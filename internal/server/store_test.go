package server

import (
	"testing"

	"github.com/kayac/go-katsubushi"
	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestStore_CreateJob(t *testing.T) {
	store := testStore(t, NewQueueManager(10))

	job := &structs.Job{}
	job.ID = 109192606348480512
	job.CreatedAt = katsubushi.ToTime(job.ID)
	job.Name = "test"
	job.Comment = "test comment"
	job.URL = "http://example.com"
	job.Headers = map[string]string{
		"X-Custom-Header": "aaa",
	}
	job.Payload = []byte(`{"foo": "bar"}`)
	job.Timeout = 0

	err := store.CreateJob(job)
	assert.NoError(t, err)
}

func TestStore_GetJob(t *testing.T) {
	store := testStore(t, NewQueueManager(10))

	// create the job
	job := &structs.Job{}
	job.ID = 109192606348480512
	job.CreatedAt = katsubushi.ToTime(job.ID)
	job.Name = "test"
	job.Comment = "test comment"
	job.URL = "http://example.com"
	job.Headers = map[string]string{
		"X-Custom-Header": "aaa",
	}
	job.Payload = []byte(`{"foo": "bar"}`)
	job.Timeout = 0

	err := store.CreateJob(job)
	assert.NoError(t, err)

	// get the job
	job2, err := store.GetJob(109192606348480512)
	assert.NoError(t, err)
	assert.Equal(t, job, job2)
}

func TestStore_UpdateJob(t *testing.T) {
	store := testStore(t, NewQueueManager(10))

	// create the job
	job := &structs.Job{}
	job.ID = 109192606348480512
	job.CreatedAt = katsubushi.ToTime(job.ID)
	job.Name = "test"
	job.Comment = "test comment"
	job.URL = "http://example.com"
	job.Headers = map[string]string{
		"X-Custom-Header": "aaa",
	}
	job.Payload = []byte(`{"foo": "bar"}`)
	job.Timeout = 0

	err := store.CreateJob(job)
	assert.NoError(t, err)

	// get the job
	job2, err := store.GetJob(109192606348480512)
	assert.NoError(t, err)
	assert.Equal(t, job, job2)

	// update the job
	job2.Name = "updated name"
	job2.Comment = "updated comment"
	err = store.UpdateJob(job2)
	assert.NoError(t, err)

	// get the job again
	job3, err := store.GetJob(109192606348480512)
	assert.NoError(t, err)
	assert.Equal(t, "updated name", job3.Name)
	assert.Equal(t, "updated comment", job3.Comment)
}

func TestStore_DeleteJob(t *testing.T) {
	store := testStore(t, NewQueueManager(10))

	// create the job
	job := &structs.Job{}
	job.ID = 109192606348480512
	job.CreatedAt = katsubushi.ToTime(job.ID)
	job.Name = "test"
	job.Comment = "test comment"
	job.URL = "http://example.com"
	job.Headers = map[string]string{
		"X-Custom-Header": "aaa",
	}
	job.Payload = []byte(`{"foo": "bar"}`)
	job.Timeout = 0

	err := store.CreateJob(job)
	assert.NoError(t, err)

	// delete the job
	err = store.DeleteJob(109192606348480512)
	assert.NoError(t, err)

	// check the deletion
	_, err = store.GetJob(109192606348480512)
	assert.IsType(t, &ErrJobNotFound{}, err)
}
