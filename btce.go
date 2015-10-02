/*
An incomplete client for the btce api.
*/
package btce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	prefix        = "/api/3"
	timeout       = 5 * time.Second
	tickerPath    = "ticker/%v"
	orderBookPath = "depth/%v"
	limit         = 5
)

// Client communicates with the API
type Client struct {
	Host       string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		Host:       "https://btc-e.com",
		HTTPClient: &http.Client{Timeout: timeout},
	}
}

// GetTicker returns a Ticker from the API
//
// Example
//
//   $ curl "https://btc-e.com/api/3/ticker/btc_usd?limit=5"
//   {
//     "btc_usd":{
//       "high":233.5,
//       "low":230,
//       "avg":231.75,
//       "vol":1624577.12558,
//       "vol_cur":6998.05547,
//       "last":232.819,
//       "buy":232.867,
//       "sell":232.819,
//       "updated":1443277438
//     }
//   }
func (c *Client) GetTicker(symbol string) (*Ticker, error) {
	var body []byte
	var err error

	path := fmt.Sprintf(tickerPath, symbol)

	if body, err = c.get(c.url(path)); err != nil {
		return nil, err
	}

	tickers := map[string]Ticker{}

	if err = json.Unmarshal(body, &tickers); err != nil {
		return nil, err
	}

	ticker := tickers[symbol]
	return &ticker, nil
}

// GetOrderBook returns the OrderBook from the API
//
// Example
//
//   $ curl "https://btc-e.com/api/3/depth/btc_usd?limit=2"
//   {
//     "btc_usd":{
//       "asks":[
//         [233.567, 0.01104],
//         [233.751, 22.0454485]
//       ],
//       "bids":[
//         [233.2, 1.66890233],
//         [233.072, 0.05070952]
//       ]
//     }
//   }
func (c *Client) GetOrderBook(symbol string) (*OrderBook, error) {
	var body []byte
	var err error

	opts := fmt.Sprintf("?limit=%v", limit)
	path := fmt.Sprintf(orderBookPath, symbol) + opts

	if body, err = c.get(c.url(path)); err != nil {
		return nil, err
	}

	res := make(map[string]map[string][][]float32)

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	book := OrderBook{
		Asks: buildOrders(res[symbol]["asks"]),
		Bids: buildOrders(res[symbol]["bids"]),
	}

	return &book, nil
}

// newOrder builds a new Order from an array of two floats.
//
// The first float is the price, the second float is the quantity.
func newOrder(order []float32) Order {
	return Order{
		Price:    order[0],
		Quantity: order[1],
	}

}

// buildOrders returns an array of Order structs from the [][]float32 arrays
// of bid and ask data.
func buildOrders(orderData [][]float32) []Order {
	orders := make([]Order, len(orderData))

	for i, order := range orderData {
		orders[i] = newOrder(order)
	}

	return orders
}

// url returns a full URL for a resource
func (c *Client) url(resource string) string {
	return fmt.Sprintf("%v%v/%v", c.Host, prefix, resource)
}

// get a response from a URL.
//
// This method will handle closing off the body.
func (c *Client) get(url string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(url)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
