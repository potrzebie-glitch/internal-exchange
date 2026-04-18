package exchange

const MaxPrice = 1000000
const MinPrice = 0

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

func NewOrderBook() *OrderBook {
	ob := OrderBook{
		BestBid:   MinPrice,
		BestOffer: MaxPrice,
		bids:      make([]*PriceLevel, MaxPrice),
		asks:      make([]*PriceLevel, MaxPrice),
	}
	for i := 0; i < MaxPrice; i++ {
		ob.bids[i] = &PriceLevel{}
		ob.asks[i] = &PriceLevel{}
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
	for o := pl.Head; o != nil; o = o.Next {
		sum += o.Volume
	}
	return sum
}

func (pl *PriceLevel) RemoveOrder(orderId int) bool {
	if pl.Head == nil {
		return false
	}

	if pl.Head.Id == orderId {
		pl.Head = pl.Head.Next
		if pl.Head == nil {
			pl.Tail = nil
		}
		return true
	}

	for cur := pl.Head; cur.Next != nil; cur = cur.Next {
		if cur.Next.Id == orderId {
			cur.Next = cur.Next.Next
			if cur.Next == nil {
				pl.Tail = cur
			}
			return true
		}
	}

	return false
}

// Start from current best bid, decrement the price and check
// if there is any bid, if so return that price, if not, check next.
// If no valid bid found, return MinPrice.
func nextBestBid(ob *OrderBook) int {
	for price := ob.BestBid - 1; price > MinPrice; price-- {
		if level := ob.bids[price]; level != nil && level.Head != nil {
			return price
		}
	}
	return MinPrice
}

// Start from current best offer, increment the price and check
// if there is any offer, if so return that price, if not, check next.
// If no valid offer found, return MaxPrice.
func nextBestOffer(ob *OrderBook) int {
	for price := ob.BestOffer + 1; price < MaxPrice; price++ {
		if level := ob.asks[price]; level != nil && level.Head != nil {
			return price
		}
	}
	return MaxPrice
}
