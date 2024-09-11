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
	tradeChan := make(chan Trade, 10)
	return MatchingEngine{
		OrderBook:   ob,
		TradeAction: tradeChan,
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
	currentOrder := pl.Head
	for currentOrder != nil && o.Volume > 0 {
		if o.IsBuy {
			if o.Price >= currentOrder.Price {
				qty := min(o.Volume, currentOrder.Volume)
				o.Volume -= qty
				currentOrder.Volume -= qty
				trade := Trade{
					OrderId:  currentOrder.Id,
					Price:    currentOrder.Price,
					Volume:   qty,
					FillTime: 1,
				}

				select {
				case engine.TradeAction <- trade:
				default:
					//fmt.Println("TradeAction channel is full, dropping trade")
				}

				if currentOrder.Volume == 0 {
					success := pl.RemoveOrder(currentOrder.Id)
					if !success {
						fmt.Printf("Warning - could not remove order %v", currentOrder.Id)
					}
				}
			}
		}

		if !o.IsBuy {
			if o.Price <= currentOrder.Price {
				qty := min(o.Volume, currentOrder.Volume)
				o.Volume -= qty
				currentOrder.Volume -= qty
				trade := Trade{
					OrderId:  currentOrder.Id,
					Price:    currentOrder.Price,
					Volume:   qty,
					FillTime: 1,
				}
				select {
				case engine.TradeAction <- trade:
				default:
					//fmt.Println("TradeAction channel is full, dropping trade")
				}
				if currentOrder.Volume == 0 {
					success := pl.RemoveOrder(currentOrder.Id)
					if !success {
						fmt.Printf("Warning - could not remove order %v", currentOrder.Id)
					}
				}
			}
		}
		if pl.Head == nil {
			if o.IsBuy {
				engine.OrderBook.BestOffer = nextBestOffer(engine.OrderBook)
			} else {
				engine.OrderBook.BestBid = nextBestBid(engine.OrderBook)
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

func (engine *MatchingEngine) StartTradeProcessor() {
	go func() {
		for trade := range engine.TradeAction {
			// Pipe the trade to stdout
			fmt.Printf("Trade executed: OrderId: %d, Price: %d, Volume: %d, FillTime: %d\n",
				trade.OrderId, trade.Price, trade.Volume, trade.FillTime)
		}
	}()
}
