package routes

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/humbertovnavarro/krat/pkg/client"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	fs, err := fs.Sub(client.ClientFS, "build")
	if err != nil {
		panic(err)
	}
	fsHandler := gin.WrapH(http.FileServer(http.FS(fs)))
	if os.Getenv("NO_TOR") != "" {
		logrus.Info("Proxying webpack dev server")
		r.NoRoute(ProxyWebpack)
	} else {
		logrus.Info("Serving static files")
		r.NoRoute(fsHandler)
	}
	return r
}

func ProxyWebpack(c *gin.Context) {
	godotenv.Load()
	remoteURL := os.Getenv("WEBPACK_DEV_SERVER")
	remote, err := url.Parse(remoteURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
