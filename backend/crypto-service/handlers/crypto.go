package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"lfg/backend/common/crypto"
	"lfg/backend/common/response"
	"lfg/backend/crypto-service/blockchain"
	"lfg/backend/crypto-service/exchange"
	"lfg/backend/crypto-service/models"
)

// CryptoHandler handles cryptocurrency operations
type CryptoHandler struct {
	db          *sqlx.DB
	rateService *exchange.RateService
	btcClient   blockchain.BlockchainClient
	ethClient   blockchain.BlockchainClient
	ltcClient   blockchain.BlockchainClient
	usdcClient  blockchain.BlockchainClient
	usdtClient  blockchain.BlockchainClient
}

// NewCryptoHandler creates a new crypto handler
func NewCryptoHandler(db *sqlx.DB, rateService *exchange.RateService) *CryptoHandler {
	// Initialize blockchain clients
	// For MVP, using testnet=false but APIs may need keys in production
	return &CryptoHandler{
		db:          db,
		rateService: rateService,
		btcClient:   blockchain.NewBitcoinClient(false),
		ethClient:   blockchain.NewEthereumClient(false, ""), // Add API key in production
		ltcClient:   blockchain.NewLitecoinClient(false),
		usdcClient:  blockchain.NewERC20Client("USDC", false, ""),
		usdtClient:  blockchain.NewERC20Client("USDT", false, ""),
	}
}

// GetWallets returns all crypto wallets for a user
func (h *CryptoHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var wallets []models.CryptoWallet
	err := h.db.Select(&wallets, "SELECT * FROM crypto_wallets WHERE user_id = $1 ORDER BY currency", userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch wallets", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"wallets": wallets,
	})
}

// CreateWallet generates a new crypto wallet for a user
func (h *CryptoHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	currency := vars["currency"]

	// Validate currency
	if !models.IsSupportedCurrency(currency) {
		response.Error(w, http.StatusBadRequest, "Unsupported currency", nil)
		return
	}

	// Check if wallet already exists
	var existing models.CryptoWallet
	err := h.db.Get(&existing, "SELECT * FROM crypto_wallets WHERE user_id = $1 AND currency = $2", userID, currency)
	if err == nil {
		response.Error(w, http.StatusConflict, "Wallet already exists for this currency", nil)
		return
	}

	// Generate new wallet using appropriate blockchain client
	var wallet *blockchain.Wallet
	switch currency {
	case models.CurrencyBTC:
		wallet, err = h.btcClient.GenerateWallet()
	case models.CurrencyETH:
		wallet, err = h.ethClient.GenerateWallet()
	case models.CurrencyUSDC:
		wallet, err = h.usdcClient.GenerateWallet()
	case models.CurrencyUSDT:
		wallet, err = h.usdtClient.GenerateWallet()
	case models.CurrencyLTC:
		wallet, err = h.ltcClient.GenerateWallet()
	default:
		response.Error(w, http.StatusBadRequest, "Unsupported currency", nil)
		return
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate wallet", err)
		return
	}

	// Encrypt private key
	encryptedPrivKey, err := crypto.Encrypt(wallet.PrivateKey)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to encrypt private key", err)
		return
	}

	// Save to database
	cryptoWallet := models.CryptoWallet{
		ID:                  uuid.New().String(),
		UserID:              userID,
		Currency:            currency,
		Address:             wallet.Address,
		PrivateKeyEncrypted: encryptedPrivKey,
		BalanceCrypto:       0,
		DerivationPath:      wallet.DerivationPath,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	_, err = h.db.NamedExec(`
		INSERT INTO crypto_wallets (id, user_id, currency, address, private_key_encrypted, balance_crypto, derivation_path, created_at, updated_at)
		VALUES (:id, :user_id, :currency, :address, :private_key_encrypted, :balance_crypto, :derivation_path, :created_at, :updated_at)
	`, cryptoWallet)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to save wallet", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"wallet": cryptoWallet,
		// Include mnemonic only on initial creation (user should save it)
		"mnemonic": wallet.Mnemonic,
	})
}

// GetDeposits returns deposit history for a user
func (h *CryptoHandler) GetDeposits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var deposits []models.CryptoDeposit
	err := h.db.Select(&deposits, `
		SELECT * FROM crypto_deposits
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch deposits", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"deposits": deposits,
	})
}

// GetPendingDeposits returns pending deposits for a user
func (h *CryptoHandler) GetPendingDeposits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var deposits []models.CryptoDeposit
	err := h.db.Select(&deposits, `
		SELECT * FROM crypto_deposits
		WHERE user_id = $1 AND status IN ($2, $3)
		ORDER BY created_at DESC
	`, userID, models.DepositStatusPending, models.DepositStatusConfirming)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch pending deposits", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"deposits": deposits,
	})
}

// SimulateDeposit simulates a deposit for testing (development only)
func (h *CryptoHandler) SimulateDeposit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Currency     string  `json:"currency"`
		AmountCrypto float64 `json:"amount_crypto"`
		TxHash       string  `json:"tx_hash,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Validate currency
	if !models.IsSupportedCurrency(req.Currency) {
		response.Error(w, http.StatusBadRequest, "Unsupported currency", nil)
		return
	}

	// Get user's wallet
	var wallet models.CryptoWallet
	err := h.db.Get(&wallet, "SELECT * FROM crypto_wallets WHERE user_id = $1 AND currency = $2", userID, req.Currency)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Wallet not found for this currency", err)
		return
	}

	// Get exchange rate
	credits, rate, err := h.rateService.ConvertCryptoToCredits(req.Currency, req.AmountCrypto)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to convert amount", err)
		return
	}

	// Generate random tx hash if not provided
	txHash := req.TxHash
	if txHash == "" {
		txHash = uuid.New().String()
	}

	// Create deposit record
	deposit := models.CryptoDeposit{
		ID:                    uuid.New().String(),
		UserID:                userID,
		WalletID:              wallet.ID,
		Currency:              req.Currency,
		AmountCrypto:          req.AmountCrypto,
		AmountCredits:         credits,
		ExchangeRate:          rate,
		Address:               wallet.Address,
		TxHash:                txHash,
		Confirmations:         models.GetRequiredConfirmations(req.Currency),
		RequiredConfirmations: models.GetRequiredConfirmations(req.Currency),
		Status:                models.DepositStatusConfirmed,
		DetectedAt:            time.Now(),
		CreatedAt:             time.Now(),
	}

	now := time.Now()
	deposit.ConfirmedAt = &now

	_, err = h.db.NamedExec(`
		INSERT INTO crypto_deposits (
			id, user_id, wallet_id, currency, amount_crypto, amount_credits, exchange_rate,
			address, tx_hash, confirmations, required_confirmations, status, detected_at,
			confirmed_at, created_at
		)
		VALUES (
			:id, :user_id, :wallet_id, :currency, :amount_crypto, :amount_credits, :exchange_rate,
			:address, :tx_hash, :confirmations, :required_confirmations, :status, :detected_at,
			:confirmed_at, :created_at
		)
	`, deposit)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create deposit", err)
		return
	}

	// Credit user's account
	err = h.creditUserAccount(userID, deposit.ID, credits, req.Currency, txHash)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to credit account", err)
		return
	}

	// Update deposit status
	deposit.Status = models.DepositStatusCredited
	deposit.CreditedAt = &now

	_, err = h.db.Exec(`
		UPDATE crypto_deposits
		SET status = $1, credited_at = $2
		WHERE id = $3
	`, deposit.Status, deposit.CreditedAt, deposit.ID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update deposit status", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"deposit": deposit,
		"message": "Deposit simulated successfully",
	})
}

// RequestWithdrawal processes a withdrawal request
func (h *CryptoHandler) RequestWithdrawal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Currency      string  `json:"currency"`
		AmountCredits float64 `json:"amount_credits"`
		ToAddress     string  `json:"to_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Validate currency
	if !models.IsSupportedCurrency(req.Currency) {
		response.Error(w, http.StatusBadRequest, "Unsupported currency", nil)
		return
	}

	// Validate address
	var client blockchain.BlockchainClient
	switch req.Currency {
	case models.CurrencyBTC:
		client = h.btcClient
	case models.CurrencyETH:
		client = h.ethClient
	case models.CurrencyUSDC:
		client = h.usdcClient
	case models.CurrencyUSDT:
		client = h.usdtClient
	case models.CurrencyLTC:
		client = h.ltcClient
	}

	if !client.ValidateAddress(req.ToAddress) {
		response.Error(w, http.StatusBadRequest, "Invalid destination address", nil)
		return
	}

	// Get user's wallet
	var wallet models.CryptoWallet
	err := h.db.Get(&wallet, "SELECT * FROM crypto_wallets WHERE user_id = $1 AND currency = $2", userID, req.Currency)
	if err != nil {
		response.Error(w, http.StatusNotFound, "Wallet not found for this currency", err)
		return
	}

	// Check minimum withdrawal
	minWithdrawal, err := h.getConfigFloat(fmt.Sprintf("%s_min_withdrawal", strings.ToLower(req.Currency)))
	if err == nil && req.AmountCredits < minWithdrawal {
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("Minimum withdrawal is %.4f credits", minWithdrawal), nil)
		return
	}

	// Check user balance
	var userWallet struct {
		BalanceCredits float64 `db:"balance_credits"`
	}
	err = h.db.Get(&userWallet, "SELECT balance_credits FROM wallets WHERE user_id = $1", userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch user balance", err)
		return
	}

	if userWallet.BalanceCredits < req.AmountCredits {
		response.Error(w, http.StatusBadRequest, "Insufficient balance", nil)
		return
	}

	// Estimate network fee
	networkFee, err := client.EstimateFee()
	if err != nil {
		networkFee = 0 // Use default if estimation fails
	}

	// Calculate withdrawal amount after fees
	amountCrypto, platformFee, totalFeeCredits, rate, err := h.rateService.CalculateWithdrawalAmount(
		req.Currency,
		req.AmountCredits,
		networkFee,
	)

	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Check if user has enough balance including fees
	if userWallet.BalanceCredits < (req.AmountCredits + totalFeeCredits) {
		response.Error(w, http.StatusBadRequest, "Insufficient balance for withdrawal including fees", nil)
		return
	}

	// Create withdrawal record
	withdrawal := models.CryptoWithdrawal{
		ID:              uuid.New().String(),
		UserID:          userID,
		WalletID:        wallet.ID,
		Currency:        req.Currency,
		AmountCredits:   req.AmountCredits,
		AmountCrypto:    amountCrypto,
		ExchangeRate:    rate,
		NetworkFee:      networkFee,
		PlatformFee:     platformFee,
		TotalFeeCredits: totalFeeCredits,
		ToAddress:       req.ToAddress,
		Status:          models.WithdrawalStatusPending,
		Requires2FA:     true,
		Verified2FA:     false, // Would verify 2FA in production
		RequestedAt:     time.Now(),
		CreatedAt:       time.Now(),
	}

	_, err = h.db.NamedExec(`
		INSERT INTO crypto_withdrawals (
			id, user_id, wallet_id, currency, amount_credits, amount_crypto, exchange_rate,
			network_fee, platform_fee, total_fee_credits, to_address, status,
			requires_2fa, verified_2fa, requested_at, created_at
		)
		VALUES (
			:id, :user_id, :wallet_id, :currency, :amount_credits, :amount_crypto, :exchange_rate,
			:network_fee, :platform_fee, :total_fee_credits, :to_address, :status,
			:requires_2fa, :verified_2fa, :requested_at, :created_at
		)
	`, withdrawal)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create withdrawal", err)
		return
	}

	// Debit user's account
	err = h.debitUserAccount(userID, withdrawal.ID, req.AmountCredits+totalFeeCredits, req.Currency)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to debit account", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"withdrawal": withdrawal,
		"message":    "Withdrawal request created successfully",
	})
}

// GetWithdrawals returns withdrawal history for a user
func (h *CryptoHandler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var withdrawals []models.CryptoWithdrawal
	err := h.db.Select(&withdrawals, `
		SELECT * FROM crypto_withdrawals
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch withdrawals", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"withdrawals": withdrawals,
	})
}

// GetExchangeRates returns current exchange rates
func (h *CryptoHandler) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	rates, err := h.rateService.GetAllRates()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch exchange rates", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"rates": rates,
	})
}

// ConvertAmount converts between crypto and credits
func (h *CryptoHandler) ConvertAmount(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	amountStr := r.URL.Query().Get("amount")
	direction := r.URL.Query().Get("direction") // "crypto_to_credits" or "credits_to_crypto"

	if currency == "" || amountStr == "" || direction == "" {
		response.Error(w, http.StatusBadRequest, "Missing required parameters", nil)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid amount", err)
		return
	}

	var result float64
	var rate float64

	if direction == "crypto_to_credits" {
		result, rate, err = h.rateService.ConvertCryptoToCredits(currency, amount)
	} else if direction == "credits_to_crypto" {
		result, rate, err = h.rateService.ConvertCreditsToC rypto(currency, amount)
	} else {
		response.Error(w, http.StatusBadRequest, "Invalid direction", nil)
		return
	}

	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Conversion failed", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"amount":        amount,
		"result":        result,
		"exchange_rate": rate,
		"currency":      currency,
		"direction":     direction,
	})
}

// Helper functions

func (h *CryptoHandler) creditUserAccount(userID string, depositID string, credits float64, currency string, txHash string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user wallet balance
	_, err = tx.Exec(`
		UPDATE wallets
		SET balance_credits = balance_credits + $1
		WHERE user_id = $2
	`, credits, userID)

	if err != nil {
		return err
	}

	// Create wallet transaction record
	_, err = tx.Exec(`
		INSERT INTO wallet_transactions (id, user_id, type, amount, balance, reference, reference_id, description, created_at)
		SELECT gen_random_uuid(), $1, 'DEPOSIT', $2, balance_credits, $3, $4, $5, NOW()
		FROM wallets WHERE user_id = $1
	`, userID, credits, "CRYPTO_DEPOSIT", depositID, fmt.Sprintf("Crypto deposit: %s %s", strconv.FormatFloat(credits, 'f', 2, 64), currency))

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (h *CryptoHandler) debitUserAccount(userID string, withdrawalID string, credits float64, currency string) error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update user wallet balance
	result, err := tx.Exec(`
		UPDATE wallets
		SET balance_credits = balance_credits - $1
		WHERE user_id = $2 AND balance_credits >= $1
	`, credits, userID)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("insufficient balance")
	}

	// Create wallet transaction record
	_, err = tx.Exec(`
		INSERT INTO wallet_transactions (id, user_id, type, amount, balance, reference, reference_id, description, created_at)
		SELECT gen_random_uuid(), $1, 'WITHDRAWAL', $2, balance_credits, $3, $4, $5, NOW()
		FROM wallets WHERE user_id = $1
	`, userID, -credits, "CRYPTO_WITHDRAWAL", withdrawalID, fmt.Sprintf("Crypto withdrawal: %s %s", strconv.FormatFloat(credits, 'f', 2, 64), currency))

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (h *CryptoHandler) getConfigFloat(key string) (float64, error) {
	var config models.CryptoConfig
	err := h.db.Get(&config, "SELECT * FROM crypto_config WHERE key = $1", key)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(config.Value, 64)
}

// Import required package
import "strings"
