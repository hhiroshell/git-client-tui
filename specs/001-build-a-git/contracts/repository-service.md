# Repository Service Unit Tests

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

## Unit Test Plan

### DiscoverRepository() Unit Tests

1. **TestDiscoverRepositorySuccess**
   - **Setup**: Create a test directory with a .git folder
   - **Mock**: Mock the filesystem operations to return a valid git repository path
   - **Assertions**: 
     - Verify Repository struct is returned with correct path
     - Verify no error is returned

2. **TestDiscoverRepositoryNotFound**
   - **Setup**: Create a test directory without a .git folder
   - **Mock**: Mock the filesystem operations to simulate no git repository
   - **Assertions**: 
     - Verify ErrNotARepository is returned
     - Verify nil Repository is returned

3. **TestDiscoverRepositoryPermissionDenied**
   - **Setup**: Create a test directory with permissions issue
   - **Mock**: Mock filesystem operations to return permission denied error
   - **Assertions**: 
     - Verify ErrPermissionDenied is returned
     - Verify nil Repository is returned

### GetStatus() Unit Tests

1. **TestGetStatusSuccess**
   - **Setup**: Create a mock repository with known state
   - **Mock**: Mock go-git worktree status to return predictable status
   - **Assertions**:
     - Verify RepositoryStatus contains expected staged/unstaged changes
     - Verify no error is returned

2. **TestGetStatusOperationFailure**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git worktree operations to fail with specific error
   - **Assertions**:
     - Verify ErrGitOperationFailed is returned
     - Verify nil RepositoryStatus is returned

3. **TestGetStatusWithFileChanges**
   - **Setup**: Create a mock repository with specific file changes
   - **Mock**: Mock go-git worktree status to return modified, added, deleted files
   - **Assertions**:
     - Verify RepositoryStatus correctly categorizes each file change
     - Verify file paths and change types match expectations

### Refresh() Unit Tests

1. **TestRefreshSuccess**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git repository operations to simulate refresh
   - **Assertions**:
     - Verify no error is returned
     - Verify repository state is updated

2. **TestRefreshFailure**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git operations to fail during refresh
   - **Assertions**:
     - Verify appropriate error is returned

## Error Handling

### Standard Errors
```go
var (
    ErrNotARepository      = errors.New("not a git repository")
    ErrPermissionDenied    = errors.New("permission denied")
    ErrGitOperationFailed  = errors.New("git operation failed")
    ErrInvalidReference    = errors.New("invalid git reference")
    ErrMergeConflict      = errors.New("merge conflict detected")
)
```

### Error Handling Tests

1. **TestRepositoryServiceErrorWrapping**
   - **Setup**: Create scenarios that trigger different errors
   - **Assertions**: 
     - Verify errors are properly wrapped with context
     - Verify error messages are descriptive and helpful
     - Verify error types can be identified through errors.Is()