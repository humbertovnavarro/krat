package torbin

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/humbertovnavarro/krat/pkg/helpers"
)

//go:embed tor.exe
var Exe []byte

var torDataDir string
var torBinary string

// createTorExecutable creates the tor.exe file in temp folder
func GetTorBinary() string {
	if torBinary != "" {
		return torBinary
	}
	// Create the file
	f, err := os.Create(os.TempDir() + helpers.RandomString(10) + ".exe")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.Write(Exe)
	if err2 != nil {
		log.Fatal(err2)
	}
	return f.Name()
}

func GetTorDataDir() string {
	if torDataDir != "" {
		return torDataDir
	}
	torDataDir := fmt.Sprintf("%s/%s", os.TempDir(), helpers.RandomString(10))
	os.MkdirAll(torDataDir, os.ModePerm)
	return torDataDir
}

func Cleanup() {
	os.RemoveAll(torDataDir)
	os.Remove(torBinary)
}
