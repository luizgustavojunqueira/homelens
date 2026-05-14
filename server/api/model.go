package api

import (
	"time"

	"homelens/shared"
)

type Agent struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
	Online   bool      `json:"online"`
}

type Snapshot struct {
	ID        int64             `json:"id"`
	AgentID   string            `json:"agent_id"`
	Timestamp time.Time         `json:"timestamp"`
	Data      shared.SystemInfo `json:"data"`
}

type GetSnapshotsResponse struct {
	Snapshots []Snapshot `json:"snapshots"`
}
