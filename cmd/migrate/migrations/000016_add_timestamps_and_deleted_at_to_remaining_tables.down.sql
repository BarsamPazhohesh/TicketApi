-- 000016_add_timestamps_and_deleted_at_to_remaining_tables.down.sql

ALTER TABLE users DROP COLUMN deleted_at;

ALTER TABLE departments DROP COLUMN created_at;
ALTER TABLE departments DROP COLUMN updated_at;
ALTER TABLE departments DROP COLUMN deleted_at;

ALTER TABLE ticket_types DROP COLUMN created_at;
ALTER TABLE ticket_types DROP COLUMN updated_at;
ALTER TABLE ticket_types DROP COLUMN deleted_at;

ALTER TABLE ticket_priorities DROP COLUMN created_at;
ALTER TABLE ticket_priorities DROP COLUMN updated_at;
ALTER TABLE ticket_priorities DROP COLUMN deleted_at;

ALTER TABLE ticket_statuses DROP COLUMN created_at;
ALTER TABLE ticket_statuses DROP COLUMN updated_at;
ALTER TABLE ticket_statuses DROP COLUMN deleted_at;

ALTER TABLE api_keys DROP COLUMN created_at;
ALTER TABLE api_keys DROP COLUMN updated_at;
ALTER TABLE api_keys DROP COLUMN deleted_at;
