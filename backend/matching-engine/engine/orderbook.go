package engine

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/google/uuid"
)

// Order represents a single order in the order book
type Order struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	ContractID uuid.UUID `json:"contract_id"`
	Side       string    `json:"side"` // "BUY" or "SELL"
	Price      float64   `json:"price"`
	Quantity   int       `json:"quantity"`
	Filled     int       `json:"filled"`
	Timestamp  time.Time `json:"timestamp"`
}

// Trade represents an executed trade
type Trade struct {
	ID           uuid.UUID `json:"id"`
	ContractID   uuid.UUID `json:"contract_id"`
	MakerOrderID uuid.UUID `json:"maker_order_id"`
	TakerOrderID uuid.UUID `json:"taker_order_id"`
	MakerUserID  uuid.UUID `json:"maker_user_id"`
	TakerUserID  uuid.UUID `json:"taker_user_id"`
	Price        float64   `json:"price"`
	Quantity     int       `json:"quantity"`
	ExecutedAt   time.Time `json:"executed_at"`
}

// OrderBook represents the in-memory order book for a single contract
type OrderBook struct {
	ContractID uuid.UUID
	Bids       []*Order // Buy orders, sorted by price descending, then time ascending
	Asks       []*Order // Sell orders, sorted by price ascending, then time ascending
	mu         sync.RWMutex
}

// NewOrderBook creates a new OrderBook
func NewOrderBook(contractID uuid.UUID) *OrderBook {
	return &OrderBook{
		ContractID: contractID,
		Bids:       make([]*Order, 0),
		Asks:       make([]*Order, 0),
	}
}

// AddOrder adds a new order to the order book and attempts to match it
func (ob *OrderBook) AddOrder(order *Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := make([]Trade, 0)

	if order.Side == "BUY" {
		// Try to match with asks
		trades = ob.matchBuyOrder(order)

		// If order not fully filled, add to bids
		if order.Filled < order.Quantity {
			ob.Bids = append(ob.Bids, order)
			ob.sortBids()
		}
	} else {
		// Try to match with bids
		trades = ob.matchSellOrder(order)

		// If order not fully filled, add to asks
		if order.Filled < order.Quantity {
			ob.Asks = append(ob.Asks, order)
			ob.sortAsks()
		}
	}

	return trades
}

// matchBuyOrder matches a buy order against existing ask orders
func (ob *OrderBook) matchBuyOrder(buyOrder *Order) []Trade {
	trades := make([]Trade, 0)

	for len(ob.Asks) > 0 && buyOrder.Filled < buyOrder.Quantity {
		askOrder := ob.Asks[0]

		// Check if prices cross (buy price >= ask price)
		if buyOrder.Price < askOrder.Price {
			break
		}

		// Calculate trade quantity
		remainingBuy := buyOrder.Quantity - buyOrder.Filled
		remainingAsk := askOrder.Quantity - askOrder.Filled
		tradeQuantity := remainingBuy
		if remainingAsk < tradeQuantity {
			tradeQuantity = remainingAsk
		}

		// Create trade (maker gets their price)
		trade := Trade{
			ID:           uuid.New(),
			ContractID:   ob.ContractID,
			MakerOrderID: askOrder.ID,
			TakerOrderID: buyOrder.ID,
			MakerUserID:  askOrder.UserID,
			TakerUserID:  buyOrder.UserID,
			Price:        askOrder.Price, // Maker price
			Quantity:     tradeQuantity,
			ExecutedAt:   time.Now(),
		}
		trades = append(trades, trade)

		// Update filled quantities
		buyOrder.Filled += tradeQuantity
		askOrder.Filled += tradeQuantity

		// Remove fully filled ask order
		if askOrder.Filled >= askOrder.Quantity {
			ob.Asks = ob.Asks[1:]
		}
	}

	return trades
}

// matchSellOrder matches a sell order against existing bid orders
func (ob *OrderBook) matchSellOrder(sellOrder *Order) []Trade {
	trades := make([]Trade, 0)

	for len(ob.Bids) > 0 && sellOrder.Filled < sellOrder.Quantity {
		bidOrder := ob.Bids[0]

		// Check if prices cross (sell price <= bid price)
		if sellOrder.Price > bidOrder.Price {
			break
		}

		// Calculate trade quantity
		remainingSell := sellOrder.Quantity - sellOrder.Filled
		remainingBid := bidOrder.Quantity - bidOrder.Filled
		tradeQuantity := remainingSell
		if remainingBid < tradeQuantity {
			tradeQuantity = remainingBid
		}

		// Create trade (maker gets their price)
		trade := Trade{
			ID:           uuid.New(),
			ContractID:   ob.ContractID,
			MakerOrderID: bidOrder.ID,
			TakerOrderID: sellOrder.ID,
			MakerUserID:  bidOrder.UserID,
			TakerUserID:  sellOrder.UserID,
			Price:        bidOrder.Price, // Maker price
			Quantity:     tradeQuantity,
			ExecutedAt:   time.Now(),
		}
		trades = append(trades, trade)

		// Update filled quantities
		sellOrder.Filled += tradeQuantity
		bidOrder.Filled += tradeQuantity

		// Remove fully filled bid order
		if bidOrder.Filled >= bidOrder.Quantity {
			ob.Bids = ob.Bids[1:]
		}
	}

	return trades
}

// CancelOrder removes an order from the order book
func (ob *OrderBook) CancelOrder(orderID uuid.UUID) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Try to find and remove from bids
	for i, order := range ob.Bids {
		if order.ID == orderID {
			ob.Bids = append(ob.Bids[:i], ob.Bids[i+1:]...)
			return true
		}
	}

	// Try to find and remove from asks
	for i, order := range ob.Asks {
		if order.ID == orderID {
			ob.Asks = append(ob.Asks[:i], ob.Asks[i+1:]...)
			return true
		}
	}

	return false
}

// GetSnapshot returns a snapshot of the order book
func (ob *OrderBook) GetSnapshot() ([]Order, []Order) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	bids := make([]Order, len(ob.Bids))
	for i, order := range ob.Bids {
		bids[i] = *order
	}

	asks := make([]Order, len(ob.Asks))
	for i, order := range ob.Asks {
		asks[i] = *order
	}

	return bids, asks
}

// sortBids sorts bids by price descending (highest first), then by time ascending
func (ob *OrderBook) sortBids() {
	sort.Slice(ob.Bids, func(i, j int) bool {
		if ob.Bids[i].Price != ob.Bids[j].Price {
			return ob.Bids[i].Price > ob.Bids[j].Price
		}
		return ob.Bids[i].Timestamp.Before(ob.Bids[j].Timestamp)
	})
}

// sortAsks sorts asks by price ascending (lowest first), then by time ascending
func (ob *OrderBook) sortAsks() {
	sort.Slice(ob.Asks, func(i, j int) bool {
		if ob.Asks[i].Price != ob.Asks[j].Price {
			return ob.Asks[i].Price < ob.Asks[j].Price
		}
		return ob.Asks[i].Timestamp.Before(ob.Asks[j].Timestamp)
	})
}

// PersistTrades saves trades to the database
func (ob *OrderBook) PersistTrades(ctx context.Context, trades []Trade) error {
	if len(trades) == 0 {
		return nil
	}

	tx, err := database.GetDB().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, trade := range trades {
		// Insert trade
		_, err := tx.Exec(ctx, `
			INSERT INTO trades (id, contract_id, maker_order_id, taker_order_id, quantity, price_credits, executed_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, trade.ID, trade.ContractID, trade.MakerOrderID, trade.TakerOrderID, trade.Quantity, trade.Price, trade.ExecutedAt)
		if err != nil {
			log.Printf("Failed to insert trade: %v", err)
			return err
		}

		// Update maker order
		_, err = tx.Exec(ctx, `
			UPDATE orders
			SET quantity_filled = quantity_filled + $1,
			    status = CASE WHEN quantity_filled + $1 >= quantity THEN 'FILLED' ELSE 'PARTIAL' END
			WHERE id = $2
		`, trade.Quantity, trade.MakerOrderID)
		if err != nil {
			log.Printf("Failed to update maker order: %v", err)
			return err
		}

		// Update taker order
		_, err = tx.Exec(ctx, `
			UPDATE orders
			SET quantity_filled = quantity_filled + $1,
			    status = CASE WHEN quantity_filled + $1 >= quantity THEN 'FILLED' ELSE 'PARTIAL' END
			WHERE id = $2
		`, trade.Quantity, trade.TakerOrderID)
		if err != nil {
			log.Printf("Failed to update taker order: %v", err)
			return err
		}
	}

	return tx.Commit(ctx)
}
