-- 000016_add_timestamps_and_deleted_at_to_remaining_tables.up.sql

-- users
ALTER TABLE users ADD COLUMN deleted_at TEXT DEFAULT NULL;

-- departments
ALTER TABLE departments ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE departments ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE departments ADD COLUMN deleted_at TEXT DEFAULT NULL;

-- ticket_types
ALTER TABLE ticket_types ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE ticket_types ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE ticket_types ADD COLUMN deleted_at TEXT DEFAULT NULL;

-- ticket_priorities
ALTER TABLE ticket_priorities ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE ticket_priorities ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE ticket_priorities ADD COLUMN deleted_at TEXT DEFAULT NULL;

-- ticket_statuses
ALTER TABLE ticket_statuses ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE ticket_statuses ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE ticket_statuses ADD COLUMN deleted_at TEXT DEFAULT NULL;

-- api_keys
ALTER TABLE api_keys ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE api_keys ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE api_keys ADD COLUMN deleted_at TEXT DEFAULT NULL;

UPDATE departments SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE ticket_types SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE ticket_priorities SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE ticket_statuses SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE api_keys SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
