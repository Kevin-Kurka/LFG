package engine

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Engine manages multiple order books (one per contract)
type Engine struct {
	orderBooks map[uuid.UUID]*OrderBook
	mu         sync.RWMutex
}

// NewEngine creates a new matching engine
func NewEngine() *Engine {
	return &Engine{
		orderBooks: make(map[uuid.UUID]*OrderBook),
	}
}

// GetOrCreateOrderBook returns an existing order book or creates a new one
func (e *Engine) GetOrCreateOrderBook(contractID uuid.UUID) *OrderBook {
	e.mu.Lock()
	defer e.mu.Unlock()

	if ob, exists := e.orderBooks[contractID]; exists {
		return ob
	}

	ob := NewOrderBook(contractID)
	e.orderBooks[contractID] = ob
	log.Printf("Created new order book for contract %s", contractID)
	return ob
}

// SubmitOrder submits an order to the matching engine
func (e *Engine) SubmitOrder(order *Order) ([]Trade, error) {
	ob := e.GetOrCreateOrderBook(order.ContractID)
	trades := ob.AddOrder(order)

	// Persist trades to database
	if len(trades) > 0 {
		ctx := context.Background()
		if err := ob.PersistTrades(ctx, trades); err != nil {
			log.Printf("Failed to persist trades: %v", err)
			return trades, err
		}
	}

	return trades, nil
}

// CancelOrder cancels an order
func (e *Engine) CancelOrder(contractID, orderID uuid.UUID) bool {
	e.mu.RLock()
	ob, exists := e.orderBooks[contractID]
	e.mu.RUnlock()

	if !exists {
		return false
	}

	return ob.CancelOrder(orderID)
}

// GetOrderBook returns a snapshot of the order book for a contract
func (e *Engine) GetOrderBook(contractID uuid.UUID) ([]Order, []Order, bool) {
	e.mu.RLock()
	ob, exists := e.orderBooks[contractID]
	e.mu.RUnlock()

	if !exists {
		return nil, nil, false
	}

	bids, asks := ob.GetSnapshot()
	return bids, asks, true
}

// Global engine instance
var GlobalEngine *Engine

func init() {
	GlobalEngine = NewEngine()
	log.Println("Matching engine initialized")
}
