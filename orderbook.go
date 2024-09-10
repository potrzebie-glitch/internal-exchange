package main

const MAX_PRICE = 1000

// Order represents a buy or sell order in the order book.
type Order struct {
	Id     int
	Side   string
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
	return &OrderBook{BestBid: 0,
		BestOffer:   MAX_PRICE,
		PriceLevels: make([]*PriceLevel, size)}

}

func (ob *OrderBook) AddOrder(o *Order) {
	ob.PriceLevels[o.Price].Tail.Next = o
	ob.PriceLevels[o.Price].Tail = o
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
