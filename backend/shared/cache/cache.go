package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"lfg/shared/metrics"
)

// Cache wraps Redis client with convenience methods
type Cache struct {
	client *redis.Client
}

// Config holds cache configuration
type Config struct {
	Addr     string
	Password string
	DB       int
}

// NewCache creates a new Redis cache client
func NewCache(cfg Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     50,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Get retrieves a value from cache and unmarshals it
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		metrics.RecordCacheMiss("redis", keyPrefix(key))
		return ErrCacheMiss
	}
	if err != nil {
		metrics.RecordCacheMiss("redis", keyPrefix(key))
		return fmt.Errorf("cache get failed: %w", err)
	}

	metrics.RecordCacheHit("redis", keyPrefix(key))

	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set failed: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("cache delete failed: %w", err)
	}

	return nil
}

// DeletePattern deletes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("cache scan failed: %w", err)
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		return c.Delete(ctx, keys...)
	}

	return nil
}

// Exists checks if a key exists
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache exists check failed: %w", err)
	}

	return count > 0, nil
}

// Increment increments a counter
func (c *Cache) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("cache increment failed: %w", err)
	}

	return val, nil
}

// Expire sets TTL on an existing key
func (c *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if err := c.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("cache expire failed: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// Ping checks if Redis is accessible
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// keyPrefix extracts prefix from cache key for metrics
func keyPrefix(key string) string {
	// Extract first segment before ":"
	for i, ch := range key {
		if ch == ':' {
			return key[:i]
		}
	}
	return "unknown"
}

// ErrCacheMiss is returned when key is not found
var ErrCacheMiss = fmt.Errorf("cache miss")
