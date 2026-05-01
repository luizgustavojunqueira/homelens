package client

import (
	"context"
	"fmt"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentClient struct {
	logf    func(f string, v ...any)
	conn    *websocket.Conn
	token   string
	agentID string
}

func NewAgentClient(logf func(f string, v ...any), token, agentID string) *AgentClient {
	return &AgentClient{
		logf:    logf,
		token:   token,
		agentID: agentID,
		conn:    nil,
	}
}

func (ac *AgentClient) Connect(addr string) error {
	if ac.conn != nil {
		return fmt.Errorf("already connected to a websocket server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://"+addr+"/ws?token="+ac.token+"&agent_id="+ac.agentID, nil)
	if err != nil {
		ac.logf("websocket dial error: %v", err)
		return err
	}

	fmt.Printf("Connected to ws://%s/ws?token=%s&agent_id=%s\n", addr, ac.token, ac.agentID)
	ac.conn = c
	return nil
}

func (ac *AgentClient) Disconnect() {
	if ac.conn == nil {
		return // No connection to close
	}

	err := ac.conn.Close(websocket.StatusNormalClosure, "agent disconnecting")
	if err != nil {
		ac.logf("websocket close error: %v", err)
		return
	}
}

func (ac *AgentClient) SendSnapshot(snapshot shared.SystemInfo) error {
	if ac.conn == nil {
		return fmt.Errorf("websocket connection is nil, cannot send snapshot")
	}
	return wsjson.Write(context.Background(), ac.conn, snapshot)
}
