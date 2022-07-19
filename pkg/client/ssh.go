package client

import (
	"fmt"
	"io"
	"net"

	"github.com/gliderlabs/ssh"
)

func StartSSHServer(l net.Listener, c chan error) {
	fmt.Printf("started ssh server on %s", l.Addr())
	ssh.Handle(func(s ssh.Session) {
		io.WriteString(s, "Hello world\n")
	})
	c <- ssh.Serve(l, ssh.DefaultHandler)
}
