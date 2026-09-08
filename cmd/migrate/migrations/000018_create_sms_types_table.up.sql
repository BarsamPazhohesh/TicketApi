CREATE TABLE IF NOT EXISTS sms_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL
);

INSERT OR IGNORE INTO sms_types (id, title, description) VALUES
    (1, 'otp', 'کد تایید یکبار مصرف');
