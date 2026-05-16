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
	addr    string
	conn    *websocket.Conn
	token   string
	agentID string
}

func NewAgentClient(logf func(f string, v ...any), token, agentID string, addr string) *AgentClient {
	return &AgentClient{
		logf:    logf,
		token:   token,
		agentID: agentID,
		conn:    nil,
		addr:    addr,
	}
}

func (ac *AgentClient) Connect() error {
	if ac.conn != nil {
		return fmt.Errorf("already connected to a websocket server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://"+ac.addr+"/ws?token="+ac.token+"&agent_id="+ac.agentID, nil)
	if err != nil {
		ac.logf("websocket dial error: %v", err)
		return err
	}

	fmt.Printf("Connected to ws://%s/ws?token=%s&agent_id=%s\n", ac.addr, ac.token, ac.agentID)
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

func (ac *AgentClient) Run(ctx context.Context, interval time.Duration) error {
	ac.logf("Starting agent with ID %s", ac.agentID)
	if err := ac.Connect(); err != nil {
		ac.logf("Failed to connect to server: %v", err)
		return err
	}

	out := make(chan shared.SystemInfo)
	go Collect(ctx, interval, out)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down...")
			ac.Disconnect()
			return nil
		case snapshot := <-out:
			if err := ac.SendSnapshot(snapshot); err != nil {
				fmt.Println("send error:", err)
				continue
			}
		}
	}
}
