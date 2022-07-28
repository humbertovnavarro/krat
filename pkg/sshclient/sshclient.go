package sshclient

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	"github.com/humbertovnavarro/krat/pkg/config"
	"github.com/humbertovnavarro/krat/pkg/kfs"
	"github.com/humbertovnavarro/krat/pkg/models"
	"github.com/humbertovnavarro/krat/pkg/onion"
	"github.com/humbertovnavarro/krat/pkg/tor"
	"github.com/sirupsen/logrus"
)

var PublicKey ssh.PublicKey
var keypath = config.FilePath("id_rsa")

func init() {
	if !kfs.FileExists(keypath) {
		logrus.Info("generating new host key")
		err := GenerateHostKey()
		if err != nil {
			logrus.Fatal(err)
		}
	}
	ssh.Handle(func(s ssh.Session) {
		shell, err := FindShell()
		if err != nil {
			logrus.Error(err)
		}
		cmd := exec.Command(shell)
		ptyReq, winCh, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()
			go func() {
				io.Copy(f, s)
			}()
			io.Copy(s, f)
			cmd.Wait()
		} else {
			io.WriteString(s, "No PTY requested.\n")
			s.Exit(1)
		}
	})
}

func Start(e *tor.TorEngine, eCh chan error) {
	fmt.Println("starting ssh client")
	dbKey := &models.CryptoKey{}
	config.DB.Find("WHERE For = ?", "ssh_authorized_key").Find(dbKey)
	sshOnion, err := onion.New(e, &onion.OnionServiceConfig{
		Port: 22,
		ID:   "ssh",
	})
	if err != nil {
		eCh <- err
		return
	}
	ssh.Serve(sshOnion, ssh.DefaultHandler, ssh.HostKeyFile(keypath))
}

func StartDebug(eCh chan error) {
	logrus.Info("starting ssh client in debug mode on port 2222")
	eCh <- ssh.ListenAndServe("127.0.0.1:2222", ssh.DefaultHandler, ssh.HostKeyFile(keypath))
}

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func FindShell() (string, error) {
	var shells = []string{"/usr/bin/bash", "/usr/bin/zsh", "/usr/bin/sh"}
	for _, shell := range shells {
		if kfs.FileExists(shell) {
			return shell, nil
		}
	}
	return "", errors.New("no shell found")
}
