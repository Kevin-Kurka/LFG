package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "endpoint", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "endpoint", "status"},
	)

	// Database metrics
	DBQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"service", "operation", "table"},
	)

	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"service", "operation", "table"},
	)

	// Business metrics
	OrdersPlaced = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_placed_total",
			Help: "Total number of orders placed",
		},
		[]string{"market_id", "order_type", "side"},
	)

	OrdersMatched = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "orders_matched_total",
			Help: "Total number of orders matched",
		},
		[]string{"market_id"},
	)

	MatchingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "matching_duration_seconds",
			Help:    "Duration of order matching in seconds",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"market_id"},
	)

	UsersRegistered = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total number of users registered",
		},
	)

	ActiveWebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Number of active WebSocket connections",
		},
	)

	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type", "key_prefix"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type", "key_prefix"},
	)
)

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(service, method, endpoint, status string, duration time.Duration) {
	HTTPRequestsTotal.WithLabelValues(service, method, endpoint, status).Inc()
	HTTPRequestDuration.WithLabelValues(service, method, endpoint, status).Observe(duration.Seconds())
}

// RecordDBQuery records a database query metric
func RecordDBQuery(service, operation, table string, duration time.Duration) {
	DBQueriesTotal.WithLabelValues(service, operation, table).Inc()
	DBQueryDuration.WithLabelValues(service, operation, table).Observe(duration.Seconds())
}

// RecordOrderPlaced records an order placement metric
func RecordOrderPlaced(marketID, orderType, side string) {
	OrdersPlaced.WithLabelValues(marketID, orderType, side).Inc()
}

// RecordOrderMatched records an order match metric
func RecordOrderMatched(marketID string, duration time.Duration) {
	OrdersMatched.WithLabelValues(marketID).Inc()
	MatchingDuration.WithLabelValues(marketID).Observe(duration.Seconds())
}

// RecordCacheHit records a cache hit
func RecordCacheHit(cacheType, keyPrefix string) {
	CacheHits.WithLabelValues(cacheType, keyPrefix).Inc()
}

// RecordCacheMiss records a cache miss
func RecordCacheMiss(cacheType, keyPrefix string) {
	CacheMisses.WithLabelValues(cacheType, keyPrefix).Inc()
}
