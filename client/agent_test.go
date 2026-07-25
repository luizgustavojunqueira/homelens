package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestAgentClient_ConnectAndSend(t *testing.T) {
	token := "test-token"
	machineID := "test-machine-1"

	snapshotReceived := make(chan shared.SystemInfo, 1)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if auth != token {
			t.Errorf("Expected token %s, got %s", token, auth)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.URL.Query().Get("machine_id") != machineID {
			t.Errorf("Expected machine_id %s", machineID)
		}

		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to accept: %v", err)
		}
		defer c.Close(websocket.StatusNormalClosure, "")

		var snap shared.SystemInfo
		if err := wsjson.Read(r.Context(), c, &snap); err == nil {
			snapshotReceived <- snap
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	logf := func(f string, v ...any) {}

	addr := strings.TrimPrefix(server.URL, "http://")
	client := NewAgentClient(logf, machineID, token, addr)

	err := client.Connect()
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	testSnapshot := shared.SystemInfo{
		AgentIP: "192.168.1.100",
	}

	err = client.SendSnapshot(testSnapshot)
	if err != nil {
		t.Fatalf("SendSnapshot failed: %v", err)
	}

	select {
	case snap := <-snapshotReceived:
		if snap.AgentIP != testSnapshot.AgentIP {
			t.Errorf("Expected AgentIP %s, got %s", testSnapshot.AgentIP, snap.AgentIP)
		}
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for snapshot to be received")
	}

	client.Disconnect()
}

func TestAgentClient_ReconnectBackoff(t *testing.T) {
	logf := func(f string, v ...any) {}

	client := NewAgentClient(logf, "test-machine", "token", "localhost:12345")
	client.reconnectDelay = time.Millisecond * 10

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	err := client.reconnect(ctx)

	if err != context.DeadlineExceeded {
		t.Fatalf("Expected DeadlineExceeded, got %v", err)
	}

	if client.reconnectAttempts < 2 {
		t.Errorf("Expected multiple reconnect attempts, got %d", client.reconnectAttempts)
	}
}

func TestReadHostInfo(t *testing.T) {
	info := readHostInfo()
	if info.OS == "" {
		t.Errorf("Expected non-empty OS string from readHostInfo()")
	}
}
