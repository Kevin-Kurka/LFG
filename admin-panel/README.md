# LFG Admin Panel

A comprehensive admin panel for managing the LFG prediction market and sportsbook platform.

## Features

### Authentication
- Admin login with JWT authentication
- Protected routes
- Token-based session management

### Dashboard
- Platform statistics overview
- Total users, markets, bets, and volume
- Active users and markets tracking
- Quick action links

### Markets Management
- View all prediction markets with filtering
- Create new markets with custom outcomes
- Edit existing markets
- Resolve markets and set winners
- Cancel markets when needed

### Sportsbook Management
- View and manage sportsbook providers
- Enable/disable providers
- Sync odds data
- View sports events (scheduled, live, finished)
- Monitor odds across providers

### Users Management
- List all platform users
- View detailed user profiles
- Check user wallet balances (available and locked)
- Monitor user activity history
- Search users by username or email

### Orders & Trades
- View all orders across markets
- Filter by status, type, and market
- Monitor order book activity
- Track trade history
- Real-time order status updates

### Bets Management
- View all sportsbook bets
- Filter by status, user, and event
- Monitor bet outcomes (pending, won, lost)
- Track betting volume and payouts
- Bet settlement management

### Arbitrage & Hedges
- View arbitrage opportunities across providers
- Monitor hedge opportunities
- Execute profitable arbitrage trades
- Track guaranteed profits
- Filter by active/executed/expired status

## Tech Stack

- **React 19** with TypeScript
- **Tailwind CSS** for styling
- **React Router** for navigation
- **Axios** for API calls
- **React Hook Form** for form management
- **date-fns** for date handling
- **Recharts** for analytics charts

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend services running (User, Market, Order, Sportsbook, Arbitrage services)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Configure environment variables:
```env
REACT_APP_API_BASE_URL=http://localhost:8000
```

### Development

Run the development server:
```bash
npm start
```

The admin panel will be available at `http://localhost:3000`

### Production Build

Build for production:
```bash
npm run build
```

## Docker Deployment

### Build and run with Docker:

```bash
docker build -t lfg-admin-panel .
docker run -p 3001:80 -e REACT_APP_API_BASE_URL=http://your-api-url lfg-admin-panel
```

### Using Docker Compose:

```bash
docker-compose up -d
```

The admin panel will be available at `http://localhost:3001`

## Project Structure

```
admin-panel/
├── src/
│   ├── components/
│   │   ├── common/          # Reusable UI components
│   │   └── layout/          # Layout components (Sidebar, Header)
│   ├── contexts/            # React contexts (Auth, Theme)
│   ├── pages/               # Page components
│   │   ├── auth/           # Login
│   │   ├── dashboard/      # Dashboard
│   │   ├── markets/        # Markets management
│   │   ├── sportsbook/     # Sportsbook management
│   │   ├── users/          # Users management
│   │   ├── orders/         # Orders & Trades
│   │   ├── bets/           # Bets management
│   │   └── arbitrage/      # Arbitrage opportunities
│   ├── services/           # API service layer
│   ├── types/              # TypeScript type definitions
│   └── utils/              # Utility functions
├── public/                 # Static assets
├── Dockerfile             # Production Docker build
├── nginx.conf             # Nginx configuration
└── docker-compose.yml     # Docker Compose configuration
```

## API Integration

The admin panel integrates with the following backend services:

- **User Service**: User management and authentication
- **Market Service**: Prediction market management
- **Order Service**: Order book and trading
- **Sportsbook Service**: Sports events, odds, and bets
- **Arbitrage Service**: Arbitrage and hedge opportunities

All API calls are handled through service classes in `src/services/`.

## Features Overview

### UI Components

- **Button**: Customizable buttons with variants (primary, secondary, danger, success)
- **Card**: Container component for content sections
- **Table**: Data table with sorting and pagination
- **Badge**: Status indicators with color variants
- **Input/Select**: Form input components
- **StatsCard**: Statistics display cards
- **Pagination**: Pagination controls

### Authentication

- JWT token storage in localStorage
- Automatic token refresh
- Protected routes with redirect to login
- Logout functionality

### Dark Mode

- Full dark mode support
- Toggle between light and dark themes
- Persistent theme preference

### Responsive Design

- Mobile-friendly interface
- Responsive layouts for all screen sizes
- Touch-optimized controls

## Default Admin Credentials

For initial login, use the admin credentials configured in your backend.

## Security Features

- JWT token authentication
- Protected API routes
- HTTPS support in production
- XSS protection headers
- CSRF token handling

## Performance

- Code splitting with React Router
- Lazy loading of components
- Optimized bundle size
- Nginx caching for static assets
- Gzip compression

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Contributing

1. Create feature branches from `main`
2. Follow TypeScript best practices
3. Use Tailwind CSS for styling
4. Add proper error handling
5. Test all features before committing

## License

Private - LFG Platform

## Support

For issues or questions, contact the development team.
