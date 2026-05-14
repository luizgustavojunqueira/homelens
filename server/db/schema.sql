CREATE TABLE IF NOT EXISTS agents (
    id        TEXT PRIMARY KEY,
    name      TEXT NOT NULL,
    last_seen DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS snapshots (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_id  TEXT NOT NULL REFERENCES agents(id),
    timestamp DATETIME NOT NULL,
    data      TEXT NOT NULL -- JSON
);

CREATE INDEX IF NOT EXISTS idx_snapshots_agent_timestamp
    ON snapshots (agent_id, timestamp);
