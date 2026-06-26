// Package api serves http/websocket endpoints for retrieving agent snapshots data
package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	"homelens/server"
	"homelens/server/alert"
	"homelens/server/db"
	"homelens/shared"
	"homelens/ui"

	"github.com/coder/websocket"
)

type Querier interface {
	ListAgents(ctx context.Context) ([]db.Agent, error)
	GetLatestSnapshot(ctx context.Context, agentGUID string) (db.Snapshot, error)
	ListSnapshotsByRange(ctx context.Context, arg db.ListSnapshotsByRangeParams) ([]db.Snapshot, error)
	UpdateAgentName(ctx context.Context, arg db.UpdateAgentNameParams) error
	UpsertAlertConfig(ctx context.Context, arg db.UpsertAlertConfigParams) (db.AlertConfig, error)
	GetAlertConfig(ctx context.Context) (db.AlertConfig, error)
}

type AlertController interface {
	UpdateConfig(config alert.AlertConfig)
}

type API struct {
	registry *server.AgentRegistry
	db       Querier
	alerts   AlertController
	logf     func(f string, v ...any)
}

func NewAPI(logf func(f string, v ...any), registry *server.AgentRegistry, db Querier, alerts AlertController) *API {
	return &API{
		registry: registry,
		logf:     logf,
		db:       db,
		alerts:   alerts,
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
			GUID:           agent.Guid,
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

func (api API) UpdateAgentName(w http.ResponseWriter, r *http.Request) {
	var updateNameRequest shared.UpdateNameRequest
	err := json.NewDecoder(r.Body).Decode(&updateNameRequest)
	if err != nil {
		http.Error(w, "Failed to decode response", http.StatusBadRequest)
		return
	}

	err = api.db.UpdateAgentName(context.Background(), db.UpdateAgentNameParams{
		Name: sql.NullString{String: updateNameRequest.Name, Valid: true},
		Guid: updateNameRequest.GUID,
	})
	if err != nil {
		http.Error(w, "Failed to list agents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(true)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api API) SaveAlertConfig(w http.ResponseWriter, r *http.Request) {
	var req shared.UpdateAlertConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	_, err := api.db.UpsertAlertConfig(r.Context(), db.UpsertAlertConfigParams{
		ID:            1,
		CpuThreshold:  sql.NullInt64{Int64: req.CPUThreshold, Valid: true},
		MemThreshold:  sql.NullInt64{Int64: req.MemThreshold, Valid: true},
		DiskThreshold: sql.NullInt64{Int64: req.DiskThreshold, Valid: true},
		OfflineMins:   sql.NullInt64{Int64: req.OfflineThreshold, Valid: true},
		ToleranceMins: sql.NullInt64{Int64: req.ToleranceMinutes, Valid: true},
		WebhookUrl:    sql.NullString{String: req.WebhookURL, Valid: true},
	})
	if err != nil {
		http.Error(w, "failed to save to database", http.StatusInternalServerError)
		return
	}

	api.alerts.UpdateConfig(alert.AlertConfig{
		CPUThreshold:     req.CPUThreshold,
		MemThreshold:     req.MemThreshold,
		DiskThreshold:    req.DiskThreshold,
		OfflineMinutes:   time.Duration(req.OfflineThreshold) * time.Minute,
		ToleranceMinutes: time.Duration(req.ToleranceMinutes) * time.Minute,
	})

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(true)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (api API) GetAlertConfig(w http.ResponseWriter, r *http.Request) {
	config, err := api.db.GetAlertConfig(r.Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			config = db.AlertConfig{
				CpuThreshold:  sql.NullInt64{Int64: 90, Valid: true},
				MemThreshold:  sql.NullInt64{Int64: 90, Valid: true},
				DiskThreshold: sql.NullInt64{Int64: 90, Valid: true},
				OfflineMins:   sql.NullInt64{Int64: 5, Valid: true},
				ToleranceMins: sql.NullInt64{Int64: 5, Valid: true},
			}
		} else {
			api.logf("GetAlertConfig error: %v", err)
			http.Error(w, "Failed to get alert config", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(shared.GetAlertConfigResponse{
		CPUThreshold:     config.CpuThreshold.Int64,
		MemThreshold:     config.MemThreshold.Int64,
		DiskThreshold:    config.DiskThreshold.Int64,
		OfflineThreshold: config.OfflineMins.Int64,
		ToleranceMinutes: config.ToleranceMins.Int64,
		WebhookURL:       config.WebhookUrl.String,
	})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
