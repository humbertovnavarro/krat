package onion

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
	"github.com/sirupsen/logrus"
)

type OnionServiceConfig struct {
	Port    int    `json:"port"`
	Onion   string `json:"onion"`
	Tag     string `json:"tag"`
	KeyFile string
}

type Onion struct {
	*tor.OnionService
}

func generateOnion(e *tor_engine.TorEngine, config OnionServiceConfig) (*Onion, error) {
	listenConf := &tor.ListenConf{
		RemotePorts: []int{config.Port},
		Version3:    true,
		Detach:      true,
	}
	if e.Tor == nil || !e.Open {
		return nil, errors.New("tried to create a new onion service but tor is not ready")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	onion, err := e.Tor.Listen(ctx, listenConf)
	if err != nil {
		return nil, err
	}
	service := &Onion{
		onion,
	}
	err = service.serialize()
	if err != nil {
		logrus.Error(err)
		return service, nil
	}
	onions.OnionStartConfigs[config.Tag] = &config
	onions.serialize()
	return service, nil
}

func loadExisting(config OnionServiceConfig) (*Onion, error) {
	return nil, nil
}

func New(e *tor_engine.TorEngine, config OnionServiceConfig) (*Onion, error) {
	exists := onions.OnionStartConfigs[config.Tag] != nil
	if exists {
		config = *onions.OnionStartConfigs[config.Tag]
		return loadExisting(config)
	}
	return generateOnion(e, config)
}

func (o *Onion) AssertKeyIsED25519KeyPair() ed25519.KeyPair {
	if !o.Version3 {
		panic("onion is not set to v3, something went horribly wrong")
	}
	key := o.Key.(ed25519.KeyPair)
	return key
}

func (o *Onion) serialize() error {
	keyPair := o.AssertKeyIsED25519KeyPair()
	err := os.MkdirAll(config.NewFilePath("tor", "onions"), 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(config.NewFilePath("tor", "onions", o.ID))
	if err != nil {
		return err
	}
	_, err = file.Write(keyPair.PrivateKey())
	if err != nil {
		return err
	}
	return nil
}
