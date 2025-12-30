# Plan: Add Code Host Links to Tin Commits

## Goal
Display GitHub commit URLs in `tin log` output, linking each tin commit to its underlying git commit on the code host.

## Configuration
- **Primary**: Auto-detect from `git remote get-url origin`
- **Override**: Allow explicit `code_host_url` in `.tin/config`
- **Scope**: GitHub only (SSH and HTTPS formats)

## Example Output

**Before:**
```
Git:    d6d83b2
```

**After:**
```
Git:    https://github.com/dadlerj/tin/commit/d6d83b2abc...
```

## Implementation Steps

### 1. Create `internal/git/codehost.go`
New file with URL parsing logic:
- `ParseGitRemoteURL(url string) (*CodeHostURL, error)` - parses GitHub SSH/HTTPS URLs
- `CodeHostURL.CommitURL(hash string) string` - generates full commit URL

Supported formats:
- `git@github.com:owner/repo.git` → `https://github.com/owner/repo`
- `https://github.com/owner/repo.git` → `https://github.com/owner/repo`

### 2. Update `internal/storage/repository.go`
Add to `Config` struct (line 38):
```go
CodeHostURL string `json:"code_host_url,omitempty"`
```

Add methods:
- `GetGitRemoteURL(name string) (string, error)` - runs `git remote get-url <name>`
- `GetCodeHostURL() (string, error)` - returns config override or parsed git remote

### 3. Update `internal/commands/log.go`
Modify line 74 to show full URL when available:
```go
// Current:
fmt.Printf("Git:    %s\n", commit.GitCommitHash[:min(8, len(commit.GitCommitHash))])

// New: Show full URL if available, fall back to short hash
```

### 4. Add tests `internal/git/codehost_test.go`
Test cases for URL parsing (SSH, HTTPS with/without .git suffix, unsupported formats).

## Files to Modify

| File | Change |
|------|--------|
| `internal/git/codehost.go` | CREATE - URL parsing logic |
| `internal/git/codehost_test.go` | CREATE - Unit tests |
| `internal/storage/repository.go` | ADD `CodeHostURL` to Config, add getter methods |
| `internal/commands/log.go` | MODIFY line 74 to display URL |
