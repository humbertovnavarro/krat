package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/humbertovnavarro/krat/pkg/routes"
	"github.com/humbertovnavarro/krat/pkg/torbin"
	"github.com/humbertovnavarro/krat/pkg/torhttp"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	if os.Getenv("NO_TOR") != "" {
		fmt.Println("Tor disabled, running in clearnet mode")
		r := routes.NewRouter()
		r.Run("127.0.0.1:" + os.Getenv("NO_TOR_PORT"))
	} else {
		fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
		t, err := tor.Start(context.Background(), &tor.StartConf{
			ExePath:         torbin.GetTorBinary(),
			TempDataDirBase: torbin.GetTorDataDir(),
		})

		if err != nil {
			log.Panicf("Unable to start Tor: %v", err)
		}

		torhttp.UseTorClient(t)

		defer t.Close()
		defer torbin.Cleanup()

		// Wait at most a few minutes to publish the service
		listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer listenCancel()

		// Create a v3 onion service to listen on any port but show as 80
		onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, RemotePorts: []int{80}})
		if err != nil {
			log.Panicf("Unable to create onion service: %v", err)
		}
		// Notify control server of the onion address
		defer onion.Close()
		fmt.Printf("Open Tor browser and navigate to http://%v.onion\n", onion.ID)
		fmt.Println("Press enter to exit")
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		fmt.Println("Blocking, press ctrl+c to continue...")
		r := routes.NewRouter()
		go func() {
			r.RunListener(onion)
		}()
		<-done // Will block here until user hits ctrl+c
	}
}
