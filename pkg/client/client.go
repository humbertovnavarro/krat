package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/krat/pkg/config"
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
	fmt.Println("tor connected")
	sshOnion, err := onion.New(e, onion.OnionServiceConfig{
		Port: 22,
		Tag:  "ssh",
	})
	if err != nil {
		panic(err)
	}
	eCh := make(chan error)
	sCh := make(chan os.Signal, 1)
	signal.Notify(sCh, os.Interrupt)
	go StartSSHServer(sshOnion, eCh)
	go func() {
		for err := range eCh {
			logrus.Info(err)
			time.Sleep(time.Second)
		}
	}()
	func() {
		for sig := range sCh {
			log.Printf("captured %v, stopping tor gracefully..", sig)
			err := e.Tor.Close()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}()
}
