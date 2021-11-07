package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func NewValidationError(message string) *echo.HTTPError {
	return &echo.HTTPError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

func ErrorHandler(err error, c echo.Context) {
	hErr := transformToHTTPError(err)
	if hErr.Code >= 500 {
		c.Logger().Error(err)
	}

	if c.Response().Committed {
		return
	}

	var message string
	if msg, ok := hErr.Message.(string); ok {
		message = msg
	} else {
		message = http.StatusText(hErr.Code)
	}

	if err := c.JSON(hErr.Code, &structs.ErrorResponse{
		Status: hErr.Code,
		Error:  message,
	}); err != nil {
		c.Logger().Error(err)
	}
}

func transformToHTTPError(err error) *echo.HTTPError {
	if hErr, ok := err.(*echo.HTTPError); ok {
		return hErr
	}

	code := http.StatusInternalServerError
	hErr := echo.NewHTTPError(code, http.StatusText(code))
	hErr.Internal = err
	return hErr
}
