package tor_engine

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cretz/bine/tor"
)

type TorEngineConf struct {
	Start *tor.StartConf
}

type OnConnect = func(e *TorEngine)

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
	t.ProcessCancelFunc = func() {
		e.Open = false
	}
	e.Tor = t
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	e.onConnect(e)
	return err
}

func (e *TorEngine) NewOnionListener(listenConf *tor.ListenConf) (*tor.OnionService, error) {
	if e.Tor == nil {
		return nil, errors.New("tor: the tor service is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	onion, err := e.Tor.Listen(ctx, listenConf)
	if err != nil {
		return nil, err
	}
	return onion, nil
}

func (e *TorEngine) NewHTTPClient() (*http.Client, error) {
	dialer, err := e.NewDialer()
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}, nil
}

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
