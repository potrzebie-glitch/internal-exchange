// matchingengine_test.go
package exchange

import (
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func BenchmarkProcessOrder(b *testing.B) {

	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal("could not create CPU profile:", err)
	}
	defer cpuFile.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		b.Fatal("could not start CPU profile:", err)
	}
	defer pprof.StopCPUProfile()

	// Seed the random number generator for randomness
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Create a sample OrderBook and MatchingEngine
	orderBook := NewOrderBook()
	engine := NewMatchingEngine(orderBook)

	// Benchmarking loop - runs b.N times
	for i := 0; i < b.N; i++ {
		// Randomize whether it's a buy or sell order (true = buy, false = sell)
		isBuy := r.Intn(2) == 1

		// Randomize the price between 90 and 110
		price := r.Intn(21) + 90 // Generates a price between 90 and 110

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
