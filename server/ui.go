package server

import (
	"github.com/labstack/echo"
	"net/http"
)

func UIIndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func UIFallbackHandler(c echo.Context) error {
	return UIIndexHandler(c)
}
