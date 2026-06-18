package shared

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	SnapshotType     EventType = "snapshot"
	StatusChangeType EventType = "status_change"
)

type Agent struct {
	GUID           string        `json:"guid"`
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
	AgentGUID string        `json:"agent_guid"`
	AgentName string        `json:"agent_name"`
	Snapshot  SnapshotEntry `json:"snapshot"`
}

type StatusChangeEvent struct {
	AgentGUID string `json:"agent_guid"`
	Online    bool   `json:"online"`
}

type BroadcastMessage struct {
	Type    EventType `json:"type"`
	Payload any       `json:"payload"`
}

type SnapshotEntryRaw struct {
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

type GetSnapshotsResponse struct {
	Snapshots []SnapshotEntryRaw `json:"snapshots"`
}

type UpdateNameRequest struct {
	Name string `json:"name"`
	GUID string `json:"guid"`
}
