package engine

import (
	"container/heap"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Order represents a trading order
type Order struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MarketID        uuid.UUID
	ContractID      uuid.UUID
	Side            string // "BUY" or "SELL"
	Type            string // "MARKET" or "LIMIT"
	Quantity        int
	FilledQuantity  int
	Price           float64
	CreatedAt       time.Time
	Index           int // heap index
}

// RemainingQuantity returns unfilled quantity
func (o *Order) RemainingQuantity() int {
	return o.Quantity - o.FilledQuantity
}

// IsFilled checks if order is completely filled
func (o *Order) IsFilled() bool {
	return o.FilledQuantity >= o.Quantity
}

// Trade represents a matched trade
type Trade struct {
	ID          uuid.UUID
	MarketID    uuid.UUID
	BuyOrderID  uuid.UUID
	SellOrderID uuid.UUID
	Quantity    int
	Price       float64
	Timestamp   time.Time
}

// OrderBook implements an efficient order book using priority queues
type OrderBook struct {
	mu         sync.RWMutex
	buyOrders  *BuyHeap  // Max heap (highest price first)
	sellOrders *SellHeap // Min heap (lowest price first)
	orderIndex map[uuid.UUID]*Order
	marketID   uuid.UUID
}

// NewOrderBook creates a new order book for a market
func NewOrderBook(marketID uuid.UUID) *OrderBook {
	buyHeap := &BuyHeap{}
	sellHeap := &SellHeap{}
	heap.Init(buyHeap)
	heap.Init(sellHeap)

	return &OrderBook{
		buyOrders:  buyHeap,
		sellOrders: sellHeap,
		orderIndex: make(map[uuid.UUID]*Order),
		marketID:   marketID,
	}
}

// AddOrder adds an order to the book
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

// CancelOrder removes an order from the book
func (ob *OrderBook) CancelOrder(orderID uuid.UUID) bool {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	order, exists := ob.orderIndex[orderID]
	if !exists {
		return false
	}

	delete(ob.orderIndex, orderID)

	if order.Side == "BUY" {
		heap.Remove(ob.buyOrders, order.Index)
	} else {
		heap.Remove(ob.sellOrders, order.Index)
	}

	return true
}

// Match executes order matching and returns resulting trades
func (ob *OrderBook) Match() []*Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := make([]*Trade, 0)

	// Keep matching while there's overlap in prices
	for ob.buyOrders.Len() > 0 && ob.sellOrders.Len() > 0 {
		topBuy := (*ob.buyOrders)[0]
		topSell := (*ob.sellOrders)[0]

		// Check if prices cross (buy price >= sell price)
		if topBuy.Price < topSell.Price {
			break // No more matches possible
		}

		// Execute trade at seller's price (price-time priority)
		tradeQuantity := min(topBuy.RemainingQuantity(), topSell.RemainingQuantity())
		tradePrice := topSell.Price

		trade := &Trade{
			ID:          uuid.New(),
			MarketID:    ob.marketID,
			BuyOrderID:  topBuy.ID,
			SellOrderID: topSell.ID,
			Quantity:    tradeQuantity,
			Price:       tradePrice,
			Timestamp:   time.Now(),
		}
		trades = append(trades, trade)

		// Update filled quantities
		topBuy.FilledQuantity += tradeQuantity
		topSell.FilledQuantity += tradeQuantity

		// Remove fully filled orders
		if topBuy.IsFilled() {
			heap.Pop(ob.buyOrders)
			delete(ob.orderIndex, topBuy.ID)
		} else {
			// Update heap to reflect changes
			heap.Fix(ob.buyOrders, topBuy.Index)
		}

		if topSell.IsFilled() {
			heap.Pop(ob.sellOrders)
			delete(ob.orderIndex, topSell.ID)
		} else {
			// Update heap to reflect changes
			heap.Fix(ob.sellOrders, topSell.Index)
		}
	}

	return trades
}

// GetTopBuy returns the best buy order without removing it
func (ob *OrderBook) GetTopBuy() *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.buyOrders.Len() == 0 {
		return nil
	}

	return (*ob.buyOrders)[0]
}

// GetTopSell returns the best sell order without removing it
func (ob *OrderBook) GetTopSell() *Order {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.sellOrders.Len() == 0 {
		return nil
	}

	return (*ob.sellOrders)[0]
}

// GetSpread returns the bid-ask spread
func (ob *OrderBook) GetSpread() float64 {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	if ob.buyOrders.Len() == 0 || ob.sellOrders.Len() == 0 {
		return 0
	}

	return (*ob.sellOrders)[0].Price - (*ob.buyOrders)[0].Price
}

// BuyHeap implements a max heap for buy orders (highest price first)
type BuyHeap []*Order

func (h BuyHeap) Len() int { return len(h) }

func (h BuyHeap) Less(i, j int) bool {
	// Primary: price (higher is better for buys)
	if h[i].Price != h[j].Price {
		return h[i].Price > h[j].Price
	}
	// Secondary: time priority (earlier is better)
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h BuyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *BuyHeap) Push(x interface{}) {
	n := len(*h)
	order := x.(*Order)
	order.Index = n
	*h = append(*h, order)
}

func (h *BuyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	order := old[n-1]
	old[n-1] = nil
	order.Index = -1
	*h = old[0 : n-1]
	return order
}

// SellHeap implements a min heap for sell orders (lowest price first)
type SellHeap []*Order

func (h SellHeap) Len() int { return len(h) }

func (h SellHeap) Less(i, j int) bool {
	// Primary: price (lower is better for sells)
	if h[i].Price != h[j].Price {
		return h[i].Price < h[j].Price
	}
	// Secondary: time priority (earlier is better)
	return h[i].CreatedAt.Before(h[j].CreatedAt)
}

func (h SellHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

func (h *SellHeap) Push(x interface{}) {
	n := len(*h)
	order := x.(*Order)
	order.Index = n
	*h = append(*h, order)
}

func (h *SellHeap) Pop() interface{} {
	old := *h
	n := len(old)
	order := old[n-1]
	old[n-1] = nil
	order.Index = -1
	*h = old[0 : n-1]
	return order
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
