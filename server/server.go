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

	logf func(f string, v ...any)
}

func NewAgentServer(logf func(f string, v ...any)) *AgentServer {
	return &AgentServer{
		agents: make(map[string]*websocket.Conn),
		logf:   logf,
	}
}

func (as AgentServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		fmt.Printf("Received snapshot from agent: %s\n", snapshot.AgentID)
		fmt.Printf("CPU AVG: %f\n", snapshot.CPUUsage.CPUAvg)

	}

	c.Close(websocket.StatusNormalClosure, "")
}
