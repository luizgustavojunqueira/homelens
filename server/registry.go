package server

import (
	"sync"

	"github.com/coder/websocket"
)

type AgentRegistry struct {
	agents map[string]*websocket.Conn
	mutex  sync.RWMutex
}

func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]*websocket.Conn),
		mutex:  sync.RWMutex{},
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
