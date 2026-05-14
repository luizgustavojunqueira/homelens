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

type SnapshotEntry struct {
	Timestamp string            `json:"timestamp"`
	Data      shared.SystemInfo `json:"data"`
}

type AgentSnapshots struct {
	AgentID   string          `json:"agent_id"`
	Snapshots []SnapshotEntry `json:"snapshots"`
}

type GetSnapshotsResponse struct {
	Agents []AgentSnapshots `json:"agents"`
}
