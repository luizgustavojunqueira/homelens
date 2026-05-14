package api

import (
	"context"
	"encoding/json"
	"net/http"

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
