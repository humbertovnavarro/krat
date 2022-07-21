package client

import (
	"io"
	"net"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

func StartSSHServer(l net.Listener, c chan error) {
	logrus.Infof("started ssh server on %s", l.Addr())
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})
	c <- ssh.Serve(l, ssh.DefaultHandler)
}
