CREATE TABLE IF NOT EXISTS commands (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    command TEXT NOT NULL,
    description TEXT NOT NULL,
    parameters JSONB NOT NULL DEFAULT (jsonb_array()),
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS secrets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT NOT NULL UNIQUE,
    value TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_commands_name ON commands(name);
CREATE INDEX IF NOT EXISTS idx_secrets_key ON secrets(key);

-- Trigger to automatically update updated_at timestamp for commands
CREATE TRIGGER IF NOT EXISTS update_commands_timestamp
AFTER UPDATE ON commands
FOR EACH ROW
BEGIN
    UPDATE commands SET updated_at = datetime('now') WHERE id = OLD.id;
END;

-- Trigger to automatically update updated_at timestamp for secrets
CREATE TRIGGER IF NOT EXISTS update_secrets_timestamp
AFTER UPDATE ON secrets
FOR EACH ROW
BEGIN
    UPDATE secrets SET updated_at = datetime('now') WHERE id = OLD.id;
END;
