package command

import (
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/urfave/cli"
	"strconv"
)

var RestartCommand = cli.Command{
	Name:      "restart",
	Usage:     `Restarts a job`,
	ArgsUsage: `<job_id...>`,
	Action:    restartAction,
	Flags: []cli.Flag{
		addressFlag,
		cli.BoolFlag{
			Name:  "copy, c",
			Usage: "Restarts the copied job instead of updating the existed job",
		},
	},
}

func restartAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() < 1 {
		return fmt.Errorf("require one id at least")
	}

	copy := ctx.Bool("copy")

	t := newTabby()

	args := ctx.Args()
	for _, idstr := range args {
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			return err
		}

		req := &hq.RestartJobRequest{
			Copy: copy,
		}
		job, err := c.RestartJob(id, req)
		if err != nil {
			return err
		}

		t.AddLine(fmt.Sprintf("%d", job.ID))
	}
	t.Print()
	return nil
}
