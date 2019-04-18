package command

import (
	"github.com/kohkimakimoto/hq/server"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"os"
)

var ServeCommand = cli.Command{
	Name:   "serve",
	Usage:  "Start the hq server process",
	Action: serverAction,
	Flags: []cli.Flag{
		configFileFlag,
		logLevelFlag,
	},
}

func serverAction(ctx *cli.Context) error {
	config := server.NewConfig()

	if err := loadServerConfigFiles(ctx, config); err != nil {
		return err
	}

	applyLogLevel(ctx, config)

	app := server.NewApp(config)
	defer app.Close()

	return app.ListenAndServe()
}

func loadServerConfigFiles(ctx *cli.Context, config *server.Config) error {
	if v := ctx.String("config-file"); v != "" {
		// config file specified by CLI option
		if err := config.LoadConfigFile(v); err != nil {
			return errors.Wrapf(err, "failed to load config from the file '%s'", v)
		}
	} else {
		// default system config
		configFile := "/etc/hq/hq.toml"

		if _, err := os.Stat(configFile); err == nil {
			if err := config.LoadConfigFile(configFile); err != nil {
				return errors.Wrapf(err, "failed to load config from the file '%s'", configFile)
			}
		}
	}

	return nil
}
