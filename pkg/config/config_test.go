package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilePath(t *testing.T) {
	AppName = "app"
	homeDir, err := os.UserHomeDir()
	assert.Nil(t, err)
	dir := FilePath("foo", "bar", "buzz", "bazz")
	expected := fmt.Sprintf("%s/%s/%s/%s/%s/%s", homeDir, AppName, "foo", "bar", "buzz", "bazz")
	assert.Equal(t, expected, dir)
	dir = FilePath("foo")
	expected = fmt.Sprintf("%s/%s/%s", homeDir, AppName, "foo")
	assert.Equal(t, expected, dir)
}
