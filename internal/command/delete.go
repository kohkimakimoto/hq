package command

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"
)

var DeleteCommand = &cli.Command{
	Name:      "delete",
	Usage:     `Deletes a job`,
	ArgsUsage: `<job_id...>`,
	Action:    deleteAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func deleteAction(ctx *cli.Context) error {
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

		deletedJob, err := c.DeleteJob(id)
		if err != nil {
			return err
		}

		t.AddLine(fmt.Sprintf("%d", deletedJob.ID))
	}
	t.Print()
	return nil
}
