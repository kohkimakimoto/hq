package command

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v2"

	"github.com/kohkimakimoto/hq/internal/server"
)

var ServeCommand = &cli.Command{
	Name:   "serve",
	Usage:  "Starts the HQ server process",
	Action: serverAction,
	Flags: []cli.Flag{
		configFileFlag,
		logLevelFlag,
	},
}

func serverAction(ctx *cli.Context) error {
	config := server.NewConfig()
	// Load config file
	if path := getConfigFilePath(ctx); path != "" {
		if _, err := toml.DecodeFile(path, config); err != nil {
			return err
		}
	}

	// override log level
	if v := ctx.String("log-level"); v != "" {
		config.LogLevelString = v
	}

	return server.Run(config)
}

func getConfigFilePath(ctx *cli.Context) string {
	if configFile := ctx.String("config-file"); configFile != "" {
		// config file specified by CLI option
		return configFile
	}

	// Use default system config if it exists.
	configFile := "/etc/hq/hq.toml"
	if _, err := os.Stat(configFile); err == nil {
		return configFile
	}

	// It means that there is no config file.
	return ""
}
