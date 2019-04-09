package structs

import "time"

type Info struct {
	ServerId   uint   `json:"serverId"`
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	DataDir    string `json:"dataDir"`
}

type CreateJobRequest struct {
	Name    string `json:"name" form:"name" query:"name"`
	Comment string `json:"comment" form:"comment" query:"comment"`
	Code    string `json:"code" form:"code" query:"code"`
	Timeout int64  `json:"timeout" form:"timeout" query:"timeout"`
}

type ListJobsQuery struct {
	Name     string `json:"-" form:"-" query:"-"`
	Begin    uint64 `json:"-" form:"-" query:"-"`
	HasBegin bool   `json:"-" form:"-" query:"-"`
	Reverse  bool   `json:"-" form:"-" query:"-"`
	Limit    int    `json:"-" form:"-" query:"-"`
}

var (
	ListJobsRequestDefaultLimit = 100
)

type Job struct {
	ID         uint64     `json:"id,string" gluamapper:"id"`
	Name       string     `json:"name" gluamapper:"name"`
	Comment    string     `json:"comment" gluamapper:"comment"`
	Code       string     `json:"code" gluamapper:"code"`
	Timeout    int64      `json:"timeout"`
	CreatedAt  time.Time  `json:"created_at" gluamapper:"created_at"`
	FinishedAt *time.Time `json:"finished_at" gluamapper:"finished_at"`
	Failure    bool       `json:"failure" gluamapper:"failure"`
	Success    bool       `json:"success" gluamapper:"success"`
	Err        string     `json:"err" gluamapper:"err"`
	Output     string     `json:"output" gluamapper:"output"`
	active     bool
	lockWait   bool
}

type DeletedJob struct {
	ID uint64 `json:"id,string"`
}

// J is internal representation of a job in the boltdb.
type J struct {
	ID         uint64
	Name       string
	Comment    string
	Code       string
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
