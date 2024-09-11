package exchange

import "fmt"

const MAX_PRICE = 1000000

// Order represents a buy or sell order in the order book.
type Order struct {
	Id     int
	IsBuy  bool
	Price  int
	Volume int
	Next   *Order
}

// PriceLevel keeps track of not yet executed orders for a given price level.
type PriceLevel struct {
	Head *Order
	Tail *Order
}

// OrderBook stores all orders for all price levels of a given book. It also keeps track of best bid and offer.
type OrderBook struct {
	BestBid   int
	BestOffer int
	bids      []*PriceLevel
	asks      []*PriceLevel
}

func NewOrderBook(size int) *OrderBook {
	ob := OrderBook{
		BestBid:   0,
		BestOffer: MAX_PRICE,
		bids:      make([]*PriceLevel, size+1),
		asks:      make([]*PriceLevel, size+1),
	}
	for i := 0; i <= size; i++ {
		ob.bids[i] = &PriceLevel{
			Head: nil,
			Tail: nil,
		}
		ob.asks[i] = &PriceLevel{
			Head: nil,
			Tail: nil,
		}
	}
	return &ob

}

func (ob *OrderBook) InsertOrder(o *Order) {
	var pl *PriceLevel
	if o.IsBuy {
		pl = ob.bids[o.Price]
		if o.Price > ob.BestBid {
			ob.BestBid = o.Price
		}
	} else {
		pl = ob.asks[o.Price]
		if o.Price < ob.BestOffer {
			ob.BestOffer = o.Price
		}
	}

	if pl.Head == nil {
		pl.Head = o
		pl.Tail = o
	} else {
		pl.Tail.Next = o
		pl.Tail = o
	}
}

// Note that there can never be both bids and offers resting at the same price level
func getLevelVolume(pl *PriceLevel) int {
	sum := 0
	currentOrder := pl.Head

	for currentOrder != nil {
		sum += currentOrder.Volume
		currentOrder = currentOrder.Next
	}

	return sum

}

func (pl *PriceLevel) RemoveOrder(orderId int) bool {
	// If the list is empty, there's nothing to remove
	if pl.Head == nil {
		return false
	}

	// If the order to remove is the head
	if pl.Head.Id == orderId {
		pl.Head = pl.Head.Next
		// If the list becomes empty, update the tail
		if pl.Head == nil {
			pl.Tail = nil
		}
		return true
	}

	// Traverse the list to find the order to remove
	current := pl.Head
	for current.Next != nil {
		if current.Next.Id == orderId {
			// Remove the order by skipping it in the linked list
			current.Next = current.Next.Next
			// If we removed the tail, update the tail reference
			if current.Next == nil {
				pl.Tail = current
			}
			return true
		}
		current = current.Next
	}

	// If the order is not found
	return false
}

// Start from current best bid, decrement the price and check
// if there is any bid, if so return that price, if not, check next
func nextBestBid(ob *OrderBook) int {
	for price := ob.BestOffer - 1; price >= 0; price-- {
		if ob.bids[price].Head != nil {
			return price
		}
	}
	fmt.Println("Panic mode: no bids found")
	return 0
}

// Start from current best offer, increment the price and check
// if there is any offer, if so return that price, if not, check next
func nextBestOffer(ob *OrderBook) int {
	for price := ob.BestOffer + 1; price <= MAX_PRICE; price++ {
		if ob.asks[price].Head != nil {
			return price
		}
	}
	fmt.Println("Panic mode: no offers found")
	return MAX_PRICE
}
