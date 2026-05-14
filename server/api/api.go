package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"homelens/server"
	"homelens/server/db"
)

type API struct {
	registry *server.AgentRegistry
	db       *db.Queries

	logf func(f string, v ...any)
}

func NewAPI(logf func(f string, v ...any), registry *server.AgentRegistry, db *db.Queries) *API {
	return &API{
		registry: registry,
		logf:     logf,
		db:       db,
	}
}

func (api API) GetAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := api.db.ListAgents(context.Background())
	if err != nil {
		http.Error(w, "Failed to list agents", http.StatusInternalServerError)
		return
	}

	agentsResult := make([]Agent, len(agents))
	for i, agent := range agents {

		agentLatestSnapshot, err := api.db.GetLatestSnapshot(context.Background(), agent.ID)
		if err != nil {
			api.logf("GetLatestSnapshot error: %v", err)
			http.Error(w, "Failed to get latest snapshot", http.StatusInternalServerError)
			return
		}

		var entry SnapshotEntry
		if err = json.Unmarshal([]byte(agentLatestSnapshot.Data), &entry.Data); err != nil {
			api.logf("unmarshal snapshot %d error: %v", agentLatestSnapshot.ID, err)
			http.Error(w, "Failed to unmarshal snapshot", http.StatusInternalServerError)
			return
		}

		entry.Timestamp = agentLatestSnapshot.Timestamp.Format(time.RFC3339)

		agentsResult[i] = Agent{
			ID:             agent.ID,
			Name:           agent.Name,
			LastSeen:       agent.LastSeen,
			Online:         api.registry.IsOnline(agent.ID),
			LatestSnapshot: entry,
		}
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(agentsResult)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api API) GetSnapshots(w http.ResponseWriter, r *http.Request) {
	agentID := r.PathValue("id")

	rows, err := api.db.ListSnapshotsByRange(context.Background(), db.ListSnapshotsByRangeParams{
		AgentID:     agentID,
		Timestamp:   time.Now().Add(-24 * time.Hour),
		Timestamp_2: time.Now(),
	})
	if err != nil {
		api.logf("ListSnapshotsByRange error: %v", err)
		http.Error(w, "Failed to list snapshots", http.StatusInternalServerError)
		return
	}

	entries := make([]SnapshotEntry, 0, len(rows))
	for _, row := range rows {
		var entry SnapshotEntry
		if err = json.Unmarshal([]byte(row.Data), &entry.Data); err != nil {
			api.logf("unmarshal snapshot %d error: %v", row.ID, err)
			http.Error(w, "Failed to unmarshal snapshot", http.StatusInternalServerError)
			return
		}
		entry.Timestamp = row.Timestamp.Format(time.RFC3339)
		entries = append(entries, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(GetSnapshotsResponse{Snapshots: entries})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
