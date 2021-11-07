package command

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

var StopCommand = &cli.Command{
	Name:      "stop",
	Usage:     `Stops a job`,
	ArgsUsage: `<job_id...>`,
	Action:    stopAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func stopAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() < 1 {
		return fmt.Errorf("require one id at least")
	}

	t := newTabby(ctx.App.Writer)

	args := ctx.Args()
	for _, idstr := range args.Slice() {
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			return err
		}

		job, err := c.StopJob(id)
		if err != nil {
			return err
		}

		t.AddLine(fmt.Sprintf("%d", job.ID))
	}
	t.Print()
	return nil
}
