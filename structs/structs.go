package structs

import (
	"encoding/json"
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
	ID         uint64          `json:"id,string"`
	Name       string          `json:"name"`
	Comment    string          `json:"comment"`
	URL        string          `json:"url"`
	Payload    json.RawMessage `json:"payload"`
	Timeout    int64           `json:"timeout"`
	CreatedAt  time.Time       `json:"created_at"`
	FinishedAt *time.Time      `json:"finished_at"`
	Failure    bool            `json:"failure"`
	Success    bool            `json:"success"`
	Err        string          `json:"err"`
	Output     string          `json:"output"`
	active     bool
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
	Err        string
	Output     string
}

type JobList struct {
	Jobs    []*Job  `json:"jobs"`
	HasNext bool    `json:"hasNext"`
	NextJob *uint64 `json:"nextJob,omitempty"`
	Count   int     `json:"count"`
}
