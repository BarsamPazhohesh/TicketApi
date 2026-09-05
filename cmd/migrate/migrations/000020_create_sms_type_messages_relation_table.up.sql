CREATE TABLE IF NOT EXISTS sms_type_messages_relation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sms_type_id INTEGER NOT NULL,
    sms_message_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    FOREIGN KEY (sms_type_id) REFERENCES sms_types(id) ON DELETE CASCADE,
    FOREIGN KEY (sms_message_id) REFERENCES sms_messages(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO sms_type_messages_relation (id, sms_type_id, sms_message_id) VALUES
    (1, 1, 1);
