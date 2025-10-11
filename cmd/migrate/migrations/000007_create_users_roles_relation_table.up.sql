CREATE TABLE IF NOT EXISTS users_roles_relation (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

INSERT INTO users_roles_relation (user_id, role_id) VALUES (1,1);
