package command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestListCommand(t *testing.T) {
	app := testApp(t)
	testRegisterTestClient(t, app, func(req *http.Request) *http.Response {
		b, _ := json.Marshal(&structs.JobList{
			Jobs: []*structs.Job{
				{ID: 1234, Name: "test1"},
				{ID: 1235, Name: "test2"},
				{ID: 1236, Name: "test3"},
			},
			HasNext: false,
		})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	err := app.Run([]string{"hq", "list", "--quiet"})
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(app.Writer.(*bytes.Buffer))
	assert.NoError(t, err)
	assert.Equal(t, "1234\n1235\n1236\n", string(b))
}
