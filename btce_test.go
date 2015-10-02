package btce

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

// Mock out HTTP requests.
//
// Pinched from http://keighl.com/post/mocking-http-responses-in-golang/
// Thanks, Kyle Truscott (@keighl)!
func httpMock(code int,
	body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport, Timeout: 5 * time.Second}

	// setting the host in the Client so I don't need to totally fake out
	// the TLS config
	client := &Client{
		Host:       "http://btce-e.com",
		HTTPClient: httpClient,
	}

	return server, client
}

// Test helper. Thanks again, @keighl
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)",
			b,
			reflect.TypeOf(b),
			a,
			reflect.TypeOf(a))
	}
}

func TestGetTicker(t *testing.T) {
	body := `{"btc_usd" : {
      "avg" : 234.6085,
      "high" : 237,
      "sell" : 234.1,
      "vol_cur" : 7389.03437,
      "vol" : 1732138.1951,
      "updated" : 1443706785,
      "buy" : 234.112,
      "low" : 232.217,
      "last" : 234.106
   }
}`

	server, client := httpMock(200, body)
	defer server.Close()

	ticker, err := client.GetTicker("btc_usd")

	if err != nil {
		t.Errorf("GetTicker : %v", err)
	}

	expect(t, ticker.LastTrade, float32(234.106))
	expect(t, ticker.Bid, float32(234.112))
	expect(t, ticker.Ask, float32(234.1))
}

func TestGetOrderBook(t *testing.T) {
	body := `{"btc_usd":{"asks":[[233.309,0.01104],[233.363,0.022089]],"bids":[[233.2,0.90891808],[233.08,0.01109033]]}}`
	server, client := httpMock(200, body)
	defer server.Close()

	orderBook, err := client.GetOrderBook("btc_usd")

	if err != nil {
		t.Errorf("GetDepth : %v", err)
	}

	expect(t, len(orderBook.Asks), 2)
	expect(t, orderBook.Asks[0].Price, float32(233.309))
	expect(t, orderBook.Asks[0].Quantity, float32(0.01104))
	expect(t, len(orderBook.Bids), 2)
	expect(t, orderBook.Bids[0].Price, float32(233.2))
	expect(t, orderBook.Bids[0].Quantity, float32(0.90891808))

}
