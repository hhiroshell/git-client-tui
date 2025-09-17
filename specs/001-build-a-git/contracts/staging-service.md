# Staging Service Contract

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

### Input/Output Contracts

#### StageFile(filePath string)
**Input**: Relative path to file from repository root
**Output**:
- Success: `nil`
- Error: `ErrFileNotFound`, `ErrInvalidPath`, `ErrGitOperationFailed`

#### StageHunks(filePath string, hunkIndexes []int)
**Input**:
- `filePath`: Relative path to file
- `hunkIndexes`: Zero-based indexes of hunks to stage
**Output**:
- Success: `nil`
- Error: `ErrInvalidHunkIndex`, `ErrPatchGenerationFailed`

#### StageLines(filePath string, lineSelections []LineSelection)
**Input**:
- `filePath`: Relative path to file
- `lineSelections`: Array of line selection ranges
**Output**:
- Success: `nil`
- Error: `ErrInvalidLineSelection`, `ErrPatchApplicationFailed`

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