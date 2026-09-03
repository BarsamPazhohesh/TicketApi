CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

INSERT INTO roles(id, title) VALUES
(1, 'Admin'),
(2, 'Agent'),
(3, 'Customer');

