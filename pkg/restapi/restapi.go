package restapi

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/onion"
	"github.com/humbertovnavarro/krat/pkg/tor"
)

func Start(e *tor.TorEngine, eCh chan error) {
	portString := config.GetEnv("HTTP_PORT")
	if portString == "" {
		portString = "8080"
	}
	port, err := strconv.ParseInt(portString, 10, 32)
	if err != nil {
		eCh <- err
	}
	apiOnion, err := onion.New(e, &onion.OnionServiceConfig{
		Port: int(port), ID: "http",
	})
	if err != nil {
		eCh <- err
	}
	api := New()
	eCh <- api.RunListener(apiOnion)
}

func StartDebug(eCh chan error) {
	api := New()
	eCh <- api.Run("127.0.0.1:8080")
}

func New() *gin.Engine {
	r := gin.Default()
	ApplyHandlers(r)
	ApplyMiddlewareStack(r)
	return r
}
