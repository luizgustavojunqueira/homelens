package server

import (
	"context"
	"sync"

	"homelens/shared"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type AgentRegistry struct {
	agents          map[string]*websocket.Conn
	subsConnections []*websocket.Conn
	mutex           sync.RWMutex
	subsMutex       sync.RWMutex
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents:          make(map[string]*websocket.Conn),
		subsConnections: make([]*websocket.Conn, 0),
		mutex:           sync.RWMutex{},
		subsMutex:       sync.RWMutex{},
	}
}

func (ar *AgentRegistry) Add(agentID string, conn *websocket.Conn) {
	ar.mutex.Lock()
	ar.agents[agentID] = conn
	ar.mutex.Unlock()
}

func (ar *AgentRegistry) Remove(agentID string) {
	ar.mutex.Lock()
	delete(ar.agents, agentID)
	ar.mutex.Unlock()
}

func (ar *AgentRegistry) Get(agentID string) (*websocket.Conn, bool) {
	ar.mutex.RLock()
	conn, exists := ar.agents[agentID]
	ar.mutex.RUnlock()
	return conn, exists
}

func (ar *AgentRegistry) IsOnline(agentID string) bool {
	ar.mutex.RLock()
	_, exists := ar.agents[agentID]
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

func (ar *AgentRegistry) Broadcast(event shared.SnapshotEvent) {
	for _, conn := range ar.subsConnections {
		wsjson.Write(context.Background(), conn, event)
	}
}
