package server

import (
	"testing"
	"time"
)

func TestBackgroundCleaner_Start(t *testing.T) {
	// just run.
	// TODO: test the background cleaner.
	bg := testBackgroundCleaner(t, NewQueueManager(10), 1*time.Second, 0)
	bg.Start()
	time.Sleep(3 * time.Second)
	bg.Stop()
}
