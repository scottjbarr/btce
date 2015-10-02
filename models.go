package btce

// Ticker representing basic details of where the current price of a symbol is.
type Ticker struct {
	LastTrade float32 `json:"last"`
	Bid       float32 `json:"buy"`
	Ask       float32 `json:"sell"`
}

// OrderBook represents the full order book returned by the API.
type OrderBook struct {
	Asks []Order `json:"asks"`
	Bids []Order `json:"bids"`
}

// Order represents the price and quanty of an individual Order, or the summary
// of multiple Orders (as in the case of an Order Book)
type Order struct {
	Price    float32 `json:"price,string"`
	Quantity float32 `json:"amount,string"`
}
