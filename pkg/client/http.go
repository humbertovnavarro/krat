package client

import (
	"fmt"
	"net"

	"github.com/humbertovnavarro/tor-reverse-shell/pkg/tor_engine"
)

func StartHTTPServer(l net.Listener, e *tor_engine.TorEngine, c chan error) {
	fmt.Println("starting http server")
	r := router(e)
	fmt.Println("creating new onion for http service")
	fmt.Println("done creating new onion for http service")
	c <- r.RunListener(l)
}
