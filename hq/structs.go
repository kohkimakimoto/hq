package hq

import (
	"encoding/json"
	"fmt"
	"time"
)

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
}

type Stats struct {
	// config
	ServerId            uint   `json:"serverId"`
	Queues              int64  `json:"queues"`
	Dispatchers         int64  `json:"dispatchers"`
	MaxWorkers          int64  `json:"maxWorkers"`
	ShutdownTimeout     int64  `json:"shutdownTimeout"`
	JobLifetime         int64  `json:"jobLifetime"`
	JobLifetimeStr      string `json:"jobLifetimeStr"`
	JobListDefaultLimit int    `json:"jobListDefaultLimit"`
	// queue stats
	QueueMax       int   `json:"queueMax"`
	QueueUsage     int   `json:"queueUsage"`
	NumWaitingJobs int   `json:"numWaitingJobs"`
	NumRunningJobs int   `json:"numRunningJobs"`
	NumWorkers     int64 `json:"numWorkers"`
	NumJobs        int   `json:"numJobs"`
}

type CreateJobRequest struct {
	Name    string          `json:"name" form:"name" query:"name"`
	Comment string          `json:"comment" form:"comment" query:"comment"`
	URL     string          `json:"url" form:"url" query:"url"`
	Payload json.RawMessage `json:"payload" form:"payload" query:"payload"`
	Timeout int64           `json:"timeout" form:"timeout" query:"timeout"`
}

type ListJobsRequest struct {
	Name    string  `query:"name"`
	Begin   *uint64 `query:"begin"`
	Reverse bool    `query:"reverse"`
	Limit   int     `query:"limit"`
	Status  string  `query:"status"`
}

type Job struct {
	ID         uint64          `json:"id,string"`
	Name       string          `json:"name"`
	Comment    string          `json:"comment"`
	URL        string          `json:"url"`
	Payload    json.RawMessage `json:"payload"`
	Timeout    int64           `json:"timeout"`
	CreatedAt  time.Time       `json:"createdAt"`
	StartedAt  *time.Time      `json:"startedAt"`
	FinishedAt *time.Time      `json:"finishedAt"`
	Failure    bool            `json:"failure"`
	Success    bool            `json:"success"`
	Canceled   bool            `json:"canceled"`
	StatusCode *int            `json:"statusCode"`
	Err        string          `json:"err"`
	Output     string          `json:"output"`
	// status properties.
	Waiting bool `json:"waiting"`
	Running bool `json:"running"`
}

func (j *Job) Status() string {
	if j.Running {
		if j.Canceled {
			return "canceling"
		} else {
			return "running"
		}
	} else if j.Waiting {
		if j.Canceled {
			return "canceling"
		} else {
			return "waiting"
		}
	} else if j.Failure {
		return "failure"
	} else if j.Success {
		return "success"
	} else if j.Canceled {
		return "canceled"
	} else if j.FinishedAt == nil {
		return "unfinished"
	} else {
		return "unknown"
	}
}

func (j *Job) MarshalJSON() ([]byte, error) {
	jobMap := map[string]interface{}{
		"id":         fmt.Sprintf("%d", j.ID),
		"name":       j.Name,
		"comment":    j.Comment,
		"url":        j.URL,
		"payload":    j.Payload,
		"timeout":    j.Timeout,
		"createdAt":  j.CreatedAt,
		"startedAt":  j.StartedAt,
		"finishedAt": j.FinishedAt,
		"failure":    j.Failure,
		"success":    j.Success,
		"canceled":   j.Canceled,
		"statusCode": j.StatusCode,
		"err":        j.Err,
		"output":     j.Output,
		"waiting":    j.Waiting,
		"running":    j.Running,
		"status":     j.Status(),
	}

	return json.Marshal(jobMap)
}

type DeletedJob struct {
	ID uint64 `json:"id,string"`
}

type StoppedJob struct {
	ID uint64 `json:"id,string"`
}

// J is internal representation of a job in the boltdb.
type J struct {
	ID         uint64
	Name       string
	Comment    string
	URL        string
	Payload    json.RawMessage
	Timeout    int64
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Failure    bool
	Success    bool
	Canceled   bool
	StatusCode *int
	Err        string
	Output     string
}

type JobList struct {
	Jobs    []*Job  `json:"jobs"`
	HasNext bool    `json:"hasNext"`
	Next    *uint64 `json:"next,string,omitempty"`
	Count   int     `json:"count"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
