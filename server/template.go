package server

import (
	"github.com/kohkimakimoto/hq/hq"
	"github.com/kohkimakimoto/hq/res/views"
	"github.com/labstack/echo"
	"html/template"
	"io"
	"strings"
)

// see https://echo.labstack.com/guide/templates

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate(app *App) *Template {
	return &Template{
		templates: getCompiledViewTemplate(app),
	}
}

func getCompiledViewTemplate(app *App) *template.Template {
	funcMap := template.FuncMap{
		"Version": func() string {
			return hq.Version
		},
		"CommitHash": func() string {
			return hq.CommitHash
		},
		"Basename": func() string {
			return app.Config.UIBasename
		},
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	t := template.New("views").Funcs(funcMap)
	for _, name := range views.AssetNames() {
		t = template.Must(t.New(name).Parse(string(views.MustAsset(name))))
	}

	return t
}
