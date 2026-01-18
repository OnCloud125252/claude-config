---
name: Fully Optimize Code
description: Comprehensive code optimization - removes unused code, eliminates duplication, enforces best practices, and reorganizes for maximum clarity
model: claude-opus-4-5
---

Analyze the project directory to identify and remove unused code:

1. First, scan all files and create an organized inventory of:
   - All exported functions, components, and hooks
   - All imports across files
   - Entry points (main files, routes, etc.)
   - Commented-out code blocks
   - Duplicate implementations
   
   Group the inventory by:
   - File/module location
   - Type (function, component, hook, utility, etc.)
   - Usage status (used, unused, potentially dynamic)

2. Identify unused items:
   - Functions/components/hooks that are defined but never imported elsewhere
   - Files with no imports from other files (excluding entry points)
   - Dead code that can be safely removed

3. For each file or module in the inventory, call the code-simplifier agent:
   - Use the Skill tool with skill: "code-simplifier:code-simplifier"
   - Instruct it to ruthlessly and aggressively simplify the code
   - Let it remove unused code, eliminate duplication, enforce best practices
   - The code-simplifier will preserve all functionality while maximizing simplicity

4. Reorganize files:
   - Move related code (functions/components/hooks etc.) into dedicated files or folders
   - Ensure a logical structure that follows best practices for the specific framework/language
   - Rename files for clarity and consistency

5. Run any formatters or linters to ensure code quality:
   - Use tools like Biome, ESLint, Prettier, or similar for project's language/framework

Important: Don't remove:
- Code that's called dynamically (check for string-based imports)
- Entry points and configuration files
