---
name: Project Cleaner
description: Deep cleanup - enforces best practices, removes unused code, formats, and cleans
model: claude-opus-4-5
---

Analyze the project directory to identify and remove unused code:

1. First, scan all files and create an inventory of:
   - All exported functions, components, and hooks
   - All imports across files
   - Entry points (main files, routes, etc.)

2. Identify unused items:
   - Functions/components/hooks that are defined but never imported elsewhere
   - Files with no imports from other files (excluding entry points)
   - Commented-out code blocks
   - Duplicate implementations

3. Reorganize files:
   - Move related code (functions/components/hooks etc.) into dedicated files or folders
   - Ensure a logical structure that follows best practices for the specific framework/language
   - Rename files for clarity and consistency

4. Run any formatters or linters to ensure code quality:
   - Use tools like Biome, ESLint, Prettier, or similar for project's language/framework.

Important: Don't remove:
- Code that's called dynamically (check for string-based imports)
