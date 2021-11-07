package command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestDeleteCommand(t *testing.T) {
	app := testApp(t)
	testRegisterTestClient(t, app, func(req *http.Request) *http.Response {
		id, err := strconv.ParseUint(path.Base(req.URL.Path), 10, 64)
		assert.NoError(t, err)

		b, _ := json.Marshal(&structs.DeletedJob{
			ID: id,
		})

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	err := app.Run([]string{"hq", "delete", "1234", "1235", "1236"})
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(app.Writer.(*bytes.Buffer))
	assert.NoError(t, err)

	assert.Equal(t, "1234\n1235\n1236\n", string(b))
}
