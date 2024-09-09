package main

import (
	"testing"
)

func TestNewOrderBook(t *testing.T) {
	desired_size := 10000
	ob := NewOrderBook(desired_size)
	if len(ob.PriceLevels) != desired_size {
		t.Errorf("Expected size %v but got %v ", desired_size, len(ob.PriceLevels))
	}

}
