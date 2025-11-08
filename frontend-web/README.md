# LFG Frontend Web Application

Production-ready React web application for the LFG prediction marketplace and sportsbook platform.

## Features

### Prediction Markets
- Browse and search prediction markets
- View detailed market information with order books
- Place buy/sell orders with limit pricing
- Real-time order matching and trade execution
- Live price updates via WebSocket

### Sportsbook Integration
- Compare odds across multiple sportsbooks
- View events with best odds highlighted
- Link sportsbook accounts securely
- Track bets across all platforms

### Arbitrage Detection
- Automatic arbitrage opportunity discovery
- Calculated stake recommendations
- Guaranteed profit calculations
- Real-time opportunity alerts

### Hedge Opportunities
- Find hedge opportunities for existing bets
- Calculate optimal hedge stakes
- Lock in guaranteed profits
- Risk-free betting strategies

### Bet Tracking
- Track all bets in one dashboard
- View profit & loss analytics
- Filter by status and provider
- Historical bet performance

### Wallet Management
- View balance and transaction history
- Buy/sell platform credits
- Secure payment integration ready
- Complete transaction audit trail

## Technology Stack

- **React 18** - Modern React with hooks
- **TypeScript** - Type-safe development
- **Vite** - Lightning-fast build tool
- **Tailwind CSS** - Utility-first styling
- **React Router** - Client-side routing
- **React Query** - Server state management
- **Axios** - HTTP client
- **WebSocket** - Real-time updates
- **Recharts** - Data visualization

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend services running (see `/backend`)

### Installation

```bash
# Install dependencies
npm install

# Copy environment variables
cp .env.example .env

# Edit .env with your backend service URLs
nano .env
```

### Development

```bash
# Start development server
npm run dev

# Open http://localhost:3000
```

The development server includes:
- Hot module replacement
- TypeScript type checking
- Tailwind CSS compilation
- API proxy to backend services

### Building for Production

```bash
# Build optimized production bundle
npm run build

# Preview production build
npm run preview
```

### Running with Docker

```bash
# Build Docker image
docker build -t lfg-frontend .

# Run container
docker run -p 3000:80 lfg-frontend

# Or use docker-compose
docker-compose up -d
```

## Project Structure

```
frontend-web/
├── public/              # Static assets
├── src/
│   ├── components/      # Reusable UI components
│   │   ├── Navbar.tsx
│   │   ├── OrderBook.tsx
│   │   ├── OddsComparison.tsx
│   │   ├── ArbitrageCard.tsx
│   │   ├── HedgeCalculator.tsx
│   │   ├── BetTracker.tsx
│   │   └── ...
│   ├── pages/           # Route pages
│   │   ├── Home.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Markets.tsx
│   │   ├── Sportsbook.tsx
│   │   ├── Arbitrage.tsx
│   │   ├── Hedges.tsx
│   │   └── ...
│   ├── services/        # API services
│   │   ├── api.ts
│   │   ├── auth.service.ts
│   │   ├── market.service.ts
│   │   ├── sportsbook.service.ts
│   │   └── websocket.service.ts
│   ├── context/         # React contexts
│   │   ├── AuthContext.tsx
│   │   ├── ThemeContext.tsx
│   │   └── OddsFormatContext.tsx
│   ├── hooks/           # Custom hooks
│   │   └── useWebSocket.ts
│   ├── types/           # TypeScript types
│   │   └── index.ts
│   ├── utils/           # Utility functions
│   │   ├── odds.ts
│   │   └── format.ts
│   ├── App.tsx          # Main app component
│   ├── main.tsx         # Entry point
│   └── index.css        # Global styles
├── Dockerfile           # Production container
├── nginx.conf           # Nginx configuration
└── package.json         # Dependencies
```

## Key Components

### OrderBook
Visual order book display with bid/ask spreads, price bars, and real-time updates.

### OddsComparison
Side-by-side odds comparison across sportsbooks with best odds highlighted.

### ArbitrageCard
Arbitrage opportunity card with stake calculator and ROI display.

### HedgeCalculator
Interactive hedge calculator showing guaranteed profit scenarios.

### BetTracker
Comprehensive bet tracking table with P&L, status, and filtering.

## API Integration

All API calls are made through service layers in `/src/services/`:

- `auth.service.ts` - User authentication
- `wallet.service.ts` - Wallet and transactions
- `market.service.ts` - Prediction markets
- `order.service.ts` - Order placement
- `sportsbook.service.ts` - Sportsbook data
- `websocket.service.ts` - Real-time updates

### Environment Variables

```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WALLET_SERVICE_URL=http://localhost:8081
VITE_MARKET_SERVICE_URL=http://localhost:8082
VITE_ORDER_SERVICE_URL=http://localhost:8085
VITE_CREDIT_SERVICE_URL=http://localhost:8086
VITE_SPORTSBOOK_SERVICE_URL=http://localhost:8088
VITE_WS_URL=ws://localhost:8087
```

## Features in Detail

### Authentication
- JWT-based authentication
- Secure token storage
- Protected routes
- Auto-redirect on expiration

### Real-Time Updates
- WebSocket connection with auto-reconnect
- Live market prices
- Trade notifications
- Arbitrage alerts

### Dark Mode
- System preference detection
- Manual toggle
- Persistent user preference
- Tailwind dark mode classes

### Responsive Design
- Mobile-first approach
- Tablet and desktop optimized
- Touch-friendly interfaces
- Accessible components

### Odds Format Toggle
- American odds (+150, -200)
- Decimal odds (2.50)
- Fractional odds (3/2)
- Persistent user preference

## Performance Optimizations

- Code splitting by route
- Lazy loading of components
- Image optimization
- Gzip compression
- Browser caching
- React Query caching
- WebSocket connection pooling

## Security Features

- JWT token authentication
- Secure credential encryption
- XSS protection headers
- CSRF protection
- Input sanitization
- API error handling

## Browser Support

- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)
- Mobile browsers (iOS Safari, Chrome Android)

## Development Guidelines

### Code Style
- TypeScript strict mode
- ESLint for linting
- Prettier for formatting
- Conventional commits

### Component Guidelines
- Functional components with hooks
- TypeScript interfaces for props
- Proper error boundaries
- Loading and error states

### State Management
- React Context for global state
- React Query for server state
- Local state for UI state
- WebSocket for real-time data

## Deployment

### Production Build
The production build is optimized for performance:
- Minified JavaScript
- CSS purging with Tailwind
- Tree shaking
- Asset optimization

### Docker Deployment
Multi-stage Dockerfile creates a minimal production image:
- Node.js for building
- Nginx for serving
- ~50MB final image size

### Environment Configuration
Update environment variables for production:
- API endpoints
- WebSocket URL
- Feature flags

## Troubleshooting

### Build Errors
```bash
# Clear cache and rebuild
rm -rf node_modules dist
npm install
npm run build
```

### API Connection Issues
- Check backend services are running
- Verify environment variables
- Check CORS configuration
- Inspect network tab in DevTools

### WebSocket Issues
- Verify WebSocket URL format
- Check firewall/proxy settings
- Monitor connection status in console

## Contributing

1. Follow TypeScript and React best practices
2. Write clean, documented code
3. Test all features thoroughly
4. Ensure responsive design
5. Update README for new features

## License

Proprietary - All rights reserved

## Support

For issues or questions, contact the development team.
