# Claude Code System Instructions

## Project Context
This is my development environment. Follow these preferences and conventions.

## Code Style Preferences
- Use 2 spaces for indentation (not tabs)
- Prefer async/await over promises
- Use meaningful variable names, avoid single letters except for loop indices
- Always use const/let, never var in JavaScript/TypeScript

## Node.js Development
- **Always use `bun` by default** for Node.js projects
  - `bun install` - install dependencies
  - `bun run <script>` - run package.json scripts
  - `bun add <package>` - add dependencies
  - `bun remove <package>` - remove dependencies
  - `bun <file.ts>` - run TypeScript files directly
- Only use npm or yarn if the project explicitly requires it (e.g., existing lock files)

## Git Commit Conventions
- Follow Conventional Commits format (feat, fix, docs, style, refactor, test, chore)
- Keep commit messages under 72 characters
- Write commit messages in imperative mood

## Testing Requirements
- Always run tests after making changes if test scripts exist
- Run linting and type checking before completing tasks
- Check for these common scripts: `npm run test`, `npm run lint`, `npm run typecheck`

## Development Workflow
- Always read existing code before making changes
- Follow existing patterns and conventions in the codebase
- Prefer modifying existing files over creating new ones
- Never create documentation files unless explicitly requested

## File Operations
- **ALWAYS use `trash` instead of `rm`** for file deletion
  - `trash file.txt` - moves to trash (recoverable)
  - `trash -r directory/` - recursively trash directory
  - `trash-list` - view trashed files
  - `trash-restore` - recover deleted files
- Never use `rm -rf` as it permanently deletes files

## Security Guidelines
- Never commit sensitive information (API keys, passwords, tokens)
- Always use environment variables for configuration
- Check .gitignore before adding new files

## Tool Usage Preferences
- Use `trash` instead of `rm` for safer file deletion (files can be recovered from trash)
- Use ripgrep (`rg`) instead of grep for searching
- Prefer using Glob and Grep tools over bash find/grep commands
- Always use absolute paths when working with files

## Project-Specific Notes
<!-- Add your project-specific requirements here -->
<!-- Example: This project uses Prettier for formatting -->
<!-- Example: API endpoints should follow REST conventions -->

## Frequently Used Commands
<!-- Add commands you use often for quick reference -->
<!-- Example: Build: npm run build -->
<!-- Example: Deploy: npm run deploy -->