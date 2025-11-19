package engine

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	pb "lfg/matching-engine/proto"
)

// MatchingEngine manages the order books for all contracts
type MatchingEngine struct {
	OrderBooks map[string]*OrderBook // Map of contract ID to OrderBook
	mu         sync.RWMutex
	natsConn   *nats.Conn
	pb.UnimplementedMatchingEngineServer
}

// NewMatchingEngine creates a new matching engine
func NewMatchingEngine(natsConn *nats.Conn) *MatchingEngine {
	return &MatchingEngine{
		OrderBooks: make(map[string]*OrderBook),
		natsConn:   natsConn,
	}
}

// GetOrCreateOrderBook retrieves an existing order book or creates a new one
func (me *MatchingEngine) GetOrCreateOrderBook(contractID string) *OrderBook {
	me.mu.Lock()
	defer me.mu.Unlock()

	if ob, ok := me.OrderBooks[contractID]; ok {
		return ob
	}

	newOB := NewOrderBook(contractID)
	me.OrderBooks[contractID] = newOB
	return newOB
}

// PlaceOrder implements the gRPC PlaceOrder method
func (me *MatchingEngine) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	// Get or create order book for contract
	orderBook := me.GetOrCreateOrderBook(req.ContractId)

	// Create order
	order := &Order{
		ID:         req.OrderId,
		UserID:     req.UserId,
		ContractID: req.ContractId,
		Type:       req.Type,
		Side:       req.Side,
		Quantity:   int(req.Quantity),
		LimitPrice: req.LimitPrice,
		Timestamp:  time.Now(),
	}

	// Add order to book and match
	trades, quantityFilled, status := orderBook.AddOrder(order)

	// Convert trades to protobuf format
	pbTrades := make([]*pb.Trade, len(trades))
	totalValue := 0.0
	for i, trade := range trades {
		pbTrades[i] = &pb.Trade{
			TradeId:       trade.ID,
			MakerOrderId:  trade.MakerOrderID,
			TakerOrderId:  trade.TakerOrderID,
			Quantity:      int32(trade.Quantity),
			Price:         trade.Price,
			ExecutedAt:    trade.ExecutedAt.Unix(),
		}
		totalValue += float64(trade.Quantity) * trade.Price
	}

	// Calculate average price
	averagePrice := 0.0
	if quantityFilled > 0 {
		averagePrice = totalValue / float64(quantityFilled)
	}

	// Publish trade events to NATS
	if me.natsConn != nil && len(trades) > 0 {
		for _, trade := range trades {
			tradeEvent := map[string]interface{}{
				"trade_id":       trade.ID,
				"contract_id":    trade.ContractID,
				"maker_order_id": trade.MakerOrderID,
				"taker_order_id": trade.TakerOrderID,
				"maker_user_id":  trade.MakerUserID,
				"taker_user_id":  trade.TakerUserID,
				"quantity":       trade.Quantity,
				"price":          trade.Price,
				"executed_at":    trade.ExecutedAt.Unix(),
			}

			eventJSON, err := json.Marshal(tradeEvent)
			if err != nil {
				log.Printf("Failed to marshal trade event: %v", err)
				continue
			}

			// Publish to trades topic
			if err := me.natsConn.Publish("trades", eventJSON); err != nil {
				log.Printf("Failed to publish trade event: %v", err)
			} else {
				log.Printf("Published trade event: %s", trade.ID)
			}
		}
	}

	return &pb.PlaceOrderResponse{
		OrderId:        req.OrderId,
		Status:         status,
		QuantityFilled: int32(quantityFilled),
		AveragePrice:   averagePrice,
		Trades:         pbTrades,
	}, nil
}

// CancelOrder implements the gRPC CancelOrder method
func (me *MatchingEngine) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	me.mu.RLock()
	orderBook, exists := me.OrderBooks[req.ContractId]
	me.mu.RUnlock()

	if !exists {
		return &pb.CancelOrderResponse{
			Success: false,
			Message: "Order book not found",
		}, nil
	}

	success := orderBook.CancelOrder(req.OrderId)
	message := "Order cancelled successfully"
	if !success {
		message = "Order not found or already filled"
	}

	return &pb.CancelOrderResponse{
		Success: success,
		Message: message,
	}, nil
}

// GetOrderBook implements the gRPC GetOrderBook method
func (me *MatchingEngine) GetOrderBook(ctx context.Context, req *pb.GetOrderBookRequest) (*pb.GetOrderBookResponse, error) {
	me.mu.RLock()
	orderBook, exists := me.OrderBooks[req.ContractId]
	me.mu.RUnlock()

	if !exists {
		return &pb.GetOrderBookResponse{
			Bids: []*pb.OrderBookLevel{},
			Asks: []*pb.OrderBookLevel{},
		}, nil
	}

	depth := int(req.Depth)
	if depth == 0 {
		depth = 10 // Default depth
	}

	bids, asks := orderBook.GetAggregatedBook(depth)

	// Convert to protobuf format
	pbBids := make([]*pb.OrderBookLevel, len(bids))
	for i, level := range bids {
		pbBids[i] = &pb.OrderBookLevel{
			Price:      level.Price,
			Quantity:   int32(level.Quantity),
			OrderCount: int32(level.OrderCount),
		}
	}

	pbAsks := make([]*pb.OrderBookLevel, len(asks))
	for i, level := range asks {
		pbAsks[i] = &pb.OrderBookLevel{
			Price:      level.Price,
			Quantity:   int32(level.Quantity),
			OrderCount: int32(level.OrderCount),
		}
	}

	return &pb.GetOrderBookResponse{
		Bids: pbBids,
		Asks: pbAsks,
	}, nil
}

// Trade represents a matched trade
type Trade struct {
	ID            string
	ContractID    string
	MakerOrderID  string
	TakerOrderID  string
	MakerUserID   string
	TakerUserID   string
	Quantity      int
	Price         float64
	ExecutedAt    time.Time
}

// NewTrade creates a new trade
func NewTrade(contractID, makerOrderID, takerOrderID, makerUserID, takerUserID string, quantity int, price float64) *Trade {
	return &Trade{
		ID:           uuid.New().String(),
		ContractID:   contractID,
		MakerOrderID: makerOrderID,
		TakerOrderID: takerOrderID,
		MakerUserID:  makerUserID,
		TakerUserID:  takerUserID,
		Quantity:     quantity,
		Price:        price,
		ExecutedAt:   time.Now(),
	}
}
