CREATE TABLE IF NOT EXISTS sms_warehouse (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    receiver_phone_number TEXT NOT NULL,
    message TEXT NOT NULL,
    status INTEGER NOT NULL DEFAULT 0, -- 0: pending, 1: sent, 2: failed
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL
);
