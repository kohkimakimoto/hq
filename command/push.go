package command

import (
	"encoding/json"
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

var PushCommand = cli.Command{
	Name:  "push",
	Usage: "Pushes a new job.",
	Description: `Pushes a new job.
If you specify '-', it will use stdin as a job JSON .`,
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

	if len(args) == 1 && args[0] == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		payload := &hq.PushJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.PushJob(payload)
		if err != nil {
			return err
		}

		fmt.Println(job.ID)

		return nil
	}

	for _, file := range args {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}

		payload := &hq.PushJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.PushJob(payload)
		if err != nil {
			return err
		}

		fmt.Println(job.ID)
	}

	return nil
}
