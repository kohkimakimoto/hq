package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kayac/go-katsubushi"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/kohkimakimoto/hq/internal/structs"
	"github.com/kohkimakimoto/hq/internal/version"
	"github.com/kohkimakimoto/hq/pkg/boltutil"
)

func registerAPIHandlers(e *echo.Echo, prefix string) {
	e.Any(prefix, InfoHandler)
	e.GET(prefix+"stats", StatsHandler)
	e.POST(prefix+"job", PushJobHandler)
	e.GET(prefix+"job", ListJobsHandler)
	e.GET(prefix+"job/:id", GetJobHandler)
	e.DELETE(prefix+"job/:id", DeleteJobHandler)
	e.POST(prefix+"job/:id/stop", StopJobHandler)
	e.POST(prefix+"job/:id/restart", RestartJobHandler)
}

func InfoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &structs.Info{
		Version:    version.Version,
		CommitHash: version.CommitHash,
	})
}

func StatsHandler(c echo.Context) error {
	stats, err := getStats()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

var (
	DefaultJobName = "default"
)

func PushJobHandler(c echo.Context) error {
	req := &structs.PushJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	if req.URL == "" {
		return NewValidationError("'url' is required")
	}

	if req.Name == "" {
		req.Name = DefaultJobName
	}

	id, err := g.IdGen.NextID()
	if err != nil {
		return errors.Wrap(err, "failed to generate uniq id")
	}

	job := &structs.Job{}
	job.ID = id
	job.CreatedAt = katsubushi.ToTime(id)
	job.Name = req.Name
	job.Comment = req.Comment
	job.URL = req.URL
	job.Payload = req.Payload
	job.Headers = req.Headers
	job.Timeout = req.Timeout

	if err := g.Store.CreateJob(job); err != nil {
		return err
	}

	g.QueueManager.EnqueueAsync(job)

	return c.JSON(http.StatusOK, job)
}

func GetJobHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewValidationError("The job id must be a number but '" + c.Param("id") + "'.")
	}

	job, err := g.Store.GetJob(id)
	if err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewValidationError(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, job)
}

func RestartJobHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewValidationError("The job id must be a number but '" + c.Param("id") + "'.")
	}

	req := &structs.RestartJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	job, err := g.Store.GetJob(id)
	if err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewValidationError(err.Error())
		} else {
			return err
		}
	}

	if job.Running {
		return NewValidationError(fmt.Sprintf("The job %d is running now", job.ID))
	}

	if job.Waiting {
		return NewValidationError(fmt.Sprintf("The job %d is waiting now", job.ID))
	}

	if req.Copy {
		id, err := g.IdGen.NextID()
		if err != nil {
			return errors.Wrap(err, "failed to generate uniq id")
		}

		job.ID = id
		job.CreatedAt = katsubushi.ToTime(id)
		job.StartedAt = nil
		job.FinishedAt = nil
		job.Failure = false
		job.Success = false
		job.Canceled = false
		job.StatusCode = nil
		job.Err = ""
		job.Output = ""

		if err := g.Store.CreateJob(job); err != nil {
			return err
		}
	} else {
		job.StartedAt = nil
		job.FinishedAt = nil
		job.Failure = false
		job.Success = false
		job.Canceled = false
		job.StatusCode = nil
		job.Err = ""
		job.Output = ""

		if err := g.Store.UpdateJob(job); err != nil {
			return err
		}
	}

	g.QueueManager.EnqueueAsync(job)

	return c.JSON(http.StatusOK, job)
}

func StopJobHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewValidationError("The job id must be a number but '" + c.Param("id") + "'.")
	}

	job, err := g.Store.GetJob(id)
	if err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewValidationError(err.Error())
		} else {
			return err
		}
	}

	if !job.Running && !job.Waiting {
		return NewValidationError(fmt.Sprintf("The job %d is not active", job.ID))
	}

	g.QueueManager.CancelJob(job.ID)

	return c.JSON(http.StatusOK, &structs.StoppedJob{
		ID: id,
	})
}

func DeleteJobHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewValidationError("The job id must be a number but '" + c.Param("id") + "'.")
	}

	job, err := g.Store.GetJob(id)
	if err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewValidationError(err.Error())
		} else {
			return err
		}
	}

	if job.Running {
		return NewValidationError(fmt.Sprintf("The job %d is running now", job.ID))
	}

	if job.Waiting {
		return NewValidationError(fmt.Sprintf("The job %d is waiting now", job.ID))
	}

	if err := g.Store.DeleteJob(id); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewValidationError(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, &structs.DeletedJob{
		ID: id,
	})
}

func ListJobsHandler(c echo.Context) error {
	req := &structs.ListJobsRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	if req.Limit == 0 {
		req.Limit = g.Config.JobListDefaultLimit
	}

	query := &ListJobsQuery{
		Name:    req.Name,
		Term:    req.Term,
		Begin:   req.Begin,
		Reverse: req.Reverse,
		Limit:   req.Limit,
		Status:  req.Status,
	}

	list, err := g.Store.ListJobs(query)
	if err != nil {
		if err != boltutil.ErrNotFound {
			return errors.Wrap(err, "failed to fetch objects")
		}
	}

	return c.JSON(http.StatusOK, list)
}

func UIIndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func UIFallbackHandler(c echo.Context) error {
	return UIIndexHandler(c)
}

func UIDashboardApiHandler(c echo.Context) error {
	req := &structs.ListJobsRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	if req.Limit == 0 {
		req.Limit = g.Config.JobListDefaultLimit
	}

	query := &ListJobsQuery{
		Name:    req.Name,
		Term:    req.Term,
		Begin:   req.Begin,
		Reverse: req.Reverse,
		Limit:   req.Limit,
		Status:  req.Status,
	}

	list, err := g.Store.ListJobs(query)
	if err != nil {
		if err != boltutil.ErrNotFound {
			return errors.Wrap(err, "failed to fetch objects")
		}
	}

	stats, err := getStats()
	if err != nil {
		return err
	}

	dashboard := &structs.Dashboard{
		Stats:   stats,
		JobList: list,
	}

	return c.JSON(http.StatusOK, dashboard)
}

func getStats() (*structs.Stats, error) {
	var numAllWorkers int64 = 0
	for _, d := range g.Dispatchers {
		numAllWorkers = numAllWorkers + d.NumWorkers()
	}

	numJobs, err := g.Store.CountJobs()
	if err != nil {
		return nil, err
	}

	tt := time.Now().Add(time.Duration(-1) * time.Minute)
	begin := katsubushi.ToID(tt)
	query := &ListJobsQuery{
		Begin: &begin,
	}
	list, err := g.Store.ListJobs(query)
	if err != nil {
		return nil, err
	}

	return &structs.Stats{
		Queues:              g.Config.Queues,
		Dispatchers:         g.Config.Dispatchers,
		MaxWorkers:          g.Config.MaxWorkers,
		NumWorkers:          numAllWorkers,
		NumJobsInQueue:      g.QueueManager.NumJobsInQueue(),
		NumJobsWaiting:      g.QueueManager.NumJobsWaiting(),
		NumJobsRunning:      g.QueueManager.NumJobsRunning(),
		NumStoredJobs:       numJobs,
		NumJobsInLastMinute: list.Count,
	}, nil
}

func bindRequest(req interface{}, c echo.Context) error {
	payload := c.FormValue("payload")

	if payload != "" {
		if err := json.Unmarshal([]byte(payload), req); err != nil {
			return err
		}
	} else {
		httpReq := c.Request()
		if httpReq.ContentLength != 0 || httpReq.Method == http.MethodGet || httpReq.Method == http.MethodDelete {
			if err := c.Bind(req); err != nil {
				return err
			}
		}
	}

	return nil
}
