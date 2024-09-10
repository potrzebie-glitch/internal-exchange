package main

const MAX_PRICE = 1000

// Order represents a buy or sell order in the order book.
type Order struct {
	Id     int
	IsBuy  bool
	Price  int
	Volume int
	Next   *Order
}

// PriceLevel keeps track of not yet executed orders for a given price level.
type PriceLevel struct {
	Head *Order
	Tail *Order
}

// OrderBook stores all orders for all price levels of a given book. It also keeps track of best bid and offer.
type OrderBook struct {
	BestBid     int
	BestOffer   int
	PriceLevels []*PriceLevel
}

func NewOrderBook(size int) *OrderBook {
	ob := OrderBook{
		BestBid:     0,
		BestOffer:   MAX_PRICE,
		PriceLevels: make([]*PriceLevel, size),
	}
	for i := 0; i < size; i++ {
		ob.PriceLevels[i] = &PriceLevel{
			Head: nil,
			Tail: nil,
		}
	}
	return &ob

}

func (ob *OrderBook) AddOrder(o *Order) {
	pl := ob.PriceLevels[o.Price]
	if pl.Head == nil {
		pl.Head = o
		pl.Tail = o
	} else {
		pl.Tail.Next = o
		pl.Tail = o
	}
}

func getLevelVolume(pl *PriceLevel) int {
	sum := 0
	currentOrder := pl.Head

	// Traverse the linked list
	for currentOrder != nil {
		sum += currentOrder.Volume
		currentOrder = currentOrder.Next
	}

	return sum

}
