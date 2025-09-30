# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

This is a personal Claude Code configuration repository that provides custom tools, slash commands, and system-wide settings. The `config/CLAUDE.md` file is symlinked to `~/CLAUDE.md` to provide system-wide defaults, while project-specific CLAUDE.md files override these settings.

## Architecture

### Status Line System (tools/statusline.go)

A custom Go tool that generates the CLI status line with real-time information:

**Core Components:**
- **Concurrent data fetching**: Uses goroutines and channels to fetch git branch, session time, context usage, and user message in parallel
- **Session tracking**: Maintains session files in `session-tracker/sessions/` with heartbeat updates and interval tracking
- **Git branch caching**: 5-second cache to avoid excessive git calls
- **Context analysis**: Parses transcript JSONL to extract token usage from API responses, filtering side-chain messages

**Key Functions:**
- `main()`: Orchestrates parallel data fetching via goroutines
- `getGitBranch()`: Cached git branch detection with mutex-protected cache
- `updateSession()`: Heartbeat-based session tracking with 10-minute timeout
- `calculateTotalHours()`: Aggregates all today's session data across multiple sessions
- `analyzeContext()`: Parses transcript JSONL to find latest non-sidechain usage data
- `extractUserMessage()`: Finds most recent user message for current session

**Data Flow:**
1. Claude Code passes JSON input via stdin containing model, session ID, workspace, and transcript path
2. Four goroutines fetch data concurrently (git, hours, context, message)
3. Results collected via channel and formatted with ANSI colors
4. Session file updated synchronously to prevent race conditions

### Slash Commands (commands/)

Commands are defined using markdown files with YAML frontmatter:

**Frontmatter Configuration:**
- `allowed-tools`: Restricts tool access (e.g., `Bash(git diff --cached)`)
- `disable-model-invocation`: Prevents streaming responses for clipboard operations
- `description`: Short description shown in command list
- `model`: Specify which Claude model to use

**Example: commit-message.md**
- Analyzes staged changes with `git diff --cached`
- Auto-detects Conventional Commit type based on file changes
- Derives scope from recent commits and directory structure
- Copies message to clipboard via `pbcopy`
- Outputs colored confirmation with ANSI escape codes

### Configuration Files

**settings.json:**
- `permissions.additionalDirectories`: Grants access to home directory (`~/`)
- `statusLine.type`: Set to "command" to run custom tool
- `statusLine.command`: Path to Go-based status line tool
- `outputStyle`: Set to "Explanatory" for educational responses
- `alwaysThinkingEnabled`: Enables detailed reasoning

**config/CLAUDE.md:**
- System-wide instructions symlinked to `~/CLAUDE.md`
- Defines code style, git conventions, testing requirements
- Security guidelines and tool preferences
- Project-specific CLAUDE.md files in working directories override these settings

## Development Workflow

### Testing Status Line Changes

1. Make changes to `tools/statusline.go`
2. Test directly:
   ```bash
   echo '{"model":{"display_name":"Sonnet 4"},"session_id":"test","workspace":{"current_dir":"'$(pwd)'"}}' | go run tools/statusline.go
   ```
3. Test with actual Claude Code session by starting a new conversation

### Adding New Slash Commands

1. Create `commands/your-command.md`
2. Add frontmatter with required tools and description
3. Write instructions in markdown body
4. Test by running `/your-command` in Claude Code

### Modifying System-Wide Configuration

Edit `config/CLAUDE.md` (not `~/CLAUDE.md` which is a symlink). Changes apply to all projects unless overridden by project-specific CLAUDE.md.

## Key Implementation Details

### Session Tracking Logic

Sessions use interval-based tracking:
- Each heartbeat (status line render) updates `LastHeartbeat`
- If gap < 10 minutes: extend current interval
- If gap ≥ 10 minutes: create new interval
- `TotalSeconds` calculated by summing all interval durations

### Context Token Calculation

Token usage extracted from transcript JSONL:
- Reads last 100 lines from transcript file
- Parses each line as JSON, skips if `isSidechain: true`
- Finds most recent message with `message.usage` object
- Sums `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`
- Returns 0 if no usage data found (start of conversation)

### Git Branch Caching Strategy

Cache prevents excessive git calls during rapid status line updates:
- Cache stored in global variables with mutex protection
- 5-second TTL balances freshness vs performance
- Returns empty string if not in git repository

## File Structure

```
.claude/
├── commands/              # Slash commands (markdown with frontmatter)
├── config/
│   └── CLAUDE.md          # System-wide instructions (symlinked to ~/CLAUDE.md)
├── tools/
│   └── statusline.go      # Custom status line implementation
├── session-tracker/       # Session data (gitignored)
│   └── sessions/*.json    # Per-session tracking files
├── settings.json          # Claude Code configuration
└── .gitignore            # Excludes session data, todos, statsig cache
```

## Common Tasks

**Run status line locally:**
```bash
go run ~/.claude/tools/statusline.go
```

**Check session data:**
```bash
cat ~/.claude/session-tracker/sessions/*.json | jq
```

**Test commit message generation:**
```bash
# Stage some changes first
git add .
# Then run the command
/commit-message
```

**Update system-wide config:**
```bash
claude edit ~/.claude/config/CLAUDE.md
```

**Sync configuration to git:**
```bash
cd ~/.claude
git add -A
git commit -m "Update configuration"
git push
```
