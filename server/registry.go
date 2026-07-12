package server

import (
	"context"
	"log"
	"maps"
	"sync"
	"time"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentRegistry struct {
	agents          map[string]*websocket.Conn
	latestSnapshots map[string]shared.SnapshotEvent
	subsConnections []*websocket.Conn
	mutex           sync.RWMutex
	subsMutex       sync.RWMutex
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents:          make(map[string]*websocket.Conn),
		latestSnapshots: make(map[string]shared.SnapshotEvent),
		subsConnections: make([]*websocket.Conn, 0),
	}
}

func (ar *AgentRegistry) Add(machineID string, conn *websocket.Conn) {
	ar.mutex.Lock()
	ar.agents[machineID] = conn
	ar.mutex.Unlock()
}

func (ar *AgentRegistry) Remove(machineID string) {
	ar.mutex.Lock()
	delete(ar.agents, machineID)
	ar.mutex.Unlock()
}

func (ar *AgentRegistry) IsOnline(machineID string) bool {
	ar.mutex.RLock()
	_, exists := ar.agents[machineID]
	ar.mutex.RUnlock()
	return exists
}

func (ar *AgentRegistry) Subscribe(conn *websocket.Conn) {
	ar.subsMutex.Lock()
	ar.subsConnections = append(ar.subsConnections, conn)
	ar.subsMutex.Unlock()
}

func (ar *AgentRegistry) Unsubscribe(conn *websocket.Conn) {
	ar.subsMutex.Lock()
	for i, c := range ar.subsConnections {
		if c == conn {
			ar.subsConnections = append(ar.subsConnections[:i], ar.subsConnections[i+1:]...)
			break
		}
	}
	ar.subsMutex.Unlock()
}

// Broadcast sends event to all subscribed frontend connections.
// It takes a snapshot of subscribers under the lock then releases it before
// writing, so individual slow/dead subscribers don't hold the lock.
// Dead connections are cleaned up after the write loop.
func (ar *AgentRegistry) Broadcast(event shared.BroadcastMessage) error {
	ar.subsMutex.RLock()
	subs := make([]*websocket.Conn, len(ar.subsConnections))
	copy(subs, ar.subsConnections)
	ar.subsMutex.RUnlock()

	var dead []*websocket.Conn
	for _, conn := range subs {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := wsjson.Write(ctx, conn, event)
		cancel()
		if err != nil {
			log.Printf("broadcast write failed, dropping subscriber: %v", err)
			_ = conn.CloseNow()
			dead = append(dead, conn)
		}
	}

	// Remove dead connections after releasing the read lock.
	for _, conn := range dead {
		ar.Unsubscribe(conn)
	}

	return nil
}

func (ar *AgentRegistry) UpsertSnapshot(machineID string, snap shared.SnapshotEvent) {
	ar.mutex.Lock()
	ar.latestSnapshots[machineID] = snap
	ar.mutex.Unlock()
}

func (ar *AgentRegistry) GetAllSnapshots() map[string]shared.SnapshotEvent {
	ar.mutex.RLock()
	defer ar.mutex.RUnlock()

	copyMap := make(map[string]shared.SnapshotEvent)
	maps.Copy(copyMap, ar.latestSnapshots)
	return copyMap
}
