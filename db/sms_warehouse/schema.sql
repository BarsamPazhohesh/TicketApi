CREATE TABLE sms_warehouse (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    sms_type_id INTEGER NOT NULL DEFAULT 1,
    receiver_phone_number TEXT NOT NULL,
    message TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    FOREIGN KEY (sms_type_id) REFERENCES sms_types(id)
);
