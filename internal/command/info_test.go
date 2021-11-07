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

func TestInfoCommand(t *testing.T) {
	app := testApp(t)
	testRegisterTestClient(t, app, func(req *http.Request) *http.Response {
		b, _ := json.Marshal(&structs.Job{
			ID:   1234,
			Name: "test job",
		})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	err := app.Run([]string{"hq", "info", "1234"})
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(app.Writer.(*bytes.Buffer))
	assert.NoError(t, err)
	ret := structs.Job{}
	err = json.Unmarshal(b, &ret)
	assert.NoError(t, err)

	assert.Equal(t, uint64(1234), ret.ID)
	assert.Equal(t, "test job", ret.Name)
}
