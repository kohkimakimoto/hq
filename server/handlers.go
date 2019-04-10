package server

import (
	"encoding/json"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/boltutil"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/kohkimakimoto/hq/structs"
	"github.com/labstack/echo"
)

func InfoHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, &structs.Info{
		Version:    hq.Version,
		CommitHash: hq.CommitHash,
	})
}

func CreateJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	req := &structs.CreateJobRequest{}
	if err := bindRequest(req, c); err != nil {
		c.Logger().Warn(errors.Wrap(err, "failed to bind request"))
		return NewHttpErrorBadRequest()
	}

	id, err := app.Gen.NextID()
	if err != nil {
		return errors.Wrap(err, "failed to generate uniq id")
	}

	if req.Name == "" {
		return NewErrorValidationFailed("'name' is required")
	}

	if req.URL == "" {
		return NewErrorValidationFailed("'url' is required")
	}

	job := &structs.Job{}
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

	go app.QueueManager.Enqueue(job)

	return c.JSON(http.StatusOK, job)
}

func GetJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	job := &structs.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, job)
}

func DeleteJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '" + c.Param("id") + "'.")
	}

	if err := app.Store.DeleteJob(id); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else {
			return err
		}
	}

	return c.JSON(http.StatusOK, &structs.DeletedJob{
		ID: id,
	})
}

var (
	ListJobsRequestDefaultLimit = 100
)

func ListJobsHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	query := &structs.ListJobsQuery{}

	// Parse query strings
	query.HasBegin = false
	if c.QueryParam("begin") != "" {
		i, err := strconv.ParseUint(c.QueryParam("begin"), 10, 64)
		if err != nil {
			return NewErrorValidationFailed("The 'begin' must be a number but '" + c.QueryParam("begin") + "'.")
		}
		query.Begin = i
		query.HasBegin = true
	}

	query.Reverse = false
	if c.QueryParam("reverse") != "" {
		query.Reverse = true
	}

	query.Limit = ListJobsRequestDefaultLimit
	if c.QueryParam("limit") != "" {
		l, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			return NewErrorValidationFailed("The 'limit' must be a number but '" + c.QueryParam("limit") + "'.")
		}
		query.Limit = l
	}

	query.Name = c.QueryParam("name")

	list := &structs.JobList{
		Jobs:    []*structs.Job{},
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
