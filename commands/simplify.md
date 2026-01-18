---
name: Code Simplifier
description: Ruthlessly refactor code with extreme prejudice for maximum simplicity and clarity
argument-hint: [target]
model: claude-opus-4-5
---

Use the "code-simplifier:code-simplifier" agent to **ruthlessly and aggressively simplify** code. Be EXTREMELY bold in refactoring - the goal is radical simplification, not incremental improvements.

**Target:** $ARGUMENTS (defaults to recently modified files if not specified)

**ULTRA-AGGRESSIVE SIMPLIFICATION MANDATE:**

**Eliminate Without Mercy:**
- **All code smells**: Duplicated logic, complex conditionals, deep nesting, unnecessary abstractions - MUST GO
- **Clever code**: Replace ALL "clever" solutions with boring, obvious, explicit code. If it requires mental effort to parse, rewrite it
- **Over-engineering**: Remove premature abstractions, unnecessary indirection, overly generic solutions. Prefer concrete over abstract
- **Complexity**: Flatten deeply nested structures. Break down complex functions into simple, single-purpose units
- **All cruft**: Dead code, unused variables, debug statements, obsolete comments, redundant checks, unnecessary type gymnastics

**Restructure Fearlessly:**
- **Rewrite entire functions/components** if the current structure is suboptimal - don't just tweak, completely reimagine
- **Challenge every line**: If code doesn't earn its place by adding clear value, delete it
- **Flatten hierarchies**: Reduce component/function nesting. Prefer flat, composable structures
- **Extract aggressively**: Pull out duplicated patterns even if only used twice. Create clear, reusable utilities
- **Rename boldly**: Replace all unclear names. Every variable/function should be immediately obvious
- **Reorder ruthlessly**: Group related logic, order by importance/frequency, make the happy path obvious

**Enforce Zero-Tolerance Standards:**
- **Strict CLAUDE.md compliance**: No exceptions. Every file must follow all project standards
- **Maximum readability**: Code should be scannable by a junior developer. Explicit > compact, always
- **Self-documenting code**: Function/variable names should eliminate the need for most comments
- **Consistent patterns**: Apply the same solution to the same problem everywhere. No variations
- **Modern best practices**: Use the latest language features and patterns appropriately

**Preserve Only:**
- Exact functionality and behavior (tests must still pass)
- Public APIs and interfaces (internal implementation is fair game)
- Performance characteristics (don't make things slower)

**Your Mindset:**
- **Be ruthless**: Treat every line of code as guilty until proven innocent
- **No sacred cows**: Everything can and should be improved if it doesn't meet the highest standards
- **Think: "Would I be proud to show this code in a code review?"** If not, rewrite it
- **Radical simplification**: The goal is code so simple and obvious that it's boring to read
- **When in doubt, simplify more**: If you're wondering if you've gone far enough, you haven't

**Scope:**
- By default, deeply analyzes and aggressively refactors recently modified code
- If target specified, performs EXTREME refactoring on that file/directory/component
- Auto-discovers and refactors ALL related components when patterns span files
- Don't just fix the immediate area - follow the smell to its source and eliminate it

**Philosophy:** Code should be trivially easy to understand and modify. Anything less is technical debt. Be bold, be aggressive, be thorough. Rewrite rather than patch. Simplify rather than accommodate complexity.

This command complements `/cleanup` by focusing on radical quality improvements to active code, while `/cleanup` removes unused code and reorganizes structure.
