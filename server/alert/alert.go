// Package alert checks for agents exceeding thresholds and send alerts
package alert

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"homelens/server/db"
	"homelens/shared"
)

const alertCheckInterval = 10 * time.Second

// webhookClient is shared across all webhook calls to enable connection reuse.
var webhookClient = &http.Client{Timeout: 10 * time.Second}

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

	ticker := time.NewTicker(alertCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			e.mu.RLock()
			cfg := e.configCache
			e.mu.RUnlock()

			currentSnapshots := e.registry.GetAllSnapshots()

			for machineID, event := range currentSnapshots {
				snap := event.Snapshot

				var avgCPU float64
				if len(snap.Data.CPU) > 0 {
					var totalCPU float64
					for _, item := range snap.Data.CPU {
						totalCPU += item.UsagePercent
					}
					avgCPU = totalCPU / float64(len(snap.Data.CPU))
				}

				var memUsagePct float64
				if snap.Data.Memory.Total > 0 {
					memUsagePct = (float64(snap.Data.Memory.Used) / float64(snap.Data.Memory.Total)) * 100.0
				}

				diskUsagePct := snap.Data.Disk.DiskSpace.UsagePercent

				lastSeen := time.Since(time.UnixMilli(snap.Timestamp))

				e.evaluateMetric(machineID, "CPU", event.AgentName, avgCPU, float64(cfg.CPUThreshold), cfg)
				e.evaluateMetric(machineID, "MEM", event.AgentName, memUsagePct, float64(cfg.MemThreshold), cfg)
				e.evaluateMetric(machineID, "DISK", event.AgentName, diskUsagePct, float64(cfg.DiskThreshold), cfg)
				e.evaluateMetric(machineID, "OFFLINE", event.AgentName, lastSeen.Minutes(), cfg.OfflineMinutes.Minutes(), cfg)
			}
		}
	}
}

func (e *AlertEngine) UpdateConfig(config AlertConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.configCache = config
}

func (e *AlertEngine) evaluateMetric(machineID, metricName, agentName string, currentValue, threshold float64, cfg AlertConfig) {
	e.mu.Lock()

	stateKey := fmt.Sprintf("%s_%s", machineID, metricName)
	tolerance := cfg.ToleranceMinutes

	var broadcastMsg *shared.BroadcastMessage
	var webhookPayload *shared.AlertPayload

	if currentValue > threshold {
		agentState := e.state[stateKey]
		if agentState == nil {
			e.state[stateKey] = &AlertState{
				StartTime: time.Now(),
				HaveFired: false,
			}
		} else if !agentState.HaveFired && time.Since(agentState.StartTime) > tolerance {
			agentState.HaveFired = true

			p := shared.AlertPayload{
				AgentName: agentName,
				Metric:    metricName,
				Value:     math.Floor(currentValue*100) / 100,
				Active:    true,
			}
			broadcastMsg = &shared.BroadcastMessage{Type: shared.AlertType, Payload: p}
			webhookPayload = &p
		}
	} else {
		if agentState, exists := e.state[stateKey]; exists {
			if agentState.HaveFired {
				p := shared.AlertPayload{
					AgentName: agentName,
					Metric:    metricName,
					Value:     math.Floor(currentValue*100) / 100,
					Active:    false,
				}
				broadcastMsg = &shared.BroadcastMessage{Type: shared.AlertType, Payload: p}
				webhookPayload = &p
			}
			delete(e.state, stateKey)
		}
	}

	e.mu.Unlock()

	if broadcastMsg != nil {
		if err := e.registry.Broadcast(*broadcastMsg); err != nil {
			log.Printf("failed to broadcast alert for %s/%s: %v", machineID, metricName, err)
		}
	}
	if webhookPayload != nil {
		e.triggerWebhook(cfg.WebhookURL, *webhookPayload)
	}
}

func ValidateWebhookURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https scheme, got %q", parsed.Scheme)
	}
	if parsed.Host == "" {
		return fmt.Errorf("webhook URL must have a host")
	}
	return nil
}

func (e *AlertEngine) triggerWebhook(webhookURL string, payload shared.AlertPayload) {
	if webhookURL == "" {
		return
	}
	go func() {
		if err := ValidateWebhookURL(webhookURL); err != nil {
			log.Printf("invalid webhook URL: %v", err)
			return
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("failed to marshal webhook payload: %v", err)
			return
		}

		resp, err := webhookClient.Post(webhookURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("webhook request failed: %v", err)
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
