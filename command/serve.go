package command

import (
	"github.com/urfave/cli"
)

var ServeCommand = cli.Command{
	Name:   "serve",
	Usage:  "Start the hq server process",
	Action: serverAction,
	Flags: []cli.Flag{
		configFileFlag,
	},
}

func serverAction(ctx *cli.Context) error {

	return nil
}