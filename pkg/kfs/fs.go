package kfs

import (
	"errors"
	"os"
)

func FileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
