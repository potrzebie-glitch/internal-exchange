package main

import (
	"fmt"
	"internal-exchange/exchange"
	"math/rand"
	"time"
)

const ITERATIONS = 10000000

func main() {
	rand.Seed(time.Now().UnixNano())
	orderBook := exchange.NewOrderBook(exchange.MAX_PRICE)
	engine := exchange.NewMatchingEngine(orderBook)
	start := time.Now()
	for i := 0; i < ITERATIONS; i++ {
		isBuy := rand.Intn(2) == 1
		price := rand.Intn(21) + 90
		order := &exchange.Order{
			Id:     i,
			IsBuy:  isBuy,
			Price:  price,
			Volume: 10,
		}
		engine.ProcessOrder(order)
	}
	duration := time.Since(start)
	fmt.Printf("Processing %v orders takes %v ", ITERATIONS, duration)

}
