-- Drop triggers
DROP TRIGGER IF EXISTS update_secrets_timestamp;
DROP TRIGGER IF EXISTS update_commands_timestamp;

-- Drop indexes
DROP INDEX IF EXISTS idx_secrets_key;
DROP INDEX IF EXISTS idx_commands_name;

-- Drop tables
DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS commands;
