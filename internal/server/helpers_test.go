package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func testStore(t *testing.T, qm *QueueManager) *Store {
	t.Helper()

	logger := testLogger(t)

	s := NewStore("", logger, qm)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		s.Close()
	})

	return s
}

func testLogger(t *testing.T) echo.Logger {
	t.Helper()

	logger := log.New("test")
	logger.SetLevel(log.OFF)

	return logger
}

func testDispatcher(t *testing.T, qm *QueueManager) *Dispatcher {
	t.Helper()

	return &Dispatcher{
		queueManager:      qm,
		store:             testStore(t, qm),
		logger:            testLogger(t),
		httpClientFactory: defaultHttpClientFactory,
		maxWorkers:        0,
		numWorkers:        0,
	}
}

func testBackgroundCleaner(t *testing.T, qm *QueueManager, tickerDuration time.Duration, jobLifetime int64) *BackgroundCleaner {
	t.Helper()
	return NewBackgroundCleaner(testLogger(t), qm, testStore(t, qm), tickerDuration, jobLifetime)
}

var allocatedServerId uint = 0

func testInitApp(t *testing.T) {
	c := NewConfig()
	// workaround: katsubushi does not allow you to set same worker id for different instance.
	allocatedServerId = allocatedServerId + 1

	c.ServerId = allocatedServerId
	c.LogLevelString = "OFF"

	if err := InitApp(c); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		g.Close()
	})
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
