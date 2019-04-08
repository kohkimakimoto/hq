package server

import (
	"github.com/kohkimakimoto/hq/hq"
	"net/http"

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
