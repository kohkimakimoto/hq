package server

import (
	"github.com/kohkimakimoto/govalidator-report"
	"github.com/kohkimakimoto/hq/util/stringutil"
	"github.com/labstack/echo"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ErrorValidationFailedResponse struct {
	Status  int                 `json:"status"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

// NewErrorResponseWithValidatorReport creates new error response with ValidationReport.
func NewErrorValidationFailedResponse(r *report.Report) *ErrorValidationFailedResponse {
	resp := &ErrorValidationFailedResponse{
		Status:  http.StatusUnprocessableEntity,
		Message: "Validation Failed",
	}

	errors := map[string][]string{}
	for _, err := range r.Errors {
		name := stringutil.LowerFirst(err.Name)

		e := errors[name]
		if e == nil {
			e = []string{}
		}

		e = append(e, err.Err.Error())
		errors[name] = e
	}

	resp.Errors = errors
	return resp
}

func NewHttpErrorBadRequest() *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusBadRequest,
		Message: http.StatusText(http.StatusBadRequest),
	}
}

func NewHttpErrorNotFound() *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusNotFound,
		Message: http.StatusText(http.StatusNotFound),
	}
}

func ErrorHandler(err error, c echo.Context) {
	e := c.Echo()

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
		message = http.StatusText(statusCode)
	}

	if err2 := c.JSON(statusCode, &ErrorResponse{
		Status:  statusCode,
		Message: message,
	}); err2 != nil {
		e.Logger.Error(err2)
	}
}
