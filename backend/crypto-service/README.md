# LFG Crypto Service

Complete cryptocurrency deposit and withdrawal system for the LFG platform.

## Features

### Supported Cryptocurrencies
- **Bitcoin (BTC)** - 3 confirmations required
- **Ethereum (ETH)** - 12 confirmations required
- **USD Coin (USDC)** - 12 confirmations required (ERC-20)
- **Tether (USDT)** - 12 confirmations required (ERC-20)
- **Litecoin (LTC)** - 6 confirmations required

### Core Functionality

1. **Wallet Management**
   - HD wallet generation (BIP39/BIP44)
   - Secure private key encryption (AES-256-GCM)
   - Unique address per user per currency
   - Balance tracking

2. **Deposit System**
   - Automatic blockchain monitoring (30-second intervals)
   - Transaction detection and confirmation tracking
   - Automatic credit to user account when confirmed
   - Real-time exchange rate conversion
   - Comprehensive audit trail

3. **Withdrawal System**
   - Secure withdrawal requests
   - Fee calculation (1% platform fee + network fees)
   - Address validation
   - Automatic transaction broadcasting
   - Status tracking (pending, processing, sent, confirmed)
   - Refund on failure

4. **Exchange Rate Management**
   - Multiple provider support (CoinGecko, CryptoCompare)
   - Automatic rate updates (60-second intervals)
   - Fallback mechanism if primary provider fails
   - Rate caching for performance

5. **Security Features**
   - Encrypted private key storage
   - JWT authentication on all endpoints
   - Rate limiting
   - 2FA support (simulated for MVP)
   - Minimum deposit/withdrawal limits
   - Transaction audit logging

## Architecture

### Database Schema

#### crypto_wallets
- Stores user cryptocurrency wallets
- Encrypted private keys
- HD wallet derivation paths
- Balance tracking

#### crypto_deposits
- Tracks all incoming deposits
- Confirmation counting
- Status tracking (pending → confirming → confirmed → credited)
- Exchange rate snapshot

#### crypto_withdrawals
- Manages withdrawal requests
- Fee calculation and tracking
- Destination address validation
- Transaction hash tracking

#### exchange_rates
- Current exchange rates for all currencies
- Provider tracking
- Last update timestamps

#### crypto_transactions
- Comprehensive audit trail
- All crypto operations logged
- Links to deposits/withdrawals

#### crypto_monitoring_jobs
- Blockchain scanning job status
- Last scanned block tracking
- Error monitoring

#### crypto_config
- System configuration
- Minimum amounts
- Fee settings
- Confirmation requirements

### Services

#### Exchange Rate Service (`exchange/rates.go`)
- Fetches rates from multiple providers
- Automatic failover
- Rate caching
- Conversion calculations
- Fee calculations

#### Blockchain Integration (`blockchain/blockchain.go`)
- Wallet generation for all supported currencies
- Address validation
- Balance checking (via APIs)
- Transaction fetching
- Transaction signing and broadcasting

**Blockchain Clients:**
- `BitcoinClient` - Bitcoin operations
- `EthereumClient` - Ethereum operations
- `ERC20Client` - USDC/USDT operations
- `LitecoinClient` - Litecoin operations

#### Deposit Monitor Worker (`workers/deposit_monitor.go`)
- Scans blockchains for new transactions
- Detects deposits to user addresses
- Tracks confirmations
- Credits user accounts when confirmed
- Runs continuously in background

#### Withdrawal Processor Worker (`workers/withdrawal_processor.go`)
- Processes pending withdrawals
- Signs and broadcasts transactions
- Monitors transaction status
- Handles failures and refunds
- Runs continuously in background

## API Endpoints

### User Endpoints (Require JWT Auth)

#### GET /crypto/wallets
Get all cryptocurrency wallets for authenticated user.

**Response:**
```json
{
  "wallets": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "currency": "BTC",
      "address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
      "balance_crypto": 0.5,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /crypto/wallets/{currency}
Create a new wallet for specified currency.

**Parameters:**
- `currency` - BTC, ETH, USDC, USDT, or LTC

**Response:**
```json
{
  "wallet": { /* wallet object */ },
  "mnemonic": "word1 word2 ... word12"
}
```

**Note:** Mnemonic is only returned once. User must save it securely.

#### GET /crypto/deposits
Get deposit history for authenticated user.

**Response:**
```json
{
  "deposits": [
    {
      "id": "uuid",
      "currency": "BTC",
      "amount_crypto": 0.1,
      "amount_credits": 5000.00,
      "exchange_rate": 50000.00,
      "tx_hash": "abc123...",
      "confirmations": 3,
      "required_confirmations": 3,
      "status": "CREDITED",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### GET /crypto/deposits/pending
Get pending deposits for authenticated user.

#### POST /crypto/deposits/simulate
Simulate a deposit for testing (development only).

**Request:**
```json
{
  "currency": "BTC",
  "amount_crypto": 0.1,
  "tx_hash": "optional-custom-hash"
}
```

#### POST /crypto/withdraw
Request a cryptocurrency withdrawal.

**Request:**
```json
{
  "currency": "BTC",
  "amount_credits": 5000.00,
  "to_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
}
```

**Response:**
```json
{
  "withdrawal": {
    "id": "uuid",
    "currency": "BTC",
    "amount_credits": 5000.00,
    "amount_crypto": 0.099,
    "exchange_rate": 50000.00,
    "network_fee": 0.0001,
    "platform_fee": 50.00,
    "total_fee_credits": 55.00,
    "to_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "status": "PENDING",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /crypto/withdrawals
Get withdrawal history for authenticated user.

#### GET /crypto/rates
Get current exchange rates for all cryptocurrencies.

**Response:**
```json
{
  "rates": {
    "BTC": {
      "currency": "BTC",
      "usd_rate": 50000.00,
      "provider": "CoinGecko",
      "last_updated": "2024-01-01T00:00:00Z"
    },
    ...
  }
}
```

#### GET /crypto/convert
Calculate conversion between crypto and credits.

**Query Parameters:**
- `currency` - BTC, ETH, etc.
- `amount` - Amount to convert
- `direction` - "crypto_to_credits" or "credits_to_crypto"

**Response:**
```json
{
  "amount": 0.1,
  "result": 5000.00,
  "exchange_rate": 50000.00,
  "currency": "BTC",
  "direction": "crypto_to_credits"
}
```

## Configuration

### Environment Variables

```bash
DATABASE_URL=postgresql://lfguser:lfgpassword@postgres:5432/lfg?sslmode=disable
ENCRYPTION_KEY=your-32-byte-encryption-key-here
PORT=8087
```

### Database Configuration

Configuration stored in `crypto_config` table:

| Key | Default | Description |
|-----|---------|-------------|
| btc_min_deposit | 0.001 | Minimum BTC deposit |
| eth_min_deposit | 0.01 | Minimum ETH deposit |
| usdc_min_deposit | 10 | Minimum USDC deposit |
| usdt_min_deposit | 10 | Minimum USDT deposit |
| ltc_min_deposit | 0.1 | Minimum LTC deposit |
| btc_min_withdrawal | 0.001 | Minimum BTC withdrawal |
| eth_min_withdrawal | 0.01 | Minimum ETH withdrawal |
| usdc_min_withdrawal | 10 | Minimum USDC withdrawal |
| usdt_min_withdrawal | 10 | Minimum USDT withdrawal |
| ltc_min_withdrawal | 0.1 | Minimum LTC withdrawal |
| platform_fee_percentage | 1.0 | Platform fee percentage |
| btc_confirmations_required | 3 | Required confirmations for BTC |
| eth_confirmations_required | 12 | Required confirmations for ETH |
| usdc_confirmations_required | 12 | Required confirmations for USDC |
| usdt_confirmations_required | 12 | Required confirmations for USDT |
| ltc_confirmations_required | 6 | Required confirmations for LTC |
| deposit_scan_interval_seconds | 30 | Deposit scanning interval |
| withdrawal_processing_enabled | true | Enable/disable withdrawals |

## Development vs Production

### MVP/Development Mode
- Uses blockchain APIs (Blockchain.info, Etherscan)
- Simulated transaction broadcasting
- No full node required
- CoinGecko free API (rate limited)
- Suitable for testing and development

### Production Requirements

1. **Blockchain Infrastructure**
   - Run full nodes or use premium API services (Infura, Alchemy)
   - Implement proper UTXO management for Bitcoin
   - Gas price optimization for Ethereum
   - Transaction fee estimation

2. **Security**
   - Hardware security modules (HSM) for key storage
   - Multi-signature hot wallets
   - Cold wallet integration for large amounts
   - Proper 2FA implementation
   - IP whitelisting for admin operations

3. **Monitoring**
   - Blockchain reorg handling
   - Failed transaction retry logic
   - Alert system for stuck transactions
   - Balance reconciliation
   - Audit logging

4. **APIs**
   - Premium API subscriptions
   - Multiple provider redundancy
   - WebSocket connections for real-time updates
   - Rate limit handling

5. **Compliance**
   - KYC/AML integration
   - Transaction limits
   - Suspicious activity monitoring
   - Regulatory reporting

## Building and Running

### Local Development

```bash
# Install dependencies
cd /home/user/LFG/backend/crypto-service
go mod download

# Run the service
go run main.go
```

### Docker

```bash
# Build
docker build -t lfg-crypto-service .

# Run
docker run -p 8087:8087 \
  -e DATABASE_URL=postgresql://... \
  -e ENCRYPTION_KEY=your-key \
  lfg-crypto-service
```

### Docker Compose

```bash
# From project root
docker-compose up crypto-service
```

## Testing

### Simulate Deposits

Use the simulate deposit endpoint for testing:

```bash
curl -X POST http://localhost:8087/crypto/deposits/simulate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "currency": "BTC",
    "amount_crypto": 0.1
  }'
```

### Monitor Logs

```bash
# Docker
docker logs -f lfg-crypto-service

# Local
# Logs print to stdout
```

## Monitoring

The service logs important events:
- Exchange rate updates
- Deposit detection
- Confirmation tracking
- User account credits
- Withdrawal processing
- Errors and failures

Monitor the `crypto_monitoring_jobs` table for blockchain scanner health.

## Troubleshooting

### Deposits not detected
1. Check `crypto_monitoring_jobs` table for errors
2. Verify blockchain API connectivity
3. Check `last_scanned_block` is updating
4. Review service logs for errors

### Withdrawals stuck
1. Check `crypto_withdrawals` table for error messages
2. Verify withdrawal processing is enabled in config
3. Check blockchain API connectivity
4. Review withdrawal processor logs

### Exchange rates not updating
1. Check CoinGecko/CryptoCompare API status
2. Review rate service logs
3. Verify network connectivity
4. Check `exchange_rates` table last_updated

## Security Considerations

1. **Never expose private keys** - Always encrypted in database
2. **Secure the ENCRYPTION_KEY** - Use strong 32-byte key, store securely
3. **Enable 2FA in production** - Current implementation simulates 2FA
4. **Monitor for anomalies** - Set up alerts for unusual activity
5. **Regular security audits** - Review code and configurations
6. **Backup private keys** - Encrypted backup of wallet keys
7. **Rate limiting** - Already implemented in API gateway
8. **Input validation** - Verify all user inputs

## Future Enhancements

1. **Additional Cryptocurrencies**
   - More ERC-20 tokens
   - Other blockchains (Solana, Polygon, etc.)

2. **Advanced Features**
   - Batch withdrawals
   - Scheduled withdrawals
   - Multi-signature wallets
   - Payment forwarding

3. **Optimization**
   - Redis caching for rates
   - WebSocket real-time updates
   - Bulk transaction processing

4. **Analytics**
   - Volume tracking
   - User analytics
   - Fee collection reports

## License

Proprietary - LFG Platform

## Support

For issues or questions, contact the development team.
