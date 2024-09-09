package main

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
