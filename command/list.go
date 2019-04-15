package command

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli"
	"strconv"
	"strings"
)

var ListCommand = cli.Command{
	Name:      "list",
	Usage:     `List jobs`,
	ArgsUsage: `[<job_id...>]`,
	Action:    listAction,
	Flags: []cli.Flag{
		addressFlag,
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Only display IDs",
		},
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Specifies a regular expression `STRING` to filter the jobs with job's name",
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
		cli.StringFlag{
			Name:  "status, s",
			Usage: "Specifies `STATUS` to filter the jobs with job's status ('running|waiting|canceling|failure|success|canceled|unfinished|unknown')",
		},
	},
}

func listAction(ctx *cli.Context) error {
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
			Status:  ctx.String("status"),
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
