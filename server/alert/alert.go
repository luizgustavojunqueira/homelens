// Package alert checks for agents exceeding thresholds and send alerts
package alert

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
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
	WebhookURL       string
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
			WebhookURL:       config.WebhookUrl.String,
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

				memUsed := float64(snap.Data.Memory.Used)
				memTotal := float64(snap.Data.Memory.Total)
				memUsagePct := (memUsed / memTotal) * 100.0

				diskUsagePct := snap.Data.Disk.DiskSpace.UsagePercent

				lastSeen := time.Since(time.UnixMilli(snap.Timestamp))

				e.evaluateMetric(machineID, "CPU", event.AgentName, avgCPU, float64(e.configCache.CPUThreshold))
				e.evaluateMetric(machineID, "MEM", event.AgentName, memUsagePct, float64(e.configCache.MemThreshold))
				e.evaluateMetric(machineID, "DISK", event.AgentName, diskUsagePct, float64(e.configCache.DiskThreshold))
				e.evaluateMetric(machineID, "OFFLINE", event.AgentName, lastSeen.Minutes(), e.configCache.OfflineMinutes.Minutes())
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

			payload := shared.AlertPayload{
				AgentName: agentName,
				Metric:    metricName,
				Value:     math.Floor(currentValue*100) / 100,
				Active:    true,
			}

			_ = e.registry.Broadcast(shared.BroadcastMessage{
				Type:    shared.AlertType,
				Payload: payload,
			})

			e.triggerWebhook(payload)
		}
	} else {
		if agentState, exists := e.state[stateKey]; exists {
			if agentState.HaveFired {

				payload := shared.AlertPayload{
					AgentName: agentName,
					Metric:    metricName,
					Value:     math.Floor(currentValue*100) / 100,
					Active:    false,
				}

				_ = e.registry.Broadcast(shared.BroadcastMessage{
					Type:    shared.AlertType,
					Payload: payload,
				})

				e.triggerWebhook(payload)
			}
			delete(e.state, stateKey)
		}
	}
}

func (e *AlertEngine) triggerWebhook(payload shared.AlertPayload) {
	go func() {
		if e.configCache.WebhookURL == "" {
			return
		}
		client := &http.Client{Timeout: 10 * time.Second}

		body, err := json.Marshal(payload)
		if err != nil {
			return
		}

		resp, err := client.Post(e.configCache.WebhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return
		}

		_ = resp.Body.Close()
	}()
}

func (e *AlertEngine) ClearAlertsForAgent(machineID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for key := range e.state {
		if strings.HasPrefix(key, machineID+"_") {
			delete(e.state, key)
		}
	}
}
