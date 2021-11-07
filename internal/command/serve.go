package command

import "github.com/urfave/cli/v2"

var ServeCommand = &cli.Command{
	Name:   "serve",
	Usage:  "Starts the HQ server process",
	Action: serverAction,
	Flags:  []cli.Flag{},
}

func serverAction(ctx *cli.Context) error {
	return nil
}
