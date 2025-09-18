package repository

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"github.com/hhiroshell/git-client-tui/src/models"
)

// RepositoryService defines the interface for repository operations
type RepositoryService interface {
	// Discover repository starting from current directory
	DiscoverRepository() (*models.Repository, error)

	// Get current repository status
	GetStatus() (*models.RepositoryStatus, error)
}

// GitRepositoryService implements the RepositoryService interface
type GitRepositoryService struct {
	currentRepo *models.Repository
	gitRepo     *git.Repository
}

// NewGitRepositoryService creates a new GitRepositoryService
func NewGitRepositoryService() *GitRepositoryService {
	return &GitRepositoryService{
		currentRepo: nil,
	}
}

// DiscoverRepository finds a git repository from the current directory upwards
func (s *GitRepositoryService) DiscoverRepository() (*models.Repository, error) {
	// In a real implementation, this would start from the current directory
	// and search upwards for a .git directory

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "DiscoverRepository",
			Cause:     err,
			Context:   map[string]interface{}{"workingDir": wd},
		}
	}

	// Check if the directory is a git repository
	repo, err := s.findGitRepository(wd)
	if err != nil {
		return nil, err
	}

	// Store the repository for future use
	s.currentRepo = repo

	return repo, nil
}

// findGitRepository finds a git repository in the given path or its parents
func (s *GitRepositoryService) findGitRepository(startPath string) (*models.Repository, error) {
	// Open the repository using go-git
	repo, err := git.PlainOpenWithOptions(startPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, &models.ServiceError{
				Operation: "findGitRepository",
				Cause:     models.ErrNotARepository,
				Context:   map[string]interface{}{"path": startPath},
			}
		}
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     err,
			Context:   map[string]interface{}{"path": startPath},
		}
	}

	// Store the git repository for future use
	s.gitRepo = repo

	// Use the traditional location for .git directory path
	gitPath := filepath.Join(startPath, ".git")

	// Get the current branch
	head, err := repo.Head()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": startPath, "error": "unable to get HEAD"},
		}
	}

	// Extract branch name from the reference
	currentBranch := ""
	if head.Name().IsBranch() {
		currentBranch = head.Name().Short()
	} else {
		// We're in detached HEAD state
		currentBranch = head.Hash().String()[:7] // Short SHA
	}

	// Get the worktree status to check if repo is clean
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     err,
			Context:   map[string]interface{}{"path": startPath},
		}
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": startPath, "error": "unable to get status"},
		}
	}

	isClean := status.IsClean()

	// Get the remotes
	remotes, err := s.getRemotes(repo)
	if err != nil {
		return nil, err
	}

	return &models.Repository{
		Path:          gitPath,
		WorkingDir:    startPath,
		CurrentBranch: currentBranch,
		IsClean:       isClean,
		Remotes:       remotes,
	}, nil
}

// getRemotes gets the remotes of the repository using go-git
func (s *GitRepositoryService) getRemotes(repo *git.Repository) ([]models.Remote, error) {
	gitRemotes, err := repo.Remotes()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "getRemotes",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"error": "unable to get remotes"},
		}
	}

	remotes := make([]models.Remote, 0, len(gitRemotes))
	for _, remote := range gitRemotes {
		config := remote.Config()
		if len(config.URLs) > 0 {
			remotes = append(remotes, models.Remote{
				Name: config.Name,
				URL:  config.URLs[0], // Use the first URL
			})
		}
	}

	return remotes, nil
}

// GetStatus gets the current status of the repository
func (s *GitRepositoryService) GetStatus() (*models.RepositoryStatus, error) {
	// Make sure we have a repository
	if s.currentRepo == nil || s.gitRepo == nil {
		_, err := s.DiscoverRepository()
		if err != nil {
			return nil, err
		}
	}

	// Get the current branch
	currentBranch := s.currentRepo.CurrentBranch

	// Get the worktree status
	worktree, err := s.gitRepo.Worktree()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "GetStatus",
			Cause:     err,
			Context:   map[string]interface{}{"path": s.currentRepo.WorkingDir},
		}
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "GetStatus",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": s.currentRepo.WorkingDir, "error": "unable to get status"},
		}
	}

	// Parse the status to extract staged and unstaged changes
	stagedChanges, unstagedChanges := s.parseGoGitStatus(status)

	// Create the repository status
	isClean := status.IsClean()

	return &models.RepositoryStatus{
		Branch:          currentBranch,
		IsClean:         isClean,
		StagedChanges:   stagedChanges,
		UnstagedChanges: unstagedChanges,
	}, nil
}

// parseGoGitStatus parses the go-git Status object into our models
func (s *GitRepositoryService) parseGoGitStatus(status git.Status) ([]models.Change, []models.Change) {
	stagedChanges := make([]models.Change, 0)
	unstagedChanges := make([]models.Change, 0)

	// Iterate through all file statuses
	for filePath, fileStatus := range status {
		// Check if the file is staged
		if fileStatus.Staging != git.Unmodified && fileStatus.Staging != git.Untracked {
			changeType := s.mapGoGitStatusToChangeType(fileStatus.Staging)
			stagedChanges = append(stagedChanges, models.NewChange(filePath, changeType, true))
		}

		// Check if the file is unstaged
		if fileStatus.Worktree != git.Unmodified {
			changeType := s.mapGoGitStatusToChangeType(fileStatus.Worktree)
			unstagedChanges = append(unstagedChanges, models.NewChange(filePath, changeType, false))
		}
	}

	// In a real implementation, we would also fetch the diffs for each change
	// and create hunks and lines, but for this implementation we'll skip that

	return stagedChanges, unstagedChanges
}

// mapGoGitStatusToChangeType maps go-git status to our ChangeType
func (s *GitRepositoryService) mapGoGitStatusToChangeType(status git.StatusCode) models.ChangeType {
	switch status {
	case git.Added, git.Untracked:
		return models.Added
	case git.Modified, git.UpdatedButUnmerged:
		return models.Modified
	case git.Deleted:
		return models.Deleted
	case git.Renamed:
		return models.Renamed
	default:
		// For simplicity, treat everything else as modified
		return models.Modified
	}
}
