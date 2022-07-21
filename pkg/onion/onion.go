package onion

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
	"github.com/sirupsen/logrus"
)

type OnionServiceConfig struct {
	Port  int
	Onion string `json:"onion"`
	Tag   string `json:"tag"`
}

type Onion struct {
	*tor.OnionService
}

// Creates a "github.com/cretz/bine/tor".OnionService
func createOnionSocket(e *tor_engine.TorEngine, conf *tor.ListenConf) (*tor.OnionService, error) {
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

func fetchPrivateKey(onion string) (ed25519.KeyPair, error) {
	privKey, err := ioutil.ReadFile(config.NewFilePath("tor", "onions", onion))
	if err != nil {
		return nil, err
	}
	keyPair := ed25519.PrivateKey.KeyPair(privKey)
	return keyPair, nil
}

func generateOnion(e *tor_engine.TorEngine, config OnionServiceConfig) (*Onion, error) {
	listenConf := &tor.ListenConf{
		RemotePorts: []int{config.Port},
		Version3:    true,
		Detach:      true,
	}
	onion, err := createOnionSocket(e, listenConf)
	if err != nil {
		logrus.Error(err)
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
	config.Onion = onion.ID
	onions.TagToOnionMap[config.Tag] = &config
	return service, nil
}

func loadOnionFromConfig(e *tor_engine.TorEngine, config OnionServiceConfig) (*Onion, error) {
	keyPair, err := fetchPrivateKey(config.Onion)
	if err != nil {
		return nil, err
	}
	listenConf := tor.ListenConf{
		RemotePorts: []int{config.Port},
		Version3:    true,
		Detach:      true,
		Key:         keyPair,
	}
	onion, err := createOnionSocket(e, &listenConf)
	if err != nil {
		return nil, err
	}
	service := &Onion{
		onion,
	}
	config.Onion = onion.ID
	return service, nil
}

func New(e *tor_engine.TorEngine, config OnionServiceConfig) (*Onion, error) {
	if onions.TagToOnionMap[config.Tag] != nil {
		logrus.Infof("found existing onion service %s running on port %d", config.Tag, config.Port)
		return loadOnionFromConfig(e, *onions.TagToOnionMap[config.Tag])
	}
	logrus.Infof("creating new onion service %s running on port %d", config.Tag, config.Port)
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
	err := os.MkdirAll(config.NewFilePath("tor", "onions"), 0600)
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
