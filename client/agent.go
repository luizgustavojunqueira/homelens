// Package client reads data from the current system
package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentClient struct {
	logf              func(f string, v ...any)
	addr              string
	conn              *websocket.Conn
	token             string
	machineID         string
	reconnectAttempts int
	reconnectDelay    time.Duration
}

func GetMachineID() (string, error) {
	machineIDFile, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}
	machineID := strings.TrimSpace(string(machineIDFile))

	return machineID, nil
}

func NewAgentClient(logf func(f string, v ...any), machineID, token, addr string) *AgentClient {
	return &AgentClient{
		logf:              logf,
		token:             token,
		machineID:         machineID,
		conn:              nil,
		addr:              addr,
		reconnectAttempts: 0,
		reconnectDelay:    time.Second * 5,
	}
}

func (ac *AgentClient) Connect() error {
	if ac.conn != nil {
		return fmt.Errorf("already connected to a websocket server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	c, _, err := websocket.Dial(ctx, "ws://"+ac.addr+"/ws?token="+ac.token+"&machine_id="+ac.machineID, nil)
	if err != nil {
		ac.logf("websocket dial error: %v", err)
		return err
	}

	fmt.Printf("Connected to ws://%s/ws?token=%s&machine_id=%s\n", ac.addr, ac.token, ac.machineID)
	ac.conn = c
	return nil
}

func (ac *AgentClient) Disconnect() {
	if ac.conn == nil {
		return
	}

	err := ac.conn.Close(websocket.StatusNormalClosure, "agent disconnecting")
	if err != nil {
		ac.logf("websocket close error: %v", err)
		ac.conn = nil
		return
	}

	ac.conn = nil
}

func (ac *AgentClient) SendSnapshot(snapshot shared.SystemInfo) error {
	if ac.conn == nil {
		return fmt.Errorf("websocket connection is nil, cannot send snapshot")
	}
	return wsjson.Write(context.Background(), ac.conn, snapshot)
}

func (ac *AgentClient) Run(ctx context.Context, interval time.Duration) error {
	ac.logf("Starting agent with ID %s", ac.machineID)
	if err := ac.Connect(); err != nil {
		ac.logf("Failed to connect to server: %v", err)
		return err
	}

	out := make(chan shared.SystemInfo)
	go func() {
		if err := Collect(ctx, interval, out); err != nil {
			ac.logf("Error on collect routine: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down...")
			ac.Disconnect()
			return nil
		case snapshot := <-out:
			if err := ac.SendSnapshot(snapshot); err != nil {
				ac.logf("Connection lost, reconnecting...")
				ac.conn = nil
				if err := ac.reconnect(ctx); err != nil {
					return err
				}
			}
		}
	}
}

func (ac *AgentClient) reconnect(ctx context.Context) error {
	delay := ac.reconnectDelay
	for {
		ac.reconnectAttempts++
		ac.logf("Reconnect attempt %d in %v", ac.reconnectAttempts, delay)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		if err := ac.Connect(); err == nil {
			ac.reconnectAttempts = 0
			return nil
		}

		delay *= 2
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
	}
}
