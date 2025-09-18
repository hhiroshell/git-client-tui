package repository

import (
	"errors"
	"testing"
)

func TestDiscoverRepository_Success(t *testing.T) {
	// Arrange
	mockRepo := &Repository{
		Path:         "/path/to/repo",
		WorkingDir:   "/path/to/repo",
		CurrentBranch: "main",
		IsClean:      true,
		Remotes: []Remote{
			{Name: "origin", URL: "https://github.com/user/repo.git"},
		},
	}
	
	svc := &MockRepositoryService{
		ShouldFail:       false,
		RepositoryToReturn: mockRepo,
	}
	
	// Act
	repo, err := svc.DiscoverRepository()
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if repo == nil {
		t.Fatal("Expected repository to be returned, got nil")
	}
	
	if repo.Path != mockRepo.Path {
		t.Errorf("Expected Path to be %s, got %s", mockRepo.Path, repo.Path)
	}
	
	if repo.WorkingDir != mockRepo.WorkingDir {
		t.Errorf("Expected WorkingDir to be %s, got %s", mockRepo.WorkingDir, repo.WorkingDir)
	}
	
	if repo.CurrentBranch != mockRepo.CurrentBranch {
		t.Errorf("Expected CurrentBranch to be %s, got %s", mockRepo.CurrentBranch, repo.CurrentBranch)
	}
	
	if repo.IsClean != mockRepo.IsClean {
		t.Errorf("Expected IsClean to be %v, got %v", mockRepo.IsClean, repo.IsClean)
	}
	
	if len(repo.Remotes) != len(mockRepo.Remotes) {
		t.Errorf("Expected %d remotes, got %d", len(mockRepo.Remotes), len(repo.Remotes))
	}
}

func TestDiscoverRepository_NotARepository(t *testing.T) {
	// Arrange
	svc := &MockRepositoryService{
		ShouldFail:    true,
		ErrorToReturn: ErrNotARepository,
	}
	
	// Act
	repo, err := svc.DiscoverRepository()
	
	// Assert
	if !errors.Is(err, ErrNotARepository) {
		t.Errorf("Expected error to be %v, got %v", ErrNotARepository, err)
	}
	
	if repo != nil {
		t.Errorf("Expected repository to be nil, got: %v", repo)
	}
}

func TestDiscoverRepository_PermissionDenied(t *testing.T) {
	// Arrange
	svc := &MockRepositoryService{
		ShouldFail:    true,
		ErrorToReturn: ErrPermissionDenied,
	}
	
	// Act
	repo, err := svc.DiscoverRepository()
	
	// Assert
	if !errors.Is(err, ErrPermissionDenied) {
		t.Errorf("Expected error to be %v, got %v", ErrPermissionDenied, err)
	}
	
	if repo != nil {
		t.Errorf("Expected repository to be nil, got: %v", repo)
	}
}

func TestDiscoverRepository_DirectoryNavigation(t *testing.T) {
	t.Skip("This test requires a real filesystem implementation")
	
	// This test would verify that the repository is correctly discovered
	// when called from a subdirectory of the repository.
	// Since this is just a contract test, we're skipping the implementation
	// but the real implementation should handle this case.
}