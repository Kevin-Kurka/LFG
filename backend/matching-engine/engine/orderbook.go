package engine

import (
	"sort"
	"sync"
	"time"

	pb "lfg/matching-engine/proto"
)

// Order represents a single order in the order book
type Order struct {
	ID         string
	UserID     string
	ContractID string
	Type       pb.OrderType
	Side       pb.OrderSide
	Quantity   int
	Filled     int
	LimitPrice float64
	Timestamp  time.Time
}

// OrderBook represents the in-memory order book for a single contract
type OrderBook struct {
	ContractID string
	Bids       []*Order // Buy orders (sorted high to low)
	Asks       []*Order // Sell orders (sorted low to high)
	mu         sync.Mutex
}

// NewOrderBook creates a new OrderBook
func NewOrderBook(contractID string) *OrderBook {
	return &OrderBook{
		ContractID: contractID,
		Bids:       make([]*Order, 0),
		Asks:       make([]*Order, 0),
	}
}

// AddOrder adds a new order to the order book and attempts to match it
func (ob *OrderBook) AddOrder(order *Order) ([]*Trade, int, string) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := []*Trade{}
	quantityFilled := 0

	// Market orders execute at any price
	// Limit orders execute at limit price or better
	if order.Side == pb.OrderSide_BUY {
		// Match against asks (sell orders)
		trades, quantityFilled = ob.matchBuyOrder(order)
	} else {
		// Match against bids (buy orders)
		trades, quantityFilled = ob.matchSellOrder(order)
	}

	order.Filled = quantityFilled

	// Determine order status
	status := "ACTIVE"
	if quantityFilled == order.Quantity {
		status = "FILLED"
	} else if quantityFilled > 0 {
		status = "PARTIALLY_FILLED"
	}

	// If not fully filled and it's a limit order, add to book
	if quantityFilled < order.Quantity && order.Type == pb.OrderType_LIMIT {
		order.Quantity = order.Quantity - quantityFilled // Remaining quantity
		if order.Side == pb.OrderSide_BUY {
			ob.Bids = append(ob.Bids, order)
			ob.sortBids()
		} else {
			ob.Asks = append(ob.Asks, order)
			ob.sortAsks()
		}
	}

	return trades, quantityFilled, status
}

// matchBuyOrder matches a buy order against the ask side
func (ob *OrderBook) matchBuyOrder(order *Order) ([]*Trade, int) {
	trades := []*Trade{}
	quantityFilled := 0
	remaining := order.Quantity

	for i := 0; i < len(ob.Asks) && remaining > 0; i++ {
		ask := ob.Asks[i]

		// For limit orders, only match at limit price or better
		if order.Type == pb.OrderType_LIMIT && ask.LimitPrice > order.LimitPrice {
			break // No more matches possible
		}

		// Calculate match quantity
		matchQty := remaining
		if ask.Quantity-ask.Filled < matchQty {
			matchQty = ask.Quantity - ask.Filled
		}

		// Create trade at the maker's price (ask price)
		trade := NewTrade(
			order.ContractID,
			ask.ID,        // Maker
			order.ID,      // Taker
			ask.UserID,    // Maker user
			order.UserID,  // Taker user
			matchQty,
			ask.LimitPrice, // Execute at ask price
		)
		trades = append(trades, trade)

		// Update quantities
		ask.Filled += matchQty
		remaining -= matchQty
		quantityFilled += matchQty

		// If ask is fully filled, mark for removal
		if ask.Filled >= ask.Quantity {
			ob.Asks[i] = nil
		}
	}

	// Remove fully filled orders
	ob.cleanupAsks()

	return trades, quantityFilled
}

// matchSellOrder matches a sell order against the bid side
func (ob *OrderBook) matchSellOrder(order *Order) ([]*Trade, int) {
	trades := []*Trade{}
	quantityFilled := 0
	remaining := order.Quantity

	for i := 0; i < len(ob.Bids) && remaining > 0; i++ {
		bid := ob.Bids[i]

		// For limit orders, only match at limit price or better
		if order.Type == pb.OrderType_LIMIT && bid.LimitPrice < order.LimitPrice {
			break // No more matches possible
		}

		// Calculate match quantity
		matchQty := remaining
		if bid.Quantity-bid.Filled < matchQty {
			matchQty = bid.Quantity - bid.Filled
		}

		// Create trade at the maker's price (bid price)
		trade := NewTrade(
			order.ContractID,
			bid.ID,        // Maker
			order.ID,      // Taker
			bid.UserID,    // Maker user
			order.UserID,  // Taker user
			matchQty,
			bid.LimitPrice, // Execute at bid price
		)
		trades = append(trades, trade)

		// Update quantities
		bid.Filled += matchQty
		remaining -= matchQty
		quantityFilled += matchQty

		// If bid is fully filled, mark for removal
		if bid.Filled >= bid.Quantity {
			ob.Bids[i] = nil
		}
	}

	// Remove fully filled orders
	ob.cleanupBids()

	return trades, quantityFilled
}

// CancelOrder removes an order from the book
func (ob *OrderBook) CancelOrder(orderID string) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// Search in bids
	for i, order := range ob.Bids {
		if order != nil && order.ID == orderID {
			ob.Bids[i] = nil
			ob.cleanupBids()
			return true
		}
	}

	// Search in asks
	for i, order := range ob.Asks {
		if order != nil && order.ID == orderID {
			ob.Asks[i] = nil
			ob.cleanupAsks()
			return true
		}
	}

	return false
}

// GetAggregatedBook returns aggregated price levels
func (ob *OrderBook) GetAggregatedBook(depth int) ([]PriceLevel, []PriceLevel) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	bidLevels := ob.aggregateSide(ob.Bids, depth)
	askLevels := ob.aggregateSide(ob.Asks, depth)

	return bidLevels, askLevels
}

// PriceLevel represents an aggregated price level
type PriceLevel struct {
	Price      float64
	Quantity   int
	OrderCount int
}

// aggregateSide aggregates orders by price level
func (ob *OrderBook) aggregateSide(orders []*Order, depth int) []PriceLevel {
	priceMap := make(map[float64]*PriceLevel)

	for _, order := range orders {
		if order == nil {
			continue
		}

		remaining := order.Quantity - order.Filled
		if remaining <= 0 {
			continue
		}

		if level, exists := priceMap[order.LimitPrice]; exists {
			level.Quantity += remaining
			level.OrderCount++
		} else {
			priceMap[order.LimitPrice] = &PriceLevel{
				Price:      order.LimitPrice,
				Quantity:   remaining,
				OrderCount: 1,
			}
		}
	}

	// Convert map to slice
	levels := make([]PriceLevel, 0, len(priceMap))
	for _, level := range priceMap {
		levels = append(levels, *level)
	}

	// Limit to depth
	if len(levels) > depth {
		levels = levels[:depth]
	}

	return levels
}

// sortBids sorts bids by price (high to low), then by time (earlier first)
func (ob *OrderBook) sortBids() {
	sort.Slice(ob.Bids, func(i, j int) bool {
		if ob.Bids[i] == nil {
			return false
		}
		if ob.Bids[j] == nil {
			return true
		}

		// Higher price first
		if ob.Bids[i].LimitPrice != ob.Bids[j].LimitPrice {
			return ob.Bids[i].LimitPrice > ob.Bids[j].LimitPrice
		}

		// Earlier timestamp first (price-time priority)
		return ob.Bids[i].Timestamp.Before(ob.Bids[j].Timestamp)
	})
}

// sortAsks sorts asks by price (low to high), then by time (earlier first)
func (ob *OrderBook) sortAsks() {
	sort.Slice(ob.Asks, func(i, j int) bool {
		if ob.Asks[i] == nil {
			return false
		}
		if ob.Asks[j] == nil {
			return true
		}

		// Lower price first
		if ob.Asks[i].LimitPrice != ob.Asks[j].LimitPrice {
			return ob.Asks[i].LimitPrice < ob.Asks[j].LimitPrice
		}

		// Earlier timestamp first (price-time priority)
		return ob.Asks[i].Timestamp.Before(ob.Asks[j].Timestamp)
	})
}

// cleanupBids removes nil entries from bids
func (ob *OrderBook) cleanupBids() {
	cleaned := make([]*Order, 0, len(ob.Bids))
	for _, order := range ob.Bids {
		if order != nil && order.Filled < order.Quantity {
			cleaned = append(cleaned, order)
		}
	}
	ob.Bids = cleaned
}

// cleanupAsks removes nil entries from asks
func (ob *OrderBook) cleanupAsks() {
	cleaned := make([]*Order, 0, len(ob.Asks))
	for _, order := range ob.Asks {
		if order != nil && order.Filled < order.Quantity {
			cleaned = append(cleaned, order)
		}
	}
	ob.Asks = cleaned
}
