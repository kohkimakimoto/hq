package server

import (
	"github.com/kohkimakimoto/hq/test"
	"os"
	"testing"
)

func TestConfig_SetLogLevel(t *testing.T) {
	c := &Config{}
	// just run
	c.SetLogLevel("info")
}


func TestConfig_LoadConfigFile(t *testing.T) {
	c := &Config{}

	tmpFile, err := test.CreateTempfile([]byte(`
server_id = 1
addr = "localhost:1234"
`))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	c.LoadConfigFile(tmpFile.Name())

	if c.ServerId != 1 {
		t.Errorf("c.ServerId must be 1 but %d", c.ServerId)
	}

	if c.Addr != "localhost:1234" {
		t.Errorf("c.Addr must be 1 but %s", c.Addr)
	}
}

