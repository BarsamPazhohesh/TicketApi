CREATE TABLE api_handlers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    handler TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    UNIQUE(handler, method)
);