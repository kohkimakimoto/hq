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

type CreateJobRequest struct {
	Name    string          `json:"name" form:"name" query:"name"`
	Comment string          `json:"comment" form:"comment" query:"comment"`
	URL     string          `json:"url" form:"url" query:"url"`
	Payload json.RawMessage `json:"payload" form:"payload" query:"payload"`
	Timeout int64           `json:"timeout" form:"timeout" query:"timeout"`
}

type ListJobsRequest struct {
	Name    string  `json:"name" form:"name" query:"name"`
	Begin   string  `json:"begin"`
	Reverse bool    `json:"reverse"`
	Limit   uint64  `json:"limit"`
}

type ListJobsQuery struct {
	Name     string `json:"-" form:"-" query:"-"`
	Begin    uint64 `json:"-" form:"-" query:"-"`
	HasBegin bool   `json:"-" form:"-" query:"-"`
	Reverse  bool   `json:"-" form:"-" query:"-"`
	Limit    int    `json:"-" form:"-" query:"-"`
}

type Job struct {
	ID         uint64          `json:"id,string"`
	Name       string          `json:"name"`
	Comment    string          `json:"comment"`
	URL        string          `json:"url"`
	Payload    json.RawMessage `json:"payload"`
	Timeout    int64           `json:"timeout"`
	CreatedAt  time.Time       `json:"createdAt"`
	FinishedAt *time.Time      `json:"finishedAt"`
	Failure    bool            `json:"failure"`
	Success    bool            `json:"success"`
	StatusCode int             `json:"statusCode"`
	Err        string          `json:"err"`
	Output     string          `json:"output"`
	Running    bool            `json:"running"`
}

func (j *Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":         fmt.Sprintf("%d", j.ID),
		"name":       j.Name,
		"comment":    j.Comment,
		"url":        j.URL,
		"payload":    j.Payload,
		"timeout":    j.Timeout,
		"createdAt":  j.CreatedAt,
		"finishedAt": j.FinishedAt,
		"failure":    j.Failure,
		"success":    j.Success,
		"statusCode": j.StatusCode,
		"err":        j.Err,
		"output":     j.Output,
		"running":    j.Running,
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
	NextJob *uint64 `json:"nextJob,string,omitempty"`
	Count   int     `json:"count"`
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
