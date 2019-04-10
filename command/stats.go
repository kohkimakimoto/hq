package command

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
)

var StatsCommand = cli.Command{
	Name:   "stats",
	Usage:  "Display HQ server statistics",
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

	fmt.Println(string(b))

	return nil
}
