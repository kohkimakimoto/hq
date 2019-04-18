package server

import (
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/labstack/echo"
	"net/http"
	"strings"
)

// NewErrorResponseWithValidatorReport creates new error response with ValidationReport.
func NewErrorValidationFailed(msgs ...string) *echo.HTTPError {
	msg := http.StatusText(http.StatusUnprocessableEntity)
	if len(msg) > 0 {
		msg = strings.Join(msgs, "\n")
	}

	return &echo.HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: msg,
	}
}

func NewHttpErrorBadRequest(msgs ...string) *echo.HTTPError {
	msg := http.StatusText(http.StatusBadRequest)
	if len(msg) > 0 {
		msg = strings.Join(msgs, "\n")
	}

	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NewHttpErrorNotFound(msgs ...string) *echo.HTTPError {
	msg := http.StatusText(http.StatusNotFound)
	if len(msg) > 0 {
		msg = strings.Join(msgs, "\n")
	}

	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: msg,
	}
}

func ErrorHandler(err error, c echo.Context) {
	e := c.Echo()

	e.Logger.Error(fmt.Sprintf("%+v", err))

	var statusCode int
	var message string

	if httperr, ok := err.(*echo.HTTPError); ok {
		statusCode = httperr.Code
		if msg, ok := httperr.Message.(string); ok {
			message = msg
		} else {
			message = http.StatusText(statusCode)
		}
	} else {
		statusCode = http.StatusInternalServerError

		message = err.Error()
		if message == "" {
			message = http.StatusText(statusCode)
		}
	}

	if err2 := c.JSON(statusCode, &hq.ErrorResponse{
		Status: statusCode,
		Error:  message,
	}); err2 != nil {
		e.Logger.Error(fmt.Sprintf("%+v", err2))
	}
}

func errorHandler(app *App) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		e := c.Echo()
		if c.Response().Committed {
			goto ERROR
		}

		ErrorHandler(err, c)
	ERROR:
		e.Logger.Error(err)
	}
}
