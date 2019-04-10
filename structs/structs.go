package structs

import (
	"encoding/json"
	"fmt"
	"time"
)

type Info struct {
	ServerId   uint   `json:"serverId"`
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	DataDir    string `json:"dataDir"`
}

type CreateJobRequest struct {
	Name    string          `json:"name" form:"name" query:"name"`
	Comment string          `json:"comment" form:"comment" query:"comment"`
	URL     string          `json:"url" form:"url" query:"url"`
	Payload json.RawMessage `json:"payload" form:"payload" query:"payload"`
	Timeout int64           `json:"timeout" form:"timeout" query:"timeout"`
}

type ListJobsQuery struct {
	Name     string `json:"-" form:"-" query:"-"`
	Begin    uint64 `json:"-" form:"-" query:"-"`
	HasBegin bool   `json:"-" form:"-" query:"-"`
	Reverse  bool   `json:"-" form:"-" query:"-"`
	Limit    int    `json:"-" form:"-" query:"-"`
}

type Job struct {
	ID         uint64
	Name       string
	Comment    string
	URL        string
	Payload    json.RawMessage
	Timeout    int64
	CreatedAt  time.Time
	FinishedAt *time.Time
	Failure    bool
	Success    bool
	StatusCode int
	Err        string
	Output     string
	Running    bool
}

func (j *Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         fmt.Sprintf("%d", j.ID),
		"name":       j.Name,
		"comment":    j.Comment,
		"url":        j.URL,
		"payload":    j.Payload,
		"createdAt":  j.CreatedAt,
		"finishedAt": j.FinishedAt,
		"failure":    j.Failure,
		"success":    j.Success,
		"statusCode": j.StatusCode,
		"err":        j.Err,
		"output":     j.Output,
		"status":     j.Status(),
	})
}

func (j *Job) Status() string {
	if j.Running {
		return "running"
	} else {
		if j.Failure {
			return "failure"
		} else if j.Success {
			return "success"
		} else {
			if j.FinishedAt == nil {
				return "waiting"
			} else {
				return "unknown"
			}
		}
	}
}

type DeletedJob struct {
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
	FinishedAt *time.Time
	Failure    bool
	Success    bool
	StatusCode int
	Err        string
	Output     string
}

type JobList struct {
	Jobs    []*Job  `json:"jobs"`
	HasNext bool    `json:"hasNext"`
	NextJob *uint64 `json:"nextJob,omitempty"`
	Count   int     `json:"count"`
}
