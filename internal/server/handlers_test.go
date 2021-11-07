package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
	"github.com/kohkimakimoto/hq/internal/version"
)

func TestInfoHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()
		testInitApp(t)
		g.Echo.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		ret := &structs.Info{}
		if err := json.Unmarshal(res.Body.Bytes(), ret); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, version.Version, ret.Version)
		assert.Equal(t, version.CommitHash, ret.CommitHash)
	})
}

func TestCreateJobHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/job", bytes.NewBufferString(`
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
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		testInitApp(t)
		g.Echo.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		job := &structs.Job{}
		if err := json.Unmarshal(res.Body.Bytes(), job); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "https://your-worker-app-server/example", job.URL)
		assert.Equal(t, "example", job.Name)
		assert.Equal(t, "This is an example job!", job.Comment)

		payload := map[string]interface{}{}
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "Hello world!", payload["message"])
		assert.Equal(t, "xxxxxxx", job.Headers["X-Custom-Token"])
		assert.Equal(t, int64(0), job.Timeout)
		assert.Equal(t, false, job.Failure)
		assert.Equal(t, false, job.Success)
	})
}

func TestStatsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/stats", nil)
		res := httptest.NewRecorder()
		testInitApp(t)
		g.Echo.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)

		stats := &structs.Stats{}
		if err := json.Unmarshal(res.Body.Bytes(), stats); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(8192), stats.Queues)
		assert.Equal(t, int64(runtime.NumCPU()), stats.Dispatchers)
		assert.Equal(t, int64(0), stats.MaxWorkers)
		assert.Equal(t, 0, stats.NumJobsInQueue)
		assert.Equal(t, 0, stats.NumJobsWaiting)
		assert.Equal(t, 0, stats.NumJobsRunning)
		assert.Equal(t, int64(0), stats.NumWorkers)
		assert.Equal(t, 0, stats.NumStoredJobs)
	})
}
