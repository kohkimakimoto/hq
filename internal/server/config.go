package server

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type Config struct {
	ServerId            uint   `toml:"server_id"`
	LogLevelString      string `toml:"log_level"`
	Addr                string `toml:"addr"`
	Logfile             string `toml:"log_file"`
	DataDir             string `toml:"data_dir"`
	AccessLogfile       string `toml:"access_log_file"`
	Queues              int64  `toml:"queues"`
	Dispatchers         int64  `toml:"dispatchers"`
	MaxWorkers          int64  `toml:"max_workers"`
	ShutdownTimeout     int64  `toml:"shutdown_timeout"`
	JobLifetime         int64  `toml:"job_lifetime"`
	JobListDefaultLimit int    `toml:"job_list_default_limit"`
	UI                  bool   `toml:"ui"`
	UIBasename          string `toml:"ui_basename"`
	IDEpoch             []int  `toml:"id_epoch"`
}

func NewConfig() *Config {
	c := &Config{
		ServerId:            0,
		LogLevelString:      "info",
		Addr:                "0.0.0.0:19900",
		Logfile:             "",
		DataDir:             "",
		AccessLogfile:       "",
		Queues:              8192,
		Dispatchers:         int64(runtime.NumCPU()),
		MaxWorkers:          0,
		ShutdownTimeout:     10,
		JobLifetime:         60 * 60 * 24 * 28, // JobLifetime's unit is second
		JobListDefaultLimit: 0,
		UI:                  true,
		UIBasename:          "/ui",
		IDEpoch:             []int{2019, 1, 1},
	}

	return c
}

func (c *Config) LogLevel() (log.Lvl, error) {
	if c.LogLevelString == "" {
		// default
		return log.INFO, nil
	}

	switch strings.ToLower(c.LogLevelString) {
	case "debug":
		return log.DEBUG, nil
	case "info":
		return log.INFO, nil
	case "warn":
		return log.WARN, nil
	case "error":
		return log.ERROR, nil
	case "off":
		return log.OFF, nil
	}

	return log.INFO, fmt.Errorf("invalid log level %s", c.LogLevelString)
}

func (c *Config) IDEpochTime() (time.Time, error) {
	if len(c.IDEpoch) != 3 {
		return time.Now(), fmt.Errorf("id_epoch must be 3 int values")
	}

	return time.Date(c.IDEpoch[0], time.Month(c.IDEpoch[1]), c.IDEpoch[2], 0, 0, 0, 0, time.UTC), nil
}
