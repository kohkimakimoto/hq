package command

import (
	"github.com/urfave/cli"
)

var InfoCommand = cli.Command{
	Name:   "info",
	Usage:  "Display HQ server info",
	Action: infoAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func infoAction(ctx *cli.Context) error {


	return nil
}
