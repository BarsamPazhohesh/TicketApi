# Ticket API

A simple Ticket API built with **Go**, **Gin**, and **sqlc**,
Includes Swagger documentation for API endpoints.

---

## Features

- Versioned API (`/api/v1`, `/api/v2`, etc.)
- SQL queries managed with **sqlc** (type-safe Go models)
- Database migrations supported
- Swagger documentation with **swaggo/gin-swagger**

---

## Prerequisites

- [Go](https://go.dev/doc/install) >= 1.21
- [Docker & Docker Compose](https://docs.docker.com/compose/install/) (optional, for local DB)
- [sqlc](https://sqlc.dev/) for generating type-safe Go queries
- [swag](https://github.com/swaggo/swag) for generating Swagger docs
- [migrate](https://github.com/golang-migrate/migrate) for DB migrations
- [Air](https://github.com/cosmtrek/air) for live-reload during development

---

## Installation

1. **Clone the repository**

```bash
git clone git@github.com:BarsamPazhohesh/TicketApi.git
cd TicketApi
```

````

2. **Install dependencies**

```bash
go mod tidy
```

3. **Install development tools**
```bash
# Air
go install github.com/air-verse/air@latest
```


```bash
# sqlc for type-safe SQL queries
brew install sqlc

# swag for Swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest

# migrate for DB migrations
brew install golang-migrate
```

4. **Generate SQL code and Swagger docs**
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

```bash
# Generate Go code from SQL queries
sqlc generate

# Generate Swagger documentation
swag init
```



---

## Database Setup

- **Using Docker (PostgreSQL example)**

```bash
docker-compose up -d
```

- **Run migrations**

```bash
go run ./cmd/mingrate/main.go up
```

---

## Run API

```bash
go run main.go
```

- API endpoints are available at: `http://localhost:8080/api/v1`
- Swagger UI is available at: `http://localhost:8080/swagger/index.html`

---

## Example Endpoints

- `GET /api/v1` – Get current version of API v1
- `POST /api/v1/version` – Create a new version

All responses use JSON camelCase fields:

```json
{
  "apiVersion": "v1",
  "version": "1.0.0",
  "notes": "Initial release",
  "isCurrent": true
}
```

---

## Development Notes

- Repository pattern is used: all DB logic is in `internal/repository`.
- DTOs in `internal/dto` handle JSON mapping and conversions.
- Handlers in `internal/handler` include Swagger annotations.

---

## Recommended Tools for macOS

- **Postgres / SQLite** → TablePlus (GUI) + pgcli (CLI)
- **MongoDB** → MongoDB Compass + mongosh
- **SQLite** → DB Browser for SQLite or TablePlus

---

## License

MIT License

```
[MIT License](./LICENSE)
```
````
