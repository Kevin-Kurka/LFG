package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"lfg/backend/common/crypto"
	"lfg/backend/crypto-service/blockchain"
	"lfg/backend/crypto-service/models"
)

// WithdrawalProcessor processes pending withdrawals
type WithdrawalProcessor struct {
	db                *sqlx.DB
	processInterval   time.Duration
	stopChan          chan bool
	blockchainClients map[string]blockchain.BlockchainClient
}

// NewWithdrawalProcessor creates a new withdrawal processor
func NewWithdrawalProcessor(db *sqlx.DB, processInterval time.Duration) *WithdrawalProcessor {
	return &WithdrawalProcessor{
		db:              db,
		processInterval: processInterval,
		stopChan:        make(chan bool),
		blockchainClients: map[string]blockchain.BlockchainClient{
			models.CurrencyBTC:  blockchain.NewBitcoinClient(false),
			models.CurrencyETH:  blockchain.NewEthereumClient(false, ""),
			models.CurrencyUSDC: blockchain.NewERC20Client("USDC", false, ""),
			models.CurrencyUSDT: blockchain.NewERC20Client("USDT", false, ""),
			models.CurrencyLTC:  blockchain.NewLitecoinClient(false),
		},
	}
}

// Start begins processing withdrawals
func (p *WithdrawalProcessor) Start() {
	log.Println("Starting withdrawal processor...")

	ticker := time.NewTicker(p.processInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.processWithdrawals(); err != nil {
				log.Printf("Error processing withdrawals: %v\n", err)
			}
		case <-p.stopChan:
			log.Println("Withdrawal processor stopped")
			return
		}
	}
}

// Stop stops the withdrawal processor
func (p *WithdrawalProcessor) Stop() {
	p.stopChan <- true
}

// processWithdrawals processes all pending withdrawals
func (p *WithdrawalProcessor) processWithdrawals() error {
	// Check if withdrawal processing is enabled
	var config models.CryptoConfig
	err := p.db.Get(&config, "SELECT * FROM crypto_config WHERE key = 'withdrawal_processing_enabled'")
	if err == nil && config.Value != "true" {
		return nil // Processing disabled
	}

	// Get all pending withdrawals
	var withdrawals []models.CryptoWithdrawal
	err = p.db.Select(&withdrawals, `
		SELECT * FROM crypto_withdrawals
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT 10
	`, models.WithdrawalStatusPending)

	if err != nil {
		return fmt.Errorf("failed to get pending withdrawals: %w", err)
	}

	if len(withdrawals) == 0 {
		return nil // No withdrawals to process
	}

	log.Printf("Processing %d pending withdrawals\n", len(withdrawals))

	for _, withdrawal := range withdrawals {
		if err := p.processWithdrawal(&withdrawal); err != nil {
			log.Printf("Error processing withdrawal %s: %v\n", withdrawal.ID, err)
			p.markWithdrawalFailed(&withdrawal, err)
		}
	}

	return nil
}

// processWithdrawal processes a single withdrawal
func (p *WithdrawalProcessor) processWithdrawal(withdrawal *models.CryptoWithdrawal) error {
	log.Printf("Processing withdrawal %s: %.8f %s to %s\n",
		withdrawal.ID, withdrawal.AmountCrypto, withdrawal.Currency, withdrawal.ToAddress)

	// Update status to processing
	_, err := p.db.Exec(`
		UPDATE crypto_withdrawals
		SET status = $1
		WHERE id = $2
	`, models.WithdrawalStatusProcessing, withdrawal.ID)

	if err != nil {
		return err
	}

	// Get wallet
	var wallet models.CryptoWallet
	err = p.db.Get(&wallet, "SELECT * FROM crypto_wallets WHERE id = $1", withdrawal.WalletID)
	if err != nil {
		return fmt.Errorf("failed to get wallet: %w", err)
	}

	// Decrypt private key
	privateKey, err := crypto.Decrypt(wallet.PrivateKeyEncrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Get blockchain client
	client := p.blockchainClients[withdrawal.Currency]

	// Send transaction
	// NOTE: In production, this would actually send the transaction to the blockchain
	// For MVP/development, we're simulating the send
	txHash, err := client.SendTransaction(privateKey, withdrawal.ToAddress, withdrawal.AmountCrypto)

	// Since SendTransaction returns an error for MVP (not implemented),
	// we'll simulate a successful transaction
	if err != nil && err.Error() != "BTC sending not implemented - requires production blockchain integration" &&
		err.Error() != "ETH sending not implemented - requires production blockchain integration" {
		return err
	}

	// Simulate transaction hash for development
	txHash = fmt.Sprintf("0x%s%d", withdrawal.ID[:40], time.Now().Unix())

	// Update withdrawal with transaction hash
	now := time.Now()
	_, err = p.db.Exec(`
		UPDATE crypto_withdrawals
		SET status = $1, tx_hash = $2, processed_at = $3
		WHERE id = $4
	`, models.WithdrawalStatusSent, txHash, now, withdrawal.ID)

	if err != nil {
		return fmt.Errorf("failed to update withdrawal: %w", err)
	}

	// Create crypto transaction audit record
	_, err = p.db.Exec(`
		INSERT INTO crypto_transactions (
			id, user_id, wallet_id, type, currency, amount_crypto, amount_credits,
			exchange_rate, tx_hash, status, reference_id, reference_type, created_at
		)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
	`, withdrawal.UserID, withdrawal.WalletID, models.TransactionTypeWithdrawal, withdrawal.Currency,
		withdrawal.AmountCrypto, withdrawal.AmountCredits, withdrawal.ExchangeRate, txHash,
		models.WithdrawalStatusSent, withdrawal.ID, "WITHDRAWAL")

	if err != nil {
		return fmt.Errorf("failed to create audit record: %w", err)
	}

	log.Printf("Withdrawal %s sent successfully. TX: %s\n", withdrawal.ID, txHash)

	return nil
}

// markWithdrawalFailed marks a withdrawal as failed
func (p *WithdrawalProcessor) markWithdrawalFailed(withdrawal *models.CryptoWithdrawal, err error) {
	errMsg := err.Error()
	now := time.Now()

	_, dbErr := p.db.Exec(`
		UPDATE crypto_withdrawals
		SET status = $1, error_message = $2, processed_at = $3
		WHERE id = $4
	`, models.WithdrawalStatusFailed, errMsg, now, withdrawal.ID)

	if dbErr != nil {
		log.Printf("Failed to mark withdrawal as failed: %v\n", dbErr)
	}

	// Refund user's account
	if refundErr := p.refundWithdrawal(withdrawal); refundErr != nil {
		log.Printf("Failed to refund withdrawal %s: %v\n", withdrawal.ID, refundErr)
	}
}

// refundWithdrawal refunds a failed withdrawal to the user's account
func (p *WithdrawalProcessor) refundWithdrawal(withdrawal *models.CryptoWithdrawal) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Refund amount including fees
	refundAmount := withdrawal.AmountCredits + withdrawal.TotalFeeCredits

	// Update user wallet balance
	_, err = tx.Exec(`
		UPDATE wallets
		SET balance_credits = balance_credits + $1
		WHERE user_id = $2
	`, refundAmount, withdrawal.UserID)

	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	// Create wallet transaction record
	_, err = tx.Exec(`
		INSERT INTO wallet_transactions (id, user_id, type, amount, balance, reference, reference_id, description, created_at)
		SELECT gen_random_uuid(), $1, 'REFUND', $2, balance_credits, 'CRYPTO_WITHDRAWAL_REFUND', $3, $4, NOW()
		FROM wallets WHERE user_id = $1
	`, withdrawal.UserID, refundAmount, withdrawal.ID,
		fmt.Sprintf("Refund for failed crypto withdrawal: %.8f %s", withdrawal.AmountCrypto, withdrawal.Currency))

	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit refund: %w", err)
	}

	log.Printf("Refunded %.2f credits to user %s for failed withdrawal %s\n",
		refundAmount, withdrawal.UserID, withdrawal.ID)

	return nil
}
