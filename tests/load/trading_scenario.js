// k6 load testing scenario for LFG Platform
// Run with: k6 run tests/load/trading_scenario.js

import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const loginSuccessRate = new Rate('login_success_rate');
const orderPlacementRate = new Rate('order_placement_rate');
const orderLatency = new Trend('order_placement_latency');
const tradesExecuted = new Counter('trades_executed');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 },   // Ramp up to 100 users
    { duration: '5m', target: 1000 },  // Ramp up to 1000 users
    { duration: '5m', target: 1000 },  // Stay at 1000 users
    { duration: '2m', target: 5000 },  // Spike to 5000 users
    { duration: '5m', target: 5000 },  // Stay at 5000 users
    { duration: '2m', target: 0 },     // Ramp down to 0
  ],
  thresholds: {
    'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% under 500ms, 99% under 1s
    'http_req_failed': ['rate<0.01'],                  // Error rate under 1%
    'login_success_rate': ['rate>0.95'],               // Login success rate > 95%
    'order_placement_rate': ['rate>0.90'],             // Order placement success rate > 90%
    'order_placement_latency': ['p(95)<300'],          // Order latency p95 < 300ms
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';

// Generate random email
function randomEmail() {
  return `test-${__VU}-${__ITER}-${Date.now()}@example.com`;
}

// Generate random password
function randomPassword() {
  return `Test123!${__VU}${__ITER}`;
}

// Test scenario
export default function () {
  const email = randomEmail();
  const password = randomPassword();
  let token = '';

  // Group 1: User Registration
  group('user_registration', function () {
    const registerPayload = JSON.stringify({
      email: email,
      password: password,
    });

    const registerParams = {
      headers: { 'Content-Type': 'application/json' },
    };

    const registerRes = http.post(
      `${BASE_URL}/register`,
      registerPayload,
      registerParams
    );

    check(registerRes, {
      'registration status is 201': (r) => r.status === 201,
      'registration response has token': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.token !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    if (registerRes.status === 201) {
      try {
        const body = JSON.parse(registerRes.body);
        token = body.token;
      } catch (e) {
        console.error('Failed to parse register response');
      }
    }

    sleep(0.5);
  });

  // Group 2: User Login
  group('user_login', function () {
    const loginPayload = JSON.stringify({
      email: email,
      password: password,
    });

    const loginParams = {
      headers: { 'Content-Type': 'application/json' },
    };

    const loginRes = http.post(
      `${BASE_URL}/login`,
      loginPayload,
      loginParams
    );

    const loginSuccess = check(loginRes, {
      'login status is 200': (r) => r.status === 200,
      'login response has token': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.token !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    loginSuccessRate.add(loginSuccess);

    if (loginRes.status === 200) {
      try {
        const body = JSON.parse(r.body);
        token = body.token;
      } catch (e) {
        // Use existing token
      }
    }

    sleep(0.5);
  });

  if (!token) {
    console.error('No token available, skipping market operations');
    return;
  }

  // Group 3: Browse Markets
  group('browse_markets', function () {
    const headers = {
      'Authorization': `Bearer ${token}`,
    };

    const marketsRes = http.get(
      `${BASE_URL}/markets?status=OPEN&page=1&page_size=20`,
      { headers }
    );

    check(marketsRes, {
      'markets status is 200': (r) => r.status === 200,
      'markets response has data': (r) => {
        try {
          const body = JSON.parse(r.body);
          return Array.isArray(body.markets);
        } catch (e) {
          return false;
        }
      },
    });

    sleep(1);
  });

  // Group 4: Place Orders
  group('place_orders', function () {
    const headers = {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    };

    // Place buy order
    const buyOrderPayload = JSON.stringify({
      contract_id: 'test-contract-uuid',  // Replace with actual contract ID
      type: 'LIMIT',
      quantity: Math.floor(Math.random() * 100) + 1,
      limit_price: (Math.random() * 0.5 + 0.3).toFixed(4), // 0.30-0.80
    });

    const startTime = Date.now();
    const buyOrderRes = http.post(
      `${BASE_URL}/orders`,
      buyOrderPayload,
      { headers }
    );
    const orderTime = Date.now() - startTime;

    const orderSuccess = check(buyOrderRes, {
      'order status is 201': (r) => r.status === 201,
      'order response has id': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.id !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    orderPlacementRate.add(orderSuccess);
    orderLatency.add(orderTime);

    sleep(0.5);

    // Place sell order
    const sellOrderPayload = JSON.stringify({
      contract_id: 'test-contract-uuid',
      type: 'LIMIT',
      quantity: Math.floor(Math.random() * 100) + 1,
      limit_price: (Math.random() * 0.3 + 0.5).toFixed(4), // 0.50-0.80
    });

    http.post(`${BASE_URL}/orders`, sellOrderPayload, { headers });

    sleep(0.5);
  });

  // Group 5: Check Order Status
  group('check_orders', function () {
    const headers = {
      'Authorization': `Bearer ${token}`,
    };

    const ordersRes = http.get(
      `${BASE_URL}/orders?status=ACTIVE`,
      { headers }
    );

    check(ordersRes, {
      'orders status is 200': (r) => r.status === 200,
    });

    sleep(0.5);
  });

  // Group 6: Check Wallet
  group('check_wallet', function () {
    const headers = {
      'Authorization': `Bearer ${token}`,
    };

    const walletRes = http.get(
      `${BASE_URL}/wallet`,
      { headers }
    );

    check(walletRes, {
      'wallet status is 200': (r) => r.status === 200,
      'wallet has balance': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.balance_credits !== undefined;
        } catch (e) {
          return false;
        }
      },
    });

    sleep(0.5);
  });

  sleep(1);
}

// Setup function runs once per VU
export function setup() {
  console.log('Starting load test...');
  console.log(`Target: ${BASE_URL}`);
}

// Teardown function runs once at the end
export function teardown(data) {
  console.log('Load test completed');
}
