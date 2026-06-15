// Package api serves http/websocket endpoints for retrieving agent snapshots data
package api

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"homelens/server"
	"homelens/server/db"
	"homelens/shared"
	"homelens/ui"

	"github.com/coder/websocket"
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

	agentsResult := make([]shared.Agent, len(agents))
	for i, agent := range agents {

		agentLatestSnapshot, err := api.db.GetLatestSnapshot(context.Background(), agent.Guid)
		if err != nil {
			api.logf("GetLatestSnapshot error: %v", err)
			http.Error(w, "Failed to get latest snapshot", http.StatusInternalServerError)
			return
		}

		var entry shared.SnapshotEntry
		if err = json.Unmarshal([]byte(agentLatestSnapshot.Data), &entry.Data); err != nil {
			api.logf("unmarshal snapshot %d error: %v", agentLatestSnapshot.ID, err)
			http.Error(w, "Failed to unmarshal snapshot", http.StatusInternalServerError)
			return
		}

		entry.Timestamp = agentLatestSnapshot.Timestamp.UnixMilli()

		agentsResult[i] = shared.Agent{
			Guid:           agent.Guid,
			Name:           agent.Name.String,
			LastSeen:       agent.LastSeen,
			Online:         api.registry.IsOnline(agent.MachineID),
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
	agentGUID := r.PathValue("guid")

	rows, err := api.db.ListSnapshotsByRange(context.Background(), db.ListSnapshotsByRangeParams{
		AgentGuid:   agentGUID,
		Timestamp:   time.Now().Add(-24 * time.Hour),
		Timestamp_2: time.Now(),
	})
	if err != nil {
		api.logf("ListSnapshotsByRange error: %v", err)
		http.Error(w, "Failed to list snapshots", http.StatusInternalServerError)
		return
	}

	entries := make([]shared.SnapshotEntryRaw, 0, len(rows))
	for _, row := range rows {
		entries = append(entries, shared.SnapshotEntryRaw{
			Timestamp: row.Timestamp.UnixMilli(),
			Data:      []byte(row.Data),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(shared.GetSnapshotsResponse{Snapshots: entries})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api API) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		api.logf("websocket accept error: %v", err)
		return
	}
	defer func() { _ = c.CloseNow() }()

	api.logf("websocket client connected: %s", r.RemoteAddr)

	api.registry.Subscribe(c)

	for {
		_, _, err := c.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				api.logf("websocket client disconnected: %s", r.RemoteAddr)
			} else {
				api.logf("websocket read error: %v", err)
			}
			api.registry.Unsubscribe(c)
			break
		}
	}
}

func (api API) ServeFrontend() http.Handler {
	strippedFS, err := fs.Sub(ui.Assets, "dist")
	if err != nil {
		log.Fatal("erro ao ler o frontend embutido:", err)
	}

	fileServer := http.FileServer(http.FS(strippedFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "."
		}

		_, err := fs.Stat(strippedFS, path)
		if err != nil {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}
