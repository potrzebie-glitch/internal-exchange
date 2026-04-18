package exchange

const MaxPrice = 65536
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
	Head        *Order
	Tail        *Order
	TotalVolume int
}

// OrderBook stores all orders for all price levels of a given book. It also keeps track of best bid and offer.
type OrderBook struct {
	BestBid   int
	BestOffer int
	bids      []*PriceLevel
	asks      []*PriceLevel
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		BestBid:   MinPrice,
		BestOffer: MaxPrice,
		bids:      make([]*PriceLevel, MaxPrice),
		asks:      make([]*PriceLevel, MaxPrice),
	}
}

// level returns the PriceLevel for the given price, allocating it lazily on first use.
func (ob *OrderBook) level(price int, isBuy bool) *PriceLevel {
	if isBuy {
		if ob.bids[price] == nil {
			ob.bids[price] = &PriceLevel{}
		}
		return ob.bids[price]
	}
	if ob.asks[price] == nil {
		ob.asks[price] = &PriceLevel{}
	}
	return ob.asks[price]
}

func (ob *OrderBook) InsertOrder(o *Order) {
	pl := ob.level(o.Price, o.IsBuy)
	if o.IsBuy {
		if o.Price > ob.BestBid {
			ob.BestBid = o.Price
		}
	} else {
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
	pl.TotalVolume += o.Volume
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
