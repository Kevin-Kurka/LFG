package engine

import (
	"sync"
)

// MatchingEngine manages the order books for all markets.
type MatchingEngine struct {
	OrderBooks map[string]*OrderBook // Map of market ticker to OrderBook
	mu         sync.RWMutex
	// NATS connection would go here
}

// NewMatchingEngine creates a new MatchingEngine.
func NewMatchingEngine() *MatchingEngine {
	return &MatchingEngine{
		OrderBooks: make(map[string]*OrderBook),
	}
}

// GetOrCreateOrderBook retrieves an existing order book for a market or creates a new one.
func (me *MatchingEngine) GetOrCreateOrderBook(marketTicker string) *OrderBook {
	me.mu.Lock()
	defer me.mu.Unlock()

	if ob, ok := me.OrderBooks[marketTicker]; ok {
		return ob
	}

	newOB := NewOrderBook()
	me.OrderBooks[marketTicker] = newOB
	return newOB
}

// ProcessOrder is the main entry point for processing an order.
func (me *MatchingEngine) ProcessOrder(marketTicker string, order *Order) {
	orderBook := me.GetOrCreateOrderBook(marketTicker)
	orderBook.AddOrder(order)
	// Further logic will be more complex, involving publishing events.
}
