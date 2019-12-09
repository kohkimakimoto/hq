package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/boltutil"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func InfoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &hq.Info{
		Version:    hq.Version,
		CommitHash: hq.CommitHash,
	})
}

var (
	DefaultJobName = "default"
)

func CreateJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	req := &hq.PushJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	if req.URL == "" {
		return NewErrorValidationFailed("'url' is required")
	}

	if req.Name == "" {
		req.Name = DefaultJobName
	}

	id, err := app.Gen.NextID()
	if err != nil {
		return errors.Wrap(err, "failed to generate uniq id")
	}

	job := &hq.Job{}
	job.ID = id
	job.CreatedAt = katsubushi.ToTime(id)
	job.Name = req.Name
	job.Comment = req.Comment
	job.URL = req.URL
	job.Payload = req.Payload
	job.Headers = req.Headers
	job.Timeout = req.Timeout

	if err := app.Store.CreateJob(job); err != nil {
		return err
	}

	app.QueueManager.EnqueueAsync(job)

	return c.JSON(http.StatusOK, job)
}

func GetJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	job := &hq.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, job)
}

func RestartJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	req := &hq.RestartJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	job := &hq.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	if job.Running {
		return NewErrorValidationFailed(fmt.Sprintf("The job %d is running now", job.ID))
	}

	if job.Waiting {
		return NewErrorValidationFailed(fmt.Sprintf("The job %d is waiting now", job.ID))
	}

	if req.Copy {
		id, err := app.Gen.NextID()
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

		if err := app.Store.CreateJob(job); err != nil {
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

		if err := app.Store.UpdateJob(job); err != nil {
			return err
		}
	}

	app.QueueManager.EnqueueAsync(job)

	return c.JSON(http.StatusOK, job)
}

func StopJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	job := &hq.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	if !job.Running && !job.Waiting {
		return NewErrorValidationFailed(fmt.Sprintf("The job %d is not active", job.ID))
	}

	if err := app.QueueManager.CancelJob(job.ID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &hq.StoppedJob{
		ID: id,
	})
}

func DeleteJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	job := &hq.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	if job.Running {
		return NewErrorValidationFailed(fmt.Sprintf("The job %d is running now", job.ID))
	}

	if job.Waiting {
		return NewErrorValidationFailed(fmt.Sprintf("The job %d is waiting now", job.ID))
	}

	if err := app.Store.DeleteJob(id); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, &hq.DeletedJob{
		ID: id,
	})
}

func ListJobsHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	req := &hq.ListJobsRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return err
	}

	if req.Limit == 0 {
		req.Limit = app.Config.JobListDefaultLimit
	}

	query := &ListJobsQuery{
		Name:    req.Name,
		Term:    req.Term,
		Begin:   req.Begin,
		Reverse: req.Reverse,
		Limit:   req.Limit,
		Status:  req.Status,
	}

	list := &hq.JobList{
		Jobs:    []*hq.Job{},
		HasNext: false,
	}
	if err := app.Store.ListJobs(query, list); err != nil {
		if err == boltutil.ErrNotFound {
			return NewHttpErrorNotFound()
		} else {
			return errors.Wrap(err, "failed to fetch objects")
		}
	}

	return c.JSON(http.StatusOK, list)
}

func StatsHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	config := app.Config
	queueManger := app.QueueManager

	var numAllWorkers int64 = 0
	for _, d := range queueManger.Dispatchers {
		numAllWorkers = numAllWorkers + atomic.LoadInt64(&d.NumWorkers)
	}

	jobStats, err := app.Store.JobsStats()
	if err != nil {
		return err
	}

	stats := &hq.Stats{
		Queues:         config.Queues,
		Dispatchers:    config.Dispatchers,
		MaxWorkers:     config.MaxWorkers,
		QueueUsage:     len(queueManger.Queue),
		NumWaitingJobs: len(queueManger.WaitingJobs),
		NumRunningJobs: len(queueManger.RunningJobs),
		NumWorkers:     numAllWorkers,
		NumJobs:        jobStats.KeyN,
	}

	return c.JSON(http.StatusOK, stats)
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
