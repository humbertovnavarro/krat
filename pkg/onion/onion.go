package onion

// extends "github.com/cretz/bine/tor" OnionService, enforces one port per onion, and enforces v3 addressing
// https://www.youtube.com/watch?v=-FtCTW2rVFM

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/cretz/bine/torutil/ed25519"
	"github.com/humbertovnavarro/krat/pkg/db"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
	"github.com/sirupsen/logrus"
)

type OnionServiceConfig struct {
	Port int
	ID   string
}

type Onion struct {
	*tor.OnionService
}

func createBineOnionSocket(e *tor_engine.TorEngine, conf *tor.ListenConf) (*tor.OnionService, error) {
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

func tryGetPrivateKey(id string) (ed25519.PrivateKey, error) {
	stmt := `SELECT privateKey FROM OnionServices WHERE onionServiceID = ?`
	rows, err := db.DB.Query(stmt, id)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	var privateKeyHex string
	for rows.Next() {
		err = rows.Scan(&privateKeyHex)
		if err != nil {
			logrus.Error(err)
		}
	}
	privateKey := make([]byte, 64)
	num, err := hex.Decode(privateKey, []byte(privateKeyHex))
	if err != nil {
		return nil, err
	}
	if num != 64 {
		return nil, errors.New("wrong private key size for private key")
	}
	return privateKey, nil
}

// loads the onion service specified in config.ID if it exists, otherwise create a new onion service
func New(e *tor_engine.TorEngine, config *OnionServiceConfig) (*Onion, error) {
	listenConf := &tor.ListenConf{
		RemotePorts: []int{config.Port},
		Version3:    true,
		Detach:      true,
	}
	var keyPair ed25519.KeyPair
	existingKey, err := tryGetPrivateKey(config.ID)
	if err != nil && existingKey == nil {
		logrus.Error(err)
	}
	if existingKey == nil {
		logrus.Infof("generating new key for service %s", config.ID)
		generatedKey, err := ed25519.GenerateKey(nil)
		if err != nil {
			logrus.Fatal(err)
		}
		keyPair = generatedKey
	} else {
		keyPair = existingKey.KeyPair()
	}
	listenConf.Key = keyPair
	bineOnion, err := createBineOnionSocket(e, listenConf)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	onion := &Onion{
		bineOnion,
	}
	if existingKey == nil {
		err = serialize(onion, config)
		if err != nil {
			logrus.Error(err)
		}
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

func serialize(o *Onion, s *OnionServiceConfig) error {
	keyPair := o.AssertKeyIsED25519KeyPair()
	stmt := `INSERT INTO OnionServices (onionServiceID, onionUrl, port, privateKey) VALUES (?, ?, ?, ?)`
	_, err := db.DB.Exec(stmt, s.ID, o.ID, o.RemotePorts[0], hex.EncodeToString(keyPair.PrivateKey()))
	return err
}
