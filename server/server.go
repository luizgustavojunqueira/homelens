// Package server receives and handles websocket connections from agents
package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"homelens/server/db"
	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
)

type AlertCleaner interface {
	ClearAlertsForAgent(machineID string)
}

type AgentServer struct {
	registry     *AgentRegistry
	db           *db.Queries
	alertCleaner AlertCleaner
	logf         func(f string, v ...any)
	token        string
}

func NewAgentServer(logf func(f string, v ...any), token string, registry *AgentRegistry, db *db.Queries, alertCleaner AlertCleaner) *AgentServer {
	return &AgentServer{
		registry:     registry,
		logf:         logf,
		token:        token,
		db:           db,
		alertCleaner: alertCleaner,
	}
}

func (as *AgentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	if token == "" || token != as.token {
		as.logf("unauthorized agent connection attempt from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	machineID := r.URL.Query().Get("machine_id")
	if machineID == "" {
		http.Error(w, "Missing agent_id", http.StatusBadRequest)
		return
	}

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		as.logf("websocket accept error: %v", err)
		return
	}
	defer func() { _ = c.CloseNow() }()

	as.logf("agent connected: %s", machineID)
	as.registry.Add(machineID, c)
	defer as.registry.Remove(machineID)

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				pingCtx, pingCancel := context.WithTimeout(ctx, 10*time.Second)
				if err := c.Ping(pingCtx); err != nil {
					as.logf("agent %s ping failed: %v", machineID, err)
					pingCancel()
					cancel()
					return
				}
				pingCancel()
			}
		}
	}()

	agent, err := as.db.UpsertAgent(r.Context(), db.UpsertAgentParams{
		Guid:      uuid.New().String(),
		MachineID: machineID,
		LastSeen:  time.Now(),
	})
	if err != nil {
		as.logf("failed to upsert agent in database: %v", err)
		return
	}
	agentGUID := agent.Guid

	for {
		var snapshot shared.SystemInfo
		err := wsjson.Read(ctx, c, &snapshot)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				as.logf("agent disconnected: %s", machineID)
			} else {
				as.logf("agent %s connection lost: %v", machineID, err)
			}
			break
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		snapshot.AgentIP = ip

		if updatedAgent, err := as.db.UpsertAgent(r.Context(), db.UpsertAgentParams{
			Guid:      agentGUID,
			MachineID: machineID,
			LastSeen:  time.Now(),
		}); err != nil {
			as.logf("failed to update agent last_seen: %v", err)
		} else {
			agent = updatedAgent
		}

		event := shared.SnapshotEvent{
			AgentName: agent.Name.String,
			AgentGUID: agentGUID,
			Snapshot: shared.SnapshotEntry{
				Timestamp: time.Now().UnixMilli(),
				Data:      snapshot,
			},
		}

		as.registry.UpsertSnapshot(machineID, event)

		data, err := json.Marshal(snapshot)
		if err != nil {
			as.logf("failed to marshal snapshot: %v", err)
			continue
		}

		dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
		dbErr := as.db.InsertSnapshot(dbCtx, db.InsertSnapshotParams{
			AgentGuid: agentGUID,
			Timestamp: time.Now(),
			Data:      string(data),
		})
		dbCancel()

		if dbErr != nil {
			as.logf("failed to insert snapshot into database: %v", dbErr)
			continue
		}

		if err := as.registry.Broadcast(shared.BroadcastMessage{
			Type:    shared.SnapshotType,
			Payload: event,
		}); err != nil {
			as.logf("failed to broadcast agent %s data: %v", machineID, err)
		}
	}

	as.alertCleaner.ClearAlertsForAgent(machineID)

	if err := as.registry.Broadcast(shared.BroadcastMessage{
		Type: shared.StatusChangeType,
		Payload: shared.StatusChangeEvent{
			AgentGUID: agentGUID,
			Online:    false,
		},
	}); err != nil {
		as.logf("failed to broadcast agent %s disconnect: %v", machineID, err)
	}
}
