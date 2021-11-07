package ui

// see https://echo.labstack.com/cookbook/embed-resources/

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist
var embeddedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

var AssetHandler = http.FileServer(getFileSystem())
