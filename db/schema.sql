--! Update `cmd/mingrate/migrations/*` files after change 

CREATE TABLE app_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_version TEXT NOT NULL,
    version TEXT NOT NULL,
    release_date TEXT NOT NULL DEFAULT (datetime('now')),
    notes TEXT,
    is_current INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

CREATE TABLE api_handlers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    handler TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    UNIQUE(handler, method)
);

CREATE TABLE api_handlers_roles_relation(
    api_handler_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (api_handler_id, role_id),
    FOREIGN KEY (api_handler_id) REFERENCES api_handlers(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    department_id INTEGER NOT NULL UNIQUE,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    Foreign Key (department_id) REFERENCES departments(id)
);

CREATE TABLE users_roles_relation (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

CREATE TABLE ticket_types(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

CREATE TABLE ticket_types_roles_relation (
    ticket_type_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (ticket_type_id, role_id),
    FOREIGN KEY (ticket_type_id) REFERENCES Ticket_Types(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
);

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

CREATE TABLE departments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);
