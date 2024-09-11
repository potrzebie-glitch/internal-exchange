// matchingengine_test.go
package exchange

import (
	"internal-exchange/exchange"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkProcessOrder(b *testing.B) {
	// Seed the random number generator for randomness
	rand.Seed(time.Now().UnixNano())

	// Create a sample OrderBook and MatchingEngine
	orderBook := NewOrderBook(exchange.MAX_PRICE)
	engine := NewMatchingEngine(orderBook)

	// Benchmarking loop - runs b.N times
	for i := 0; i < b.N; i++ {
		// Randomize whether it's a buy or sell order (true = buy, false = sell)
		isBuy := rand.Intn(2) == 1

		// Randomize the price between 90 and 110
		price := rand.Intn(21) + 90 // Generates a price between 90 and 110

		// Create a randomized order
		order := &Order{
			Id:     i, // Use loop index as the order ID
			IsBuy:  isBuy,
			Price:  price,
			Volume: 10,
		}

		// Process the randomized order
		engine.ProcessOrder(order)
	}
}
