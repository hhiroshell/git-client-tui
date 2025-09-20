# Commit Service Unit Tests

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

## Unit Test Plan

### GetCommitHistory() Unit Tests

1. **TestGetCommitHistorySuccess**
   - **Setup**: Create a mock repository with known commit history
   - **Mock**: Mock go-git repository.Log to return predefined commit history
   - **Assertions**:
     - Verify returned commits match expected count and order
     - Verify commit details (hash, message, author) are correctly parsed
     - Verify no error is returned

2. **TestGetCommitHistoryInvalidLimit**
   - **Setup**: Create a mock repository
   - **Mock**: Mock validation to fail with invalid limit error
   - **Assertions**:
     - Verify ErrInvalidLimit is returned
     - Verify nil commit array is returned

3. **TestGetCommitHistoryEmptyRepository**
   - **Setup**: Create a mock empty repository with no commits
   - **Mock**: Mock go-git repository.Log to return empty history
   - **Assertions**:
     - Verify empty array is returned (not nil)
     - Verify no error is returned

### SquashCommits() Unit Tests

1. **TestSquashCommitsSuccess**
   - **Setup**: Create a mock repository with multiple commits
   - **Mock**: Mock go-git repository rebase operations for squashing
   - **Assertions**:
     - Verify returned Commit has correct hash and message
     - Verify history is updated correctly
     - Verify no error is returned

2. **TestSquashCommitsInvalidHash**
   - **Setup**: Create a mock repository
   - **Mock**: Mock validation to fail with invalid commit hash
   - **Assertions**:
     - Verify ErrInvalidCommitHash is returned
     - Verify nil Commit is returned

3. **TestSquashCommitsNonContiguous**
   - **Setup**: Create a mock repository
   - **Mock**: Mock validation to fail with non-contiguous commits
   - **Assertions**:
     - Verify ErrNonContiguousCommits is returned
     - Verify nil Commit is returned

4. **TestSquashCommitsRebaseConflict**
   - **Setup**: Create a mock repository with conflicting changes
   - **Mock**: Mock go-git rebase operations to fail with conflict
   - **Assertions**:
     - Verify ErrRebaseConflict is returned
     - Verify nil Commit is returned

### EditCommitMessage() Unit Tests

1. **TestEditCommitMessageSuccess**
   - **Setup**: Create a mock repository with target commit
   - **Mock**: Mock go-git commit object operations for message editing
   - **Assertions**:
     - Verify no error is returned
     - Verify commit message is updated (via subsequent check)

2. **TestEditCommitMessageNotFound**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git repository to fail with commit not found
   - **Assertions**:
     - Verify ErrCommitNotFound is returned

3. **TestEditCommitMessagePushedCommit**
   - **Setup**: Create a mock repository with pushed commit
   - **Mock**: Mock validation to detect pushed commit
   - **Assertions**:
     - Verify ErrPushedCommit is returned

4. **TestEditCommitMessageInvalidMessage**
   - **Setup**: Create a mock repository
   - **Mock**: Mock validation to fail with invalid message
   - **Assertions**:
     - Verify ErrInvalidMessage is returned

### CreateCommit() Unit Tests

1. **TestCreateCommitSuccess**
   - **Setup**: Create a mock repository with staged changes
   - **Mock**: Mock go-git worktree.Commit operation
   - **Assertions**:
     - Verify returned Commit has correct hash and message
     - Verify no error is returned

2. **TestCreateCommitNoStagedChanges**
   - **Setup**: Create a mock repository with no staged changes
   - **Mock**: Mock go-git worktree.Commit to fail with no changes error
   - **Assertions**:
     - Verify ErrNoStagedChanges is returned
     - Verify nil Commit is returned

3. **TestCreateCommitInvalidMessage**
   - **Setup**: Create a mock repository with staged changes
   - **Mock**: Mock validation to fail with invalid message
   - **Assertions**:
     - Verify ErrInvalidMessage is returned
     - Verify nil Commit is returned

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

### Error Handling Tests

1. **TestCommitServiceErrorWrapping**
   - **Setup**: Create scenarios that trigger different commit-related errors
   - **Assertions**:
     - Verify errors are properly wrapped with context
     - Verify error messages are descriptive and helpful
     - Verify error types can be identified through errors.Is()

2. **TestCommitValidationBehavior**
   - **Setup**: Create test cases with various commit message formats
   - **Assertions**:
     - Verify validation correctly identifies valid/invalid messages
     - Verify validation correctly identifies valid/invalid commit hashes