// Package server receives and handles websocket connections from agents
package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
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

func (as AgentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token != as.token {
		as.logf("%s : %s", token, as.token)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
	defer cancel()

	agentGUID := uuid.New().String()

	agentConnected := true

	for {

		var snapshot shared.SystemInfo
		err := wsjson.Read(ctx, c, &snapshot)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				as.logf("agent disconnected")
			} else {
				as.logf("agent connection lost: %v", err)
			}
			agentConnected = false
			break
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		snapshot.AgentIP = ip

		agent, err := as.db.UpsertAgent(context.Background(), db.UpsertAgentParams{
			Guid:      agentGUID,
			MachineID: machineID,
			LastSeen:  time.Now(),
		})
		if err != nil {
			as.logf("failed to upsert agent in database: %v", err)
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

		agentGUID = agent.Guid

		dbSnapshot := db.InsertSnapshotParams{
			AgentGuid: agentGUID,
			Timestamp: time.Now(),
			Data:      string(data),
		}

		dbErr := as.db.InsertSnapshot(ctx, dbSnapshot)

		if dbErr != nil {
			as.logf("failed to insert snapshot into database: %v", dbErr)
			continue
		}

		broadcastErr := as.registry.Broadcast(shared.BroadcastMessage{
			Type:    shared.SnapshotType,
			Payload: event,
		})

		if broadcastErr != nil {
			as.logf("failed to broadcast agent %s data: %v", machineID, broadcastErr)
		}

	}

	if !agentConnected {

		as.alertCleaner.ClearAlertsForAgent(machineID)

		broadcastErr := as.registry.Broadcast(shared.BroadcastMessage{
			Type: shared.StatusChangeType,
			Payload: shared.StatusChangeEvent{
				AgentGUID: agentGUID,
				Online:    false,
			},
		})

		if broadcastErr != nil {
			as.logf("failed to broadcast agent %s data: %v", machineID, broadcastErr)
		}

		return
	}

	if err := c.Close(websocket.StatusNormalClosure, ""); err != nil {
		as.logf("error closing agent %s websocket cleanly: %v", machineID, err)
	}
}
