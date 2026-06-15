CREATE TABLE IF NOT EXISTS agents (
    guid        TEXT PRIMARY KEY,
    name        TEXT,
    machine_id  TEXT UNIQUE NOT NULL,
    last_seen   DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS snapshots (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    agent_guid  TEXT NOT NULL REFERENCES agents(guid),
    timestamp   DATETIME NOT NULL,
    data        TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_snapshots_agent_timestamp
    ON snapshots (agent_guid, timestamp);
