package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentServer struct {
	agents map[string]*websocket.Conn

	logf  func(f string, v ...any)
	token string
}

func NewAgentServer(logf func(f string, v ...any), token string) *AgentServer {
	return &AgentServer{
		agents: make(map[string]*websocket.Conn),
		logf:   logf,
		token:  token,
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour*24)
	defer cancel()

	for {
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

		as.agents[agentID] = c
	}

	c.Close(websocket.StatusNormalClosure, "")
}
