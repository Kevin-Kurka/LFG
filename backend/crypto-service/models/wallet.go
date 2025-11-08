package models

import (
	"time"
)

// CryptoWallet represents a user's cryptocurrency wallet
type CryptoWallet struct {
	ID                  string    `json:"id" db:"id"`
	UserID              string    `json:"user_id" db:"user_id"`
	Currency            string    `json:"currency" db:"currency"`
	Address             string    `json:"address" db:"address"`
	PrivateKeyEncrypted string    `json:"-" db:"private_key_encrypted"` // Never expose in JSON
	BalanceCrypto       float64   `json:"balance_crypto" db:"balance_crypto"`
	DerivationPath      string    `json:"derivation_path,omitempty" db:"derivation_path"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// CryptoDeposit represents a cryptocurrency deposit
type CryptoDeposit struct {
	ID                     string    `json:"id" db:"id"`
	UserID                 string    `json:"user_id" db:"user_id"`
	WalletID               string    `json:"wallet_id" db:"wallet_id"`
	Currency               string    `json:"currency" db:"currency"`
	AmountCrypto           float64   `json:"amount_crypto" db:"amount_crypto"`
	AmountCredits          float64   `json:"amount_credits" db:"amount_credits"`
	ExchangeRate           float64   `json:"exchange_rate" db:"exchange_rate"`
	Address                string    `json:"address" db:"address"`
	TxHash                 string    `json:"tx_hash" db:"tx_hash"`
	Confirmations          int       `json:"confirmations" db:"confirmations"`
	RequiredConfirmations  int       `json:"required_confirmations" db:"required_confirmations"`
	Status                 string    `json:"status" db:"status"`
	BlockNumber            *int64    `json:"block_number,omitempty" db:"block_number"`
	DetectedAt             time.Time `json:"detected_at" db:"detected_at"`
	ConfirmedAt            *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CreditedAt             *time.Time `json:"credited_at,omitempty" db:"credited_at"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

// CryptoWithdrawal represents a cryptocurrency withdrawal
type CryptoWithdrawal struct {
	ID               string     `json:"id" db:"id"`
	UserID           string     `json:"user_id" db:"user_id"`
	WalletID         string     `json:"wallet_id" db:"wallet_id"`
	Currency         string     `json:"currency" db:"currency"`
	AmountCredits    float64    `json:"amount_credits" db:"amount_credits"`
	AmountCrypto     float64    `json:"amount_crypto" db:"amount_crypto"`
	ExchangeRate     float64    `json:"exchange_rate" db:"exchange_rate"`
	NetworkFee       float64    `json:"network_fee" db:"network_fee"`
	PlatformFee      float64    `json:"platform_fee" db:"platform_fee"`
	TotalFeeCredits  float64    `json:"total_fee_credits" db:"total_fee_credits"`
	ToAddress        string     `json:"to_address" db:"to_address"`
	TxHash           *string    `json:"tx_hash,omitempty" db:"tx_hash"`
	Status           string     `json:"status" db:"status"`
	Confirmations    int        `json:"confirmations" db:"confirmations"`
	BlockNumber      *int64     `json:"block_number,omitempty" db:"block_number"`
	Requires2FA      bool       `json:"requires_2fa" db:"requires_2fa"`
	Verified2FA      bool       `json:"verified_2fa" db:"verified_2fa"`
	ErrorMessage     *string    `json:"error_message,omitempty" db:"error_message"`
	RequestedAt      time.Time  `json:"requested_at" db:"requested_at"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	ConfirmedAt      *time.Time `json:"confirmed_at,omitempty" db:"confirmed_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// ExchangeRate represents cryptocurrency exchange rate
type ExchangeRate struct {
	ID          string    `json:"id" db:"id"`
	Currency    string    `json:"currency" db:"currency"`
	USDRate     float64   `json:"usd_rate" db:"usd_rate"`
	Provider    string    `json:"provider" db:"provider"`
	LastUpdated time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// CryptoTransaction represents a comprehensive audit trail entry
type CryptoTransaction struct {
	ID            string                 `json:"id" db:"id"`
	UserID        string                 `json:"user_id" db:"user_id"`
	WalletID      string                 `json:"wallet_id" db:"wallet_id"`
	Type          string                 `json:"type" db:"type"`
	Currency      string                 `json:"currency" db:"currency"`
	AmountCrypto  float64                `json:"amount_crypto" db:"amount_crypto"`
	AmountCredits *float64               `json:"amount_credits,omitempty" db:"amount_credits"`
	ExchangeRate  *float64               `json:"exchange_rate,omitempty" db:"exchange_rate"`
	TxHash        *string                `json:"tx_hash,omitempty" db:"tx_hash"`
	Status        string                 `json:"status" db:"status"`
	ReferenceID   *string                `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType *string                `json:"reference_type,omitempty" db:"reference_type"`
	Metadata      map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
}

// CryptoMonitoringJob represents a blockchain monitoring job
type CryptoMonitoringJob struct {
	ID               string     `json:"id" db:"id"`
	Currency         string     `json:"currency" db:"currency"`
	LastScannedBlock int64      `json:"last_scanned_block" db:"last_scanned_block"`
	LastScanAt       time.Time  `json:"last_scan_at" db:"last_scan_at"`
	IsRunning        bool       `json:"is_running" db:"is_running"`
	ErrorCount       int        `json:"error_count" db:"error_count"`
	LastError        *string    `json:"last_error,omitempty" db:"last_error"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// CryptoConfig represents system configuration
type CryptoConfig struct {
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description *string   `json:"description,omitempty" db:"description"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Currency constants
const (
	CurrencyBTC  = "BTC"
	CurrencyETH  = "ETH"
	CurrencyUSDC = "USDC"
	CurrencyUSDT = "USDT"
	CurrencyLTC  = "LTC"
)

// Deposit status constants
const (
	DepositStatusPending    = "PENDING"
	DepositStatusConfirming = "CONFIRMING"
	DepositStatusConfirmed  = "CONFIRMED"
	DepositStatusCredited   = "CREDITED"
	DepositStatusFailed     = "FAILED"
)

// Withdrawal status constants
const (
	WithdrawalStatusPending    = "PENDING"
	WithdrawalStatusProcessing = "PROCESSING"
	WithdrawalStatusSent       = "SENT"
	WithdrawalStatusConfirmed  = "CONFIRMED"
	WithdrawalStatusFailed     = "FAILED"
	WithdrawalStatusCancelled  = "CANCELLED"
)

// Transaction type constants
const (
	TransactionTypeDeposit    = "DEPOSIT"
	TransactionTypeWithdrawal = "WITHDRAWAL"
	TransactionTypeFee        = "FEE"
)

// SupportedCurrencies returns list of all supported cryptocurrencies
func SupportedCurrencies() []string {
	return []string{CurrencyBTC, CurrencyETH, CurrencyUSDC, CurrencyUSDT, CurrencyLTC}
}

// IsSupportedCurrency checks if a currency is supported
func IsSupportedCurrency(currency string) bool {
	for _, c := range SupportedCurrencies() {
		if c == currency {
			return true
		}
	}
	return false
}

// GetRequiredConfirmations returns the required confirmations for a currency
func GetRequiredConfirmations(currency string) int {
	switch currency {
	case CurrencyBTC:
		return 3
	case CurrencyETH, CurrencyUSDC, CurrencyUSDT:
		return 12
	case CurrencyLTC:
		return 6
	default:
		return 6
	}
}
