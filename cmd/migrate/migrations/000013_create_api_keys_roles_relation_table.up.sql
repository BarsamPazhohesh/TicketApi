CREATE TABLE api_keys_roles_relation (
    api_key_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (api_key_id, role_id),
    FOREIGN KEY (api_key_id) REFERENCES api_Keys(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

INSERT INTO api_keys_roles_relation (api_key_id, role_id) VALUES (1, 1);
