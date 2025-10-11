CREATE TABLE IF NOT EXISTS departments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

INSERT INTO departments (title, description) VALUES ('SampleDep', 'its a sample department!');
