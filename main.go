package main

import "internal-exchange/exchange"

func main() {

	o1 := exchange.Order{
		Id:     1,
		IsBuy:  false,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	o2 := exchange.Order{
		Id:     2,
		IsBuy:  false,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	ob := exchange.NewOrderBook(4)
	ob.InsertOrder(&o1)
	ob.InsertOrder(&o2)
	engine := exchange.NewMatchingEngine(ob)
	o3 := exchange.Order{
		Id:     3,
		IsBuy:  true,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	engine.ProcessOrder(&o3)

}
