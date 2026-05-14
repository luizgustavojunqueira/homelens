package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"homelens/server/db"
	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentServer struct {
	registry *AgentRegistry
	db       *db.Queries

	logf  func(f string, v ...any)
	token string
}

func NewAgentServer(logf func(f string, v ...any), token string, registry *AgentRegistry, db *db.Queries) *AgentServer {
	return &AgentServer{
		registry: registry,
		logf:     logf,
		token:    token,
		db:       db,
	}
}

func (as AgentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token != as.token {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	agentID := r.URL.Query().Get("agent_id")
	if agentID == "" {
		http.Error(w, "Missing agent_id", http.StatusBadRequest)
		return
	}

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		as.logf("websocket accept error: %v", err)
		return
	}
	defer c.CloseNow()

	as.logf("agent connected: %s", agentID)
	as.registry.Add(agentID, c)
	defer as.registry.Remove(agentID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
	defer cancel()

	for {

		err = as.db.UpsertAgent(context.Background(), db.UpsertAgentParams{
			ID:       agentID,
			Name:     agentID,
			LastSeen: time.Now(),
		})
		if err != nil {
			as.logf("failed to upsert agent in database: %v", err)
		}

		var snapshot shared.SystemInfo
		err := wsjson.Read(ctx, c, &snapshot)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				as.logf("agent disconnected")
			} else {
				as.logf("agent connection lost: %v", err)
			}
			break
		}

		fmt.Printf("Received snapshot from agent: %s\n", agentID)
		fmt.Printf("CPU AVG: %f\n", snapshot.CPUUsage.CPUAvg)

		dbSnapshot := db.InsertSnapshotParams{
			AgentID:   agentID,
			Timestamp: time.Now(),
			Data:      fmt.Sprintf("%+v", snapshot),
		}

		dbErr := as.db.InsertSnapshot(ctx, dbSnapshot)

		if dbErr != nil {
			as.logf("failed to insert snapshot into database: %v", dbErr)
		}
	}

	c.Close(websocket.StatusNormalClosure, "")
}
