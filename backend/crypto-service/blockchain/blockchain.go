package blockchain

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// BlockchainClient defines the interface for blockchain operations
type BlockchainClient interface {
	GenerateWallet() (*Wallet, error)
	ValidateAddress(address string) bool
	GetBalance(address string) (float64, error)
	GetTransactions(address string, fromBlock int64) ([]Transaction, error)
	SendTransaction(fromPrivKey string, toAddress string, amount float64) (string, error)
	EstimateFee() (float64, error)
}

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	Address        string
	PrivateKey     string
	PublicKey      string
	DerivationPath string
	Mnemonic       string // Only for initial generation
}

// Transaction represents a blockchain transaction
type Transaction struct {
	Hash          string
	From          string
	To            string
	Amount        float64
	Confirmations int
	BlockNumber   int64
	Timestamp     time.Time
	Status        string
}

// BitcoinClient handles Bitcoin blockchain operations
type BitcoinClient struct {
	network    *chaincfg.Params
	apiBaseURL string
	client     *http.Client
}

// NewBitcoinClient creates a new Bitcoin client
func NewBitcoinClient(testnet bool) *BitcoinClient {
	network := &chaincfg.MainNetParams
	apiBaseURL := "https://blockchain.info"

	if testnet {
		network = &chaincfg.TestNet3Params
		apiBaseURL = "https://testnet.blockchain.info"
	}

	return &BitcoinClient{
		network:    network,
		apiBaseURL: apiBaseURL,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateWallet generates a new Bitcoin wallet using BIP39/BIP32
func (c *BitcoinClient) GenerateWallet() (*Wallet, error) {
	// Generate mnemonic
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Generate seed
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}

	// BIP44 path for Bitcoin: m/44'/0'/0'/0/0
	purpose, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinType, _ := purpose.NewChildKey(bip32.FirstHardenedChild + 0)
	account, _ := coinType.NewChildKey(bip32.FirstHardenedChild + 0)
	change, _ := account.NewChildKey(0)
	addressKey, _ := change.NewChildKey(0)

	// Convert to Bitcoin key
	privKey, pubKey := btcec.PrivKeyFromBytes(addressKey.Key)

	// Generate address
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, c.network)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	return &Wallet{
		Address:        address.EncodeAddress(),
		PrivateKey:     hex.EncodeToString(privKey.Serialize()),
		PublicKey:      hex.EncodeToString(pubKey.SerializeCompressed()),
		DerivationPath: "m/44'/0'/0'/0/0",
		Mnemonic:       mnemonic,
	}, nil
}

// ValidateAddress validates a Bitcoin address
func (c *BitcoinClient) ValidateAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, c.network)
	return err == nil
}

// GetBalance gets the balance for a Bitcoin address using blockchain.info API
func (c *BitcoinClient) GetBalance(address string) (float64, error) {
	url := fmt.Sprintf("%s/q/addressbalance/%s", c.apiBaseURL, address)

	resp, err := c.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// Balance is in satoshis
	var satoshis int64
	fmt.Sscanf(string(body), "%d", &satoshis)

	// Convert to BTC
	btc := float64(satoshis) / 100000000.0

	return btc, nil
}

// GetTransactions gets transactions for a Bitcoin address
func (c *BitcoinClient) GetTransactions(address string, fromBlock int64) ([]Transaction, error) {
	url := fmt.Sprintf("%s/rawaddr/%s", c.apiBaseURL, address)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Txs []struct {
			Hash        string `json:"hash"`
			BlockHeight int64  `json:"block_height"`
			Time        int64  `json:"time"`
			Inputs      []struct {
				PrevOut struct {
					Addr  string `json:"addr"`
					Value int64  `json:"value"`
				} `json:"prev_out"`
			} `json:"inputs"`
			Out []struct {
				Addr  string `json:"addr"`
				Value int64  `json:"value"`
			} `json:"out"`
		} `json:"txs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, tx := range result.Txs {
		// Check if this is an incoming transaction to our address
		for _, out := range tx.Out {
			if out.Addr == address {
				transactions = append(transactions, Transaction{
					Hash:        tx.Hash,
					To:          out.Addr,
					Amount:      float64(out.Value) / 100000000.0, // Convert satoshis to BTC
					BlockNumber: tx.BlockHeight,
					Timestamp:   time.Unix(tx.Time, 0),
					Status:      "confirmed",
				})
			}
		}
	}

	return transactions, nil
}

// SendTransaction sends a Bitcoin transaction (SIMPLIFIED - Production needs proper UTXO handling)
func (c *BitcoinClient) SendTransaction(fromPrivKey string, toAddress string, amount float64) (string, error) {
	// NOTE: This is a simplified placeholder
	// Production implementation would use btcd or btcwallet libraries to:
	// 1. Fetch UTXOs for the address
	// 2. Create and sign transaction
	// 3. Broadcast to network
	return "", fmt.Errorf("BTC sending not implemented - requires production blockchain integration")
}

// EstimateFee estimates the transaction fee
func (c *BitcoinClient) EstimateFee() (float64, error) {
	// Simplified fee estimation
	// Production would query fee estimation API
	return 0.0001, nil // 10,000 satoshis
}

// EthereumClient handles Ethereum blockchain operations
type EthereumClient struct {
	networkID  int64
	apiBaseURL string
	apiKey     string
	client     *http.Client
}

// NewEthereumClient creates a new Ethereum client
func NewEthereumClient(testnet bool, apiKey string) *EthereumClient {
	networkID := int64(1) // Mainnet
	apiBaseURL := "https://api.etherscan.io/api"

	if testnet {
		networkID = 5 // Goerli testnet
		apiBaseURL = "https://api-goerli.etherscan.io/api"
	}

	return &EthereumClient{
		networkID:  networkID,
		apiBaseURL: apiBaseURL,
		apiKey:     apiKey,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateWallet generates a new Ethereum wallet
func (c *EthereumClient) GenerateWallet() (*Wallet, error) {
	// Generate private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to cast public key")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &Wallet{
		Address:    address,
		PrivateKey: privateKeyHex,
		PublicKey:  publicKeyHex,
	}, nil
}

// ValidateAddress validates an Ethereum address
func (c *EthereumClient) ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// GetBalance gets the balance for an Ethereum address
func (c *EthereumClient) GetBalance(address string) (float64, error) {
	url := fmt.Sprintf("%s?module=account&action=balance&address=%s&tag=latest", c.apiBaseURL, address)
	if c.apiKey != "" {
		url += "&apikey=" + c.apiKey
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Status != "1" {
		return 0, fmt.Errorf("API error: %s", result.Message)
	}

	// Convert wei to ETH
	wei := new(big.Int)
	wei.SetString(result.Result, 10)
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
	balance, _ := eth.Float64()

	return balance, nil
}

// GetTransactions gets transactions for an Ethereum address
func (c *EthereumClient) GetTransactions(address string, fromBlock int64) ([]Transaction, error) {
	url := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=%d&sort=desc",
		c.apiBaseURL, address, fromBlock)
	if c.apiKey != "" {
		url += "&apikey=" + c.apiKey
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  []struct {
			Hash            string `json:"hash"`
			From            string `json:"from"`
			To              string `json:"to"`
			Value           string `json:"value"`
			BlockNumber     string `json:"blockNumber"`
			TimeStamp       string `json:"timeStamp"`
			Confirmations   string `json:"confirmations"`
			TxReceiptStatus string `json:"txreceipt_status"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var transactions []Transaction
	for _, tx := range result.Result {
		// Only include transactions to our address
		if strings.EqualFold(tx.To, address) {
			wei := new(big.Int)
			wei.SetString(tx.Value, 10)
			eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18))
			amount, _ := eth.Float64()

			blockNum := new(big.Int)
			blockNum.SetString(tx.BlockNumber, 10)

			timestamp := new(big.Int)
			timestamp.SetString(tx.TimeStamp, 10)

			confirmations := new(big.Int)
			confirmations.SetString(tx.Confirmations, 10)

			status := "confirmed"
			if tx.TxReceiptStatus == "0" {
				status = "failed"
			}

			transactions = append(transactions, Transaction{
				Hash:          tx.Hash,
				From:          tx.From,
				To:            tx.To,
				Amount:        amount,
				BlockNumber:   blockNum.Int64(),
				Timestamp:     time.Unix(timestamp.Int64(), 0),
				Confirmations: int(confirmations.Int64()),
				Status:        status,
			})
		}
	}

	return transactions, nil
}

// SendTransaction sends an Ethereum transaction (SIMPLIFIED)
func (c *EthereumClient) SendTransaction(fromPrivKey string, toAddress string, amount float64) (string, error) {
	// NOTE: This is a simplified placeholder
	// Production implementation would use go-ethereum client to:
	// 1. Create transaction
	// 2. Sign with private key
	// 3. Broadcast to network
	return "", fmt.Errorf("ETH sending not implemented - requires production blockchain integration")
}

// EstimateFee estimates the transaction fee
func (c *EthereumClient) EstimateFee() (float64, error) {
	// Simplified fee estimation
	// Production would query gas price and estimate gas limit
	gasPrice := 20.0    // 20 gwei
	gasLimit := 21000.0 // Standard transfer
	return (gasPrice * gasLimit) / 1e9, nil
}

// LitecoinClient handles Litecoin blockchain operations (similar to Bitcoin)
type LitecoinClient struct {
	*BitcoinClient
}

// NewLitecoinClient creates a new Litecoin client
func NewLitecoinClient(testnet bool) *LitecoinClient {
	network := &chaincfg.MainNetParams
	apiBaseURL := "https://api.blockcypher.com/v1/ltc/main"

	if testnet {
		network = &chaincfg.TestNet3Params
		apiBaseURL = "https://api.blockcypher.com/v1/ltc/test3"
	}

	return &LitecoinClient{
		BitcoinClient: &BitcoinClient{
			network:    network,
			apiBaseURL: apiBaseURL,
			client:     &http.Client{Timeout: 30 * time.Second},
		},
	}
}

// ERC20Client handles ERC20 token operations (USDC, USDT)
type ERC20Client struct {
	*EthereumClient
	contractAddress string
	tokenSymbol     string
}

// NewERC20Client creates a new ERC20 client
func NewERC20Client(tokenSymbol string, testnet bool, apiKey string) *ERC20Client {
	// Contract addresses (mainnet)
	contracts := map[string]string{
		"USDC": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		"USDT": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
	}

	return &ERC20Client{
		EthereumClient:  NewEthereumClient(testnet, apiKey),
		contractAddress: contracts[tokenSymbol],
		tokenSymbol:     tokenSymbol,
	}
}

// GetBalance gets the ERC20 token balance
func (c *ERC20Client) GetBalance(address string) (float64, error) {
	url := fmt.Sprintf("%s?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest",
		c.apiBaseURL, c.contractAddress, address)
	if c.apiKey != "" {
		url += "&apikey=" + c.apiKey
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	if result.Status != "1" {
		return 0, fmt.Errorf("API error: %s", result.Message)
	}

	// Convert to token amount (USDC/USDT have 6 decimals)
	amount := new(big.Int)
	amount.SetString(result.Result, 10)
	tokenAmount := new(big.Float).Quo(new(big.Float).SetInt(amount), big.NewFloat(1e6))
	balance, _ := tokenAmount.Float64()

	return balance, nil
}

// Helper functions

// GenerateRandomAddress generates a random address for testing
func GenerateRandomAddress(currency string) (string, error) {
	switch currency {
	case "BTC", "LTC":
		// Generate random Bitcoin-style address
		bytes := make([]byte, 20)
		rand.Read(bytes)
		hash := sha256.Sum256(bytes)
		return "1" + hex.EncodeToString(hash[:20]), nil
	case "ETH", "USDC", "USDT":
		// Generate random Ethereum address
		bytes := make([]byte, 20)
		rand.Read(bytes)
		return "0x" + hex.EncodeToString(bytes), nil
	default:
		return "", fmt.Errorf("unsupported currency: %s", currency)
	}
}
