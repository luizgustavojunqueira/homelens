package shared

import (
	"encoding/json"
	"time"
)

type Agent struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	LastSeen       time.Time     `json:"last_seen"`
	Online         bool          `json:"online"`
	LatestSnapshot SnapshotEntry `json:"latest_snapshot"`
}

type SnapshotEntry struct {
	Timestamp int64      `json:"timestamp"`
	Data      SystemInfo `json:"data"`
}

type SnapshotEvent struct {
	AgentID  string        `json:"agent_id"`
	Snapshot SnapshotEntry `json:"snapshot"`
}

type GetSnapshotsResponse struct {
	Snapshots []SnapshotEntryRaw `json:"snapshots"`
}

type SnapshotEntryRaw struct {
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}
