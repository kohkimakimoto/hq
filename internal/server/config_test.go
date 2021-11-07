package server

import (
	"testing"

	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestConfig_LogLevel(t *testing.T) {
	c := NewConfig()

	l, err := c.LogLevel()
	assert.Nil(t, err)
	assert.Equal(t, log.INFO, l, "default log level should be info")

	tests := []struct {
		Level    string
		Expected log.Lvl
	}{
		{"debug", log.DEBUG},
		{"info", log.INFO},
		{"warn", log.WARN},
		{"error", log.ERROR},
		{"off", log.OFF},
	}

	for _, test := range tests {
		c.LogLevelString = test.Level
		l, err = c.LogLevel()
		assert.Nil(t, err)
		assert.Equal(t, test.Expected, l)
	}
}
