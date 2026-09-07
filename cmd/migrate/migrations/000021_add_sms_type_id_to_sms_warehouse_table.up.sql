ALTER TABLE sms_warehouse ADD COLUMN sms_type_id INTEGER NOT NULL DEFAULT 1 REFERENCES sms_types(id);
