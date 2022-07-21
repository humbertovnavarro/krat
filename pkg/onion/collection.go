package onion

import (
	"encoding/json"
	"io/ioutil"

	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/sirupsen/logrus"
)

type OnionCollection struct {
	TagToOnionMap map[string]*OnionServiceConfig `json:"onions"`
}

var onions OnionCollection = NewOnionCollection()

func (o *OnionCollection) serialize() {
	json, err := json.Marshal(o)
	if err != nil {
		logrus.Error(err)
		return
	}
	ioutil.WriteFile(config.NewFilePath("index.json"), json, 0755)
}

func NewOnionCollection() OnionCollection {
	collection := OnionCollection{
		TagToOnionMap: make(map[string]*OnionServiceConfig),
	}
	bytes, err := ioutil.ReadFile(config.NewFilePath("index.json"))
	if err != nil {
		logrus.Error(err)
		return collection
	}
	collectionOnDisk := &OnionCollection{
		TagToOnionMap: make(map[string]*OnionServiceConfig),
	}
	err = json.Unmarshal(bytes, collectionOnDisk)
	if err != nil {
		logrus.Error(err)
		return collection
	}
	return *collectionOnDisk
}
