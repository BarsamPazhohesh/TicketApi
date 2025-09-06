CREATE TABLE ticket_priorities (
    user_id INTEGER NOT NULL,
    ticket_type_id INTEGER NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, ticket_type_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_type_id) REFERENCES Ticket_Types(id) ON DELETE CASCADE
);