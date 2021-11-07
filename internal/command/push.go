package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/kohkimakimoto/hq/internal/structs"
)

var PushCommand = &cli.Command{
	Name:  "push",
	Usage: "Pushes a new job.",
	Description: `Pushes a new job.
If you specify '-', it reads a Job JSON from STDIN.`,
	ArgsUsage: `<-|json_file...>`,
	Action:    pushAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func pushAction(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("require one JSON file at least")
	}

	c := newClient(ctx)
	args := ctx.Args()

	if args.Len() == 1 && args.First() == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		payload := &structs.PushJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.PushJob(payload)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(ctx.App.Writer, job.ID)

		return nil
	}

	for _, file := range args.Slice() {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		payload := &structs.PushJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.PushJob(payload)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(ctx.App.Writer, job.ID)
	}

	return nil
}
