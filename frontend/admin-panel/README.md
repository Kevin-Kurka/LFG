# LFG Admin Panel

React-based admin panel for managing the LFG Prediction Market Platform.

## Features

### Dashboard
- Overview of platform statistics
- Total markets, active markets, users, and orders
- Quick access to key metrics

### Market Management
- View all markets with filtering and search
- Filter by status (OPEN, CLOSED, RESOLVED)
- Search markets by question
- View market details including category, resolution date, and outcome

## Tech Stack

- **React 18.2**: UI library
- **TypeScript 5.2**: Type-safe development
- **Vite 5.0**: Fast build tool and dev server
- **React Router 6.20**: Client-side routing
- **Vanilla CSS**: Lightweight styling

## Setup

### Prerequisites
- Node.js 18+ and npm/yarn
- Running LFG backend services (API Gateway on port 8000)

### Installation

1. Install dependencies:
```bash
npm install
# or
yarn install
```

2. Start development server:
```bash
npm run dev
# or
yarn dev
```

The admin panel will be available at `http://localhost:3000`

### Build for Production

```bash
npm run build
# or
yarn build
```

The production build will be in the `dist/` directory.

## API Integration

The admin panel connects to the LFG API Gateway through a Vite proxy:
- Development: `http://localhost:3000/api` → `http://localhost:8000`
- Production: Configure reverse proxy (nginx, etc.) to route `/api` to API Gateway

## Project Structure

```
src/
├── main.tsx              # Entry point
├── App.tsx               # Main app component with routing
├── index.css             # Global styles
├── types/
│   └── index.ts          # TypeScript interfaces
├── services/
│   └── api.ts            # API service layer
└── pages/
    ├── Dashboard.tsx     # Dashboard page
    └── Markets.tsx       # Markets management page
```

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Adding New Pages

1. Create new component in `src/pages/`
2. Add route in `src/App.tsx`
3. Add navigation link in `Header` component

### API Service

The `ApiService` class in `src/services/api.ts` provides methods for:
- `getMarkets(params)` - Fetch markets with optional filters
- `getMarket(id)` - Fetch single market details
- `getDashboardStats()` - Fetch dashboard statistics

Add new methods as needed for additional features.

## Features Roadmap

- [ ] User management (view, suspend, ban users)
- [ ] Order monitoring and analytics
- [ ] Market creation and editing
- [ ] Market resolution interface
- [ ] Analytics and reporting
- [ ] Real-time updates with WebSocket
- [ ] Export data (CSV, JSON)
- [ ] Audit logs

## Styling

The admin panel uses vanilla CSS with a custom design system:
- Clean, modern interface
- Responsive grid layout
- Consistent color scheme (purple accent: #6C5CE7)
- Reusable CSS classes (.card, .btn, .badge, etc.)

## Security

- **No authentication**: Currently the admin panel has no auth. In production, add:
  - Admin authentication (separate from user auth)
  - Role-based access control
  - Session management
  - CSRF protection

## License

Proprietary - All rights reserved
