package client_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
	"github.com/sirupsen/logrus"
)

const REQUEST_TIMEOUT = time.Second * 10

type fetchPasswordResponse struct {
	Key string `json:"key"`
}

// Get the encryption key required to open the local database, using generated tor http client
func FetchDBPassword(e *tor_engine.TorEngine) string {
	if config.GetEnv("DATABASE_PASSWORD") != "" {
		return config.GetEnv("DATABASE_PASSWORD")
	}
	client, err := e.NewHTTPClient()
	client.Timeout = time.Second * 10
	if err != nil {
		logrus.Panic(err)
	}
	resp, err := client.Get(fmt.Sprintf("%s/%s/%s", config.Remote, config.UUID, "unlock"))
	if err != nil {
		logrus.Error("could not get decryption keys for host database")
		logrus.Fatal(err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("could not read response body")
		logrus.Fatal(err)
	}
	respStruct := &fetchPasswordResponse{}
	err = json.Unmarshal(respBody, respStruct)
	if err != nil {
		logrus.Fatal(err)
	}
	return respStruct.Key
}
