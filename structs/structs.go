package structs

type Info struct {
	ServerId   uint   `json:"serverId"`
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	DataDir    string `json:"dataDir"`
}
