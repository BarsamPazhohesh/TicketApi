CREATE TABLE IF NOT EXISTS api_handlers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    handler TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    UNIQUE(handler, method)
);

INSERT INTO api_handlers (handler, method, description) VALUES ('SampleHandler', 'SampleMethod', 'its a sample record!') RETURNING id;