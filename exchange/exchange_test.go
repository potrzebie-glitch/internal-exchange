package exchange

import (
	"fmt"
	"testing"
)

func TestNewOrderBookSizeCorrect(t *testing.T) {
	desired_size := 10000
	ob := NewOrderBook(desired_size)
	if len(ob.bids) != desired_size {
		t.Errorf("Expected size %v but got %v ", desired_size, len(ob.bids))
	}

}

func TestOrderBookVolumeAggregateCorrect(t *testing.T) {
	o1 := Order{
		Id:     1,
		IsBuy:  true,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	o2 := Order{
		Id:     2,
		IsBuy:  true,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	ob := NewOrderBook(4)
	ob.InsertOrder(&o1)
	ob.InsertOrder(&o2)
	pl := ob.bids[3]
	qty := getLevelVolume(pl)
	if qty != 2 {
		t.Errorf("Expected volume %v but got %v ", 2, qty)
	}
}

func TestAggressiveBidTakesOutPriceLevel(t *testing.T) {
	o1 := Order{
		Id:     1,
		IsBuy:  false,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	o2 := Order{
		Id:     2,
		IsBuy:  false,
		Price:  3,
		Volume: 1,
		Next:   nil,
	}
	o3 := Order{
		Id:     3,
		IsBuy:  true,
		Price:  1,
		Volume: 3,
		Next:   nil,
	}
	o4 := Order{
		Id:     4,
		IsBuy:  true,
		Price:  3,
		Volume: 3,
		Next:   nil,
	}

	ob := NewOrderBook(5)
	ob.InsertOrder(&o1)
	ob.InsertOrder(&o2)
	ob.InsertOrder(&o3)
	me := NewMatchingEngine(ob)
	pl := me.OrderBook.asks[3]
	currentOrder := pl.Head
	for currentOrder != nil {
		fmt.Println(currentOrder)
		currentOrder = currentOrder.Next
	}

	fmt.Println("Enter aggressive bid")
	me.ProcessOrder(&o4)
	pl3 := me.OrderBook.asks[3]
	current := pl3.Head
	fmt.Println("Printing asks")
	for current != nil {
		fmt.Println(current)
		current = current.Next
	}
	fmt.Println("Printing bids")
	pl4 := me.OrderBook.bids[3]
	currentO := pl4.Head
	for currentO != nil {
		fmt.Println(currentO)
		currentO = currentO.Next
	}
	fmt.Println(ob.BestBid)
	fmt.Println(ob.BestOffer)

}
