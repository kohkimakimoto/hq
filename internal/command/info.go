package command

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

var InfoCommand = &cli.Command{
	Name:      "info",
	Usage:     `Displays a job detail`,
	ArgsUsage: `<job_id>`,
	Action:    infoAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func infoAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() != 1 {
		return fmt.Errorf("require just one argument as a job ID")
	}

	id, err := strconv.ParseUint(ctx.Args().First(), 10, 64)
	if err != nil {
		return err
	}

	job, err := c.GetJob(id)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(job, "", "  ")
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(ctx.App.Writer, string(b))
	return nil
}
