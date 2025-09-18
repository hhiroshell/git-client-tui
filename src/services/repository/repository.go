package repository

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hhiroshell/git-client-tui/src/models"
)

// RepositoryService defines the interface for repository operations
type RepositoryService interface {
	// Discover repository starting from current directory
	DiscoverRepository() (*models.Repository, error)

	// Get current repository status
	GetStatus() (*models.RepositoryStatus, error)

	// Refresh repository state
	Refresh() error
}

// GitRepositoryService implements the RepositoryService interface
type GitRepositoryService struct {
	currentRepo *models.Repository
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
	// In a real implementation, this would:
	// 1. Check if startPath contains a .git directory
	// 2. If not, check its parent directory
	// 3. Continue until the root directory is reached
	
	// This implementation just checks if the current directory is a git repository
	gitPath := filepath.Join(startPath, ".git")
	_, err := os.Stat(gitPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &models.ServiceError{
				Operation: "findGitRepository",
				Cause:     models.ErrNotARepository,
				Context:   map[string]interface{}{"path": startPath},
			}
		}
		if os.IsPermission(err) {
			return nil, &models.ServiceError{
				Operation: "findGitRepository",
				Cause:     models.ErrPermissionDenied,
				Context:   map[string]interface{}{"path": startPath},
			}
		}
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     err,
			Context:   map[string]interface{}{"path": startPath},
		}
	}
	
	// Get the current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = startPath
	out, err := cmd.Output()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": startPath, "command": "git rev-parse --abbrev-ref HEAD"},
		}
	}
	currentBranch := strings.TrimSpace(string(out))
	
	// Check if the repository is clean
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = startPath
	out, err = cmd.Output()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "findGitRepository",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": startPath, "command": "git status --porcelain"},
		}
	}
	isClean := len(strings.TrimSpace(string(out))) == 0
	
	// Get the remotes
	remotes, err := s.getRemotes(startPath)
	if err != nil {
		return nil, err
	}
	
	return &models.Repository{
		Path:         gitPath,
		WorkingDir:   startPath,
		CurrentBranch: currentBranch,
		IsClean:      isClean,
		Remotes:      remotes,
	}, nil
}

// getRemotes gets the remotes of the repository
func (s *GitRepositoryService) getRemotes(repoPath string) ([]models.Remote, error) {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "getRemotes",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": repoPath, "command": "git remote -v"},
		}
	}
	
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	remotes := make([]models.Remote, 0)
	seen := make(map[string]bool)
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		
		name := parts[0]
		url := parts[1]
		
		// Skip duplicates (fetch and push entries)
		if seen[name] {
			continue
		}
		
		seen[name] = true
		remotes = append(remotes, models.Remote{
			Name: name,
			URL:  url,
		})
	}
	
	return remotes, nil
}

// GetStatus gets the current status of the repository
func (s *GitRepositoryService) GetStatus() (*models.RepositoryStatus, error) {
	// Make sure we have a repository
	if s.currentRepo == nil {
		_, err := s.DiscoverRepository()
		if err != nil {
			return nil, err
		}
	}
	
	// Get the current branch
	currentBranch := s.currentRepo.CurrentBranch
	
	// Get the status
	cmd := exec.Command("git", "status", "--porcelain", "-z")
	cmd.Dir = s.currentRepo.WorkingDir
	out, err := cmd.Output()
	if err != nil {
		return nil, &models.ServiceError{
			Operation: "GetStatus",
			Cause:     models.ErrGitCommandFailed,
			Context:   map[string]interface{}{"path": s.currentRepo.WorkingDir, "command": "git status --porcelain -z"},
		}
	}
	
	// Parse the status
	stagedChanges, unstagedChanges := s.parseStatus(out)
	
	// Create the repository status
	isClean := len(stagedChanges) == 0 && len(unstagedChanges) == 0
	
	return &models.RepositoryStatus{
		Branch:         currentBranch,
		IsClean:        isClean,
		StagedChanges:  stagedChanges,
		UnstagedChanges: unstagedChanges,
	}, nil
}

// parseStatus parses the output of git status --porcelain -z
func (s *GitRepositoryService) parseStatus(output []byte) ([]models.Change, []models.Change) {
	if len(output) == 0 {
		return []models.Change{}, []models.Change{}
	}
	
	entries := strings.Split(strings.TrimRight(string(output), "\x00"), "\x00")
	stagedChanges := make([]models.Change, 0)
	unstagedChanges := make([]models.Change, 0)
	
	for _, entry := range entries {
		if len(entry) < 3 {
			continue
		}
		
		statusCode := entry[0:2]
		filePath := entry[3:]
		
		// Check if the file is staged
		if statusCode[0] != ' ' && statusCode[0] != '?' {
			changeType := s.getChangeType(statusCode[0])
			stagedChanges = append(stagedChanges, models.NewChange(filePath, changeType, true))
		}
		
		// Check if the file is unstaged
		if statusCode[1] != ' ' {
			changeType := s.getChangeType(statusCode[1])
			unstagedChanges = append(unstagedChanges, models.NewChange(filePath, changeType, false))
		}
	}
	
	// In a real implementation, we would also fetch the diffs for each change
	// and create hunks and lines, but for this implementation we'll skip that
	
	return stagedChanges, unstagedChanges
}

// getChangeType converts a git status code to a ChangeType
func (s *GitRepositoryService) getChangeType(code byte) models.ChangeType {
	switch code {
	case 'A':
		return models.Added
	case 'M':
		return models.Modified
	case 'D':
		return models.Deleted
	case 'R':
		return models.Renamed
	default:
		// For simplicity, treat everything else as modified
		return models.Modified
	}
}

// Refresh updates the repository state
func (s *GitRepositoryService) Refresh() error {
	if s.currentRepo == nil {
		_, err := s.DiscoverRepository()
		return err
	}
	
	// Re-discover the repository
	repo, err := s.findGitRepository(s.currentRepo.WorkingDir)
	if err != nil {
		return err
	}
	
	s.currentRepo = repo
	return nil
}