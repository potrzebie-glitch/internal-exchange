package exchange

import (
	"testing"
)

func TestOrderBookVolumeAggregateCorrect(t *testing.T) {
	o1 := Order{Id: 1, IsBuy: true, Price: 3, Volume: 1}
	o2 := Order{Id: 2, IsBuy: true, Price: 3, Volume: 1}
	ob := NewOrderBook()
	ob.InsertOrder(&o1)
	ob.InsertOrder(&o2)
	if qty := getLevelVolume(ob.bids[3]); qty != 2 {
		t.Errorf("expected volume 2, got %d", qty)
	}
}

func TestAggressiveBidTakesOutPriceLevel(t *testing.T) {
	// Two resting sells at 3, one resting buy at 1, then an aggressive buy at 3 with volume 3.
	// The aggressive buy should consume both sells (2 vol) and rest the remainder (1 vol) at 3.
	o1 := Order{Id: 1, IsBuy: false, Price: 3, Volume: 1}
	o2 := Order{Id: 2, IsBuy: false, Price: 3, Volume: 1}
	o3 := Order{Id: 3, IsBuy: true, Price: 1, Volume: 3}
	o4 := Order{Id: 4, IsBuy: true, Price: 3, Volume: 3}

	ob := NewOrderBook()
	ob.InsertOrder(&o1)
	ob.InsertOrder(&o2)
	ob.InsertOrder(&o3)

	me := NewMatchingEngine(ob)
	me.ProcessOrder(&o4)

	if ob.BestBid != 3 {
		t.Errorf("expected BestBid 3, got %d", ob.BestBid)
	}
	if ob.BestOffer != MaxPrice {
		t.Errorf("expected BestOffer %d, got %d", MaxPrice, ob.BestOffer)
	}
	if ob.asks[3].Head != nil {
		t.Error("expected ask price level 3 to be empty")
	}
	if head := ob.bids[3].Head; head == nil || head.Volume != 1 {
		t.Errorf("expected resting bid at 3 with volume 1, got %v", head)
	}

	// Two trades should have been generated
	if got := me.Ring.Len(); got != 2 {
		t.Errorf("expected 2 trades in ring, got %d", got)
	}
	trade1, _ := me.Ring.Pop()
	if trade1.OrderId != 1 || trade1.Volume != 1 {
		t.Errorf("unexpected first trade: %+v", trade1)
	}
	trade2, _ := me.Ring.Pop()
	if trade2.OrderId != 2 || trade2.Volume != 1 {
		t.Errorf("unexpected second trade: %+v", trade2)
	}
}
