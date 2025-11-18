# LFG Platform - Optimization & Performance Guide

**Version**: 1.0
**Last Updated**: 2025-11-18
**Target**: 10,000+ concurrent users, <500ms p95 latency

This guide provides detailed optimization strategies and implementation examples for the LFG Platform.

---

## Table of Contents

1. [Database Optimization](#database-optimization)
2. [API Performance](#api-performance)
3. [Matching Engine Optimization](#matching-engine-optimization)
4. [Caching Strategy](#caching-strategy)
5. [WebSocket Optimization](#websocket-optimization)
6. [Frontend Optimization](#frontend-optimization)
7. [Infrastructure Optimization](#infrastructure-optimization)
8. [Monitoring & Profiling](#monitoring--profiling)

---

## Database Optimization

### Query Optimization

#### 1. Add Missing Indexes

```sql
-- Orders table - Critical for trading performance
CREATE INDEX CONCURRENTLY idx_orders_user_status
ON orders(user_id, status)
WHERE status IN ('ACTIVE', 'PARTIALLY_FILLED');

CREATE INDEX CONCURRENTLY idx_orders_contract_status
ON orders(contract_id, status)
WHERE status IN ('ACTIVE', 'PARTIALLY_FILLED');

CREATE INDEX CONCURRENTLY idx_orders_created_at_desc
ON orders(created_at DESC);

-- Trades table - For market history and analytics
CREATE INDEX CONCURRENTLY idx_trades_market_created
ON trades(market_id, created_at DESC);

CREATE INDEX CONCURRENTLY idx_trades_user_created
ON trades(user_id, created_at DESC);

-- Wallets table - For balance queries
CREATE INDEX CONCURRENTLY idx_wallets_balance
ON wallets(balance_credits)
WHERE balance_credits > 0;

-- Markets table - For market discovery
CREATE INDEX CONCURRENTLY idx_markets_status_expires
ON markets(status, expires_at)
WHERE status = 'OPEN';

-- Composite indexes for common queries
CREATE INDEX CONCURRENTLY idx_orders_composite
ON orders(contract_id, status, limit_price_credits, created_at);
```

#### 2. Optimize Common Queries

**Before** (N+1 Query):
```go
// Bad: N+1 query problem
func (r *OrderRepository) GetOrdersWithMarketInfo(ctx context.Context, userID uuid.UUID) ([]*Order, error) {
    orders, _ := r.GetByUserID(ctx, userID)
    for _, order := range orders {
        market, _ := r.marketRepo.GetByID(ctx, order.MarketID) // N queries!
        order.Market = market
    }
    return orders, nil
}
```

**After** (Optimized with JOIN):
```go
// Good: Single query with JOIN
func (r *OrderRepository) GetOrdersWithMarketInfo(ctx context.Context, userID uuid.UUID) ([]*Order, error) {
    query := `
        SELECT
            o.id, o.user_id, o.contract_id, o.type, o.status,
            o.quantity, o.quantity_filled, o.limit_price_credits,
            o.created_at, o.updated_at,
            m.id, m.ticker, m.question, m.status
        FROM orders o
        JOIN contracts c ON o.contract_id = c.id
        JOIN markets m ON c.market_id = m.id
        WHERE o.user_id = $1
        ORDER BY o.created_at DESC
        LIMIT 100
    `
    // Single query instead of N+1
    rows, err := r.pool.Query(ctx, query, userID)
    // ... scan results
}
```

#### 3. Use EXPLAIN ANALYZE

```go
// Add query performance logging in development
func (r *Repository) executeWithExplain(ctx context.Context, query string, args ...interface{}) {
    if os.Getenv("ENVIRONMENT") == "development" {
        explainQuery := "EXPLAIN ANALYZE " + query
        rows, _ := r.pool.Query(ctx, explainQuery, args...)
        // Log execution plan
        var plan string
        for rows.Next() {
            rows.Scan(&plan)
            log.Debug().Str("plan", plan).Msg("Query execution plan")
        }
    }
    // Execute actual query
    return r.pool.Query(ctx, query, args...)
}
```

### Connection Pool Optimization

```go
// backend/shared/db/pool.go
package db

import (
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type PoolConfig struct {
    // Connection limits
    MaxConns        int32         // Maximum connections (default: 4 per CPU)
    MinConns        int32         // Minimum idle connections (default: 2)

    // Timeouts
    MaxConnLifetime time.Duration // Max connection age (default: 1 hour)
    MaxConnIdleTime time.Duration // Max idle time (default: 30 minutes)

    // Health checks
    HealthCheckPeriod time.Duration // Health check interval (default: 1 minute)
}

func NewOptimizedPool(ctx context.Context, cfg PoolConfig) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(cfg.DSN)
    if err != nil {
        return nil, err
    }

    // Optimize connection pool settings
    config.MaxConns = cfg.MaxConns        // 25-50 for production
    config.MinConns = cfg.MinConns        // 5-10 for production
    config.MaxConnLifetime = cfg.MaxConnLifetime
    config.MaxConnIdleTime = cfg.MaxConnIdleTime
    config.HealthCheckPeriod = cfg.HealthCheckPeriod

    // Connection pool metrics
    config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
        // Set connection-level optimizations
        _, err := conn.Exec(ctx, "SET statement_timeout = '30s'")
        return err
    }

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, err
    }

    // Verify pool health
    if err := pool.Ping(ctx); err != nil {
        return nil, err
    }

    return pool, nil
}
```

### Read Replicas for Scalability

```go
// backend/shared/db/replica.go
package db

type DBCluster struct {
    primary  *pgxpool.Pool
    replicas []*pgxpool.Pool
    current  int
    mu       sync.Mutex
}

func NewDBCluster(primaryDSN string, replicaDSNs []string) (*DBCluster, error) {
    cluster := &DBCluster{
        replicas: make([]*pgxpool.Pool, 0, len(replicaDSNs)),
    }

    // Connect to primary
    primary, err := NewOptimizedPool(context.Background(), primaryDSN)
    if err != nil {
        return nil, err
    }
    cluster.primary = primary

    // Connect to replicas
    for _, dsn := range replicaDSNs {
        replica, err := NewOptimizedPool(context.Background(), dsn)
        if err != nil {
            continue // Log but don't fail
        }
        cluster.replicas = append(cluster.replicas, replica)
    }

    return cluster, nil
}

// Primary returns the primary database for writes
func (c *DBCluster) Primary() *pgxpool.Pool {
    return c.primary
}

// Replica returns a replica using round-robin load balancing
func (c *DBCluster) Replica() *pgxpool.Pool {
    if len(c.replicas) == 0 {
        return c.primary // Fallback to primary
    }

    c.mu.Lock()
    defer c.mu.Unlock()

    replica := c.replicas[c.current]
    c.current = (c.current + 1) % len(c.replicas)
    return replica
}

// Usage in repository
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Order, error) {
    // Use replica for read query
    return r.queryOrders(ctx, r.cluster.Replica(), "WHERE user_id = $1", userID)
}

func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
    // Use primary for write query
    return r.insertOrder(ctx, r.cluster.Primary(), order)
}
```

### Database Partitioning

```sql
-- Partition orders table by created_at for better performance
-- This is useful when you have millions of orders

-- 1. Create partitioned table
CREATE TABLE orders_partitioned (
    id UUID NOT NULL,
    user_id UUID NOT NULL,
    contract_id UUID NOT NULL,
    type order_type NOT NULL,
    status order_status NOT NULL,
    quantity INTEGER NOT NULL,
    quantity_filled INTEGER NOT NULL DEFAULT 0,
    limit_price_credits DECIMAL(10, 8),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- 2. Create monthly partitions
CREATE TABLE orders_2025_01 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');

CREATE TABLE orders_2025_02 PARTITION OF orders_partitioned
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');

-- 3. Create indexes on partitions
CREATE INDEX idx_orders_2025_01_user ON orders_2025_01(user_id);
CREATE INDEX idx_orders_2025_02_user ON orders_2025_02(user_id);

-- 4. Automate partition creation (using pg_cron or external script)
-- Example: Monthly job to create next month's partition
```

---

## API Performance

### Response Compression

```go
// backend/api-gateway/middleware/compression.go
package middleware

import (
    "compress/gzip"
    "net/http"
    "strings"
)

type gzipResponseWriter struct {
    http.ResponseWriter
    Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func CompressionMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if client accepts gzip
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        // Only compress if response is large enough (>1KB)
        w.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(w)
        defer gz.Close()

        gzw := gzipResponseWriter{ResponseWriter: w, Writer: gz}
        next.ServeHTTP(gzw, r)
    })
}
```

### Request Batching

```go
// backend/shared/batch/batcher.go
package batch

import (
    "context"
    "sync"
    "time"
)

type Batcher struct {
    maxSize    int
    maxWait    time.Duration
    buffer     []interface{}
    mu         sync.Mutex
    flushFunc  func([]interface{}) error
    timer      *time.Timer
}

func NewBatcher(maxSize int, maxWait time.Duration, flushFunc func([]interface{}) error) *Batcher {
    b := &Batcher{
        maxSize:   maxSize,
        maxWait:   maxWait,
        buffer:    make([]interface{}, 0, maxSize),
        flushFunc: flushFunc,
    }
    return b
}

func (b *Batcher) Add(item interface{}) {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.buffer = append(b.buffer, item)

    if len(b.buffer) >= b.maxSize {
        b.flush()
    } else if b.timer == nil {
        b.timer = time.AfterFunc(b.maxWait, func() {
            b.mu.Lock()
            defer b.mu.Unlock()
            b.flush()
        })
    }
}

func (b *Batcher) flush() {
    if len(b.buffer) == 0 {
        return
    }

    items := b.buffer
    b.buffer = make([]interface{}, 0, b.maxSize)

    if b.timer != nil {
        b.timer.Stop()
        b.timer = nil
    }

    // Flush in goroutine to avoid blocking
    go b.flushFunc(items)
}

// Usage: Batch database inserts
func (r *Repository) BatchInsertOrders(orders []*Order) error {
    batcher := NewBatcher(100, 5*time.Second, func(items []interface{}) error {
        // Insert all orders in single query
        return r.bulkInsert(items)
    })

    for _, order := range orders {
        batcher.Add(order)
    }

    return nil
}
```

### HTTP/2 and Connection Pooling

```go
// backend/shared/client/http.go
package client

import (
    "net"
    "net/http"
    "time"
)

func NewOptimizedHTTPClient() *http.Client {
    transport := &http.Transport{
        // Connection pooling
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 20,
        MaxConnsPerHost:     100,

        // Timeouts
        DialContext: (&net.Dialer{
            Timeout:   10 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   10 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
        IdleConnTimeout:       90 * time.Second,

        // HTTP/2
        ForceAttemptHTTP2: true,
    }

    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }
}
```

---

## Matching Engine Optimization

### Efficient Order Book Data Structure

```go
// backend/matching-engine/engine/order_book.go
package engine

import (
    "container/heap"
    "sync"
)

// OrderBook uses priority queues for O(log n) insertion and O(1) peek
type OrderBook struct {
    mu         sync.RWMutex
    buyOrders  *MaxHeap // Max heap by price (highest price first)
    sellOrders *MinHeap // Min heap by price (lowest price first)
    orderIndex map[string]*Order
}

// MaxHeap for buy orders (highest price has priority)
type MaxHeap []*Order

func (h MaxHeap) Len() int { return len(h) }

func (h MaxHeap) Less(i, j int) bool {
    // Primary: price (higher is better for buys)
    if h[i].Price != h[j].Price {
        return h[i].Price > h[j].Price
    }
    // Secondary: time priority (earlier is better)
    return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h MaxHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x interface{}) {
    *h = append(*h, x.(*Order))
}

func (h *MaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

// MinHeap for sell orders (lowest price has priority)
type MinHeap []*Order

func (h MinHeap) Len() int { return len(h) }

func (h MinHeap) Less(i, j int) bool {
    // Primary: price (lower is better for sells)
    if h[i].Price != h[j].Price {
        return h[i].Price < h[j].Price
    }
    // Secondary: time priority (earlier is better)
    return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h MinHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
    *h = append(*h, x.(*Order))
}

func (h *MinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}

func NewOrderBook() *OrderBook {
    buyHeap := &MaxHeap{}
    sellHeap := &MinHeap{}
    heap.Init(buyHeap)
    heap.Init(sellHeap)

    return &OrderBook{
        buyOrders:  buyHeap,
        sellOrders: sellHeap,
        orderIndex: make(map[string]*Order),
    }
}

func (ob *OrderBook) AddOrder(order *Order) {
    ob.mu.Lock()
    defer ob.mu.Unlock()

    ob.orderIndex[order.ID] = order

    if order.Side == "BUY" {
        heap.Push(ob.buyOrders, order)
    } else {
        heap.Push(ob.sellOrders, order)
    }
}

func (ob *OrderBook) Match() []*Trade {
    ob.mu.Lock()
    defer ob.mu.Unlock()

    trades := make([]*Trade, 0)

    // Match orders while there's overlap in prices
    for ob.buyOrders.Len() > 0 && ob.sellOrders.Len() > 0 {
        topBuy := (*ob.buyOrders)[0]
        topSell := (*ob.sellOrders)[0]

        // Check if prices cross
        if topBuy.Price < topSell.Price {
            break // No more matches possible
        }

        // Execute trade
        tradeQuantity := min(topBuy.RemainingQuantity(), topSell.RemainingQuantity())
        tradePrice := topSell.Price // Price-time priority: seller's price

        trade := &Trade{
            BuyOrderID:  topBuy.ID,
            SellOrderID: topSell.ID,
            Quantity:    tradeQuantity,
            Price:       tradePrice,
            Timestamp:   time.Now(),
        }
        trades = append(trades, trade)

        // Update order quantities
        topBuy.FilledQuantity += tradeQuantity
        topSell.FilledQuantity += tradeQuantity

        // Remove fully filled orders
        if topBuy.RemainingQuantity() == 0 {
            heap.Pop(ob.buyOrders)
            delete(ob.orderIndex, topBuy.ID)
        }
        if topSell.RemainingQuantity() == 0 {
            heap.Pop(ob.sellOrders)
            delete(ob.orderIndex, topSell.ID)
        }
    }

    return trades
}
```

### Batch Order Processing

```go
// Process multiple orders in batch for better throughput
func (e *MatchingEngine) ProcessOrderBatch(orders []*Order) []*Trade {
    allTrades := make([]*Trade, 0)

    // Sort orders by timestamp to maintain fairness
    sort.Slice(orders, func(i, j int) bool {
        return orders[i].CreatedAt.Before(orders[j].CreatedAt)
    })

    for _, order := range orders {
        e.orderBook.AddOrder(order)
        trades := e.orderBook.Match()
        allTrades = append(allTrades, trades...)
    }

    // Publish all trades in batch
    if len(allTrades) > 0 {
        e.publishTrades(allTrades)
    }

    return allTrades
}
```

### Lock-Free Concurrent Matching

```go
// For very high throughput, partition order books by market
type PartitionedMatchingEngine struct {
    markets map[string]*OrderBook
    mu      sync.RWMutex
}

func (e *PartitionedMatchingEngine) GetOrderBook(marketID string) *OrderBook {
    e.mu.RLock()
    ob, exists := e.markets[marketID]
    e.mu.RUnlock()

    if !exists {
        e.mu.Lock()
        ob = NewOrderBook()
        e.markets[marketID] = ob
        e.mu.Unlock()
    }

    return ob
}

// Each market can be matched concurrently
func (e *PartitionedMatchingEngine) ProcessOrder(order *Order) {
    ob := e.GetOrderBook(order.MarketID)
    ob.AddOrder(order)
    trades := ob.Match()

    if len(trades) > 0 {
        e.publishTrades(order.MarketID, trades)
    }
}
```

---

## Caching Strategy

### Multi-Level Caching

```go
// backend/shared/cache/multi_level.go
package cache

import (
    "context"
    "time"

    "github.com/allegro/bigcache/v3"
    "github.com/redis/go-redis/v9"
)

type MultiLevelCache struct {
    local  *bigcache.BigCache  // L1: In-memory cache (fast, limited size)
    redis  *redis.Client        // L2: Redis cache (slower, larger size)
}

func NewMultiLevelCache() (*MultiLevelCache, error) {
    // L1 cache: In-memory (100MB, 10 minute TTL)
    localConfig := bigcache.DefaultConfig(10 * time.Minute)
    localConfig.HardMaxCacheSize = 100 // 100MB
    local, err := bigcache.New(context.Background(), localConfig)
    if err != nil {
        return nil, err
    }

    // L2 cache: Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        PoolSize: 50,
    })

    return &MultiLevelCache{
        local: local,
        redis: redisClient,
    }, nil
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) ([]byte, error) {
    // Try L1 cache first
    if data, err := c.local.Get(key); err == nil {
        return data, nil
    }

    // Try L2 cache (Redis)
    data, err := c.redis.Get(ctx, key).Bytes()
    if err == nil {
        // Populate L1 cache for next request
        c.local.Set(key, data)
        return data, nil
    }

    return nil, err
}

func (c *MultiLevelCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Set in both caches
    c.local.Set(key, value)
    return c.redis.Set(ctx, key, value, ttl).Err()
}

func (c *MultiLevelCache) Delete(ctx context.Context, key string) error {
    c.local.Delete(key)
    return c.redis.Del(ctx, key).Err()
}
```

### Cache-Aside Pattern

```go
// backend/market-service/repository/cached_market_repository.go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "lfg/shared/cache"
    "lfg/shared/models"
)

type CachedMarketRepository struct {
    repo  *MarketRepository
    cache *cache.MultiLevelCache
}

func (r *CachedMarketRepository) GetByID(ctx context.Context, id string) (*models.Market, error) {
    cacheKey := fmt.Sprintf("market:%s", id)

    // Try cache first
    if data, err := r.cache.Get(ctx, cacheKey); err == nil {
        var market models.Market
        if err := json.Unmarshal(data, &market); err == nil {
            return &market, nil
        }
    }

    // Cache miss: fetch from database
    market, err := r.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Populate cache for next request
    if data, err := json.Marshal(market); err == nil {
        r.cache.Set(ctx, cacheKey, data, 5*time.Minute)
    }

    return market, nil
}

func (r *CachedMarketRepository) Update(ctx context.Context, market *models.Market) error {
    // Update database
    if err := r.repo.Update(ctx, market); err != nil {
        return err
    }

    // Invalidate cache
    cacheKey := fmt.Sprintf("market:%s", market.ID)
    r.cache.Delete(ctx, cacheKey)

    return nil
}
```

### Cache Warming

```go
// Pre-populate cache on startup with frequently accessed data
func (s *MarketService) WarmCache(ctx context.Context) error {
    log.Info().Msg("Warming cache...")

    // Cache all open markets
    markets, err := s.repo.GetOpenMarkets(ctx)
    if err != nil {
        return err
    }

    for _, market := range markets {
        cacheKey := fmt.Sprintf("market:%s", market.ID)
        data, _ := json.Marshal(market)
        s.cache.Set(ctx, cacheKey, data, 10*time.Minute)
    }

    log.Info().Int("count", len(markets)).Msg("Cache warmed")
    return nil
}
```

---

## WebSocket Optimization

### Connection Pooling and Load Balancing

```go
// backend/notification-service/websocket/pool.go
package websocket

import (
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type ConnectionPool struct {
    connections map[string]*Connection
    mu          sync.RWMutex
    upgrader    websocket.Upgrader
}

type Connection struct {
    UserID    string
    Conn      *websocket.Conn
    Send      chan []byte
    LastPing  time.Time
    mu        sync.Mutex
}

func NewConnectionPool() *ConnectionPool {
    return &ConnectionPool{
        connections: make(map[string]*Connection),
        upgrader: websocket.Upgrader{
            ReadBufferSize:  1024,
            WriteBufferSize: 4096, // Larger write buffer for batching
            CheckOrigin: func(r *http.Request) bool {
                return true // Configure properly in production
            },
        },
    }
}

func (p *ConnectionPool) Add(userID string, conn *websocket.Conn) *Connection {
    c := &Connection{
        UserID:   userID,
        Conn:     conn,
        Send:     make(chan []byte, 256), // Buffer for messages
        LastPing: time.Now(),
    }

    p.mu.Lock()
    p.connections[userID] = c
    p.mu.Unlock()

    // Start read/write pumps
    go c.readPump(p)
    go c.writePump()

    return c
}

func (c *Connection) writePump() {
    ticker := time.NewTicker(30 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // Write message
            if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            // Send ping
            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Connection) readPump(pool *ConnectionPool) {
    defer func() {
        pool.Remove(c.UserID)
        c.Conn.Close()
    }()

    c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.Conn.SetPongHandler(func(string) error {
        c.mu.Lock()
        c.LastPing = time.Now()
        c.mu.Unlock()
        c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        _, _, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }
    }
}

// Broadcast to specific user
func (p *ConnectionPool) SendToUser(userID string, message []byte) {
    p.mu.RLock()
    conn, exists := p.connections[userID]
    p.mu.RUnlock()

    if exists {
        select {
        case conn.Send <- message:
        default:
            // Channel full, close connection
            close(conn.Send)
            p.Remove(userID)
        }
    }
}
```

### Message Batching for WebSocket

```go
// Batch multiple messages to reduce overhead
type MessageBatcher struct {
    messages [][]byte
    mu       sync.Mutex
    timer    *time.Timer
    flush    func([][]byte)
}

func (b *MessageBatcher) Add(message []byte) {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.messages = append(b.messages, message)

    if len(b.messages) >= 10 || b.timer == nil {
        b.flushMessages()
    }
}

func (b *MessageBatcher) flushMessages() {
    if len(b.messages) == 0 {
        return
    }

    messages := b.messages
    b.messages = make([][]byte, 0, 10)

    if b.timer != nil {
        b.timer.Stop()
        b.timer = nil
    }

    // Combine messages
    combined := bytes.Join(messages, []byte("\n"))
    b.flush([][]byte{combined})
}
```

---

## Frontend Optimization

### Flutter Mobile App

```dart
// lib/services/api_service.dart
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';

class OptimizedApiService {
  late final Dio _dio;

  OptimizedApiService() {
    _dio = Dio(
      BaseOptions(
        baseUrl: 'https://api.lfg.com',
        connectTimeout: Duration(seconds: 10),
        receiveTimeout: Duration(seconds: 30),
        sendTimeout: Duration(seconds: 10),
        // Enable HTTP/2
        headers: {
          'Accept-Encoding': 'gzip',
        },
      ),
    );

    // Add response compression
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        // Add auth token
        options.headers['Authorization'] = 'Bearer ${getToken()}';
        handler.next(options);
      },
      onError: (error, handler) {
        // Retry logic
        if (error.response?.statusCode == 429) {
          // Rate limited, wait and retry
          Future.delayed(Duration(seconds: 2), () {
            handler.resolve(retry(error.requestOptions));
          });
        } else {
          handler.next(error);
        }
      },
    ));
  }

  // Implement pagination for market lists
  Future<List<Market>> getMarkets({
    int page = 1,
    int pageSize = 20,
  }) async {
    final response = await _dio.get(
      '/markets',
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
    );

    return (response.data['markets'] as List)
        .map((m) => Market.fromJson(m))
        .toList();
  }
}

// Implement lazy loading in ListView
class MarketListView extends StatefulWidget {
  @override
  _MarketListViewState createState() => _MarketListViewState();
}

class _MarketListViewState extends State<MarketListView> {
  final ScrollController _scrollController = ScrollController();
  final List<Market> _markets = [];
  int _currentPage = 1;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadMarkets();

    // Load more when scrolled to bottom
    _scrollController.addListener(() {
      if (_scrollController.position.pixels ==
          _scrollController.position.maxScrollExtent) {
        _loadMoreMarkets();
      }
    });
  }

  Future<void> _loadMarkets() async {
    setState(() => _isLoading = true);

    final markets = await apiService.getMarkets(page: _currentPage);

    setState(() {
      _markets.addAll(markets);
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      controller: _scrollController,
      itemCount: _markets.length + (_isLoading ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == _markets.length) {
          return CircularProgressIndicator();
        }
        return MarketListItem(market: _markets[index]);
      },
    );
  }
}
```

### React Admin Panel

```typescript
// src/hooks/useInfiniteScroll.ts
import { useState, useEffect, useCallback } from 'react';

export function useInfiniteScroll<T>(
  fetchFunction: (page: number) => Promise<T[]>,
  options?: { pageSize?: number }
) {
  const [items, setItems] = useState<T[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const loadMore = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    try {
      const newItems = await fetchFunction(page);

      if (newItems.length === 0) {
        setHasMore(false);
      } else {
        setItems(prev => [...prev, ...newItems]);
        setPage(prev => prev + 1);
      }
    } catch (error) {
      console.error('Failed to load items:', error);
    } finally {
      setLoading(false);
    }
  }, [page, loading, hasMore, fetchFunction]);

  useEffect(() => {
    loadMore();
  }, []);

  return { items, loading, hasMore, loadMore };
}

// Usage
function MarketList() {
  const { items: markets, loading, hasMore, loadMore } = useInfiniteScroll(
    (page) => api.getMarkets({ page, pageSize: 20 })
  );

  return (
    <div>
      {markets.map(market => (
        <MarketCard key={market.id} market={market} />
      ))}
      {hasMore && (
        <button onClick={loadMore} disabled={loading}>
          {loading ? 'Loading...' : 'Load More'}
        </button>
      )}
    </div>
  );
}

// Memoization for expensive computations
import { useMemo } from 'react';

function MarketAnalytics({ orders }: { orders: Order[] }) {
  const analytics = useMemo(() => {
    return calculateExpensiveAnalytics(orders);
  }, [orders]); // Only recalculate when orders change

  return <div>{analytics}</div>;
}

// Code splitting
const AdminPanel = React.lazy(() => import('./AdminPanel'));

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <AdminPanel />
    </Suspense>
  );
}
```

---

## Infrastructure Optimization

### Horizontal Pod Autoscaling

```yaml
# kubernetes/hpa/user-service-hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: user-service-hpa
  namespace: lfg
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: user-service
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  - type: Pods
    pods:
      metric:
        name: http_requests_per_second
      target:
        type: AverageValue
        averageValue: "1000"
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 100
        periodSeconds: 30
      - type: Pods
        value: 4
        periodSeconds: 30
      selectPolicy: Max
```

### Resource Optimization

```yaml
# kubernetes/deployments/user-service.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  template:
    spec:
      containers:
      - name: user-service
        image: lfg/user-service:latest
        resources:
          requests:
            # Minimum guaranteed resources
            memory: "128Mi"
            cpu: "100m"
          limits:
            # Maximum allowed resources
            memory: "512Mi"
            cpu: "500m"
        # Optimize startup time
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
```

---

## Monitoring & Profiling

### CPU and Memory Profiling

```go
// backend/shared/profiling/profiler.go
package profiling

import (
    "net/http"
    _ "net/http/pprof"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

// Enable pprof profiling endpoint
func EnableProfiling(port string) {
    go func() {
        http.ListenAndServe(":"+port, nil)
    }()
}

// Profile CPU usage
func ProfileCPU(duration time.Duration, filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    if err := pprof.StartCPUProfile(f); err != nil {
        return err
    }
    defer pprof.StopCPUProfile()

    time.Sleep(duration)
    return nil
}

// Profile memory usage
func ProfileMemory(filename string) error {
    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    runtime.GC() // Force GC before profiling
    return pprof.WriteHeapProfile(f)
}

// Usage in main.go
func main() {
    if os.Getenv("ENABLE_PROFILING") == "true" {
        profiling.EnableProfiling("6060")
    }

    // ... rest of main
}
```

### Continuous Benchmarking

```go
// backend/matching-engine/engine/matching_engine_bench_test.go
package engine

import (
    "testing"
    "time"
)

func BenchmarkOrderMatching(b *testing.B) {
    ob := NewOrderBook()

    // Prepare test orders
    buyOrders := generateOrders(1000, "BUY")
    sellOrders := generateOrders(1000, "SELL")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        for _, order := range buyOrders {
            ob.AddOrder(order)
        }
        for _, order := range sellOrders {
            ob.AddOrder(order)
        }
        ob.Match()
    }
}

func BenchmarkConcurrentOrderProcessing(b *testing.B) {
    engine := NewPartitionedMatchingEngine()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            order := generateRandomOrder()
            engine.ProcessOrder(order)
        }
    })
}

// Run benchmarks with: go test -bench=. -benchmem -cpuprofile=cpu.prof
// Analyze with: go tool pprof cpu.prof
```

---

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| API Latency (p50) | <100ms | Prometheus |
| API Latency (p95) | <500ms | Prometheus |
| API Latency (p99) | <1000ms | Prometheus |
| Order Matching | 1000+ orders/sec | Benchmarks |
| WebSocket Connections | 10,000+ concurrent | Load test |
| Database Queries (p95) | <100ms | pganalyze |
| Cache Hit Rate | >70% | Redis metrics |
| Error Rate | <1% | Prometheus |
| Uptime | 99.9% | Uptime monitor |

---

*Version: 1.0*
*Last Updated: 2025-11-18*
*Next Review: After Phase 3 optimization*
