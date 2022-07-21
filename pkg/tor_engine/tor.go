package tor_engine

import (
	"context"
	"net/http"
	"time"

	"github.com/cretz/bine/tor"
)

type OnConnect = func(e *TorEngine)

type TorEngineConf struct {
	Start *tor.StartConf
}

type TorEngine struct {
	onConnect OnConnect
	conf      *TorEngineConf
	Tor       *tor.Tor
	Dialer    *tor.Dialer
	Open      bool
}

func New(conf *TorEngineConf, callback OnConnect) *TorEngine {
	return &TorEngine{
		conf:      conf,
		Open:      true,
		onConnect: callback,
	}
}

func (e *TorEngine) Start() error {
	t, err := tor.Start(context.Background(), e.conf.Start)
	// set the tor Open status to closed using the ProcessCancelFunc callback
	t.ProcessCancelFunc = func() {
		e.Open = false
	}
	e.Tor = t
	if err != nil {
		return err
	}
	// call the on connect callback inside of TorEngine
	e.onConnect(e)
	return err
}

// Create a new http client proxied through tor
func (e *TorEngine) NewHTTPClient() (*http.Client, error) {
	dialer, err := e.NewDialer()
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}, nil
}

// Create a new dialer proxied through tor
func (c *TorEngine) NewDialer() (*tor.Dialer, error) {
	if c.Dialer != nil {
		return c.Dialer, nil
	}
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	_dialer, err := c.Tor.Dialer(dialCtx, nil)
	c.Dialer = _dialer
	if err != nil {
		return nil, err
	}
	return c.Dialer, nil
}
