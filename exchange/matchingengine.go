package exchange

import "fmt"

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
	TradeAction chan<- Trade
}

func NewMatchingEngine(ob *OrderBook) MatchingEngine {
	tradeChan := make(chan Trade)
	return MatchingEngine{
		OrderBook:   ob,
		TradeAction: tradeChan,
	}

}

func (engine *MatchingEngine) ProcessOrder(order *Order) {
	if order.IsBuy {
		for price := engine.OrderBook.BestOffer; price >= order.Price; price-- {
			engine.processTrades(order, price)
			if order.Volume == 0 {
				break
			}
		}
	} else {
		for price := engine.OrderBook.BestBid; price <= order.Price; price++ {
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
	currentOrder := pl.Head
	for currentOrder != nil && o.Volume > 0 {
		if o.IsBuy {
			if o.Price >= currentOrder.Price {
				qty := min(o.Volume, currentOrder.Volume)
				o.Volume -= qty
				currentOrder.Volume -= qty

				if currentOrder.Volume == 0 {
					success := pl.RemoveOrder(currentOrder.Id)
					if !success {
						fmt.Printf("Warning - could not remove order %v", currentOrder.Id)
					}
					nextBestAsk := nextBestOffer(engine.OrderBook)
					engine.OrderBook.BestOffer = nextBestAsk

				}
			}
		}

		if !o.IsBuy {
			if o.Price <= currentOrder.Price {
				qty := min(o.Volume, currentOrder.Volume)
				o.Volume -= qty
				currentOrder.Volume -= qty
				if currentOrder.Volume == 0 {
					success := pl.RemoveOrder(currentOrder.Id)
					if !success {
						fmt.Printf("Warning - could not remove order %v", currentOrder.Id)
					}
					//if there is no more order around, then update best bid
					nextBestBid := nextBestBid(engine.OrderBook)
					engine.OrderBook.BestBid = nextBestBid
				}
			}
		}
		currentOrder = currentOrder.Next
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
