-- 000015_add_timestamps_to_existing_tables.down.sql
ALTER TABLE roles DROP COLUMN deleted_at;
ALTER TABLE roles DROP COLUMN updated_at;
ALTER TABLE roles DROP COLUMN created_at;

ALTER TABLE api_routes DROP COLUMN deleted_at;
ALTER TABLE api_routes DROP COLUMN updated_at;
ALTER TABLE api_routes DROP COLUMN created_at;

ALTER TABLE api_routes_roles_relation DROP COLUMN deleted_at;
ALTER TABLE api_routes_roles_relation DROP COLUMN updated_at;
ALTER TABLE api_routes_roles_relation DROP COLUMN created_at;

ALTER TABLE users_roles_relation DROP COLUMN deleted_at;
ALTER TABLE users_roles_relation DROP COLUMN updated_at;
ALTER TABLE users_roles_relation DROP COLUMN created_at;

ALTER TABLE ticket_types_roles_relation DROP COLUMN deleted_at;
ALTER TABLE ticket_types_roles_relation DROP COLUMN updated_at;
ALTER TABLE ticket_types_roles_relation DROP COLUMN created_at;

ALTER TABLE api_keys_roles_relation DROP COLUMN deleted_at;
ALTER TABLE api_keys_roles_relation DROP COLUMN updated_at;
ALTER TABLE api_keys_roles_relation DROP COLUMN created_at;
