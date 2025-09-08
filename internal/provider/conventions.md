# Professional Git Commit Message Conventions

This file contains professional guidelines for writing clear, consistent, and informative commit messages based on industry best practices and established conventions.

## Structure

Follow the conventional commits format with proper separation:

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## The Seven Golden Rules

1. **Separate subject from body with a blank line**
2. **Limit the subject line to 50 characters**
3. **Capitalize the subject line**
4. **Do not end the subject line with a period**
5. **Use the imperative mood in the subject line**
6. **Wrap the body at 72 characters**
7. **Use the body to explain what and why vs. how**

## Commit Types

### Primary Types

- **feat**: A new feature for the user (correlates with MINOR in SemVer)
- **fix**: A bug fix for the user (correlates with PATCH in SemVer)

### Supporting Types

- **docs**: Documentation changes only
- **style**: Code style changes (formatting, missing semicolons, whitespace, etc.)
- **refactor**: Code refactoring without changing external behavior
- **test**: Adding missing tests or correcting existing tests
- **chore**: Regular maintenance tasks (updating dependencies, build tools, etc.)
- **perf**: Performance improvements
- **ci**: Changes to CI/CD configuration files and scripts
- **build**: Changes to build system or external dependencies
- **revert**: Reverting a previous commit

### Optional Scope

Provide additional context in parentheses:

- `feat(auth): add OAuth2 integration`
- `fix(parser): handle edge case in CSV processing`
- `docs(api): update authentication examples`

## Writing Guidelines

### Subject Line (≤50 characters)

- **Use imperative mood**: "Add feature" not "Added feature" or "Adds feature"
- **Capitalize first letter**: "Fix bug" not "fix bug"
- **No period at end**: "Update documentation" not "Update documentation."
- **Be specific and concise**: Focus on the essential change
- **Complete this sentence**: "If applied, this commit will ___"

### Body (72 characters per line)

- **Explain what and why**, not how (code explains how)
- **Provide context**: Why was this change necessary?
- **Describe the problem** that this commit solves
- **Mention side effects** or unintuitive consequences
- **Use proper grammar** and punctuation
- **Separate paragraphs** with blank lines
- **Use bullet points** when appropriate

### Footer

- **Issue references**: `Fixes #123`, `Closes #456`, `Refs #789`
- **Breaking changes**: `BREAKING CHANGE: description`
- **Co-authors**: `Co-authored-by: Name <email@example.com>`
- **Reviewers**: `Reviewed-by: Name <email@example.com>`

## Breaking Changes

Indicate breaking changes in two ways:

### Option 1: Exclamation mark

```text
feat!: remove deprecated API v1 endpoints
```

### Option 2: Footer (required for detailed explanation)

```text
feat: implement new authentication system

BREAKING CHANGE: The old token format is no longer supported.
All clients must upgrade to use the new JWT token format.
```

## Examples

### Good Examples

**Simple fix:**

```text
fix: prevent race condition in user session
```

**Feature with scope:**

```text
feat(api): add user profile image upload

Users can now upload and update their profile pictures.
Images are automatically resized and optimized for web delivery.
Supports JPEG, PNG, and WebP formats up to 5MB.

Closes #142
```

**Bug fix with detailed explanation:**

```text
fix: resolve memory leak in connection pooling

The connection pool was not properly releasing connections
after failed database queries, leading to memory buildup
over extended periods. Modified the cleanup logic to ensure
all connections are disposed regardless of query outcome.

This bug was causing production servers to run out of memory
after approximately 48 hours of operation under normal load.

Fixes #456
Reviewed-by: John Doe <john@example.com>
```

**Breaking change with migration info:**

```text
feat!: migrate from REST to GraphQL API

BREAKING CHANGE: All REST endpoints have been removed.
The new GraphQL endpoint is available at /graphql.
See migration guide: docs/graphql-migration.md

Closes #789
```

**Refactoring with context:**

```text
refactor(auth): extract JWT validation into middleware

Consolidate token validation logic that was duplicated across
multiple route handlers. This improves maintainability and
ensures consistent security checks throughout the application.

No functional changes to the API.
```

### Poor Examples (Avoid These)

```text
fixed stuff                           // Too vague
Update user.js                        // No context
WIP: working on authentication        // Not a completed change
Added new feature.                    // Uses past tense, ends with period
FIXED CRITICAL BUG!!!                 // Poor capitalization, excessive punctuation
added feature to do the thing that users wanted for a while now  // Too long
```

## Commit Message Anti-Patterns

**Avoid these common mistakes:**

- Vague messages: "fixed bug", "updated files"
- Past tense: "Added feature", "Fixed issue"
- Present continuous: "Adding feature", "Fixing bug"
- Excessive punctuation: "Fix bug!!!", "Update docs???"
- All caps: "URGENT FIX"
- Personal notes: "Going to lunch", "End of day commit"
- Technical jargon without context
- Commit messages longer than one screen

## Tools and Automation

To maintain consistency:

- Use commit message templates
- Set up commit-msg hooks for validation
- Consider tools like Commitizen or commitlint
- Configure your editor to wrap at 72 characters
- Use `git commit` (without -m) for multi-line messages

## Best Practices Summary

**Focus on the reader:**

- Write for your future self and teammates
- Assume the reader doesn't know the background
- Provide enough context to understand the change

**Be professional:**

- Use clear, concise language
- Maintain consistent formatting
- Follow team conventions

**Think like a journalist:**

- Answer: What changed and why?
- Provide the most important information first
- Keep the audience in mind

**Remember:**
A well-crafted commit message is a letter to your future self and your teammates. It should clearly communicate the intent and context of your changes, making the codebase easier to maintain and understand.
