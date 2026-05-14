package api

import (
	"time"

	"homelens/shared"
)

type Agent struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	LastSeen       time.Time     `json:"last_seen"`
	Online         bool          `json:"online"`
	LatestSnapshot SnapshotEntry `json:"latest_snapshot"`
}

type SnapshotEntry struct {
	Timestamp string            `json:"timestamp"`
	Data      shared.SystemInfo `json:"data"`
}

type GetSnapshotsResponse struct {
	Snapshots []SnapshotEntry `json:"snapshots"`
}
