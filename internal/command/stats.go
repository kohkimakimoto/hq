package command

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
)

var StatsCommand = &cli.Command{
	Name:   "stats",
	Usage:  "Displays the HQ server statistics.",
	Action: statsAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func statsAction(ctx *cli.Context) error {
	c := newClient(ctx)

	stats, err := c.Stats()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(ctx.App.Writer, string(b))
	return nil
}
