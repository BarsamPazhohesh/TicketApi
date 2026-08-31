-- 000014_create_permissions_table.up.sql
CREATE TABLE IF NOT EXISTS permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    status INT2 NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS roles_permissions_relation (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS api_routes_permissions_relation (
    api_route_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    status INT2 NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    deleted_at TEXT DEFAULT NULL,
    PRIMARY KEY (api_route_id, permission_id),
    FOREIGN KEY (api_route_id) REFERENCES api_routes(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

-- Seed Default Roles (1: Admin, 2: Agent, 3: Customer)
INSERT OR IGNORE INTO roles (id, title, status) VALUES
(1, 'Admin', 1),
(2, 'Agent', 1),
(3, 'Customer', 1);

-- Seed System Permissions
INSERT OR IGNORE INTO permissions (id, name, description) VALUES
(1, 'tickets:create', 'Create a new ticket'),
(2, 'tickets:read_own', 'Read own tickets'),
(3, 'tickets:read_all', 'Read all tickets in system'),
(4, 'tickets:chat', 'Send chat messages on tickets'),
(5, 'users:read', 'Read user profiles and lists'),
(6, 'users:manage', 'Create or manage users and roles'),
(7, 'files:upload', 'Upload files for tickets'),
(8, 'files:download', 'Download ticket attachments'),
(9, 'metadata:read', 'Read ticket types, statuses and departments');

-- Seed API Routes (All routes registered in system)
INSERT OR IGNORE INTO api_routes (id, route, method, description) VALUES
(1, 'tickets/CreateTicket/', 'POST', 'Create ticket'),
(2, 'tickets/GetTicketByID/', 'POST', 'Get ticket by ID'),
(3, 'tickets/:id/CreateChat/', 'POST', 'Create chat on ticket'),
(4, 'tickets/GetTicketByTrackCode/', 'POST', 'Get ticket by track code'),
(5, 'tickets/GetTicketsList/', 'POST', 'List tickets'),
(6, 'tickets/GetAllActiveTicketTypes/', 'GET', 'List active ticket types'),
(7, 'tickets/GetAllActiveTicketStatuses/', 'GET', 'List active ticket statuses'),
(8, 'auth/LoginWithNoAuth/', 'POST', 'Guest login'),
(9, 'auth/SignUp/', 'POST', 'User sign up'),
(10, 'auth/Login/', 'POST', 'User password login'),
(11, 'auth/GetSingleUseToken/', 'POST', 'API Key token minting'),
(12, 'auth/LoginWithSingleUseToken/', 'GET', 'Login with one-time token'),
(13, 'captcha/GetCaptcha/', 'GET', 'Get captcha image'),
(14, 'captcha/VerifyCaptcha/', 'POST', 'Verify captcha solution'),
(15, 'departments/GetAllActiveDepartments/', 'GET', 'List active departments'),
(16, 'users/GetUserByUsername/', 'POST', 'Find user by username'),
(17, 'users/GetUserByID/', 'POST', 'Find user by ID'),
(18, 'users/GetUsersByIDs/', 'POST', 'Batch get users by IDs'),
(19, 'files/UploadTicketFile/', 'POST', 'Upload ticket attachment'),
(20, 'files/GetDownloadLinkTicketFile/:objectName', 'POST', 'Get attachment download URL');

-- Map Role Permissions
-- Admin (Role 1): All permissions (1-9)
INSERT OR IGNORE INTO roles_permissions_relation (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8), (1, 9);

-- Agent (Role 2): tickets:read_all, tickets:chat, users:read, files:upload, files:download, metadata:read
INSERT OR IGNORE INTO roles_permissions_relation (role_id, permission_id) VALUES
(2, 3), (2, 4), (2, 5), (2, 7), (2, 8), (2, 9);

-- Customer (Role 3): tickets:create, tickets:read_own, tickets:chat, files:upload, files:download, metadata:read
INSERT OR IGNORE INTO roles_permissions_relation (role_id, permission_id) VALUES
(3, 1), (3, 2), (3, 4), (3, 7), (3, 8), (3, 9);

-- Map API Routes to Required Permissions
INSERT OR IGNORE INTO api_routes_permissions_relation (api_route_id, permission_id) VALUES
(2, 2),  -- GetTicketByID -> tickets:read_own
(2, 3),  -- GetTicketByID -> tickets:read_all
(5, 3),  -- GetTicketsList -> tickets:read_all
(16, 5), -- GetUserByUsername -> users:read
(17, 5), -- GetUserByID -> users:read
(18, 5); -- GetUsersByIDs -> users:read
