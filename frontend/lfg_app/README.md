# LFG Mobile App

Flutter mobile application for the LFG Prediction Market Trading Platform.

## Features

### Authentication
- User registration with email and password
- Secure login with JWT tokens
- Persistent authentication with shared preferences
- Profile management

### Markets
- Browse prediction markets with filtering and search
- View market details with contracts
- Real-time pricing information
- Category-based filtering (OPEN, RESOLVED, etc.)

### Trading
- Place market and limit orders
- View order history with status tracking
- Cancel active orders
- Real-time order execution

### Wallet
- View account balance
- Transaction history
- Buy credits with cryptocurrency (BTC, ETH, USDC)
- Real-time balance updates

## Architecture

### State Management
- Provider pattern for reactive state management
- Separate providers for Auth, Markets, Wallet, and Orders
- Automatic token injection for authenticated requests

### Data Layer
- RESTful API client with http package
- Model classes for type-safe data handling
- JSON serialization/deserialization

### UI Components
- Material Design 3 (Material You)
- Dark mode support
- Responsive layouts
- Pull-to-refresh on list screens

## Setup

### Prerequisites
- Flutter SDK 3.0+
- Dart SDK 3.0+
- Android Studio / Xcode for mobile development
- Running LFG backend services

### Installation

1. Install dependencies:
```bash
flutter pub get
```

2. Configure API endpoint in `lib/services/api_client.dart`:
```dart
static const String baseUrl = 'http://your-api-gateway:8000';
```

3. Run the app:
```bash
# Development
flutter run

# Production build
flutter build apk  # Android
flutter build ios  # iOS
```

## Project Structure

```
lib/
├── main.dart                 # App entry point
├── models/                   # Data models
│   ├── user.dart
│   ├── market.dart
│   ├── order.dart
│   └── wallet.dart
├── providers/                # State management
│   ├── auth_provider.dart
│   ├── market_provider.dart
│   ├── order_provider.dart
│   └── wallet_provider.dart
├── services/                 # API client
│   └── api_client.dart
└── screens/                  # UI screens
    ├── splash_screen.dart
    ├── auth/
    │   ├── login_screen.dart
    │   └── register_screen.dart
    ├── home/
    │   ├── home_screen.dart
    │   └── profile_screen.dart
    ├── market/
    │   ├── market_list_screen.dart
    │   └── market_detail_screen.dart
    ├── order/
    │   ├── orders_screen.dart
    │   └── place_order_screen.dart
    └── wallet/
        └── wallet_screen.dart
```

## API Integration

The app integrates with the following backend endpoints:

- **Auth**: `/register`, `/login`, `/profile`
- **Markets**: `/markets`, `/markets/:id`
- **Wallet**: `/wallet`, `/wallet/transactions`
- **Orders**: `/orders`, `/orders/:id/cancel`
- **Credit Exchange**: `/credit-exchange/buy`

## Testing

```bash
# Run unit tests
flutter test

# Run integration tests
flutter test integration_test/
```

## Build & Deploy

### Android
```bash
flutter build apk --release
flutter build appbundle --release
```

### iOS
```bash
flutter build ios --release
flutter build ipa
```

## Dependencies

- **provider**: State management
- **http**: HTTP client
- **shared_preferences**: Local storage
- **intl**: Internationalization and date formatting
- **web_socket_channel**: WebSocket support (for future real-time updates)

## Future Enhancements

- [ ] WebSocket integration for real-time notifications
- [ ] Push notifications
- [ ] Biometric authentication
- [ ] Charts and analytics
- [ ] Social features (leaderboard, following)
- [ ] Multi-language support

## License

Proprietary - All rights reserved
