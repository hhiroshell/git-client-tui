package repository

import (
	"errors"
)

// Define test interfaces based on contract
type RepositoryService interface {
	// Discover repository starting from current directory
	DiscoverRepository() (*Repository, error)
	// Get current repository status
	GetStatus() (*RepositoryStatus, error)
}

type Repository struct {
	Path          string
	WorkingDir    string
	CurrentBranch string
	IsClean       bool
	Remotes       []Remote
}

type Remote struct {
	Name string
	URL  string
}

type RepositoryStatus struct {
	Branch          string
	IsClean         bool
	StagedChanges   []Change
	UnstagedChanges []Change
}

type Change struct {
	FilePath   string
	ChangeType ChangeType
	IsStaged   bool
	Hunks      []Hunk
}

type Hunk struct {
	OldStart   int
	OldLines   int
	NewStart   int
	NewLines   int
	Lines      []Line
	IsSelected bool
}

type Line struct {
	Type       LineType
	Content    string
	OldLineNo  int
	NewLineNo  int
	IsSelected bool
}

type ChangeType int
type LineType int

const (
	Added ChangeType = iota
	Modified
	Deleted
	Renamed
)

const (
	Context LineType = iota
	Addition
	Deletion
)

// Define expected errors from the contract
var (
	ErrNotARepository   = errors.New("not a git repository")
	ErrPermissionDenied = errors.New("permission denied")
	ErrGitCommandFailed = errors.New("git command failed")
)

// MockRepositoryService implements the RepositoryService interface for testing
type MockRepositoryService struct {
	ShouldFail         bool
	ErrorToReturn      error
	RepositoryToReturn *Repository
	StatusToReturn     *RepositoryStatus
}

func (m *MockRepositoryService) DiscoverRepository() (*Repository, error) {
	if m.ShouldFail {
		return nil, m.ErrorToReturn
	}
	return m.RepositoryToReturn, nil
}

func (m *MockRepositoryService) GetStatus() (*RepositoryStatus, error) {
	if m.ShouldFail {
		return nil, m.ErrorToReturn
	}
	return m.StatusToReturn, nil
}