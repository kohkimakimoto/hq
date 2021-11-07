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

func TestStatsCommand(t *testing.T) {
	app := testApp(t)
	testRegisterTestClient(t, app, func(req *http.Request) *http.Response {
		b, _ := json.Marshal(&structs.Stats{
			Queues: 1234,
		})
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	err := app.Run([]string{"hq", "stats"})
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(app.Writer.(*bytes.Buffer))
	assert.NoError(t, err)

	ret := structs.Stats{}
	err = json.Unmarshal(b, &ret)
	assert.NoError(t, err)

	assert.Equal(t, int64(1234), ret.Queues)
}
