# Repository Service Contract

## Interface Definition

### Repository Discovery and Initialization
```go
type RepositoryService interface {
    // Discover repository starting from current directory
    DiscoverRepository() (*Repository, error)

    // Get current repository status
    GetStatus() (*RepositoryStatus, error)

    // Refresh repository state
    Refresh() error
}
```

### Input/Output Contracts

#### DiscoverRepository()
**Input**: None (uses current working directory)
**Output**:
- Success: `Repository` with path, current branch, clean status
- Error: `ErrNotARepository` if no .git found, `ErrPermissionDenied` if access denied

#### GetStatus()
**Input**: None
**Output**:
- Success: `RepositoryStatus` with staged/unstaged changes
- Error: `ErrGitCommandFailed` if git status fails

## Error Contracts

### Standard Errors
```go
var (
    ErrNotARepository      = errors.New("not a git repository")
    ErrPermissionDenied    = errors.New("permission denied")
    ErrGitCommandFailed    = errors.New("git command failed")
    ErrInvalidReference    = errors.New("invalid git reference")
    ErrMergeConflict      = errors.New("merge conflict detected")
)
```

### Error Response Format
```go
type ServiceError struct {
    Operation string
    Cause     error
    Context   map[string]interface{}
}
```