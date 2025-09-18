package models

import (
	"errors"
	"path/filepath"
)

// Standard error definitions
var (
	ErrNotARepository   = errors.New("not a git repository")
	ErrPermissionDenied = errors.New("permission denied")
	ErrGitCommandFailed = errors.New("git command failed")
	ErrInvalidReference = errors.New("invalid git reference")
	ErrMergeConflict    = errors.New("merge conflict detected")
)

// ServiceError wraps operation errors with context
type ServiceError struct {
	Operation string
	Cause     error
	Context   map[string]interface{}
}

func (e *ServiceError) Error() string {
	return e.Operation + ": " + e.Cause.Error()
}

func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// Remote represents a git remote
type Remote struct {
	Name string
	URL  string
}

// Repository represents a git repository and its current state
type Repository struct {
	Path         string     // Path to the repository (.git directory)
	WorkingDir   string     // Path to the working directory
	CurrentBranch string     // Name of the current branch
	IsClean      bool       // Whether the repository has no changes
	Remotes      []Remote   // List of remotes
}

// RepositoryStatus represents the current state of the repository
type RepositoryStatus struct {
	Branch         string     // Name of the current branch
	IsClean        bool       // Whether the repository has no changes
	StagedChanges  []Change   // Files staged for commit
	UnstagedChanges []Change   // Unstaged changes
}

// IsGitRepository checks if the given path is a git repository
func IsGitRepository(path string) bool {
	// In a real implementation, this would check for the existence of .git
	// For now, we'll just return true for testing
	return true
}

// NewRepository creates a new Repository instance
func NewRepository(path string) (*Repository, error) {
	// In a real implementation, this would:
	// 1. Check if path is a git repository
	// 2. Get the current branch
	// 3. Check if the repository is clean
	// 4. Get the remotes

	// For now, we'll just return a basic repository for testing
	return &Repository{
		Path:          filepath.Join(path, ".git"),
		WorkingDir:    path,
		CurrentBranch: "main",
		IsClean:       true,
		Remotes:       []Remote{},
	}, nil
}
