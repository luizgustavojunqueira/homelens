package client

import (
	"context"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentClient struct {
	logf func(f string, v ...any)
}

func NewAgentClient(logf func(f string, v ...any)) *AgentClient {
	return &AgentClient{
		logf: logf,
	}
}

func (ac AgentClient) Connect(addr string) *websocket.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://"+addr+"/ws", nil)
	if err != nil {
		ac.logf("websocket dial error: %v", err)
		return nil
	}

	return c
}

func (ac AgentClient) Disconnect(c *websocket.Conn) {
	err := c.Close(websocket.StatusNormalClosure, "agent disconnecting")
	if err != nil {
		ac.logf("websocket close error: %v", err)
		return
	}
}

func (ac AgentClient) SendSnapshot(c *websocket.Conn, snapshot shared.SystemInfo) error {
	return wsjson.Write(context.Background(), c, snapshot)
}
