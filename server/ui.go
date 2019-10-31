package server

import (
	"github.com/labstack/echo"
	"net/http"
)

func UIHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}
