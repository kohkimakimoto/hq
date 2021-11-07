package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestDispatcher_EventLoop(t *testing.T) {
	t.Run("dispatch a job synchronously", func(t *testing.T) {
		queueManager := NewQueueManager(10)
		d := testDispatcher(t, queueManager)
		d.maxWorkers = 0
		d.httpClientFactory = func() *http.Client {
			return testHttpClient(t, func(req *http.Request) *http.Response {
				// check request headers
				assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
				assert.Equal(t, WorkerDefaultUserAgent, req.Header.Get("User-Agent"))
				assert.Equal(t, "1", req.Header.Get("X-Hq-Job-Id"))

				// check request body as a specified payload
				b, err := ioutil.ReadAll(req.Body)
				assert.Nil(t, err)
				bodyJson := map[string]interface{}{}
				err = json.Unmarshal(b, &bodyJson)
				assert.Nil(t, err)
				assert.Equal(t, "Hello World", bodyJson["message"])

				// check job status in the HQ server
				assert.Equal(t, int64(1), d.NumWorkers())

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       nil,
					Header:     make(http.Header),
				}
			})
		}

		go d.EventLoop()

		queueManager.EnqueueAsync(&structs.Job{
			ID:      1,
			URL:     "http://example.com",
			Payload: []byte(`{"message": "Hello World"}`),
		})

		d.Wait()
	})

	t.Run("dispatch jobs asynchronously", func(t *testing.T) {
		// just run. no check.
		queueManager := NewQueueManager(10)
		d := testDispatcher(t, queueManager)
		d.maxWorkers = 5
		d.httpClientFactory = func() *http.Client {
			return testHttpClient(t, func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       nil,
					Header:     make(http.Header),
				}
			})
		}

		go d.EventLoop()

		for i := 0; i < 5; i++ {
			queueManager.EnqueueAsync(&structs.Job{
				ID:      uint64(i),
				URL:     "http://example.com",
				Payload: []byte(`{"message": "Hello World"}`),
			})
		}

		d.Wait()
	})

	t.Run("dispatch jobs more than max workers", func(t *testing.T) {
		// just run. no check.
		queueManager := NewQueueManager(10)
		d := testDispatcher(t, queueManager)
		d.maxWorkers = 5
		d.httpClientFactory = func() *http.Client {
			return testHttpClient(t, func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       nil,
					Header:     make(http.Header),
				}
			})
		}

		go d.EventLoop()

		for i := 0; i < 10; i++ {
			queueManager.EnqueueAsync(&structs.Job{
				ID:      uint64(i),
				URL:     "http://example.com",
				Payload: []byte(`{"message": "Hello World"}`),
			})
		}

		d.Wait()
	})
}
