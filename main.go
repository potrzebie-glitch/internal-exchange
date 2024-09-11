package main

import (
	"fmt"
	"internal-exchange/exchange"
	"math/rand"
	"time"
)

const ITERATIONS = 10000000

func main() {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	orderBook := exchange.NewOrderBook()
	engine := exchange.NewMatchingEngine(orderBook)
	engine.StartTradeProcessor()
	start := time.Now()
	for i := 0; i < ITERATIONS; i++ {
		isBuy := r.Intn(2) == 1
		price := r.Intn(21) + 90
		volume := r.Intn(5) + 1
		order := &exchange.Order{
			Id:     i,
			IsBuy:  isBuy,
			Price:  price,
			Volume: volume,
		}
		//fmt.Printf("Market is %v @ %v", engine.OrderBook.BestBid, engine.OrderBook.BestOffer)
		engine.ProcessOrder(order)

	}
	duration := time.Since(start)
	fmt.Printf("Processing %v orders takes %v ", ITERATIONS, duration)

}
