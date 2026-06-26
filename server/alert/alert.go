// Package alert checks for agents exceeding thresholds and send alerts
package alert

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"homelens/server/db"
	"homelens/shared"
)

type AgentRegistry interface {
	GetAllSnapshots() map[string]shared.SnapshotEvent
	Broadcast(event shared.BroadcastMessage) error
}

type Querier interface {
	GetAlertConfig(ctx context.Context) (db.AlertConfig, error)
}

type AlertConfig struct {
	CPUThreshold     int64
	MemThreshold     int64
	DiskThreshold    int64
	OfflineMinutes   time.Duration
	ToleranceMinutes time.Duration
}

type AlertState struct {
	StartTime time.Time
	HaveFired bool
}

type AlertEngine struct {
	store       Querier
	registry    AgentRegistry
	configCache AlertConfig
	state       map[string]*AlertState
	mu          sync.RWMutex
}

func NewEngine(store Querier, registry AgentRegistry) *AlertEngine {
	return &AlertEngine{
		store:    store,
		registry: registry,
		state:    make(map[string]*AlertState),
		mu:       sync.RWMutex{},
	}
}

func (e *AlertEngine) Start(ctx context.Context) error {
	config, err := e.store.GetAlertConfig(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e.configCache = AlertConfig{
				CPUThreshold:     90,
				MemThreshold:     90,
				DiskThreshold:    95,
				OfflineMinutes:   5 * time.Minute,
				ToleranceMinutes: 5 * time.Minute,
			}
		} else {
			return err
		}
	} else {
		e.configCache = AlertConfig{
			CPUThreshold:     config.CpuThreshold.Int64,
			MemThreshold:     config.MemThreshold.Int64,
			DiskThreshold:    config.DiskThreshold.Int64,
			OfflineMinutes:   time.Minute * time.Duration(config.OfflineMins.Int64),
			ToleranceMinutes: time.Minute * time.Duration(config.ToleranceMins.Int64),
		}
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			currentSnapshots := e.registry.GetAllSnapshots()

			for machineID, event := range currentSnapshots {
				snap := event.Snapshot
				var totalCPU float64
				for _, item := range snap.Data.CPU {
					totalCPU += item.UsagePercent
				}
				avgCPU := totalCPU / float64(len(snap.Data.CPU))

				e.evaluateMetric(machineID, "cpu", event.AgentName, avgCPU, float64(e.configCache.CPUThreshold))

				memUsed := float64(snap.Data.Memory.Used)
				memTotal := float64(snap.Data.Memory.Total)
				memUsagePct := (memUsed / memTotal) * 100.0

				e.evaluateMetric(machineID, "mem", event.AgentName, memUsagePct, float64(e.configCache.MemThreshold))

				diskUsagePct := snap.Data.Disk.DiskSpace.UsagePercent
				e.evaluateMetric(machineID, "disk", event.AgentName, diskUsagePct, float64(e.configCache.DiskThreshold))
			}
		}
	}
}

func (e *AlertEngine) UpdateConfig(config AlertConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.configCache = config
}

func (e *AlertEngine) evaluateMetric(machineID, metricName, agentName string, currentValue, threshold float64) {
	e.mu.Lock()
	defer e.mu.Unlock()

	stateKey := fmt.Sprintf("%s_%s", machineID, metricName)
	tolerance := e.configCache.ToleranceMinutes

	if currentValue > threshold {
		agentState := e.state[stateKey]
		if agentState == nil {
			e.state[stateKey] = &AlertState{
				StartTime: time.Now(),
				HaveFired: false,
			}
		} else if !agentState.HaveFired && time.Since(agentState.StartTime) > tolerance {
			agentState.HaveFired = true

			_ = e.registry.Broadcast(shared.BroadcastMessage{
				Type: shared.AlertType,
				Payload: shared.AlertPayload{
					AgentName: agentName,
					Metric:    metricName,
					Value:     currentValue,
					Active:    true,
				},
			})
		}
	} else {
		if agentState, exists := e.state[stateKey]; exists {
			if agentState.HaveFired {
				_ = e.registry.Broadcast(shared.BroadcastMessage{
					Type: shared.AlertType,
					Payload: shared.AlertPayload{
						AgentName: agentName,
						Metric:    metricName,
						Value:     currentValue,
						Active:    false,
					},
				})
			}
			delete(e.state, stateKey)
		}
	}
}
