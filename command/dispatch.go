package command

import (
	"github.com/urfave/cli"
)

var DispatchCommand = cli.Command{
	Name:   "dispatch",
	Usage:  "Dispatch a job to the HQ server",
	Action: dispatchAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func dispatchAction(ctx *cli.Context) error {


	return nil
}
