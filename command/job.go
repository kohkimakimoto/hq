package command

import (
	"encoding/json"
	"fmt"
	"github.com/kohkimakimoto/hq/structs"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"strconv"
)

var JobCommand = cli.Command{
	Name:  "job",
	Usage: `Manages jobs.`,
	Subcommands: cli.Commands{
		JobRunCommand,
		JobListCommand,
	},
}

var JobRunCommand = cli.Command{
	Name:  "run",
	Usage: "Enqueues and runs a job.",
	Description: `Enqueues and runs a job.
If you specify '-', it will use stdin as a job JSON .`,
	ArgsUsage: `<-|json_file...>`,
	Action:    jobRunAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func jobRunAction(ctx *cli.Context) error {
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

		payload := &structs.CreateJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.CreateJob(payload)
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

		payload := &structs.CreateJobRequest{}
		if err := json.Unmarshal(b, payload); err != nil {
			return err
		}

		job, err := c.CreateJob(payload)
		if err != nil {
			return err
		}

		fmt.Println(job.ID)
	}

	return nil
}

var JobListCommand = cli.Command{
	Name:      "list",
	Usage:     `List jobs`,
	ArgsUsage: `[<job_id...>]`,
	Action:    jobListAction,
	Flags: []cli.Flag{
		addressFlag,
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Only display IDs",
		},
		cli.StringSliceFlag{
			Name:  "name, n",
			Usage: "Regular expression to filter the jobs with job's name",
		},
		cli.BoolFlag{
			Name:  "reverse, r",
			Usage: "Sort by descending ID.",
		},
		cli.StringFlag{
			Name:  "begin, b",
			Usage: "Load the jobs from `ID`.",
		},
		cli.Uint64Flag{
			Name:  "limit, l",
			Usage: "Only display `N` job(s).",
			Value: 100,
		},
	},
}

func jobListAction(ctx *cli.Context) error {
	c := newClient(ctx)

	ids := []string{}
	if ctx.NArg() > 0 {
		ids = ctx.Args()
	}

	quiet := ctx.Bool("quiet")

	jobs := []*structs.Job{}

	if len(ids) == 0 {
		payload := &structs.ListJobsRequest{
			Name: ctx.String("name"),
			Reverse: ctx.Bool("reverse"),
			Begin:   ctx.String("begin"),
			Limit:   ctx.Uint64("limit"),
		}

		list, err := c.ListJobs(payload)
		if err != nil {
			return err
		}
		jobs = list.Jobs
	} else {
		for _, idstr := range ids {
			id, err := strconv.ParseUint(idstr, 10, 64)
			job, err := c.GetJob(id)
			if err != nil {
				return err
			}

			jobs = append(jobs, job)
		}
	}

	t := newTabby()

	if !quiet {
		t.AddLine("ID", "NAME", "URL", "STATUS", "CREATED_AT", "FINISHED_AT", "ERROR")
	}

	for _, job := range jobs {
		if quiet {
			t.AddLine(job.ID)
			continue
		}

		status := job.Status()
		switch status {
		case "failure":
			status = color.Red(status)
		case "success":
			status = color.Green(status)
		case "running":
			status = color.Cyan(status)
		case "waiting":
			status = color.Dim(status)
		case "unknown":
			status = color.Yellow(status)
		}

		createdAt := fmt.Sprintf("%v", job.CreatedAt)
		finishedAt := ""
		if job.FinishedAt != nil {
			finishedAt = fmt.Sprintf("%v", job.FinishedAt)
		}

		t.AddLine(job.ID, job.Name, job.URL, status, createdAt, finishedAt, job.Err)
	}

	t.Print()
	return nil

	return nil
}
