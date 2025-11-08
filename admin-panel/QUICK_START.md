# Quick Start Guide - LFG Admin Panel

## Installation & Setup

### 1. Install Dependencies
```bash
cd /home/user/LFG/admin-panel
npm install
```

### 2. Configure Environment
```bash
# Create .env file
cp .env.example .env

# Edit .env and set your API URL
echo "REACT_APP_API_BASE_URL=http://localhost:8000" > .env
```

### 3. Start Development Server
```bash
npm start
```

The admin panel will open at `http://localhost:3000`

---

## Docker Deployment

### Option 1: Docker Build & Run
```bash
# Build the image
docker build -t lfg-admin-panel .

# Run the container
docker run -d \
  -p 3001:80 \
  -e REACT_APP_API_BASE_URL=http://your-api-url \
  --name lfg-admin \
  lfg-admin-panel
```

### Option 2: Docker Compose
```bash
# Set API URL in environment
export API_BASE_URL=http://your-api-url

# Start the service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the service
docker-compose down
```

The admin panel will be available at `http://localhost:3001`

---

## Admin Features

### 1. Login
- Navigate to `/login`
- Enter admin credentials
- JWT token is stored and used for all API calls

### 2. Dashboard
- View platform statistics
- Quick access to all admin functions
- Real-time data updates

### 3. Markets Management
- **View Markets**: Browse all prediction markets
- **Create Market**: Add new markets with custom outcomes
- **Edit Market**: Update market details
- **Resolve Market**: Set the winning outcome
- **Cancel Market**: Cancel a market if needed

### 4. Sportsbook Management
- **View Providers**: See all sportsbook data providers
- **Enable/Disable**: Turn providers on or off
- **Sync Data**: Manually trigger odds sync
- **View Events**: Monitor live and scheduled sports events

### 5. Users Management
- **View Users**: Browse all platform users
- **User Details**: See user profile, wallet, and activity
- **Search Users**: Find users by username or email
- **Wallet Info**: Check user balances and locked funds

### 6. Orders & Trades
- **View Orders**: See all buy/sell orders
- **Filter Orders**: By status, type, and market
- **Trade History**: View completed trades
- **Order Book**: Monitor order book depth

### 7. Bets Management
- **View Bets**: See all sportsbook bets
- **Bet Status**: Monitor pending, won, and lost bets
- **User Bets**: View bets by user
- **Bet Analytics**: Track betting volume and payouts

### 8. Arbitrage & Hedges
- **Arbitrage Opportunities**: View profitable arbitrage setups
- **Hedge Opportunities**: See hedge options for existing bets
- **Execute Trades**: One-click arbitrage/hedge execution
- **Profit Tracking**: Monitor guaranteed profits

---

## API Configuration

The admin panel connects to the following backend services:

### Service Endpoints
- **User Service**: `/api/users/*`
- **Market Service**: `/api/markets/*`
- **Order Service**: `/api/orders/*`
- **Sportsbook Service**: `/api/sportsbook/*`
- **Arbitrage Service**: `/api/arbitrage/*`

### Authentication
All API requests automatically include:
```
Authorization: Bearer <jwt_token>
```

### Error Handling
- 401 responses trigger automatic logout
- Network errors show user-friendly messages
- Retry logic for transient failures

---

## Dark Mode

Toggle dark mode using the sun/moon icon in the header.

Theme preference is saved to localStorage and persists across sessions.

---

## Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| Go to Dashboard | Click "Dashboard" in sidebar |
| Create Market | Navigate to Markets > "Create Market" button |
| Search Users | Use search input in Users page |
| Filter Data | Use dropdown filters on list pages |

---

## Troubleshooting

### Issue: Can't login
**Solution**:
- Verify API URL in `.env` file
- Check backend service is running
- Verify admin credentials

### Issue: No data showing
**Solution**:
- Check browser console for API errors
- Verify API endpoints are accessible
- Check CORS configuration on backend

### Issue: Dark mode not working
**Solution**:
- Clear browser localStorage
- Refresh the page
- Check browser console for errors

### Issue: Docker container not starting
**Solution**:
```bash
# Check logs
docker logs lfg-admin

# Rebuild image
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## Production Deployment

### 1. Build for Production
```bash
npm run build
```

### 2. Deploy with Docker
```bash
# Build and tag
docker build -t lfg-admin-panel:v1.0.0 .

# Push to registry
docker tag lfg-admin-panel:v1.0.0 your-registry.com/lfg-admin:v1.0.0
docker push your-registry.com/lfg-admin:v1.0.0

# Deploy on server
docker run -d \
  -p 80:80 \
  -e REACT_APP_API_BASE_URL=https://api.yourplatform.com \
  --restart unless-stopped \
  --name lfg-admin-panel \
  your-registry.com/lfg-admin:v1.0.0
```

### 3. Configure HTTPS (Optional)
Update `nginx.conf` to include SSL configuration:
```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/private/key.pem;
    # ... rest of config
}
```

---

## Development Tips

### Hot Reload
The development server supports hot module replacement. Changes to components will update instantly without full page reload.

### TypeScript
All files use TypeScript for type safety. The compiler will catch type errors before runtime.

### Styling
Use Tailwind CSS utility classes for styling. Custom classes should be added to component files.

### API Mocking
For development without backend, you can modify service files to return mock data:
```typescript
// In any service file
async getUsers() {
  // return mockData; // Uncomment for mock data
  return userServiceApi.get('/users');
}
```

---

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

**Note**: Internet Explorer is not supported.

---

## Support

For issues or questions:
1. Check the main README.md
2. Review IMPLEMENTATION_SUMMARY.md
3. Check browser console for errors
4. Review Docker logs if using containers

---

## Next Steps

After setup:
1. Test login with admin credentials
2. Explore the dashboard
3. Create a test market
4. Review all admin sections
5. Configure real backend API URL
6. Set up production deployment
7. Configure monitoring and logging

---

## Quick Reference

### Common Commands
```bash
# Development
npm start              # Start dev server
npm run build          # Production build
npm test              # Run tests

# Docker
docker-compose up -d   # Start containers
docker-compose down    # Stop containers
docker-compose logs -f # View logs

# Production
docker build -t lfg-admin .
docker run -d -p 80:80 lfg-admin
```

### Important Files
- `.env` - Environment configuration
- `src/services/api.ts` - API client configuration
- `src/App.tsx` - Main routing configuration
- `Dockerfile` - Production build configuration
- `nginx.conf` - Web server configuration

---

## Success Checklist

- [ ] Dependencies installed
- [ ] Environment configured
- [ ] Development server runs
- [ ] Can access login page
- [ ] Can authenticate
- [ ] Dashboard loads with data
- [ ] All pages accessible
- [ ] Dark mode works
- [ ] Docker build succeeds
- [ ] Production deployment works

---

**You're all set! The LFG Admin Panel is ready to use.**
