package exchange

import (
	"fmt"
	"sync/atomic"
)

// Trade represents the transaction resulting from a matching bid and offer.
type Trade struct {
	OrderId  int
	Price    int
	Volume   int
	FillTime int
}

const ringSize = 1 << 10
const ringMask = ringSize - 1

// RingBuffer is a lock-free SPSC queue. head is written only by the producer,
// tail only by the consumer — no CAS needed.
type RingBuffer struct {
	slots [ringSize]Trade
	_     [64]byte       // pad to separate head from slots cache line
	head  atomic.Uint64  // producer index
	_     [56]byte       // pad to separate head and tail cache lines
	tail  atomic.Uint64  // consumer index
}

func (r *RingBuffer) Push(t Trade) bool {
	head := r.head.Load()
	if head-r.tail.Load() >= ringSize {
		return false
	}
	r.slots[head&ringMask] = t
	r.head.Store(head + 1)
	return true
}

func (r *RingBuffer) Pop() (Trade, bool) {
	tail := r.tail.Load()
	if r.head.Load() == tail {
		return Trade{}, false
	}
	t := r.slots[tail&ringMask]
	r.tail.Store(tail + 1)
	return t, true
}

func (r *RingBuffer) Len() int {
	return int(r.head.Load() - r.tail.Load())
}

// MatchingEngine updates the state of the orderbook when new orders come in. It
// also issues trades.
type MatchingEngine struct {
	OrderBook *OrderBook
	Ring      RingBuffer
}

func NewMatchingEngine(ob *OrderBook) MatchingEngine {
	return MatchingEngine{OrderBook: ob}
}

func (engine *MatchingEngine) ProcessOrder(order *Order) {
	if order.IsBuy && order.Price >= engine.OrderBook.BestOffer {
		for price := engine.OrderBook.BestOffer; price <= order.Price; price++ {
			engine.processTrades(order, price)
			if order.Volume == 0 {
				break
			}
		}
	}
	if !order.IsBuy && order.Price <= engine.OrderBook.BestBid {
		for price := engine.OrderBook.BestBid; price >= order.Price; price-- {
			engine.processTrades(order, price)
			if order.Volume == 0 {
				break
			}
		}
	}
	if order.Volume > 0 {
		engine.OrderBook.InsertOrder(order)
	}
}

func (engine *MatchingEngine) processTrades(o *Order, p int) {
	var pl *PriceLevel
	if o.IsBuy {
		pl = engine.OrderBook.asks[p]
	} else {
		pl = engine.OrderBook.bids[p]
	}

	for currentOrder := pl.Head; currentOrder != nil && o.Volume > 0; {
		next := currentOrder.Next
		qty := min(o.Volume, currentOrder.Volume)
		o.Volume -= qty
		currentOrder.Volume -= qty
		pl.TotalVolume -= qty

		engine.Ring.Push(Trade{
			OrderId:  currentOrder.Id,
			Price:    currentOrder.Price,
			Volume:   qty,
			FillTime: 1,
		})

		if currentOrder.Volume == 0 {
			pl.Head = next
			if next == nil {
				pl.Tail = nil
			}
		}
		currentOrder = next
	}

	if pl.Head == nil {
		if o.IsBuy {
			engine.OrderBook.BestOffer = nextBestOffer(engine.OrderBook)
		} else {
			engine.OrderBook.BestBid = nextBestBid(engine.OrderBook)
		}
	}
}

func (engine *MatchingEngine) StartTradeProcessor(verbose bool) {
	go func() {
		for {
			if t, ok := engine.Ring.Pop(); ok && verbose {
				fmt.Printf("Trade executed: OrderId: %d, Price: %d, Volume: %d, FillTime: %d\n",
					t.OrderId, t.Price, t.Volume, t.FillTime)
			}
		}
	}()
}
