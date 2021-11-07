package structs

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
	Queues              int64 `json:"queues"`
	Dispatchers         int64 `json:"dispatchers"`
	MaxWorkers          int64 `json:"maxWorkers"`
	NumWorkers          int64 `json:"numWorkers"`
	NumJobsInQueue      int   `json:"numJobsInQueue"`
	NumJobsWaiting      int   `json:"numJobsWaiting"`
	NumJobsRunning      int   `json:"numJobsRunning"`
	NumStoredJobs       int   `json:"numStoredJobs"`
	NumJobsInLastMinute int   `json:"numJobsInLastMinute"`
}

type Job struct {
	ID         uint64            `json:"id,string"`
	Name       string            `json:"name"`
	Comment    string            `json:"comment"`
	URL        string            `json:"url"`
	Payload    json.RawMessage   `json:"payload"`
	Headers    map[string]string `json:"headers"`
	Timeout    int64             `json:"timeout"`
	CreatedAt  time.Time         `json:"createdAt"`
	StartedAt  *time.Time        `json:"startedAt"`
	FinishedAt *time.Time        `json:"finishedAt"`
	Failure    bool              `json:"failure"`
	Success    bool              `json:"success"`
	Canceled   bool              `json:"canceled"`
	StatusCode *int              `json:"statusCode"`
	Err        string            `json:"err"`
	Output     string            `json:"output"`
	Waiting    bool              `json:"waiting"`
	Running    bool              `json:"running"`
}

const (
	JobStatusWaiting    = "waiting"
	JobStatusRunning    = "running"
	JobStatusCanceling  = "canceling"
	JobStatusCanceled   = "canceled"
	JobStatusFailure    = "failure"
	JobStatusSuccess    = "success"
	JobStatusUnfinished = "unfinished"
	JobStatusUnknown    = "unknown"
)

func (j *Job) Status() string {
	if j.Running {
		if j.Canceled {
			return JobStatusCanceling
		} else {
			return JobStatusRunning
		}
	} else if j.Waiting {
		if j.Canceled {
			return JobStatusCanceling
		} else {
			return JobStatusWaiting
		}
	} else if j.Failure {
		return JobStatusFailure
	} else if j.Success {
		return JobStatusSuccess
	} else if j.Canceled {
		return JobStatusCanceled
	} else if j.FinishedAt == nil {
		return JobStatusUnfinished
	} else {
		return JobStatusUnknown
	}
}

func (j *Job) MarshalJSON() ([]byte, error) {
	jobMap := map[string]interface{}{
		"id":         fmt.Sprintf("%d", j.ID),
		"name":       j.Name,
		"comment":    j.Comment,
		"url":        j.URL,
		"payload":    j.Payload,
		"headers":    j.Headers,
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

type Dashboard struct {
	Stats   *Stats   `json:"stats"`
	JobList *JobList `json:"jobList"`
}
