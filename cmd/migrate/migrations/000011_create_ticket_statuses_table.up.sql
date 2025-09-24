CREATE TABLE IF NOT EXISTS ticket_statuses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL UNIQUE,
  description TEXT,
  status INT2 NOT NULL DEFAULT 1,
  deleted INT2 NOT NULL DEFAULT 0
);

INSERT INTO ticket_statuses (title, description) VALUES ('Close', 'This status is for CLOSE only!');
INSERT INTO ticket_statuses (title, description) VALUES ('Open', 'This status is for OPEN only!');
