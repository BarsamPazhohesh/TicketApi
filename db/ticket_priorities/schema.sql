CREATE TABLE ticket_priorities (
    user_id INTEGER NOT NULL,
    ticket_type_id INTEGER NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    status INT2 NOT NULL DEFAULT 1,
    PRIMARY KEY (user_id, ticket_type_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_type_id) REFERENCES ticket_types(id) ON DELETE CASCADE
);
