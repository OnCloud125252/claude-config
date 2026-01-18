# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

Personal Claude Code configuration repository providing custom slash commands, skills, plugins, and system-wide settings. The `config/CLAUDE.md` file is symlinked to `~/CLAUDE.md` for system-wide defaults; project-specific CLAUDE.md files override these.

**Official Docs:** [code.claude.com/docs](https://code.claude.com/docs)
**Custom Slash Command Docs** [code.claude.com/docs/slash-commands#custom-slash-commands](https://code.claude.com/docs/en/slash-commands#custom-slash-commands)

## Architecture

### Status Line

Powered by `ccstatusline` npm package. Shows git branch, session time, context usage, and current message preview.

### Plugin System (v2)

Plugins extend functionality with specialized agents, skills, and tools.

| Plugin | Marketplace | Description |
|--------|-------------|-------------|
| `frontend-design` | claude-code-plugins | Production-grade frontend interface generation |
| `bun` | mag-claude-plugins | Bun runtime integration |
| `frontend` | mag-claude-plugins | Frontend development tools |
| `code-analysis` | mag-claude-plugins | Codebase investigation and analysis |
| `gopls-lsp` | claude-plugins-official | Go language server integration |
| `clangd-lsp` | claude-plugins-official | C/C++ language server integration |
| `code-simplifier` | claude-plugins-official | Code simplification and refactoring |

**Marketplaces:**
- `claude-code-plugins` - anthropics/claude-code (official)
- `claude-plugins-official` - anthropics/claude-plugins-official (official)
- `mag-claude-plugins` - MadAppGang/claude-code (community)

### Skills

Self-contained automation packages in `skills/`.

**playwright-skill** - Browser automation with Playwright. Auto-detects dev servers, writes scripts to `/tmp`, visible browser by default.

### Slash Commands

| Command | Description | Model |
|---------|-------------|-------|
| `/commit-message` | Generate conventional commit message from staged changes | default |
| `/cleanup` | Deep cleanup - enforces best practices, removes unused code | claude-opus-4-5 |
| `/simplify` | Code simplification using code-simplifier agent | claude-opus-4-5 |

## Development Workflow

### Adding Slash Commands

1. Create `commands/your-command.md` with YAML frontmatter
2. Frontmatter options: `name`, `description`, `model`, `allowed-tools`, `disable-model-invocation`
3. Test with `/your-command`

### Managing Plugins

Toggle in `settings.json`:
```json
"enabledPlugins": {
  "plugin-name@marketplace": true
}
```

### Modifying System-Wide Config

Edit `config/CLAUDE.md` (not `~/CLAUDE.md` which is a symlink).

## File Structure

```
.claude/
├── commands/              # Slash commands (markdown with frontmatter)
├── config/CLAUDE.md       # System-wide instructions (symlinked to ~/CLAUDE.md)
├── plugins/               # Plugin system
│   ├── installed_plugins.json
│   ├── known_marketplaces.json
│   ├── cache/
│   └── marketplaces/
├── skills/                # Automation packages
│   └── playwright-skill/
├── settings.json          # Configuration
└── settings.local.json    # Local overrides (gitignored)
```
