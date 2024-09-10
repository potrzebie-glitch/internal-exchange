package main

import (
	"testing"
)

func TestNewOrderBookSizeCorrect(t *testing.T) {
	desired_size := 10000
	ob := NewOrderBook(desired_size)
	if len(ob.PriceLevels) != desired_size {
		t.Errorf("Expected size %v but got %v ", desired_size, len(ob.PriceLevels))
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
	ob.AddOrder(&o1)
	ob.AddOrder(&o2)
	pl := ob.PriceLevels[3]
	qty := getLevelVolume(pl)
	if qty != 2 {
		t.Errorf("Expected volume %v but got %v ", 2, qty)
	}
}
