CREATE TABLE IF NOT EXISTS ticket_types(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

INSERT INTO ticket_types (title, description) VALUES ('SampleType', 'its a sample ticket type!');
