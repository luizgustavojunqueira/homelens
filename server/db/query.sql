-- name: UpsertAgent :one
INSERT INTO agents (guid, machine_id, last_seen)
VALUES (?, ?, ?)
ON CONFLICT(machine_id) DO UPDATE SET
    last_seen = excluded.last_seen
RETURNING *;

-- name: UpdateAgentName :exec
UPDATE agents
SET name = ?
WHERE guid = ?;

-- name: GetAgent :one
SELECT * FROM agents WHERE guid = ? LIMIT 1;

-- name: ListAgents :many
SELECT * FROM agents ORDER BY last_seen DESC;

-- name: InsertSnapshot :exec
INSERT INTO snapshots (agent_guid, timestamp, data)
VALUES (?, ?, ?);

-- name: GetLatestSnapshot :one
SELECT * FROM snapshots
WHERE agent_guid = ?
ORDER BY timestamp DESC
LIMIT 1;

-- name: ListSnapshotsByRange :many
SELECT id, agent_guid, timestamp, data
FROM snapshots
WHERE agent_guid = ? AND timestamp >= ? AND timestamp <= ?
ORDER BY timestamp ASC;

-- name: DeleteSnapshotsOlderThan :exec
DELETE FROM snapshots WHERE timestamp < ?;

-- name: GetAlertConfig :one
SELECT * FROM alert_configs LIMIT 1;

-- name: UpsertAlertConfig :one
INSERT INTO alert_configs (id, cpu_threshold, mem_threshold, disk_threshold, offline_mins, tolerance_mins, webhook_url)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    cpu_threshold = excluded.cpu_threshold,
    mem_threshold = excluded.mem_threshold,
    disk_threshold = excluded.disk_threshold,
    offline_mins = excluded.offline_mins,
    tolerance_mins = excluded.tolerance_mins,
    webhook_url = excluded.webhook_url
RETURNING *;
