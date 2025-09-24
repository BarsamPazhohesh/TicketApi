CREATE TABLE api_routes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    route TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    UNIQUE(route, method)
);
