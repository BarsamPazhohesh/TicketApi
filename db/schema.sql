--! Update `cmd/mingrate/migrations/*` files after change 

CREATE TABLE app_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_version TEXT NOT NULL,
    version TEXT NOT NULL,
    release_date TEXT NOT NULL DEFAULT (datetime('now')),
    notes TEXT,
    is_current INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE Roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0
);

CREATE TABLE API_Handlers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    handler TEXT NOT NULL,
    method TEXT NOT NULL,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    UNIQUE(handler, method),
)

CREATE TABLE API_Handlers_Roles_Relation(
    api_handler_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (api_handler_id, role_id),
    FOREIGN KEY (api_handler_id) REFERENCES API_Handlers(id) ON DELETE CASCADE
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE,
)

CREATE TABLE Users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
);

CREATE TABLE Users_Roles_Relation (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, role_id)
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
)

CREATE TABLE Ticket_Types(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
)

CREATE TABLE Ticket_Types_Roles_Relation (
    ticket_type_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (ticket_type_id, role_id),
    FOREIGN KEY (ticket_type_id) REFERENCES Ticket_Types(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES Roles(id) ON DELETE CASCADE
)

CREATE TABLE Priority_Of_Tickets (
    user_id INTEGER NOT NULL,
    ticket_type_id INTEGER NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, ticket_type_id),
    FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_type_id) REFERENCES Ticket_Types(id) ON DELETE CASCADE
)

CREATE TABLE Departments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    deleted INT2 NOT NULL DEFAULT 0,
)