CREATE TABLE IF NOT EXISTS ticket_types_roles_relation (
    ticket_type_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (ticket_type_id, role_id),
    FOREIGN KEY (ticket_type_id) REFERENCES Ticket_Types(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

INSERT INTO ticket_types_roles_relation (ticket_type_id, role_id) VALUES (1,1);