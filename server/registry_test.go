package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestAgentRegistry_AddRemove(t *testing.T) {
	registry := NewAgentRegistry()
	machineID := "agent-123"

	var dummyConn websocket.Conn

	registry.Add(machineID, &dummyConn)
	if !registry.IsOnline(machineID) {
		t.Errorf("Expected agent to be online")
	}

	registry.Remove(machineID)
	if registry.IsOnline(machineID) {
		t.Errorf("Expected agent to be offline")
	}
}

func TestAgentRegistry_Snapshots(t *testing.T) {
	registry := NewAgentRegistry()
	machineID := "agent-123"

	event := shared.SnapshotEvent{
		AgentName: "Test Agent",
		AgentGUID: "guid-123",
		Snapshot: shared.SnapshotEntry{
			Timestamp: 123456789,
		},
	}

	registry.UpsertSnapshot(machineID, event)

	snapshots := registry.GetAllSnapshots()
	if len(snapshots) != 1 {
		t.Fatalf("Expected 1 snapshot, got %d", len(snapshots))
	}

	if snapshots[machineID].AgentName != "Test Agent" {
		t.Errorf("Unexpected snapshot data: %+v", snapshots[machineID])
	}
}

func TestAgentRegistry_Concurrency(t *testing.T) {
	registry := NewAgentRegistry()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			machineID := "agent-concurrency"
			var dummyConn websocket.Conn

			registry.Add(machineID, &dummyConn)
			_ = registry.IsOnline(machineID)
			registry.UpsertSnapshot(machineID, shared.SnapshotEvent{})
			_ = registry.GetAllSnapshots()
			registry.Remove(machineID)
		}(i)
	}

	wg.Wait()
}

func TestAgentRegistry_BroadcastAndSubscriptions(t *testing.T) {
	registry := NewAgentRegistry()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to accept: %v", err)
		}
		registry.Subscribe(c)

		for {
			_, _, err := c.Read(context.Background())
			if err != nil {
				break
			}
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws"+server.URL[4:], nil)
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	time.Sleep(time.Millisecond * 100)

	registry.subsMutex.RLock()
	subsCount := len(registry.subsConnections)
	registry.subsMutex.RUnlock()

	if subsCount != 1 {
		t.Fatalf("Expected 1 subscriber, got %d", subsCount)
	}

	event := shared.BroadcastMessage{
		Type: shared.StatusChangeType,
		Payload: shared.StatusChangeEvent{
			AgentGUID: "some-guid",
			Online:    true,
		},
	}
	err = registry.Broadcast(event)
	if err != nil {
		t.Fatalf("Failed to broadcast: %v", err)
	}

	var receivedEvent shared.BroadcastMessage
	err = wsjson.Read(ctx, c, &receivedEvent)
	if err != nil {
		t.Fatalf("Client failed to read broadcast: %v", err)
	}

	if receivedEvent.Type != shared.StatusChangeType {
		t.Errorf("Expected event type %s, got %s", shared.StatusChangeType, receivedEvent.Type)
	}
}
