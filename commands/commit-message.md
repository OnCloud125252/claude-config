---
name: Commit Message Creator
allowed-tools: Bash(git diff --cached), Bash(git log -n 10), Bash(git status --porcelain), Bash(pbcopy), Bash(echo:*)
disable-model-invocation: true
description: Create a Conventional Commits style message with intelligent type detection and copy to clipboard
model: claude-haiku-4-5
---

Analyze staged changes and create a Conventional Commits compliant message.

### Process:
1. Run `git status --porcelain` to verify there are staged changes
2. Run `git diff --cached` to analyze the staged changes
3. Run `git log -n 10` to match the project's existing commit style and identify common scopes
4. Copy the generated commit message to clipboard using `pbcopy`
5. Display the commit message and confirm it's copied to clipboard

### Message Generation Rules:

**If arguments provided:**
- Parse "$ARGUMENTS" intelligently:
  - If starts with conventional type (feat/fix/docs/etc), use as-is
  - Otherwise, generate appropriate type based on changes and use argument as description

**If no arguments:**
- Auto-detect the most appropriate type based on file changes:
  - `feat`: New files, new functions/methods, new features
  - `fix`: Changes to existing logic that resolve issues
  - `docs`: Changes to .md, .txt, comments, or documentation files
  - `style`: Formatting, whitespace, semicolons (no logic changes)
  - `refactor`: Code restructuring without changing functionality
  - `test`: Adding or modifying test files
  - `chore`: Build scripts, configs, dependencies, tool changes
  - `perf`: Performance improvements
  - `ci`: CI/CD configuration changes

**Scope detection:**
- Derive from common directory or module name
- Use existing scopes from recent commits when applicable
- Omit if changes span multiple areas

**Description guidelines:**
- Start with lowercase verb (add, update, fix, remove, etc.)
- Keep under 50 characters
- Be specific but concise
- No period at the end

**Body guidelines (for complex changes):**
- Include if changes affect 3+ files or have breaking changes
- Explain the what and why, not the how
- Wrap at 72 characters
- Separate from subject with blank line

**Output format:**
- Generate the commit message and copy it to clipboard using `echo "message" | pbcopy`
- Display the commit message to the user
- Add a colored confirmation note using ANSI escape codes:
  ```
  echo -e "\033[32m✓ Copied to clipboard\033[0m - paste with \033[36mCmd+V\033[0m in VSCode"
  ```
- No markdown formatting in the commit message itself
- Include body/footer only when necessary
- NEVER include "Generated with", "Co-authored-by", or any AI/tool attribution in the commit message
- No meta-information about how the message was created in the commit message

### Examples of expected output:

For a bug fix in authentication:
```
fix(auth): resolve token expiration check

Token validation was using local time instead of UTC,
causing premature session timeouts for users in different
timezones.
```

For a new feature:
```
feat(api): add pagination support for user endpoints
```

For documentation update:
```
docs: update API authentication examples
```

For breaking change:
```
feat(database)!: migrate to PostgreSQL from MySQL

BREAKING CHANGE: Database connection strings must be updated
to PostgreSQL format. See migration guide in docs/MIGRATION.md
```
