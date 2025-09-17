# Data Model: Git Client TUI

## Core Entities

### Repository
Represents a git repository and its current state.
```go
type Repository struct {
    Path         string
    WorkingDir   string
    CurrentBranch string
    IsClean      bool
    Remotes      []Remote
}
```
**Relationships**: Contains Branches, Commits, Changes
**Validation**: Path must contain .git directory, WorkingDir must be valid directory

### Branch
Represents a git branch with tracking information.
```go
type Branch struct {
    Name          string
    IsLocal       bool
    IsRemote      bool
    IsCurrent     bool
    TrackingBranch string
    AheadBy       int
    BehindBy      int
    LastCommit    Commit
}
```
**Relationships**: Belongs to Repository, contains Commits
**Validation**: Name must be valid git reference, tracking info consistent

### Commit
Represents a git commit with metadata.
```go
type Commit struct {
    Hash      string
    Message   string
    Author    Author
    Timestamp time.Time
    Parents   []string
    Files     []FileChange
}
```
**Relationships**: Belongs to Branch, contains FileChanges
**Validation**: Hash must be valid SHA, message non-empty

### Change
Represents modifications to files with diff information.
```go
type Change struct {
    FilePath    string
    ChangeType  ChangeType // Added, Modified, Deleted, Renamed
    IsStaged    bool
    Hunks       []Hunk
    OldMode     os.FileMode
    NewMode     os.FileMode
}
```
**Relationships**: Belongs to Repository working directory
**Validation**: FilePath must be relative to repository root

### Hunk
Represents a contiguous section of changes within a file.
```go
type Hunk struct {
    OldStart    int
    OldLines    int
    NewStart    int
    NewLines    int
    Lines       []Line
    IsSelected  bool
}
```
**Relationships**: Belongs to Change
**Validation**: Line counts must match actual lines, positions must be valid

### Line
Represents a single line of change with selection state.
```go
type Line struct {
    Type       LineType // Context, Addition, Deletion
    Content    string
    OldLineNo  int
    NewLineNo  int
    IsSelected bool
}
```
**Relationships**: Belongs to Hunk
**Validation**: Line numbers consistent with hunk positions

## State Transitions

### Staging Workflow
```
Untracked → Staged (git add)
Modified → Staged (git add)
Staged → Committed (git commit)
Staged → Modified (git reset)
```

### Branch Workflow
```
Current Branch → Switch Branch (git checkout)
Branch → Merged (git merge)
Commits → Squashed (git rebase -i)
```

### Selection States
```
File: Unselected → Selected → Staged
Hunk: Unselected → Selected → Staged
Line: Unselected → Selected → Staged
```

## Validation Rules

### Repository Validation
- Must contain valid .git directory
- Working directory must be accessible
- Git operations must be possible (not corrupted)

### Change Validation
- File paths must be relative to repository root
- Hunks must cover all changed lines
- Selection state must be consistent across hierarchy

### Operation Validation
- Staging operations must preserve git repository integrity
- Branch operations must not lose uncommitted changes
- Commit operations must have valid author information

## Data Flow

### Loading Repository State
```
1. Discover repository root (.git directory)
2. Load current branch and status
3. Parse staged and unstaged changes
4. Build change hierarchy (files → hunks → lines)
```

### Staging Operation Flow
```
1. User selects lines/hunks/files
2. Generate patch from selection
3. Validate patch integrity
4. Apply patch to staging area
5. Refresh repository state
```

### Branch Operation Flow
```
1. Fetch remote references
2. List all available branches
3. User selects target branch
4. Check for uncommitted changes
5. Execute checkout/merge operation
6. Refresh repository state
```