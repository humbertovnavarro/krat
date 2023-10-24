package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cretz/bine/tor"
)

// CONTROL_SERVER is the URL of the control server also running on the Tor network
var CONTROL_SERVER string

func init() {
	CONTROL_SERVER = os.Getenv("CONTROL_SERVER")
	if CONTROL_SERVER == "" {
		CONTROL_SERVER = "http://localhost:3000"
	}
}

//go:embed tor.exe
var torExe []byte

//go:embed next
var nextFS embed.FS

func RandomString(stringlen int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, stringlen)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// createTorExecutable creates the tor.exe file in temp folder
func createTorExecutable() string {
	// Create the file
	f, err := os.Create(os.TempDir() + RandomString(10) + ".exe")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.Write(torExe)
	if err2 != nil {
		log.Fatal(err2)
	}
	return f.Name()
}

// createTorHTTPClient creates a HTTP client that uses the Tor network
func createTorHTTPClient(t *tor.Tor) *http.Client {
	// Create a proxy dialer using Tor SOCKS proxy
	dialer, err := t.Dialer(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a custom transport that uses the proxy dialer
	tr := &http.Transport{DialContext: dialer.DialContext}
	// Create a custom HTTP client
	return &http.Client{Transport: tr}
}

func main() {
	distFS, err := fs.Sub(nextFS, "nextjs/dist")
	if err != nil {
		log.Fatal(err)
	}

	torExePath := createTorExecutable()
	torDataDir := fmt.Sprintf("%s/%s", os.TempDir(), RandomString(10))
	os.MkdirAll(torDataDir, 0755)

	// Start tor with default config (can set start conf's DebugWriter to os.Stdout for debug logs)
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(context.Background(), &tor.StartConf{
		ExePath:         torExePath,
		TempDataDirBase: torDataDir,
	})

	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}

	defer t.Close()

	// Wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	// Create a v3 onion service to listen on any port but show as 80
	onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, RemotePorts: []int{80}})
	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}

	// Create a HTTP client that uses the Tor network
	httpClient := createTorHTTPClient(t)
	notifyOnline(httpClient, onion.ID)

	// Notify control server of the onion address
	defer onion.Close()
	fmt.Printf("Open Tor browser and navigate to http://%v.onion\n", onion.ID)
	fmt.Println("Press enter to exit")
	http.Handle("/", http.FileServer(http.FS(distFS)))
	http.Serve(onion, http.DefaultServeMux)
	// TODO: Wait for exit
	errChan := make(chan error)
	<-errChan
	t.Close()
	// Clean up created files
	notifyOffline(httpClient, onion.ID)
	os.Remove(torDataDir)
	os.Remove(torExePath)
}

func notifyOnline(client *http.Client, onionID string) {
	body := []byte(fmt.Sprintf(`{"id": "%s"}`, onionID))
	client.Post(fmt.Sprintf("%s/node/online", CONTROL_SERVER), "application/json", bytes.NewReader(body))
}

func notifyOffline(client *http.Client, onionID string) {
	body := []byte(fmt.Sprintf(`{"id": "%s"}`, onionID))
	client.Post(fmt.Sprintf("%s/node/offline", CONTROL_SERVER), "application/json", bytes.NewReader(body))
}
