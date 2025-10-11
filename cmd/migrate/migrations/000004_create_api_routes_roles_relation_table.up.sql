CREATE TABLE IF NOT EXISTS api_routes_roles_relation(
    api_route_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (api_route_id, role_id),
    FOREIGN KEY (api_route_id) REFERENCES api_routes(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

INSERT INTO api_routes_roles_relation (api_route_id, role_id) VALUES (1,1);
