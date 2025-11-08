package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"lfg/backend/crypto-service/blockchain"
	"lfg/backend/crypto-service/exchange"
	"lfg/backend/crypto-service/models"
)

// DepositMonitor monitors blockchain for incoming deposits
type DepositMonitor struct {
	db                *sqlx.DB
	rateService       *exchange.RateService
	scanInterval      time.Duration
	stopChan          chan bool
	blockchainClients map[string]blockchain.BlockchainClient
}

// NewDepositMonitor creates a new deposit monitor
func NewDepositMonitor(db *sqlx.DB, rateService *exchange.RateService, scanInterval time.Duration) *DepositMonitor {
	return &DepositMonitor{
		db:           db,
		rateService:  rateService,
		scanInterval: scanInterval,
		stopChan:     make(chan bool),
		blockchainClients: map[string]blockchain.BlockchainClient{
			models.CurrencyBTC:  blockchain.NewBitcoinClient(false),
			models.CurrencyETH:  blockchain.NewEthereumClient(false, ""),
			models.CurrencyUSDC: blockchain.NewERC20Client("USDC", false, ""),
			models.CurrencyUSDT: blockchain.NewERC20Client("USDT", false, ""),
			models.CurrencyLTC:  blockchain.NewLitecoinClient(false),
		},
	}
}

// Start begins monitoring for deposits
func (m *DepositMonitor) Start() {
	log.Println("Starting deposit monitor...")

	// Start monitoring for each currency
	for _, currency := range models.SupportedCurrencies() {
		go m.monitorCurrency(currency)
	}

	// Keep the monitor alive
	<-m.stopChan
	log.Println("Deposit monitor stopped")
}

// Stop stops the deposit monitor
func (m *DepositMonitor) Stop() {
	m.stopChan <- true
}

// monitorCurrency monitors deposits for a specific currency
func (m *DepositMonitor) monitorCurrency(currency string) {
	ticker := time.NewTicker(m.scanInterval)
	defer ticker.Stop()

	log.Printf("Started monitoring %s deposits\n", currency)

	for {
		select {
		case <-ticker.C:
			if err := m.scanForDeposits(currency); err != nil {
				log.Printf("Error scanning %s deposits: %v\n", currency, err)
			}
		case <-m.stopChan:
			return
		}
	}
}

// scanForDeposits scans blockchain for new deposits
func (m *DepositMonitor) scanForDeposits(currency string) error {
	// Get monitoring job
	var job models.CryptoMonitoringJob
	err := m.db.Get(&job, "SELECT * FROM crypto_monitoring_jobs WHERE currency = $1", currency)
	if err != nil {
		return fmt.Errorf("failed to get monitoring job: %w", err)
	}

	// Skip if already running
	if job.IsRunning {
		return nil
	}

	// Mark as running
	_, err = m.db.Exec("UPDATE crypto_monitoring_jobs SET is_running = true WHERE currency = $1", currency)
	if err != nil {
		return err
	}
	defer m.db.Exec("UPDATE crypto_monitoring_jobs SET is_running = false WHERE currency = $1", currency)

	// Get all active wallets for this currency
	var wallets []models.CryptoWallet
	err = m.db.Select(&wallets, "SELECT * FROM crypto_wallets WHERE currency = $1", currency)
	if err != nil {
		return fmt.Errorf("failed to get wallets: %w", err)
	}

	if len(wallets) == 0 {
		return nil // No wallets to monitor
	}

	client := m.blockchainClients[currency]

	// Check each wallet for new transactions
	for _, wallet := range wallets {
		transactions, err := client.GetTransactions(wallet.Address, job.LastScannedBlock)
		if err != nil {
			log.Printf("Error getting transactions for %s wallet %s: %v\n", currency, wallet.Address, err)
			m.updateJobError(currency, err)
			continue
		}

		// Process new transactions
		for _, tx := range transactions {
			if err := m.processDeposit(wallet, tx); err != nil {
				log.Printf("Error processing deposit %s: %v\n", tx.Hash, err)
			}
		}
	}

	// Update last scanned block (simplified - in production would track actual block heights)
	_, err = m.db.Exec(`
		UPDATE crypto_monitoring_jobs
		SET last_scanned_block = last_scanned_block + 1,
		    last_scan_at = NOW(),
		    error_count = 0,
		    last_error = NULL
		WHERE currency = $1
	`, currency)

	return err
}

// processDeposit processes a new deposit transaction
func (m *DepositMonitor) processDeposit(wallet models.CryptoWallet, tx blockchain.Transaction) error {
	// Check if deposit already exists
	var existing models.CryptoDeposit
	err := m.db.Get(&existing, "SELECT * FROM crypto_deposits WHERE tx_hash = $1", tx.Hash)
	if err == nil {
		// Deposit already exists, update confirmations
		return m.updateDepositConfirmations(&existing, tx.Confirmations)
	}

	// Convert crypto amount to credits
	credits, rate, err := m.rateService.ConvertCryptoToCredits(wallet.Currency, tx.Amount)
	if err != nil {
		return fmt.Errorf("failed to convert amount: %w", err)
	}

	// Create new deposit record
	deposit := models.CryptoDeposit{
		ID:                    uuid.New().String(),
		UserID:                wallet.UserID,
		WalletID:              wallet.ID,
		Currency:              wallet.Currency,
		AmountCrypto:          tx.Amount,
		AmountCredits:         credits,
		ExchangeRate:          rate,
		Address:               wallet.Address,
		TxHash:                tx.Hash,
		Confirmations:         tx.Confirmations,
		RequiredConfirmations: models.GetRequiredConfirmations(wallet.Currency),
		Status:                models.DepositStatusPending,
		BlockNumber:           &tx.BlockNumber,
		DetectedAt:            time.Now(),
		CreatedAt:             time.Now(),
	}

	// Determine status based on confirmations
	if tx.Confirmations > 0 {
		deposit.Status = models.DepositStatusConfirming
	}

	if tx.Confirmations >= deposit.RequiredConfirmations {
		deposit.Status = models.DepositStatusConfirmed
		now := time.Now()
		deposit.ConfirmedAt = &now
	}

	// Insert deposit record
	_, err = m.db.NamedExec(`
		INSERT INTO crypto_deposits (
			id, user_id, wallet_id, currency, amount_crypto, amount_credits, exchange_rate,
			address, tx_hash, confirmations, required_confirmations, status, block_number,
			detected_at, confirmed_at, created_at
		)
		VALUES (
			:id, :user_id, :wallet_id, :currency, :amount_crypto, :amount_credits, :exchange_rate,
			:address, :tx_hash, :confirmations, :required_confirmations, :status, :block_number,
			:detected_at, :confirmed_at, :created_at
		)
	`, deposit)

	if err != nil {
		return fmt.Errorf("failed to insert deposit: %w", err)
	}

	log.Printf("New %s deposit detected: %s - %.8f %s (%.2f credits)\n",
		wallet.Currency, tx.Hash, tx.Amount, wallet.Currency, credits)

	// If confirmed, credit user account
	if deposit.Status == models.DepositStatusConfirmed {
		return m.creditUserAccount(&deposit)
	}

	return nil
}

// updateDepositConfirmations updates the confirmation count for a deposit
func (m *DepositMonitor) updateDepositConfirmations(deposit *models.CryptoDeposit, confirmations int) error {
	// Only update if confirmations increased
	if confirmations <= deposit.Confirmations {
		return nil
	}

	deposit.Confirmations = confirmations

	// Update status if needed
	if confirmations >= deposit.RequiredConfirmations && deposit.Status != models.DepositStatusConfirmed && deposit.Status != models.DepositStatusCredited {
		deposit.Status = models.DepositStatusConfirmed
		now := time.Now()
		deposit.ConfirmedAt = &now

		_, err := m.db.Exec(`
			UPDATE crypto_deposits
			SET confirmations = $1, status = $2, confirmed_at = $3
			WHERE id = $4
		`, deposit.Confirmations, deposit.Status, deposit.ConfirmedAt, deposit.ID)

		if err != nil {
			return err
		}

		log.Printf("Deposit %s confirmed with %d confirmations\n", deposit.TxHash, confirmations)

		// Credit user account
		return m.creditUserAccount(deposit)
	}

	// Just update confirmations
	_, err := m.db.Exec(`
		UPDATE crypto_deposits
		SET confirmations = $1
		WHERE id = $2
	`, deposit.Confirmations, deposit.ID)

	return err
}

// creditUserAccount credits the user's account for a confirmed deposit
func (m *DepositMonitor) creditUserAccount(deposit *models.CryptoDeposit) error {
	// Check if already credited
	if deposit.Status == models.DepositStatusCredited {
		return nil
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user wallet balance
	_, err = tx.Exec(`
		UPDATE wallets
		SET balance_credits = balance_credits + $1
		WHERE user_id = $2
	`, deposit.AmountCredits, deposit.UserID)

	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create wallet transaction record
	_, err = tx.Exec(`
		INSERT INTO wallet_transactions (id, user_id, type, amount, balance, reference, reference_id, description, created_at)
		SELECT gen_random_uuid(), $1, 'DEPOSIT', $2, balance_credits, 'CRYPTO_DEPOSIT', $3, $4, NOW()
		FROM wallets WHERE user_id = $1
	`, deposit.UserID, deposit.AmountCredits, deposit.ID,
		fmt.Sprintf("Crypto deposit: %.8f %s", deposit.AmountCrypto, deposit.Currency))

	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	// Update deposit status
	now := time.Now()
	_, err = tx.Exec(`
		UPDATE crypto_deposits
		SET status = $1, credited_at = $2
		WHERE id = $3
	`, models.DepositStatusCredited, now, deposit.ID)

	if err != nil {
		return fmt.Errorf("failed to update deposit status: %w", err)
	}

	// Create crypto transaction audit record
	_, err = tx.Exec(`
		INSERT INTO crypto_transactions (id, user_id, wallet_id, type, currency, amount_crypto, amount_credits, exchange_rate, tx_hash, status, reference_id, reference_type, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
	`, deposit.UserID, deposit.WalletID, models.TransactionTypeDeposit, deposit.Currency,
		deposit.AmountCrypto, deposit.AmountCredits, deposit.ExchangeRate, deposit.TxHash,
		models.DepositStatusCredited, deposit.ID, "DEPOSIT")

	if err != nil {
		return fmt.Errorf("failed to create audit record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Credited %.2f credits to user %s for deposit %s\n",
		deposit.AmountCredits, deposit.UserID, deposit.TxHash)

	return nil
}

// updateJobError updates the monitoring job error state
func (m *DepositMonitor) updateJobError(currency string, err error) {
	errMsg := err.Error()
	m.db.Exec(`
		UPDATE crypto_monitoring_jobs
		SET error_count = error_count + 1,
		    last_error = $1
		WHERE currency = $2
	`, errMsg, currency)
}
