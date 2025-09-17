# Branch Service Contract

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

### Input/Output Contracts

#### ListBranches()
**Input**: None
**Output**:
- Success: `BranchList` with local and remote branches
- Error: `ErrGitCommandFailed`, `ErrRemoteUnavailable`

#### CreateBranch(name string)
**Input**: Branch name (must be valid git reference)
**Output**:
- Success: `nil`
- Error: `ErrInvalidBranchName`, `ErrBranchExists`

#### SwitchBranch(name string)
**Input**: Branch name (local or remote)
**Output**:
- Success: `nil`
- Error: `ErrBranchNotFound`, `ErrUncommittedChanges`, `ErrMergeConflict`

#### MergeBranch(sourceBranch string)
**Input**: Source branch name to merge into current branch
**Output**:
- Success: `MergeResult` with commit hash and file changes
- Error: `ErrMergeConflict`, `ErrBranchNotFound`

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