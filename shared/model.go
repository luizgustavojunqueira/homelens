package shared

import "time"

type Agent struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	LastSeen       time.Time     `json:"last_seen"`
	Online         bool          `json:"online"`
	LatestSnapshot SnapshotEntry `json:"latest_snapshot"`
}

type SnapshotEntry struct {
	Timestamp string     `json:"timestamp"`
	Data      SystemInfo `json:"data"`
}

type SnapshotEvent struct {
	AgentID  string        `json:"agent_id"`
	Snapshot SnapshotEntry `json:"snapshot"`
}

type GetSnapshotsResponse struct {
	Snapshots []SnapshotEntry `json:"snapshots"`
}
