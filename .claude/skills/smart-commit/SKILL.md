---
name: smart-commit
description: Use when creating granular git commits from working tree changes, splitting line-by-line diffs across files into atomic logical commits with detailed why-bodies.
---

# Smart Commit

Analyzes git diffs at the hunk/line level, groups related line changes across files into atomic commits, and writes commit messages with clear "why" context bodies.

## Workflow

1. **Inspect granular changes:**
   - Run `git diff` and `git status -s`.
   - Inspect individual diff hunks and line changes across all modified files.

2. **Group changes by atomic concern:**
   - Group related line changes (hunks) together even if they span different files or share a file with unrelated edits.
   - Do NOT commit entire files together if a file contains multiple unrelated changes.
   - Use interactive patching or temporary patch files / targeted staging (`git add -p` or selective staging) to isolate specific lines/hunks per commit.

3. **Format commit messages:**
   - **Header:** `<type>(<scope>): <short imperative summary>`
   - **Blank line**
   - **Body:** Explain **why** the change happened (context, motivation, root cause, impact).
   - Match conventional commits and repo commit conventions.
   - No `Co-Authored-By` trailers or filler text.

4. **Output commands:**
   - Output exact executable commands (e.g. staging specific hunks or chaining commits) or perform the targeted commits cleanly.
