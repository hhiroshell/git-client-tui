package repository

import (
	"errors"
	"testing"
)

func TestGetStatus_Success_CleanRepository(t *testing.T) {
	// Arrange
	mockStatus := &RepositoryStatus{
		Branch:        "main",
		IsClean:       true,
		StagedChanges:   []Change{},
		UnstagedChanges: []Change{},
	}
	
	svc := &MockRepositoryService{
		ShouldFail:     false,
		StatusToReturn: mockStatus,
	}
	
	// Act
	status, err := svc.GetStatus()
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if status == nil {
		t.Fatal("Expected status to be returned, got nil")
	}
	
	if status.Branch != "main" {
		t.Errorf("Expected Branch to be 'main', got '%s'", status.Branch)
	}
	
	if !status.IsClean {
		t.Errorf("Expected IsClean to be true, got false")
	}
	
	if len(status.StagedChanges) != 0 {
		t.Errorf("Expected 0 staged changes, got %d", len(status.StagedChanges))
	}
	
	if len(status.UnstagedChanges) != 0 {
		t.Errorf("Expected 0 unstaged changes, got %d", len(status.UnstagedChanges))
	}
}

func TestGetStatus_Success_WithChanges(t *testing.T) {
	// Arrange
	mockStatus := &RepositoryStatus{
		Branch:  "feature/new-feature",
		IsClean: false,
		StagedChanges: []Change{
			{
				FilePath:   "file1.go",
				ChangeType: Modified,
				IsStaged:   true,
			},
		},
		UnstagedChanges: []Change{
			{
				FilePath:   "file2.go",
				ChangeType: Modified,
				IsStaged:   false,
			},
			{
				FilePath:   "new-file.go",
				ChangeType: Added,
				IsStaged:   false,
			},
		},
	}
	
	svc := &MockRepositoryService{
		ShouldFail:     false,
		StatusToReturn: mockStatus,
	}
	
	// Act
	status, err := svc.GetStatus()
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if status == nil {
		t.Fatal("Expected status to be returned, got nil")
	}
	
	if status.Branch != "feature/new-feature" {
		t.Errorf("Expected Branch to be 'feature/new-feature', got '%s'", status.Branch)
	}
	
	if status.IsClean {
		t.Errorf("Expected IsClean to be false, got true")
	}
	
	if len(status.StagedChanges) != 1 {
		t.Errorf("Expected 1 staged change, got %d", len(status.StagedChanges))
	} else {
		stagedChange := status.StagedChanges[0]
		if stagedChange.FilePath != "file1.go" {
			t.Errorf("Expected staged change for 'file1.go', got '%s'", stagedChange.FilePath)
		}
		if stagedChange.ChangeType != Modified {
			t.Errorf("Expected staged change type to be Modified")
		}
		if !stagedChange.IsStaged {
			t.Errorf("Expected IsStaged to be true for staged change")
		}
	}
	
	if len(status.UnstagedChanges) != 2 {
		t.Errorf("Expected 2 unstaged changes, got %d", len(status.UnstagedChanges))
	} else {
		// Check first unstaged change
		if status.UnstagedChanges[0].FilePath != "file2.go" {
			t.Errorf("Expected unstaged change for 'file2.go', got '%s'", status.UnstagedChanges[0].FilePath)
		}
		if status.UnstagedChanges[0].ChangeType != Modified {
			t.Errorf("Expected unstaged change type to be Modified")
		}
		
		// Check second unstaged change
		if status.UnstagedChanges[1].FilePath != "new-file.go" {
			t.Errorf("Expected unstaged change for 'new-file.go', got '%s'", status.UnstagedChanges[1].FilePath)
		}
		if status.UnstagedChanges[1].ChangeType != Added {
			t.Errorf("Expected unstaged change type to be Added")
		}
	}
}

func TestGetStatus_GitCommandFailed(t *testing.T) {
	// Arrange
	svc := &MockRepositoryService{
		ShouldFail:    true,
		ErrorToReturn: ErrGitCommandFailed,
	}
	
	// Act
	status, err := svc.GetStatus()
	
	// Assert
	if !errors.Is(err, ErrGitCommandFailed) {
		t.Errorf("Expected error to be %v, got %v", ErrGitCommandFailed, err)
	}
	
	if status != nil {
		t.Errorf("Expected status to be nil, got: %v", status)
	}
}

func TestGetStatus_WithHunkDetails(t *testing.T) {
	// This test would verify that the status correctly includes hunk details
	// For contract testing purposes, we're keeping this simple
	// but the real implementation should handle hunks and lines
	
	// Create a mock status with hunk details
	mockStatus := &RepositoryStatus{
		Branch:  "main",
		IsClean: false,
		StagedChanges: []Change{},
		UnstagedChanges: []Change{
			{
				FilePath:   "file.go",
				ChangeType: Modified,
				IsStaged:   false,
				Hunks: []Hunk{
					{
						OldStart:   10,
						OldLines:   3,
						NewStart:   10,
						NewLines:   4,
						IsSelected: false,
						Lines: []Line{
							{Type: Context, Content: " func main() {", OldLineNo: 10, NewLineNo: 10},
							{Type: Deletion, Content: "-    fmt.Println(\"Hello\")", OldLineNo: 11, NewLineNo: 0},
							{Type: Addition, Content: "+    fmt.Println(\"Hello, World!\")", OldLineNo: 0, NewLineNo: 11},
							{Type: Addition, Content: "+    fmt.Println(\"How are you?\")", OldLineNo: 0, NewLineNo: 12},
							{Type: Context, Content: " }", OldLineNo: 12, NewLineNo: 13},
						},
					},
				},
			},
		},
	}
	
	svc := &MockRepositoryService{
		ShouldFail:     false,
		StatusToReturn: mockStatus,
	}
	
	// Act
	status, err := svc.GetStatus()
	
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if status == nil {
		t.Fatal("Expected status to be returned, got nil")
	}
	
	if len(status.UnstagedChanges) != 1 {
		t.Fatalf("Expected 1 unstaged change, got %d", len(status.UnstagedChanges))
	}
	
	change := status.UnstagedChanges[0]
	if len(change.Hunks) != 1 {
		t.Fatalf("Expected 1 hunk, got %d", len(change.Hunks))
	}
	
	hunk := change.Hunks[0]
	if hunk.OldStart != 10 || hunk.OldLines != 3 {
		t.Errorf("Expected hunk to start at line 10 with 3 old lines")
	}
	
	if hunk.NewStart != 10 || hunk.NewLines != 4 {
		t.Errorf("Expected hunk to start at line 10 with 4 new lines")
	}
	
	if len(hunk.Lines) != 5 {
		t.Fatalf("Expected 5 lines in hunk, got %d", len(hunk.Lines))
	}
	
	// Check that we have the right number of each line type
	contextCount := 0
	additionCount := 0
	deletionCount := 0
	
	for _, line := range hunk.Lines {
		switch line.Type {
		case Context:
			contextCount++
		case Addition:
			additionCount++
		case Deletion:
			deletionCount++
		}
	}
	
	if contextCount != 2 {
		t.Errorf("Expected 2 context lines, got %d", contextCount)
	}
	
	if additionCount != 2 {
		t.Errorf("Expected 2 addition lines, got %d", additionCount)
	}
	
	if deletionCount != 1 {
		t.Errorf("Expected 1 deletion line, got %d", deletionCount)
	}
}