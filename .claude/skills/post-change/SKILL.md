---
name: post-change
description: Runs verification pipeline (go vet, test, build) and targeted Go reviews after changes.
---

Post-change pipeline:
1. Run verifier: `go vet ./... && go test ./... && go build -o /dev/null ./cmd/api`.
2. Inspect changed files via `git status -s`.
3. If routes, auth, middleware, or token handling changed, invoke `security-reviewer`.
4. If handlers, DTOs, or Swagger docs changed, invoke `api-doc-reviewer`.
5. If DB schema, sqlc queries, migrations, or MongoDB models changed, invoke `db-reviewer`.
6. Summarize verification and review results concisely.
