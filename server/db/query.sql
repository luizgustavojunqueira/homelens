-- name: UpsertAgent :exec
INSERT INTO agents (id, name, last_seen)
VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET
    name      = excluded.name,
    last_seen = excluded.last_seen;

-- name: GetAgent :one
SELECT * FROM agents WHERE id = ? LIMIT 1;

-- name: ListAgents :many
SELECT * FROM agents ORDER BY last_seen DESC;

-- name: InsertSnapshot :exec
INSERT INTO snapshots (agent_id, timestamp, data)
VALUES (?, ?, ?);

-- name: GetLatestSnapshot :one
SELECT * FROM snapshots
WHERE agent_id = ?
ORDER BY timestamp DESC
LIMIT 1;

-- name: ListSnapshotsByRange :many
SELECT
    agent_id,
    json_group_array(
        json_object(
            'timestamp', timestamp,
            'data', json(data)
        )
    ) as snapshots
FROM snapshots
WHERE timestamp >= ? AND timestamp <= ?
GROUP BY agent_id
ORDER BY agent_id;

-- name: DeleteSnapshotsOlderThan :exec
DELETE FROM snapshots WHERE timestamp < ?;

