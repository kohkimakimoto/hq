package server

import (
	"encoding/json"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/kohkimakimoto/hq/structs"
	"github.com/labstack/echo"
)

func InfoHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	return c.JSON(http.StatusOK, &structs.Info{
		ServerId:   app.Config.ServerId,
		Version:    hq.Version,
		CommitHash: hq.CommitHash,
		DataDir:    app.DataDir,
	})
}

func CreateJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	req := &structs.RegisterJobRequest{}
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

	if req.Code == "" {
		return NewErrorValidationFailed("'code' is required")
	}

	job := &structs.Job{}
	job.ID = id
	job.Name = req.Name
	job.Comment = req.Comment
	job.Code = req.Code
	job.CreatedAt = katsubushi.ToTime(id)

	if err := app.Store.CreateJob(job); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, job)
}

func FetchJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '"+c.Param("id")+"'.")
	}

	job := &structs.Job{}
	if err := app.Store.FetchJob(id, job); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else{
			return err
		}
	}

	return c.JSON(http.StatusOK, job)
}


func DeleteJobHandler(c echo.Context) error {
	app := c.(*AppContext).App()

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return NewErrorValidationFailed("The job id  must be a number but '"+c.Param("id")+"'.")
	}

	if err := app.Store.DeleteJob(id); err != nil {
		if _, ok := err.(*ErrJobNotFound); ok {
			return NewErrorValidationFailed(err.Error())
		} else{
			return err
		}
	}

	return c.JSON(http.StatusOK, &structs.DeletedJob{
		ID: id,
	})
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
