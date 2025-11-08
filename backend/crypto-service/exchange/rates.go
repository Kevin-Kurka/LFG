package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"lfg/backend/crypto-service/models"
)

// RateProvider defines the interface for exchange rate providers
type RateProvider interface {
	GetRates(currencies []string) (map[string]float64, error)
	Name() string
}

// RateService manages exchange rate updates
type RateService struct {
	db                *sqlx.DB
	providers         []RateProvider
	updateInterval    time.Duration
	cache             map[string]*models.ExchangeRate
	cacheMutex        sync.RWMutex
	stopChan          chan bool
	platformFeePercent float64
}

// NewRateService creates a new exchange rate service
func NewRateService(db *sqlx.DB, updateInterval time.Duration) *RateService {
	return &RateService{
		db:                db,
		providers:         []RateProvider{NewCoinGeckoProvider(), NewCryptoCompareProvider()},
		updateInterval:    updateInterval,
		cache:             make(map[string]*models.ExchangeRate),
		stopChan:          make(chan bool),
		platformFeePercent: 1.0, // 1% platform fee
	}
}

// Start begins the rate update loop
func (s *RateService) Start() {
	// Initial update
	s.updateRates()

	// Start update loop
	ticker := time.NewTicker(s.updateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.updateRates()
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the rate update loop
func (s *RateService) Stop() {
	s.stopChan <- true
}

// updateRates fetches and updates exchange rates from providers
func (s *RateService) updateRates() {
	currencies := models.SupportedCurrencies()

	var rates map[string]float64
	var provider string

	// Try each provider until one succeeds
	for _, p := range s.providers {
		fetchedRates, err := p.GetRates(currencies)
		if err == nil && len(fetchedRates) > 0 {
			rates = fetchedRates
			provider = p.Name()
			break
		}
		fmt.Printf("Provider %s failed: %v\n", p.Name(), err)
	}

	if rates == nil {
		fmt.Println("All exchange rate providers failed")
		return
	}

	// Update database and cache
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for currency, rate := range rates {
		// Update database
		_, err := s.db.Exec(`
			INSERT INTO exchange_rates (id, currency, usd_rate, provider, last_updated)
			VALUES (gen_random_uuid(), $1, $2, $3, NOW())
			ON CONFLICT (currency) DO UPDATE
			SET usd_rate = $2, provider = $3, last_updated = NOW()
		`, currency, rate, provider)

		if err != nil {
			fmt.Printf("Failed to update rate for %s: %v\n", currency, err)
			continue
		}

		// Update cache
		s.cache[currency] = &models.ExchangeRate{
			Currency:    currency,
			USDRate:     rate,
			Provider:    provider,
			LastUpdated: time.Now(),
		}

		fmt.Printf("Updated %s rate: $%.2f (from %s)\n", currency, rate, provider)
	}
}

// GetRate returns the current exchange rate for a currency
func (s *RateService) GetRate(currency string) (float64, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	if rate, ok := s.cache[currency]; ok {
		return rate.USDRate, nil
	}

	// Fallback to database if not in cache
	var rate models.ExchangeRate
	err := s.db.Get(&rate, "SELECT * FROM exchange_rates WHERE currency = $1", currency)
	if err != nil {
		return 0, fmt.Errorf("rate not found for %s: %w", currency, err)
	}

	return rate.USDRate, nil
}

// GetAllRates returns all current exchange rates
func (s *RateService) GetAllRates() (map[string]*models.ExchangeRate, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	// Return cached rates
	if len(s.cache) > 0 {
		rates := make(map[string]*models.ExchangeRate)
		for k, v := range s.cache {
			rates[k] = v
		}
		return rates, nil
	}

	// Fallback to database
	var rates []models.ExchangeRate
	err := s.db.Select(&rates, "SELECT * FROM exchange_rates")
	if err != nil {
		return nil, err
	}

	result := make(map[string]*models.ExchangeRate)
	for i := range rates {
		result[rates[i].Currency] = &rates[i]
	}

	return result, nil
}

// ConvertCryptoToCredits converts crypto amount to platform credits
func (s *RateService) ConvertCryptoToCredits(currency string, amountCrypto float64) (float64, float64, error) {
	rate, err := s.GetRate(currency)
	if err != nil {
		return 0, 0, err
	}

	// 1 Credit = 1 USD
	credits := amountCrypto * rate
	return credits, rate, nil
}

// ConvertCreditsToC rypto converts platform credits to crypto amount
func (s *RateService) ConvertCreditsToC rypto(currency string, amountCredits float64) (float64, float64, error) {
	rate, err := s.GetRate(currency)
	if err != nil {
		return 0, 0, err
	}

	if rate == 0 {
		return 0, 0, fmt.Errorf("invalid rate for %s", currency)
	}

	// 1 Credit = 1 USD
	crypto := amountCredits / rate
	return crypto, rate, nil
}

// CalculateWithdrawalAmount calculates withdrawal amount after fees
func (s *RateService) CalculateWithdrawalAmount(currency string, amountCredits float64, networkFee float64) (amountCrypto float64, platformFee float64, totalFeeCredits float64, rate float64, err error) {
	// Calculate platform fee (1% of credits)
	platformFee = amountCredits * (s.platformFeePercent / 100.0)

	// Get exchange rate
	rate, err = s.GetRate(currency)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Convert network fee to credits
	networkFeeCredits := networkFee * rate

	// Total fee in credits
	totalFeeCredits = platformFee + networkFeeCredits

	// Calculate crypto amount after fees
	netCredits := amountCredits - totalFeeCredits
	if netCredits <= 0 {
		return 0, platformFee, totalFeeCredits, rate, fmt.Errorf("amount too small after fees")
	}

	amountCrypto = netCredits / rate

	return amountCrypto, platformFee, totalFeeCredits, rate, nil
}

// CoinGeckoProvider implements RateProvider for CoinGecko API
type CoinGeckoProvider struct {
	baseURL string
	client  *http.Client
}

// NewCoinGeckoProvider creates a new CoinGecko provider
func NewCoinGeckoProvider() *CoinGeckoProvider {
	return &CoinGeckoProvider{
		baseURL: "https://api.coingecko.com/api/v3",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the provider name
func (p *CoinGeckoProvider) Name() string {
	return "CoinGecko"
}

// GetRates fetches rates from CoinGecko
func (p *CoinGeckoProvider) GetRates(currencies []string) (map[string]float64, error) {
	// Map our currency codes to CoinGecko IDs
	coinIDs := map[string]string{
		"BTC":  "bitcoin",
		"ETH":  "ethereum",
		"USDC": "usd-coin",
		"USDT": "tether",
		"LTC":  "litecoin",
	}

	ids := []string{}
	for _, currency := range currencies {
		if id, ok := coinIDs[currency]; ok {
			ids = append(ids, id)
		}
	}

	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", p.baseURL, strings.Join(ids, ","))

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CoinGecko API error: %s", string(body))
	}

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	rates := make(map[string]float64)
	for currency, id := range coinIDs {
		if data, ok := result[id]; ok {
			if usdRate, ok := data["usd"]; ok {
				rates[currency] = usdRate
			}
		}
	}

	return rates, nil
}

// CryptoCompareProvider implements RateProvider for CryptoCompare API
type CryptoCompareProvider struct {
	baseURL string
	client  *http.Client
}

// NewCryptoCompareProvider creates a new CryptoCompare provider
func NewCryptoCompareProvider() *CryptoCompareProvider {
	return &CryptoCompareProvider{
		baseURL: "https://min-api.cryptocompare.com/data",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Name returns the provider name
func (p *CryptoCompareProvider) Name() string {
	return "CryptoCompare"
}

// GetRates fetches rates from CryptoCompare
func (p *CryptoCompareProvider) GetRates(currencies []string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/pricemulti?fsyms=%s&tsyms=USD", p.baseURL, strings.Join(currencies, ","))

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CryptoCompare API error: %s", string(body))
	}

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	rates := make(map[string]float64)
	for currency, data := range result {
		if usdRate, ok := data["USD"]; ok {
			rates[currency] = usdRate
		}
	}

	return rates, nil
}
