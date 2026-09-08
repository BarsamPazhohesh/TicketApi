---
name: agents-flow
description: Runs sequential multi-agent review pipeline: api-doc-reviewer, db-reviewer, security-reviewer, refactorer, and verifier.
---

# Sequential Go Review Pipeline

Execute the following agents in exact sequence:

1. **`api-doc-reviewer`**: Audit Swagger annotations, DTO validator tags, HTTP status codes, and REST API conventions.
2. **`db-reviewer`**: Audit SQLite schema changes (`sqlc`), migrations (`cmd/migrate`), indexes, and MongoDB collections/queries.
3. **`security-reviewer`**: Audit Gin routes, authentication guards, JWT/cookie handling, input validation, and password security.
4. **`refactorer`**: Review code against clean architecture (handler -> service -> repo), error handling with `errx`, DRY, and idiomatic Go.
5. **`verifier`**: Run Go static analysis (`go vet ./...`), unit tests (`go test -v ./...`), compile checks (`go build -o /dev/null ./cmd/api`), and verify `sqlc`/`swag` sync.

Summarize findings from each agent in a structured report.
