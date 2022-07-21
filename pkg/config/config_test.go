package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppDataDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	assert.Equal(t, AppDataDir, fmt.Sprintf("%s/%s", homeDir, AppName))
}

func TestNewFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	dir := NewFilePath("foo", "bar", "buzz", "bazz")
	expected := fmt.Sprintf("%s/%s/%s/%s/%s/%s", homeDir, AppName, "foo", "bar", "buzz", "bazz")
	assert.Equal(t, expected, dir)
	dir = NewFilePath("foo")
	expected = fmt.Sprintf("%s/%s/%s", homeDir, AppName, "foo")
	assert.Equal(t, expected, dir)
}
