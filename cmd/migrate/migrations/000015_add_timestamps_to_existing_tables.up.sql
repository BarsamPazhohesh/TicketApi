-- 000015_add_timestamps_to_existing_tables.up.sql
ALTER TABLE roles ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE roles ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE roles ADD COLUMN deleted_at TEXT DEFAULT NULL;

ALTER TABLE api_routes ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE api_routes ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE api_routes ADD COLUMN deleted_at TEXT DEFAULT NULL;

ALTER TABLE api_routes_roles_relation ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE api_routes_roles_relation ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE api_routes_roles_relation ADD COLUMN deleted_at TEXT DEFAULT NULL;

ALTER TABLE users_roles_relation ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE users_roles_relation ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE users_roles_relation ADD COLUMN deleted_at TEXT DEFAULT NULL;

ALTER TABLE ticket_types_roles_relation ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE ticket_types_roles_relation ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE ticket_types_roles_relation ADD COLUMN deleted_at TEXT DEFAULT NULL;

ALTER TABLE api_keys_roles_relation ADD COLUMN created_at TEXT DEFAULT NULL;
ALTER TABLE api_keys_roles_relation ADD COLUMN updated_at TEXT DEFAULT NULL;
ALTER TABLE api_keys_roles_relation ADD COLUMN deleted_at TEXT DEFAULT NULL;

UPDATE roles SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE api_routes SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE api_routes_roles_relation SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE users_roles_relation SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE ticket_types_roles_relation SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
UPDATE api_keys_roles_relation SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL;
