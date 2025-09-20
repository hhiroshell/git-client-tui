# Git Client TUI - Quickstart Guide

## Prerequisites
- Go 1.21+ installed
- Git 2.0+ installed on system
- Terminal with ANSI color support
- Minimum 80x24 terminal size

## Installation and Basic Usage

### 1. Build and Install
```bash
# Clone the repository
git clone <repository-url>
cd my-git-client

# Build the application
go build -o gitui ./cmd/gitui

# Install system-wide (optional)
go install ./cmd/gitui
```

### 2. Quick Start Scenarios

#### Scenario 1: Granular Staging
```bash
# Navigate to a git repository with changes
cd /path/to/your/git/repo

# Launch the TUI
gitui

# Expected workflow:
# 1. See list of modified files
# 2. Press Enter on a file to view diff
# 3. Use 's' to stage individual lines/hunks
# 4. Use 'u' to unstage changes
# 5. Press Tab to switch between views
```

#### Scenario 2: Branch Management
```bash
# Launch TUI in repository
gitui

# Expected workflow:
# 1. Press 'b' for branch view
# 2. See local and remote branches automatically fetched
# 3. Press Enter to switch branches
# 4. Press 'n' to create new branch
# 5. Press 'm' to merge from another branch
```

#### Scenario 3: Commit Operations
```bash
# Launch TUI with staged changes
gitui

# Expected workflow:
# 1. Press 'c' to commit staged changes
# 2. Press 'h' for commit history
# 3. Select multiple commits and press 'q' to squash
# 4. Press 'e' to edit commit message
```

## Validation Tests

### Test 1: File-level Staging
```bash
# Setup test environment
mkdir test-repo && cd test-repo
git init
echo "line 1" > file1.txt
git add file1.txt && git commit -m "initial"
echo "line 2" >> file1.txt
echo "line 3" >> file1.txt

# Run test
gitui

# Verify:
# - File shows as modified
# - Can stage entire file
# - Staging area updates correctly
```

### Test 2: Line-level Staging
```bash
# Setup with multiple changes
echo "added line" >> file1.txt
echo "another change" >> file1.txt

# Run test
gitui

# Verify:
# - Individual lines can be selected
# - Partial staging works correctly
# - Diff view shows staged vs unstaged
```

### Test 3: Branch Creation and Switching
```bash
# Setup repository with remote
git remote add origin https://github.com/example/repo.git
git fetch

# Run test
gitui

# Verify:
# - Remote branches appear in list
# - Can create new branch
# - Can switch between branches
# - Auto-fetch works
```

### Test 4: Commit Squashing
```bash
# Setup multiple commits
echo "change 1" >> file1.txt && git add . && git commit -m "commit 1"
echo "change 2" >> file1.txt && git add . && git commit -m "commit 2"
echo "change 3" >> file1.txt && git add . && git commit -m "commit 3"

# Run test
gitui

# Verify:
# - Multiple commits can be selected
# - Squash operation combines them
# - New commit message can be set
# - History shows single commit
```

## Expected Outcomes

### Success Criteria
- ✅ TUI launches without errors in any git repository
- ✅ All file changes are visible and navigable
- ✅ Staging operations work at file, hunk, and line levels
- ✅ Branch operations complete without data loss
- ✅ Commit operations maintain git repository integrity
- ✅ UI is responsive and intuitive

### Performance Benchmarks
- ✅ Startup time < 500ms for typical repositories
- ✅ Diff loading < 100ms for files under 1000 lines
- ✅ Branch list loading < 200ms with network fetch
- ✅ Memory usage < 50MB for repositories with 10k files

### Error Handling Verification
- ✅ Graceful degradation when go-git operations fail
- ✅ Clear error messages for user mistakes
- ✅ No data corruption on unexpected exit
- ✅ Proper handling of merge conflicts

## Troubleshooting

### Common Issues
- **"Not a git repository"**: Ensure you're in a directory with .git folder
- **"Permission denied"**: Check file permissions on repository
- **"Command not found"**: Ensure git is installed and in PATH
- **UI corruption**: Check terminal size and ANSI color support

### Debug Mode
```bash
# Run with debug logging
gitui --debug

# Check logs
tail -f ~/.gitui/debug.log
```