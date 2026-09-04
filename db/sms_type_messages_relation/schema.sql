CREATE TABLE sms_type_messages_relation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sms_type_id INTEGER NOT NULL,
    sms_message_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    FOREIGN KEY (sms_type_id) REFERENCES sms_types(id) ON DELETE CASCADE,
    FOREIGN KEY (sms_message_id) REFERENCES sms_messages(id) ON DELETE CASCADE
);
