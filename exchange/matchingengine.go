package exchange

import (
	"fmt"
)

// Trade represents the transaction resulting from a matching bid and offer.
type Trade struct {
	OrderId  int
	Price    int
	Volume   int
	FillTime int
}

// MatchingEngine updates the state of the orderbook when new orders come in. It
// also issues trades.
type MatchingEngine struct {
	OrderBook   *OrderBook
	TradeAction chan Trade
}

func NewMatchingEngine(ob *OrderBook) MatchingEngine {
	return MatchingEngine{
		OrderBook:   ob,
		TradeAction: make(chan Trade, 10),
	}
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

		select {
		case engine.TradeAction <- Trade{
			OrderId:  currentOrder.Id,
			Price:    currentOrder.Price,
			Volume:   qty,
			FillTime: 1,
		}:
		default:
		}

		if currentOrder.Volume == 0 {
			if !pl.RemoveOrder(currentOrder.Id) {
				fmt.Printf("Warning: could not remove order %d\n", currentOrder.Id)
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

func (engine *MatchingEngine) StartTradeProcessor() {
	go func() {
		for trade := range engine.TradeAction {
			fmt.Printf("Trade executed: OrderId: %d, Price: %d, Volume: %d, FillTime: %d\n",
				trade.OrderId, trade.Price, trade.Volume, trade.FillTime)
		}
	}()
}
