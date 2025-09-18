# Branch Service Unit Tests

## Interface Definition

### Branch Management Operations
```go
type BranchService interface {
    // List all branches (local and remote)
    ListBranches() (*BranchList, error)

    // Create new branch from current HEAD
    CreateBranch(name string) error

    // Switch to existing branch
    SwitchBranch(name string) error

    // Merge branch into current branch
    MergeBranch(sourceBranch string) (*MergeResult, error)

    // Fetch remote branches
    FetchRemotes() error
}
```

## Unit Test Plan

### ListBranches() Unit Tests

1. **TestListBranchesSuccess**
   - **Setup**: Create a mock repository with known branches
   - **Mock**: Mock git command execution to return predefined branch list
   - **Assertions**:
     - Verify BranchList contains expected local and remote branches
     - Verify current branch is correctly identified
     - Verify no error is returned

2. **TestListBranchesCommandFailure**
   - **Setup**: Create a mock repository
   - **Mock**: Mock git command to fail with specific error
   - **Assertions**:
     - Verify ErrGitCommandFailed is returned
     - Verify nil BranchList is returned

3. **TestListBranchesRemoteUnavailable**
   - **Setup**: Create a mock repository with unreachable remote
   - **Mock**: Mock git command to succeed for local branches but fail for remote branches
   - **Assertions**:
     - Verify ErrRemoteUnavailable is returned
     - Verify BranchList contains local branches but empty remote branches

### CreateBranch() Unit Tests

1. **TestCreateBranchSuccess**
   - **Setup**: Create a mock repository
   - **Mock**: Mock git command execution for branch creation
   - **Assertions**:
     - Verify no error is returned
     - Verify branch is created (via subsequent validation)

2. **TestCreateBranchInvalidName**
   - **Setup**: Create a mock repository
   - **Mock**: Mock git command to fail with invalid reference error
   - **Assertions**:
     - Verify ErrInvalidBranchName is returned

3. **TestCreateBranchAlreadyExists**
   - **Setup**: Create a mock repository with existing branch
   - **Mock**: Mock git command to fail with branch exists error
   - **Assertions**:
     - Verify ErrBranchExists is returned

### SwitchBranch() Unit Tests

1. **TestSwitchBranchSuccess**
   - **Setup**: Create a mock repository with multiple branches
   - **Mock**: Mock git command execution for branch switching
   - **Assertions**:
     - Verify no error is returned
     - Verify current branch is changed (via subsequent validation)

2. **TestSwitchBranchNotFound**
   - **Setup**: Create a mock repository
   - **Mock**: Mock git command to fail with branch not found error
   - **Assertions**:
     - Verify ErrBranchNotFound is returned

3. **TestSwitchBranchUncommittedChanges**
   - **Setup**: Create a mock repository with uncommitted changes
   - **Mock**: Mock git command to fail with uncommitted changes error
   - **Assertions**:
     - Verify ErrUncommittedChanges is returned

### MergeBranch() Unit Tests

1. **TestMergeBranchSuccess**
   - **Setup**: Create a mock repository with source and target branches
   - **Mock**: Mock git command execution for branch merging
   - **Assertions**:
     - Verify MergeResult contains expected commit hash and file changes
     - Verify no error is returned

2. **TestMergeBranchFastForward**
   - **Setup**: Create a mock repository with fast-forward merge scenario
   - **Mock**: Mock git command for fast-forward merge
   - **Assertions**:
     - Verify MergeResult.FastForward is true
     - Verify MergeResult contains expected changes

3. **TestMergeBranchConflict**
   - **Setup**: Create a mock repository with conflicting changes
   - **Mock**: Mock git command to fail with merge conflict error
   - **Assertions**:
     - Verify ErrMergeConflict is returned
     - Verify MergeResult contains conflict information

### FetchRemotes() Unit Tests

1. **TestFetchRemotesSuccess**
   - **Setup**: Create a mock repository with remotes
   - **Mock**: Mock git command execution for fetch
   - **Assertions**:
     - Verify no error is returned

2. **TestFetchRemotesFailure**
   - **Setup**: Create a mock repository with unreachable remote
   - **Mock**: Mock git command to fail with network error
   - **Assertions**:
     - Verify ErrFetchFailed is returned

## Data Structures

### BranchList
```go
type BranchList struct {
    Current string
    Local   []Branch
    Remote  []Branch
}
```

### MergeResult
```go
type MergeResult struct {
    CommitHash    string
    FilesChanged  []string
    Conflicts     []string
    FastForward   bool
}
```

## Error Handling

### Branch-specific Errors
```go
var (
    ErrInvalidBranchName   = errors.New("invalid branch name")
    ErrBranchExists        = errors.New("branch already exists")
    ErrBranchNotFound      = errors.New("branch not found")
    ErrUncommittedChanges  = errors.New("uncommitted changes present")
    ErrRemoteUnavailable   = errors.New("remote repository unavailable")
    ErrFetchFailed         = errors.New("fetch operation failed")
)
```

### Error Handling Tests

1. **TestBranchServiceErrorWrapping**
   - **Setup**: Create scenarios that trigger different branch-related errors
   - **Assertions**:
     - Verify errors are properly wrapped with context
     - Verify error messages are descriptive and helpful
     - Verify error types can be identified through errors.Is()