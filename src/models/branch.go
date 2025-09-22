package models

import "errors"

// Branch-specific errors
var (
	ErrInvalidBranchName  = errors.New("invalid branch name")
	ErrBranchExists       = errors.New("branch already exists")
	ErrBranchNotFound     = errors.New("branch not found")
	ErrUncommittedChanges = errors.New("uncommitted changes present")
	ErrRemoteUnavailable  = errors.New("remote repository unavailable")
	ErrFetchFailed        = errors.New("fetch operation failed")
)

// Branch represents a git branch
type Branch struct {
	Name     string // Name of the branch
	IsRemote bool   // Whether this is a remote branch
	IsCurrent bool  // Whether this is the current branch
	CommitHash string // Hash of the commit this branch points to
	Remote   string // Name of the remote (for remote branches)
}

// BranchList represents a collection of branches
type BranchList struct {
	Current string   // Name of the current branch
	Local   []Branch // Local branches
	Remote  []Branch // Remote branches
}

// MergeResult represents the result of a merge operation
type MergeResult struct {
	CommitHash   string   // Hash of the merge commit
	FilesChanged []string // List of files that were changed
	Conflicts    []string // List of files with conflicts
	FastForward  bool     // Whether this was a fast-forward merge
}

// NewBranch creates a new Branch instance
func NewBranch(name string, isRemote bool, commitHash string) Branch {
	return Branch{
		Name:       name,
		IsRemote:   isRemote,
		IsCurrent:  false,
		CommitHash: commitHash,
		Remote:     "",
	}
}

// NewRemoteBranch creates a new remote Branch instance
func NewRemoteBranch(name string, remote string, commitHash string) Branch {
	return Branch{
		Name:       name,
		IsRemote:   true,
		IsCurrent:  false,
		CommitHash: commitHash,
		Remote:     remote,
	}
}

// FullName returns the full name of the branch including remote prefix if applicable
func (b *Branch) FullName() string {
	if b.IsRemote && b.Remote != "" {
		return b.Remote + "/" + b.Name
	}
	return b.Name
}

// NewBranchList creates a new BranchList instance
func NewBranchList(current string) *BranchList {
	return &BranchList{
		Current: current,
		Local:   []Branch{},
		Remote:  []Branch{},
	}
}

// AddLocalBranch adds a local branch to the list
func (bl *BranchList) AddLocalBranch(branch Branch) {
	branch.IsRemote = false
	if branch.Name == bl.Current {
		branch.IsCurrent = true
	}
	bl.Local = append(bl.Local, branch)
}

// AddRemoteBranch adds a remote branch to the list
func (bl *BranchList) AddRemoteBranch(branch Branch) {
	branch.IsRemote = true
	bl.Remote = append(bl.Remote, branch)
}

// NewMergeResult creates a new MergeResult instance
func NewMergeResult(commitHash string, fastForward bool) *MergeResult {
	return &MergeResult{
		CommitHash:   commitHash,
		FilesChanged: []string{},
		Conflicts:    []string{},
		FastForward:  fastForward,
	}
}

// AddChangedFile adds a file to the list of changed files
func (mr *MergeResult) AddChangedFile(filename string) {
	mr.FilesChanged = append(mr.FilesChanged, filename)
}

// AddConflict adds a file to the list of conflicted files
func (mr *MergeResult) AddConflict(filename string) {
	mr.Conflicts = append(mr.Conflicts, filename)
}

// HasConflicts returns whether the merge has any conflicts
func (mr *MergeResult) HasConflicts() bool {
	return len(mr.Conflicts) > 0
}