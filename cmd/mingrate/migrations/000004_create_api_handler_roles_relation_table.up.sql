CREATE TABLE IF NOT EXISTS api_handlers_roles_relation(
    api_handler_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (api_handler_id, role_id),
    FOREIGN KEY (api_handler_id) REFERENCES api_handlers(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

INSERT INTO api_handlers_roles_relation (api_handler_id, role_id) VALUES (1,1);