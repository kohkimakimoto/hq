package command

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var JobCommand = cli.Command{
	Name:  "job",
	Usage: `Manages jobs.`,
	Subcommands: cli.Commands{
		JobDeleteCommand,
		JobInfoCommand,
		JobListCommand,
		JobRestartCommand,
		JobRunCommand,
		JobStopCommand,
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

		payload := &hq.CreateJobRequest{}
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

		payload := &hq.CreateJobRequest{}
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
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Regular expression `STRING` to filter the jobs with job's name",
		},
		cli.BoolFlag{
			Name:  "reverse, r",
			Usage: "Sort by descending ID.",
		},
		cli.Uint64Flag{
			Name:  "begin, b",
			Usage: "Load the jobs from `ID`.",
		},
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Only display `N` job(s).",
		},
		cli.BoolFlag{
			Name:  "detail, d",
			Usage: "Display detail info.",
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
	detail := ctx.Bool("detail")

	jobs := []*hq.Job{}

	if len(ids) == 0 {
		payload := &hq.ListJobsRequest{
			Name:    ctx.String("name"),
			Reverse: ctx.Bool("reverse"),
			Limit:   ctx.Int("limit"),
		}

		if ctx.Uint64("begin") != 0 {
			b := ctx.Uint64("begin")
			payload.Begin = &b
		}

		list, err := c.ListJobs(payload)
		if err != nil {
			return err
		}
		jobs = list.Jobs
	} else {
		for _, idstr := range ids {
			id, err := strconv.ParseUint(idstr, 10, 64)
			if err != nil {
				return err
			}
			job, err := c.GetJob(id)
			if err != nil {
				return err
			}

			jobs = append(jobs, job)
		}
	}

	t := newTabby()

	if !quiet {
		if detail {
			t.AddLine("ID", "NAME", "COMMENT", "URL", "CREATED", "STARTED", "FINISHED", "DURATION", "STATUS")
		} else {
			t.AddLine("ID", "NAME", "COMMENT", "CREATED", "DURATION", "STATUS")
		}
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
			status = color.Reset(status)
		case "canceled":
			status = color.Grey(status)
		case "canceling":
			status = color.Grey(status)
		case "unfinished":
			status = color.Dim(status)
		case "unknown":
			status = color.Yellow(status)
		}

		createdAt := humanize.Time(job.CreatedAt)
		finishedAt := ""
		startedAt := ""
		duration := ""
		if job.StartedAt != nil {
			startedAt = humanize.Time(*job.StartedAt)
		}
		if job.FinishedAt != nil {
			finishedAt = humanize.Time(*job.FinishedAt)
		}
		if job.StartedAt != nil && job.FinishedAt != nil {
			duration = fmt.Sprintf("%v", job.FinishedAt.Sub(*job.StartedAt))
		}

		comment := strings.Replace(job.Comment, "\n", " ", -1)
		if detail {
			t.AddLine(job.ID, job.Name, comment, job.URL, createdAt, startedAt, finishedAt, duration, status)
		} else {
			t.AddLine(job.ID, job.Name, comment, createdAt, duration, status)
		}
	}

	t.Print()
	return nil
}

var JobInfoCommand = cli.Command{
	Name:      "info",
	Usage:     `Display job detail`,
	ArgsUsage: `<job_id>`,
	Action:    jobInfoAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func jobInfoAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() != 1 {
		return fmt.Errorf("require just one argument as a job ID")
	}

	id, err := strconv.ParseUint(ctx.Args()[0], 10, 64)
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

	fmt.Println(string(b))

	return nil
}

var JobDeleteCommand = cli.Command{
	Name:      "delete",
	Usage:     `Delete a job`,
	ArgsUsage: `<job_id...>`,
	Action:    jobDeleteAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func jobDeleteAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() < 1 {
		return fmt.Errorf("require one id at least")
	}

	t := newTabby()

	args := ctx.Args()
	for _, idstr := range args {
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

var JobRestartCommand = cli.Command{
	Name:      "restart",
	Usage:     `Restart a job`,
	ArgsUsage: `<job_id...>`,
	Action:    jobRestartAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func jobRestartAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() < 1 {
		return fmt.Errorf("require one id at least")
	}

	t := newTabby()

	args := ctx.Args()
	for _, idstr := range args {
		id, err := strconv.ParseUint(idstr, 10, 64)
		if err != nil {
			return err
		}

		job, err := c.RestartJob(id)
		if err != nil {
			return err
		}

		t.AddLine(fmt.Sprintf("%d", job.ID))
	}
	t.Print()
	return nil
}

var JobStopCommand = cli.Command{
	Name:      "stop",
	Usage:     `Stop a job`,
	ArgsUsage: `<job_id...>`,
	Action:    jobStopAction,
	Flags: []cli.Flag{
		addressFlag,
	},
}

func jobStopAction(ctx *cli.Context) error {
	c := newClient(ctx)

	if ctx.NArg() < 1 {
		return fmt.Errorf("require one id at least")
	}

	t := newTabby()

	args := ctx.Args()
	for _, idstr := range args {
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
