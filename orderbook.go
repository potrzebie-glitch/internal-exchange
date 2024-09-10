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

type PriceLevel struct {
	Head *Order
	Tail *Order
}

type OrderBook struct {
	PriceLevels []*PriceLevel
}

func NewOrderBook(size int) *OrderBook {
	return &OrderBook{PriceLevels: make([]*PriceLevel, size)}

}
