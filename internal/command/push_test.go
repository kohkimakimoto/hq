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

func TestPushCommand(t *testing.T) {
	t.Run("push with json file", func(t *testing.T) {
		app := testApp(t)
		testRegisterTestClient(t, app, func(req *http.Request) *http.Response {
			reqB, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)

			pushJobRequest := &structs.PushJobRequest{}
			err = json.Unmarshal(reqB, pushJobRequest)
			assert.NoError(t, err)

			b, _ := json.Marshal(&structs.Job{
				ID:      1234,
				Name:    pushJobRequest.Name,
				Comment: pushJobRequest.Comment,
			})
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
				Header:     make(http.Header),
			}
		})

		jsonFile := testTempFile(t, []byte(`
{
  "url": "https://your-worker-app-server/example",
  "name": "example",
  "comment": "This is an example job!",
  "payload": {
    "message": "Hello world!"
  },
  "headers": {
    "X-Custom-Token": "xxxxxxx"
  },
  "timeout": 0
}
`))
		err := app.Run([]string{"hq", "push", jsonFile.Name()})
		assert.NoError(t, err)

		b, err := ioutil.ReadAll(app.Writer.(*bytes.Buffer))
		assert.NoError(t, err)

		assert.Equal(t, "1234\n", string(b))
	})

}
