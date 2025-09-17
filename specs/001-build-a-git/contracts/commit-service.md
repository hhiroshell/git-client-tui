# Commit Service Contract

## Interface Definition

### Commit Operations
```go
type CommitService interface {
    // Get commit history for current branch
    GetCommitHistory(limit int) ([]Commit, error)

    // Squash multiple commits into one
    SquashCommits(commitHashes []string, newMessage string) (*Commit, error)

    // Edit commit message for specific commit
    EditCommitMessage(commitHash string, newMessage string) error

    // Create new commit with staged changes
    CreateCommit(message string) (*Commit, error)
}
```

### Input/Output Contracts

#### GetCommitHistory(limit int)
**Input**: Maximum number of commits to retrieve
**Output**:
- Success: Array of `Commit` objects in reverse chronological order
- Error: `ErrGitCommandFailed`, `ErrInvalidLimit`

#### SquashCommits(commitHashes []string, newMessage string)
**Input**:
- `commitHashes`: Array of commit hashes to squash (must be contiguous)
- `newMessage`: Commit message for the squashed commit
**Output**:
- Success: New `Commit` object representing the squashed commit
- Error: `ErrInvalidCommitHash`, `ErrNonContiguousCommits`, `ErrRebaseConflict`

#### EditCommitMessage(commitHash string, newMessage string)
**Input**:
- `commitHash`: Hash of commit to edit
- `newMessage`: New commit message
**Output**:
- Success: `nil`
- Error: `ErrCommitNotFound`, `ErrPushedCommit`, `ErrInvalidMessage`

#### CreateCommit(message string)
**Input**: Commit message
**Output**:
- Success: New `Commit` object
- Error: `ErrNoStagedChanges`, `ErrInvalidMessage`

## Data Structures

### Commit
```go
type Commit struct {
    Hash        string
    Message     string
    Author      Author
    Timestamp   time.Time
    Parents     []string
    FilesChanged []string
}
```

### Author
```go
type Author struct {
    Name  string
    Email string
}
```

## Error Handling

### Commit-specific Errors
```go
var (
    ErrInvalidLimit         = errors.New("invalid commit limit")
    ErrInvalidCommitHash    = errors.New("invalid commit hash")
    ErrNonContiguousCommits = errors.New("commits are not contiguous")
    ErrRebaseConflict       = errors.New("rebase conflict during squash")
    ErrCommitNotFound       = errors.New("commit not found")
    ErrPushedCommit         = errors.New("cannot edit pushed commit")
    ErrInvalidMessage       = errors.New("invalid commit message")
    ErrNoStagedChanges      = errors.New("no staged changes to commit")
)
```