package server

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"runtime"
	"time"
)

type Config struct {
	ServerId            uint   `toml:"server_id"`
	LogLevel            string `toml:"log_level"`
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
	IDEpoch             []int  `toml:"id_epoch"`
}

func NewConfig() *Config {
	numCPU := runtime.NumCPU()

	return &Config{
		ServerId:        0,
		LogLevel:        "info",
		Addr:            "0.0.0.0:19900",
		Logfile:         "",
		DataDir:         "",
		AccessLogfile:   "",
		Queues:          8192,
		Dispatchers:     int64(numCPU),
		MaxWorkers:      0,
		ShutdownTimeout: 10,
		// JobLifetime's unit is second
		JobLifetime:         60 * 60 * 24 * 28,
		JobListDefaultLimit: 0,
		IDEpoch:             []int{2019, 1, 1},
	}
}

func (c *Config) SetLogLevel(level string) {
	c.LogLevel = level
}

func (c *Config) LoadConfigFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) IDEpochTime() (time.Time, error) {
	if len(c.IDEpoch) != 3 {
		return time.Now(), fmt.Errorf("id_epoch must be 3 int values")
	}

	return time.Date(c.IDEpoch[0], time.Month(c.IDEpoch[1]), c.IDEpoch[2], 0, 0, 0, 0, time.UTC), nil
}
