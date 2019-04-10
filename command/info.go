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
	c := newClient(ctx)

	info, err := c.Info()
	if err != nil {
		return err
	}

	t := newTabby()
	t.AddLine("Version", info.Version)
	t.AddLine("CommitHash", info.CommitHash)
	t.Print()

	return nil
}
