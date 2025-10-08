package engine

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Order represents a single order in the order book.
type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Side      string // "BUY" or "SELL"
	Price     float64
	Quantity  int
	Timestamp time.Time
}

// OrderBook represents the in-memory order book for a single market.
type OrderBook struct {
	Bids []*Order
	Asks []*Order
	mu   sync.Mutex
}

// NewOrderBook creates a new OrderBook.
func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: make([]*Order, 0),
		Asks: make([]*Order, 0),
	}
}

// AddOrder adds a new order to the order book.
// This is a placeholder. The actual implementation will involve matching logic.
func (ob *OrderBook) AddOrder(order *Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// TODO: Implement event sourcing - persist order to NATS before processing.

	if order.Side == "BUY" {
		ob.Bids = append(ob.Bids, order)
		// TODO: Sort bids by price descending
	} else {
		ob.Asks = append(ob.Asks, order)
		// TODO: Sort asks by price ascending
	}

	// TODO: Implement matching logic here.
}
