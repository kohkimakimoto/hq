package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestClient_Info(t *testing.T) {
	c := New("http://127.0.0.1:19900")
	c.HttpClient = testHttpClient(t, func(req *http.Request) *http.Response {
		b, _ := json.Marshal(&structs.Info{
			Version:    "2.0.0",
			CommitHash: "abcdef",
		})

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	ret, err := c.Info()
	assert.NoError(t, err)
	assert.Equal(t, "2.0.0", ret.Version)
	assert.NotNil(t, "abcdef", ret.CommitHash)
}

func TestClient_PushJob(t *testing.T) {
	c := New("http://127.0.0.1:19900")
	c.HttpClient = testHttpClient(t, func(req *http.Request) *http.Response {
		reqB, err := ioutil.ReadAll(req.Body)
		assert.NoError(t, err)

		pushJobRequest := &structs.PushJobRequest{}
		err = json.Unmarshal(reqB, pushJobRequest)
		assert.NoError(t, err)

		b, _ := json.Marshal(&structs.Job{
			Name:    pushJobRequest.Name,
			Comment: pushJobRequest.Comment,
			URL:     pushJobRequest.URL,
		})

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
			Header:     make(http.Header),
		}
	})

	job, err := c.PushJob(&structs.PushJobRequest{
		Name:    "test",
		Comment: "comment",
		URL:     "http://example.com",
		Payload: nil,
		Headers: map[string]string{},
		Timeout: 0,
	})
	assert.NoError(t, err)
	assert.Equal(t, "test", job.Name)

}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func testHttpClient(t *testing.T, fn RoundTripFunc) *http.Client {
	t.Helper()
	return &http.Client{
		Transport: fn,
	}
}
