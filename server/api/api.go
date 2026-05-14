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
		agentsResult[i] = Agent{
			ID:       agent.ID,
			Name:     agent.Name,
			LastSeen: agent.LastSeen,
			Online:   api.registry.IsOnline(agent.ID),
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
	snapshots, err := api.db.ListSnapshotsByRange(context.Background(), db.ListSnapshotsByRangeParams{
		Timestamp:   time.Now().Add(-24 * time.Hour),
		Timestamp_2: time.Now(),
	})
	if err != nil {
		api.logf("ListSnapshotsByRange error: %v", err)
		http.Error(w, "Failed to list snapshots", http.StatusInternalServerError)
		return
	}

	agents := make([]AgentSnapshots, 0, len(snapshots))
	for _, row := range snapshots {
		raw, ok := row.Snapshots.(string)
		if !ok {
			api.logf("snapshots for agent %s is not a string: %T", row.AgentID, row.Snapshots)
			http.Error(w, "Failed to process snapshots", http.StatusInternalServerError)
			return
		}
		var entries []SnapshotEntry
		if err = json.Unmarshal([]byte(raw), &entries); err != nil {
			api.logf("unmarshal snapshots for agent %s error: %v | raw: %s", row.AgentID, err, raw)
			http.Error(w, "Failed to unmarshal snapshots", http.StatusInternalServerError)
			return
		}
		agents = append(agents, AgentSnapshots{
			AgentID:   row.AgentID,
			Snapshots: entries,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(GetSnapshotsResponse{Agents: agents})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
