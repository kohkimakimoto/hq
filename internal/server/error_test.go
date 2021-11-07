package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("this is a validation error")
	assert.Equal(t, 422, err.Code)
	assert.Equal(t, "this is a validation error", err.Message)
}

func TestErrorHandler(t *testing.T) {
	t.Run("handle internal server error", func(t *testing.T) {
		logBuf := &bytes.Buffer{}
		e := echo.New()
		e.Logger.SetOutput(logBuf)
		e.HTTPErrorHandler = ErrorHandler
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		e.GET("/", func(c echo.Context) error {
			return errors.New("something went wrong")
		})
		e.ServeHTTP(res, req)

		// check error response
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		b, _ := ioutil.ReadAll(res.Body)
		errResp := &structs.ErrorResponse{}
		err := json.Unmarshal(b, errResp)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusText(http.StatusInternalServerError), errResp.Error)
		assert.Equal(t, http.StatusInternalServerError, errResp.Status)

		// check error log
		assert.Regexp(t, "something went wrong", logBuf.String())
	})

	t.Run("handle ErrorValidationFailed", func(t *testing.T) {
		logBuf := &bytes.Buffer{}
		e := echo.New()
		e.Logger.SetOutput(logBuf)
		e.HTTPErrorHandler = ErrorHandler
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		e.GET("/", func(c echo.Context) error {
			return NewValidationError("validation error message that should be sent to a client")
		})
		e.ServeHTTP(res, req)

		// check error response
		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
		b, _ := ioutil.ReadAll(res.Body)
		errResp := &structs.ErrorResponse{}
		err := json.Unmarshal(b, errResp)
		assert.Nil(t, err)
		assert.Equal(t, "validation error message that should be sent to a client", errResp.Error)
		assert.Equal(t, http.StatusUnprocessableEntity, errResp.Status)

		// check that the log is empty
		assert.Equal(t, "", logBuf.String())
	})
}
