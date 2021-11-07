package server

import (
	"encoding/json"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/kohkimakimoto/hq/internal/version"
)

type UITemplateRenderer struct {
	templates *template.Template
}

func (t *UITemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewUITemplateRenderer(basename string) *UITemplateRenderer {
	return &UITemplateRenderer{
		templates: getCompiledViewTemplate(basename),
	}
}

func getCompiledViewTemplate(basename string) *template.Template {
	funcMap := template.FuncMap{
		"Version": func() string {
			return version.Version
		},
		"CommitHash": func() string {
			return version.CommitHash
		},
		"Basename": func() string {
			return basename
		},
		"appConfig": func() template.JS {
			ret, _ := json.Marshal(map[string]interface{}{
				"basename":   basename,
				"version":    version.Version,
				"commitHash": version.CommitHash,
			})
			return template.JS(ret)
		},
	}

	t := template.New("views").Funcs(funcMap)
	t = template.Must(t.New("index.html").Parse(strings.TrimSpace(indexHTML)))
	return t
}

var indexHTML = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <meta name="description" content="">
  <link rel="icon" href="{{ Basename }}/dist/favicon.svg" type="image/svg+xml">
  <title>HQ</title>
</head>
<body>
<div id="app"></div>
<script>
var appConfig = {{ appConfig }};
</script>
<script src="{{ Basename }}/dist/js/vendor.js?v={{ CommitHash }}"></script>
<script src="{{ Basename }}/dist/js/app.js?v={{ CommitHash }}"></script>
</body>
</html>
`
