package onion

import (
	"encoding/json"
	"io/ioutil"

	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/sirupsen/logrus"
)

type OnionCollection struct {
	OnionStartConfigs map[string]*OnionServiceConfig `json:"onions"`
}

var onions OnionCollection = OnionCollection{
	make(map[string]*OnionServiceConfig),
}

func (o *OnionCollection) serialize() {
	json, err := json.Marshal(o)
	if err != nil {
		logrus.Error(err)
		return
	}
	ioutil.WriteFile(config.NewFilePath("index.json"), json, 0755)
}
