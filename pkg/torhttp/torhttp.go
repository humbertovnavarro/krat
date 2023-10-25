package torhttp

import (
	"context"
	"net/http"

	"github.com/cretz/bine/tor"
	"github.com/sirupsen/logrus"
)

var client *http.Client

func UseTorClient(t *tor.Tor) *http.Client {
	if client != nil {
		return client
	}
	client = createTorHTTPClient(t)
	return client
}

// createTorHTTPClient creates a HTTP client that uses the Tor network
func createTorHTTPClient(t *tor.Tor) *http.Client {
	// Create a proxy dialer using Tor SOCKS proxy
	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		logrus.Fatal(err)
	}
	// Create a custom transport that uses the proxy dialer
	tr := &http.Transport{DialContext: dialer.DialContext}
	// Create a custom HTTP client
	return &http.Client{Transport: tr}
}
