# HomeLens — Task Checklist

## 1. Authentication & Identity

- [x] Define `HOMELENS_AGENT_ID` env var on the client
- [x] Define `HOMELENS_AUTH_TOKEN` env var on the client
- [x] Define `HOMELENS_SERVER_ADDR` env var on the client
- [x] Send agent ID + token as query params on WebSocket connect (`ws://server:8080/ws?token=xxx&agent_id=ubuntu-1`)
- [x] Server validates token before accepting the connection
- [x] Server rejects and closes connection on invalid token
- [x] Server registers agent ID and associates incoming snapshots to it

## 2. Reconnection (Client)

- [x] Detect server disconnect (write/read error)
- [x] Implement reconnect loop with exponential backoff (1s, 2s, 4s, 8s... max 30s)
- [x] Reset backoff on successful reconnection
- [x] Log reconnection attempts
- [x] Continue collecting metrics during disconnect (discard or buffer)

## 3. SQLite (Server)

- [x] Add SQLite dependency (`modernc.org/sqlite` or `mattn/go-sqlite3`)
- [x] Define schema: `agents` table (id, name, last_seen)
- [x] Define schema: `snapshots` table (id, agent_id, timestamp, data JSON)
- [x] Insert snapshot on receive (throttle: 1 per 10s or 1 per minute for storage)
- [ ] Data retention: cron/goroutine to delete snapshots older than X days
- [x] Upsert agent `last_seen` on each snapshot

## 4. Server — Agent Management

- [x] Track connected agents in memory (map of agent ID → connection info)
- [x] Remove agent from map on disconnect
- [x] Endpoint or method to list all agents with status (online/offline/last_seen)
- [x] Store latest snapshot per agent in memory for instant access

## 5. REST API (Server)

- [x] `GET /api/agents` — list all agents with status
- [x] `GET /api/agents/:id` — agent detail with latest snapshot
- [x] `GET /api/agents/:id/history?from=&to=` — historical snapshots for graphs
- [x] `GET /api/stats/:id` — aggregated metrics (avg CPU, max memory over time range)

## 6. Frontend WebSocket (Server → Browser)

- [x] WebSocket endpoint for frontend clients (`/ws/live`)
- [x] On agent snapshot received, broadcast to all connected frontend clients
- [x] Send initial state (all agents + latest snapshots) on frontend connect
- [x] Handle frontend disconnect gracefully

## 7. Frontend (Web UI)

- [x] Choose framework (plain HTML+JS, React, or Templ for Go templates)
- [x] Dashboard page: overview of all agents (name, status, CPU, memory, temp)
- [x] Agent detail page: per-agent graphs (CPU, memory, network, disk over time)
- [x] Online/offline indicator per agent
- [x] Auto-update via WebSocket (live data)
- [x] Historical graphs using REST API data

## 8. Alerts

- [ ] Define alert rules (CPU > 90% for 5 min, disk > 95%, agent offline > 2 min)
- [ ] Alert engine: goroutine evaluating rules against incoming snapshots
- [ ] Alert state management (firing, resolved, cooldown to avoid spam)
- [ ] Log alerts to SQLite
- [ ] Notification channel: Telegram bot integration (webhook/API)
- [ ] Optional: email, generic webhook

## 9. Deploy & Infrastructure

- [ ] Dockerfile for the server
- [ ] Systemd unit file for the agent
- [ ] Build script: `go build -o homelens-agent ./cmd/agent` and `go build -o homelens-server ./cmd/server`
- [ ] `.env.example` for both agent and server
- [ ] README with setup instructions

## 10. Nice to Have

- [ ] Multi-path disk space monitoring (not just `/`)
- [ ] Read thermal zone type from `/sys/class/thermal/thermal_zone*/type` for descriptive labels
- [x] Filter disk I/O to whole disks only (skip partitions, loop, zram, dm)
- [x] Network interface filtering (skip `lo`, optionally skip zero-traffic interfaces)
- [ ] Event log (agent connected, disconnected, alert fired/resolved)
- [ ] Server config file (YAML/TOML) for alert thresholds, retention period, port, etc.
- [ ] Agent config file as alternative to env vars
- [ ] Docker container metrics (via Docker socket)
