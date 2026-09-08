CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT,
    department_id INTEGER NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    Foreign Key (department_id) REFERENCES departments(id)
);

-- Seed SuperAdmin (username: 09120000000, password: Password123)
INSERT INTO users (id, username, password, department_id) VALUES
(1, '09120000000', '$2a$10$NNAVqBU3f0RgRPi76jkCreKMFkpu/Jq97OJ2t3hc7tv4fjF4pG71K', 1);

