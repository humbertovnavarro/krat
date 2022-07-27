package onion

// wraps "github.com/cretz/bine/tor" to make it easy to use and save onion services
// https://www.youtube.com/watch?v=-FtCTW2rVFM

import (
	"context"
	"errors"
	"net"
	"time"

	bine "github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/models"
	"github.com/humbertovnavarro/krat/pkg/tor"
	"github.com/sirupsen/logrus"
)

type OnionServiceConfig struct {
	Port int
	ID   string
}

type Onion struct {
	*bine.OnionService
}

func createBineOnionSocket(e *tor.TorEngine, conf *bine.ListenConf) (*bine.OnionService, error) {
	if e.Tor == nil || !e.Open {
		return nil, errors.New("tried to create a new onion service but tor is not ready")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	onion, err := e.Tor.Listen(ctx, conf)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return onion, nil
}

type OnionServiceHandler func(l net.Listener, ch chan error)

// loads the onion service specified in config.ID if it exists, otherwise create a new onion service
func New(e *tor.TorEngine, config *OnionServiceConfig) (*Onion, error) {
	service := &models.OnionService{}
	var keyPair ed25519.KeyPair
	var existingKey ed25519.PrivateKey
	if service != nil {
		existingKey = service.PrivateKey
	} else {
		logrus.Infof("generating new key for service %s", config.ID)
		generatedKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			logrus.Fatal(err)
		}
		keyPair = generatedKey
	}
	listenConf := &bine.ListenConf{
		RemotePorts: []int{config.Port},
		Version3:    true,
		Detach:      true,
		Key:         keyPair,
	}
	bineOnion, err := createBineOnionSocket(e, listenConf)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	onion := &Onion{
		bineOnion,
	}
	if existingKey == nil {
		serialize(onion, config)
	}
	return onion, nil
}

func (o *Onion) AssertKeyIsED25519KeyPair() ed25519.KeyPair {
	if !o.Version3 {
		panic("onion is not set to v3, something went horribly wrong")
	}
	key := o.Key.(ed25519.KeyPair)
	return key
}

func serialize(o *Onion, s *OnionServiceConfig) {
	keyPair := o.AssertKeyIsED25519KeyPair()
	config.DB.FirstOrCreate(&models.OnionService{
		ID:         s.ID,
		URL:        o.ID,
		Port:       s.Port,
		PrivateKey: keyPair.PrivateKey(),
	})
}
