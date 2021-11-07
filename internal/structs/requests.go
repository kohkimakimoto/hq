package structs

import "encoding/json"

type PushJobRequest struct {
	Name    string            `json:"name" form:"name" query:"name"`
	Comment string            `json:"comment" form:"comment" query:"comment"`
	URL     string            `json:"url" form:"url" query:"url"`
	Payload json.RawMessage   `json:"payload" form:"payload" query:"payload"`
	Headers map[string]string `json:"headers" form:"headers" query:"headers"`
	Timeout int64             `json:"timeout" form:"timeout" query:"timeout"`
}

type ListJobsRequest struct {
	Name    string  `query:"name"`
	Term    string  `query:"term"`
	Begin   *uint64 `query:"begin"`
	Reverse bool    `query:"reverse"`
	Limit   int     `query:"limit"`
	Status  string  `query:"status"`
}

type RestartJobRequest struct {
	Copy bool `json:"copy" form:"copy" query:"copy"`
}
