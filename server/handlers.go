package server

import (
	"encoding/json"
	"fmt"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/boltutil"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"sync/atomic"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/echo"
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

	req := &hq.CreateJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return NewHttpErrorBadRequest()
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
	job.Name = req.Name
	job.Comment = req.Comment
	job.URL = req.URL
	job.Payload = req.Payload
	job.Timeout = req.Timeout
	job.CreatedAt = katsubushi.ToTime(id)

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

func StopJobHandler(c echo.Context) error {
	return nil
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

var (
	ListJobsRequestDefaultLimit = 100
)

func ListJobsHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	req := &hq.ListJobsRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return NewHttpErrorBadRequest()
	}

	if req.Limit == 0 {
		req.Limit = ListJobsRequestDefaultLimit
	}

	query := &ListJobsQuery{
		Name:    req.Name,
		Begin:   req.Begin,
		Reverse: req.Reverse,
		Limit:   req.Limit,
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

	stats := &hq.Stats{
		ServerId:        config.ServerId,
		Queues:          config.Queues,
		Dispatchers:     config.Dispatchers,
		MaxWorkers:      config.MaxWorkers,
		ShutdownTimeout: config.ShutdownTimeout,
		JobLifetime:     config.JobLifetime,
		QueueMax:        cap(queueManger.Queue),
		QueueUsage:      len(queueManger.Queue),
		NumWorkers:      numAllWorkers,
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
		if err := c.Bind(req); err != nil {
			return err
		}
	}

	return nil
}
