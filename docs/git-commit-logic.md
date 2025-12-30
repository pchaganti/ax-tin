# Git Commit Logic in `tin commit`

This document explains how `tin commit` handles git staging and committing.

## Overview

When you run `tin commit`, tin creates both a **tin commit** (stored in `.tin/commits/`) and a **git commit** to capture file changes. The git commit is always created first, so the tin commit can include its hash.

## The Flow

```
1. Stage all changed files via git add
2. If staged changes exist:
   a. Create git commit with thread URLs
   b. Get resulting git hash
3. Else:
   a. Use current HEAD as git hash
4. Create tin commit with git hash
5. Save tin commit, update branch ref
6. Mark threads as committed, update thread.GitCommitHash
7. Clear staging index
```

## Git Commit Message Format

### Single Thread

```
[tin <thread-id>] <first line of commit message>

<rest of commit message if multi-line>

<thread-host-url>/repo/<repo-path>/thread/<thread-id>/<content-hash>
```

### Multiple Threads

```
[tin] <first line of commit message>

<rest of commit message if multi-line>

Threads:
- <thread-host-url>/repo/<repo-path>/thread/<thread-id-1>/<content-hash-1>
- <thread-host-url>/repo/<repo-path>/thread/<thread-id-2>/<content-hash-2>
```

## Key Design Decisions

### Git Commit First, Then Tin Commit

The tin commit hash calculation includes the git commit hash. This creates a clear ordering:

1. Git commit is created (if there are file changes)
2. Git hash is obtained
3. Tin commit is created, including the git hash in its computation

This means the tin commit URL cannot appear in the git commit message (chicken-and-egg problem), but thread URLs can since thread IDs are known beforehand.

### Thread GitCommitHash Updates

When `tin commit` runs, all staged threads have their `GitCommitHash` field updated to point to the new git commit (or HEAD if no changes). This ensures:

- Threads reference the git commit that captured their code changes
- Subsequent `tin commit` calls will use the latest git hash

## SessionEnd Hook

The SessionEnd hook (`internal/hooks/claude_code.go`) also creates git commits when a Claude Code session ends. These use the format `[tin <thread-id>] <message>`.

The `tin commit` flow handles this correctly:
- If SessionEnd already committed changes: no staged git changes exist, so `tin commit` uses the existing HEAD hash
- If `tin commit` runs mid-session: staged changes exist, so `tin commit` creates a git commit

## Code Location

All this logic lives in `internal/commands/commit.go`:

- `formatGitCommitMessage()`: Creates git commit message with thread URLs
- `generateCommitMessage()`: Auto-generates tin commit message from threads
