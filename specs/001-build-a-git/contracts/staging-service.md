# Staging Service Unit Tests

## Interface Definition

### File-level Operations
```go
type StagingService interface {
    // Stage entire file
    StageFile(filePath string) error

    // Unstage entire file
    UnstageFile(filePath string) error

    // Stage selected hunks from file
    StageHunks(filePath string, hunkIndexes []int) error

    // Stage selected lines from file
    StageLines(filePath string, lineSelections []LineSelection) error
}
```

## Unit Test Plan

### StageFile() Unit Tests

1. **TestStageFileSuccess**
   - **Setup**: Create a mock repository with unstaged file
   - **Mock**: Mock go-git worktree Add operation for staging file
   - **Assertions**:
     - Verify no error is returned
     - Verify file is staged (via subsequent status check)

2. **TestStageFileNotFound**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git worktree Add to fail with file not found error
   - **Assertions**:
     - Verify ErrFileNotFound is returned

3. **TestStageFileInvalidPath**
   - **Setup**: Create a mock repository
   - **Mock**: Mock validation to fail with invalid path error
   - **Assertions**:
     - Verify ErrInvalidPath is returned

4. **TestStageFileBinaryContent**
   - **Setup**: Create a mock repository with binary file
   - **Mock**: Mock go-git worktree Add for binary file staging
   - **Assertions**:
     - Verify binary file is handled correctly
     - Verify no error is returned

### UnstageFile() Unit Tests

1. **TestUnstageFileSuccess**
   - **Setup**: Create a mock repository with staged file
   - **Mock**: Mock go-git index Reset operation for unstaging file
   - **Assertions**:
     - Verify no error is returned
     - Verify file is unstaged (via subsequent status check)

2. **TestUnstageFileNotFound**
   - **Setup**: Create a mock repository
   - **Mock**: Mock go-git index Reset to fail with file not found error
   - **Assertions**:
     - Verify ErrFileNotFound is returned

3. **TestUnstageFileNotStaged**
   - **Setup**: Create a mock repository with file that isn't staged
   - **Mock**: Mock go-git index Reset operation for unstaging file
   - **Assertions**:
     - Verify appropriate handling (should not error)
     - Verify file remains unstaged

### StageHunks() Unit Tests

1. **TestStageHunksSuccess**
   - **Setup**: Create a mock repository with file containing multiple hunks
   - **Mock**: Mock go-git patch operations for hunk staging
   - **Assertions**:
     - Verify no error is returned
     - Verify specified hunks are staged (via subsequent status check)
     - Verify unspecified hunks remain unstaged

2. **TestStageHunksInvalidIndex**
   - **Setup**: Create a mock repository with file containing hunks
   - **Mock**: Mock validation to fail with invalid hunk index error
   - **Assertions**:
     - Verify ErrInvalidHunkIndex is returned

3. **TestStageHunksPatchGenerationFailure**
   - **Setup**: Create a mock repository with file
   - **Mock**: Mock patch generation to fail
   - **Assertions**:
     - Verify ErrPatchGenerationFailed is returned

### StageLines() Unit Tests

1. **TestStageLinesSuccess**
   - **Setup**: Create a mock repository with file containing multiple lines
   - **Mock**: Mock go-git patch operations for line-level staging
   - **Assertions**:
     - Verify no error is returned
     - Verify specified lines are staged (via subsequent status check)
     - Verify unspecified lines remain unstaged

2. **TestStageLinesInvalidSelection**
   - **Setup**: Create a mock repository with file
   - **Mock**: Mock validation to fail with invalid line selection error
   - **Assertions**:
     - Verify ErrInvalidLineSelection is returned

3. **TestStageLinesPatchApplicationFailure**
   - **Setup**: Create a mock repository with file
   - **Mock**: Mock patch application to fail
   - **Assertions**:
     - Verify ErrPatchApplicationFailed is returned

## Data Structures

### LineSelection
```go
type LineSelection struct {
    HunkIndex int
    StartLine int
    EndLine   int
    Type      LineType // Addition, Deletion
}
```

### Staging Result
```go
type StagingResult struct {
    FilesStaged   []string
    FilesUnstaged []string
    Conflicts     []StagingConflict
}
```

## Error Handling

### Staging-specific Errors
```go
var (
    ErrFileNotFound          = errors.New("file not found")
    ErrInvalidPath          = errors.New("invalid file path")
    ErrInvalidHunkIndex     = errors.New("invalid hunk index")
    ErrInvalidLineSelection = errors.New("invalid line selection")
    ErrPatchGenerationFailed = errors.New("patch generation failed")
    ErrPatchApplicationFailed = errors.New("patch application failed")
    ErrStagingConflict      = errors.New("staging conflict detected")
)
```

### Error Handling Tests

1. **TestStagingServiceErrorWrapping**
   - **Setup**: Create scenarios that trigger different staging-related errors
   - **Assertions**:
     - Verify errors are properly wrapped with context
     - Verify error messages are descriptive and helpful
     - Verify error types can be identified through errors.Is()

2. **TestStagingConflictHandling**
   - **Setup**: Create a mock repository with conflicting changes
   - **Mock**: Mock go-git operations to detect and report conflicts
   - **Assertions**:
     - Verify ErrStagingConflict is returned with proper context
     - Verify conflict information is correctly provided