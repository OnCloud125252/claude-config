# Claude Code Configuration

My personal Claude Code configuration, custom commands, and tools for enhanced AI-assisted development.

## 🎯 Features

- **Custom Status Line**: Beautiful Go-based status line showing model, project, git branch, context usage, and session time
- **Smart Commit Messages**: Automated Conventional Commits generator with clipboard integration
- **System-wide Instructions**: Centralized CLAUDE.md configuration with project override support

## 📁 Directory Structure

```
.claude/
├── commands/          # Custom slash commands
│   └── commit-message.md
├── config/            # Configuration files
│   └── CLAUDE.md      # System-wide instructions
├── tools/             # Custom tools and utilities
│   └── statusline.go  # Custom status line implementation
├── .gitignore         # Git ignore file
├── README.md          # This file
└── settings.json      # Claude Code settings
```

## 🚀 Quick Start

### 1. Clone this repository

```bash
# Backup your existing .claude directory first
mv ~/.claude ~/.claude.backup

# Clone this configuration
git clone https://github.com/OnCloud125252/claude-config.git ~/.claude
```

### 2. Set up the symlink for CLAUDE.md

```bash
ln -s ~/.claude/config/CLAUDE.md ~/CLAUDE.md
```

### 3. Verify the setup

```bash
# Check the symlink
ls -la ~/CLAUDE.md
```

## ⚙️ Configuration Files

### CLAUDE.md
System-wide instructions that Claude Code loads automatically. Includes:
- Code style preferences (indentation, naming conventions)
- Git commit conventions
- Testing requirements
- Security guidelines
- Tool usage preferences (trash vs rm, ripgrep vs grep)

### settings.json
Main Claude Code settings:
- Status line configuration
- Permission settings
- Output style preferences

### Custom Commands
Located in `commands/` directory:
- **commit-message.md**: Generates Conventional Commits with intelligent type detection and clipboard integration

## 🛠️ Custom Tools

### Status Line (statusline.go)
A custom Go implementation that displays:
- Current AI model (with color coding)
- Project name and git branch
- Context usage with progress bar
- Session time tracking
- User message preview

**Features:**
- Concurrent data fetching for performance
- Smart caching for git operations
- Session tracking with automatic heartbeat
- Beautiful ANSI color formatting

## 📝 Customization

### Adding Your Own Commands

1. Create a new markdown file in `commands/`:
```bash
touch ~/.claude/commands/my-command.md
```

2. Add frontmatter and instructions:
```markdown
---
allowed-tools: Bash, Read, Write
description: My custom command description
---

## Instructions
Your command logic here...
```

### Modifying CLAUDE.md

Edit `~/.claude/config/CLAUDE.md` to add your preferences:
```bash
claude edit ~/.claude/config/CLAUDE.md
```

### Project-Specific Overrides

Create a `CLAUDE.md` in any project directory to override global settings:
```bash
echo "# Project-specific instructions" > ./CLAUDE.md
```

The loading order is:
1. Global: `~/CLAUDE.md` (symlink to `~/.claude/config/CLAUDE.md`)
2. Project: `./CLAUDE.md` (if exists in current directory)

Project instructions override global ones when there's a conflict.

## 🔧 Troubleshooting

### If Claude Code doesn't recognize CLAUDE.md

1. Check the symlink exists:
   ```bash
   ls -la ~/CLAUDE.md
   ```

2. Verify it points to the correct location:
   ```bash
   readlink ~/CLAUDE.md
   ```

3. Ensure the target file exists:
   ```bash
   test -f ~/.claude/config/CLAUDE.md && echo "File exists" || echo "File missing"
   ```

### To recreate the symlink

```bash
# Remove broken symlink
trash ~/CLAUDE.md  # or rm ~/CLAUDE.md

# Create new symlink
ln -s ~/.claude/config/CLAUDE.md ~/CLAUDE.md
```

### Windows Compatibility

Windows users may need to use `mklink` or junction points instead of symlinks:
```cmd
# Run as Administrator
mklink %USERPROFILE%\CLAUDE.md %USERPROFILE%\.claude\config\CLAUDE.md
```

## 🔒 Security Notes

- Session data and tracking information are excluded from git (see `.gitignore`)
- Never commit sensitive information like API keys or passwords
- The configuration uses environment variables for sensitive data
- Files are deleted using `trash` instead of `rm` for recovery options

## 📦 Requirements

- Claude Code CLI installed
- Go runtime (for statusline.go)
- macOS/Linux (for symlinks and trash command)
- Git (for version control)

## 💾 Sync

### Syncing with Git
```bash
cd ~/.claude
git add .
git commit -m "Update Claude Code configuration"
git push
```

### Restoring on a New Machine
```bash
# Clone your configuration
git clone https://github.com/OnCloud125252/claude-config.git ~/.claude

# Set up the symlink
ln -s ~/.claude/config/CLAUDE.md ~/CLAUDE.md

# Verify setup
cat ~/CLAUDE.md
```

---

*Note: Remember to update paths and usernames when setting up on a new machine.*