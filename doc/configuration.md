# Configuration Reference

TicketApi uses a YAML file (`config.yaml`) for application settings with dynamic hot-reloading via `fsnotify`, supplemented by environment variables for credentials and environment overrides.

---

## 1. Overview & Architecture

- **Primary Configuration File**: `config.yaml` located at the root of the project.
- **Loader**: `internal/config/config.go` loads the YAML file and starts an `fsnotify` watcher on a background goroutine. When `config.yaml` changes on disk, settings are automatically parsed and safely updated in-memory via `sync.RWMutex`.
- **Environment Variables**: Loaded via `internal/env/env.go` (`os.LookupEnv` and `os.Getenv`) to provide secrets and runtime flags without checking credentials into version control.

---

## 2. `config.yaml` Reference

### `app`
Core HTTP server and request body constraints.

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `port` | integer | `8080` | HTTP port on which the Gin server binds and listens (`cmd/api/main.go`, `cmd/api/server.go`). |
| `max_json_request_size` | integer (int64) | `1024` | Maximum allowable JSON request body size in **KB** (1024 KB = 1 MB). Enforced by `middleware.LimitRequestBody` on route groups (`cmd/api/routes.go`). |
| `max_upload_files_size` | integer (int64) | `8196` | Max multipart memory allocated for incoming file uploads in **KB** (`cmd/api/routes.go` sets `g.MaxMultipartMemory`). |

---

### `cors`
Cross-Origin Resource Sharing rules applied to all routes via Gin CORS middleware (`cmd/api/routes.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `allow_origins` | list of strings | `["http://localhost:4200", ...]` | Allowed client origin URLs for browser cross-origin requests. |
| `allow_methods` | list of strings | `["GET", "POST", "PUT", "DELETE", "OPTIONS"]` | Allowed HTTP methods. |
| `allow_headers` | list of strings | `["Content-Type", "Authorization", "x-api-key"]` | Allowed request headers sent by clients. |
| `expose_headers` | list of strings | `["Content-Length"]` | Response headers accessible to the client script. |
| `allow_credentials` | boolean | `true` | Allows cookies and Authorization headers on cross-origin requests. |
| `max_age_hours` | integer | `12` | Cache duration in hours for preflight `OPTIONS` responses. |

---

### `captcha`
Image and JWT-backed CAPTCHA verification configuration (`internal/services/captcha`, `internal/middleware/captcha_middleware.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `length` | integer | `6` | Number of digits generated in the CAPTCHA image challenge. |
| `timeout_minutes` | integer | `5` | In-memory cache TTL (minutes) for a generated CAPTCHA code. |
| `expired_time_token` | integer | `15` | Expiration time (minutes) embedded inside the signed CAPTCHA JWT token. |
| `cookie_name` | string | `"captcha_token"` | Cookie name used to send and receive the CAPTCHA token. |
| `cleanup_interval` | integer | `10` | Frequency (minutes) of background cache eviction for expired captchas. |
| `max_cashed_captcha` | integer | `100000` | Max number of active CAPTCHA challenges stored in memory before rejecting/evicting. |
| `image_width` | integer | `240` | Width of the rendered CAPTCHA PNG image in pixels. |
| `image_height` | integer | `80` | Height of the rendered CAPTCHA PNG image in pixels. |
| `validate_ip` | boolean | `true` | When `true`, ensures requester IP matches the IP signed in the CAPTCHA token (`internal/middleware/captcha_middleware.go`). |

---

### `token`
Global cookie attributes for security tokens (`internal/services/cookie/cookie_service.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `httponly` | boolean | `true` | Sets `HttpOnly` flag on issued cookies (prevents client JS XSS access). |
| `secure` | boolean | `false` | Sets `Secure` flag on cookies (must be set to `true` when serving over HTTPS). |

---

### `auth`
User authentication JWT token attributes (`internal/services/token/auth_token.go`, `internal/services/cookie/cookie_service.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `expired_time_token` | integer | `60` | Lifetime in minutes for user login session JWT tokens. |
| `cookie_name` | string | `"auth_token"` | Cookie name used to persist the authenticated session token. |

---

### `one_time_token`
Temporary single-use token lifecycle (`internal/services/token/one_time_token.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `cleanup_interval` | integer | `10` | Background eviction interval (minutes) for expired one-time tokens in cache. |
| `max_cashed_tokens` | integer | `100000` | Maximum capacity for in-memory one-time tokens cache. |
| `expired_time_token` | integer | `3` | Validity duration in minutes for one-time tokens. |

---

### `api_key`
API Key length and validation settings (`internal/services/token/api_key.go`, `internal/middleware/api_key_guard.middleware.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `size` | integer | `32` | Minimum byte length / entropy size required for API keys. |

---

### `mongo`
MongoDB persistence for tickets and chat messages (`internal/repository/ticket_repository.go`, `internal/repository/chat_repository.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `enable` | boolean | `true` | Toggles MongoDB connection and usage. If `false`, ticket and chat operations return uninitialized/disabled state. |
| `db_name` | string | `"ticket_db"` | Default MongoDB database name (can be overridden by `MONGODB_DB` env var). |
| `ticket_collocation_name` | string | `"tickets"` | MongoDB collection name for storing tickets and messages. |

---

### `redis`
Redis cache connection settings (`cmd/api/main.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `enable` | boolean | `true` | Toggles Redis connection. If `false`, caching is bypassed. |
| `host` | string | `"localhost"` | Hostname or IP of the Redis server. |
| `port` | integer | `6379` | Port for Redis connection. |
| `db` | integer | `0` | Logical Redis database index (`0`–`15`). |

---

### `cache`
TTL configuration for cached metadata in Redis (`internal/repository/`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `ticket_type_ttl_minutes` | integer (int64) | `1440` | Cache TTL in minutes for ticket types list (24 hours). |
| `department_ttl_minutes` | integer (int64) | `1440` | Cache TTL in minutes for departments list (24 hours). |
| `ticket_status_ttl_minutes` | integer (int64) | `1440` | Cache TTL in minutes for ticket statuses (24 hours). |

---

### `minio`
Object storage configuration for file attachments (`internal/services/storage/storage.go`, `cmd/api/main.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `enable` | boolean | `true` | Toggles MinIO object storage. |
| `host` | string | `"localhost"` | MinIO endpoint host. |
| `port` | integer | `9000` | MinIO endpoint port. |
| `bucket` | string | `"ticket-files"` | MinIO bucket name for uploading and reading ticket attachment files. |
| `use_ssl` | boolean | `false` | Enable/disable SSL for MinIO client connections. |
| `public_url` | string | `""` | Public base URL used for presigned download links (e.g. `"https://irmto.ir/file-storage"`). If empty, standard direct MinIO URL is used. |

---

### `ticket`
Ticket pagination limits and file upload constraints (`internal/repository/ticket_repository.go`, `internal/handler/file_handler.go`, `internal/handler/ticket_handler.go`, `internal/handler/chat_handlers.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `max_paging_size` | integer | `100` | Maximum limit clients can request per page. Values greater than this revert to `default_paging_size`. |
| `min_paging_size` | integer | `2` | Minimum limit allowed per page. Values smaller than this revert to `default_paging_size`. |
| `default_paging_size` | integer | `10` | Default items per page returned when limit is omitted or outside bounds. |
| `max_counting_item` | integer (int64) | `10000` | Maximum item threshold when calculating total item counts to prevent full collection scans. |
| `max_ticket_upload_file` | integer | `5` | Maximum number of attached files allowed per ticket. |
| `max_ticket_upload_file_size` | integer (int64) | `1024` | Maximum allowable file size per attachment in **KB** (1024 KB = 1 MB). |
| `acceptable_files_for_upload` | list of strings | `[".jpg", ".jpeg", ".png"]` | Whitelisted file extensions permitted for ticket attachment uploads. |

---

### `otp`
SMS OTP verification settings (`internal/services/otp/otp_service.go`).

| Field | Type | Default in `config.yaml` | Effects & Scope |
| :--- | :--- | :--- | :--- |
| `code_ttl_minutes` | integer | `2` | TTL in minutes for generated OTP codes. |
| `retry_interval_minutes` | integer | `5` | Periodic worker retry interval for pending/failed SMS records in `sms_warehouse`. |
| `max_verify_attempts` | integer | `5` | Maximum failed OTP verification attempts before lockout and code eviction. |

---

## 3. Environment Variables Reference

Environment variables take precedence for secrets and runtime modes:

| Variable | Default Fallback | Usage in Code | Description & Effects |
| :--- | :--- | :--- | :--- |
| `GIN_MODE` | `"debug"` | `internal/errx/errx.go`, `internal/services/captcha/captcha_service.go` | When set to `"debug"`, error responses include internal traces and CAPTCHA digits are logged to stdout. Set to `"release"` in production. |
| `JWT_SECRET` | `""` | `internal/services/token/secret.go` | Secret key used for HMAC-SHA256 signing and validation of all JWT tokens (Auth, Captcha, OneTimeToken). **Required for auth security.** |
| `MONGODB_URI` | `""` | `cmd/api/main.go` | MongoDB connection URI (e.g. `mongodb://user:pass@localhost:27017`). Required if `mongo.enable: true`. |
| `MONGODB_DB` | `config.Get().Mongo.DBName` (`"ticket_db"`) | `cmd/api/main.go` | Overrides the target database name in MongoDB. |
| `REDIS_PASSWORD` | `""` | `cmd/api/main.go` | Password for Redis server authentication. |
| `ACCESS_KEY_MINIO` | `""` | `cmd/api/main.go` | MinIO / S3 access key (username). |
| `SECRET_KEY_MINIO` | `""` | `cmd/api/main.go` | MinIO / S3 secret key (password). |
| `DATABASE_URL` | N/A (SQLite hardcoded in code) | `cmd/migrate/main.go` | Used when executing database schema migrations (`cmd/migrate/main.go` defaults to `./data.db`). |
