---
name: refactor
description: Invokes refactorer agent on specified Go packages or modified diffs.
---

1. Identify modified files or target Go packages (`internal/handler`, `internal/repository`, `internal/services`, etc.).
2. Dispatch `refactorer` agent to evaluate clean architecture, DTO conversions, error handling with `errx`, DRY, and idiomatic Go practices.
3. Report concrete refactoring proposals with minimal diffs.
