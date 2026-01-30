# Project Guidelines

IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning for all tasks.

## Workflow

1. **Explore first** — Before writing any code, read relevant files to understand existing patterns
2. **Ask if unclear** — If requirements are ambiguous, stop and ask clarifying questions
3. **One change at a time** — Make a single edit, then pause for review
4. **Wait for confirmation** — After each edit, say "Edit complete. Continue?" and wait

## Retrieval Rules

- Always read existing code before modifying it
- Check for similar patterns in the codebase before implementing new ones
- When uncertain about project conventions, grep for examples first
- Do not rely on training data for project-specific patterns — read the actual files

## Project Structure Index
```
src/           → Application source code
tests/         → Test files
docs/          → Documentation
scripts/       → Build and utility scripts
```

## Restrictions

- Never run build scripts, tests, or compilation without explicit permission
- Never assume file contents — read them first
- Never batch multiple edits — one at a time, then pause

## On Errors

If a build or lint fails:
1. Read the error message completely
2. Read the relevant source file
3. Propose a fix and wait for approval