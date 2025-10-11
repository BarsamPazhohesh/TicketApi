--! Update `cmd/mingrate/migrations/*` files after change 

CREATE TABLE app_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_version TEXT NOT NULL,
    version TEXT NOT NULL,
    release_date TEXT NOT NULL DEFAULT (datetime('now')),
    notes TEXT,
    is_current INTEGER NOT NULL DEFAULT 0
);