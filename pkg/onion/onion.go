package onion

// wraps "github.com/cretz/bine/tor" to make it easy to use and save onion services
// https://www.youtube.com/watch?v=-FtCTW2rVFM

import (
	"context"
	"errors"
	"net"
	"time"

	bine "github.com/cretz/bine/tor"
	bed25519 "github.com/cretz/bine/torutil/ed25519"
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

func New(e *tor.TorEngine, onionConfig *OnionServiceConfig) (*Onion, error) {
	onionService := &models.OnionService{
		ID:   onionConfig.ID,
		Port: onionConfig.Port,
	}
	onionNotFoundError := config.DB.Take(onionService).Error
	if onionNotFoundError != nil {
		logrus.Error(onionNotFoundError)
		generatedKey, err := bed25519.GenerateKey(nil)
		if err != nil {
			logrus.Fatal(err)
		}
		onionService.PrivateKey = generatedKey.PrivateKey()
	}
	keyPair := bed25519.FromCryptoPrivateKey(onionService.PrivateKey)
	listenConf := &bine.ListenConf{
		RemotePorts: []int{onionConfig.Port},
		Version3:    true,
		Detach:      false,
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
	onionService.URL = bineOnion.ID
	if onionNotFoundError != nil {
		err = config.DB.Create(onionService).Error
		if err != nil {
			logrus.Error(err)
		} else {
			logrus.Infof("wrote service %s to database", onionService.ID)
		}
	}
	return onion, nil
}
