package models

import "os"

// ChangeType represents the type of change to a file
type ChangeType int

const (
	// Added represents a new file
	Added ChangeType = iota
	// Modified represents a changed file
	Modified
	// Deleted represents a removed file
	Deleted
	// Renamed represents a renamed file
	Renamed
)

// Change represents modifications to files with diff information
type Change struct {
	FilePath    string     // Path to the file, relative to the repository root
	ChangeType  ChangeType // Added, Modified, Deleted, Renamed
	IsStaged    bool       // Whether the change is staged
	Hunks       []Hunk     // Hunks of changes
	OldMode     os.FileMode // Old file mode
	NewMode     os.FileMode // New file mode
}

// NewChange creates a new change
func NewChange(filePath string, changeType ChangeType, isStaged bool) Change {
	return Change{
		FilePath:   filePath,
		ChangeType: changeType,
		IsStaged:   isStaged,
		Hunks:      []Hunk{},
		OldMode:    0,
		NewMode:    0,
	}
}

// AddHunk adds a hunk to the change
func (c *Change) AddHunk(hunk Hunk) {
	c.Hunks = append(c.Hunks, hunk)
}

// IsAdded returns whether the change is an addition
func (c *Change) IsAdded() bool {
	return c.ChangeType == Added
}

// IsModified returns whether the change is a modification
func (c *Change) IsModified() bool {
	return c.ChangeType == Modified
}

// IsDeleted returns whether the change is a deletion
func (c *Change) IsDeleted() bool {
	return c.ChangeType == Deleted
}

// IsRenamed returns whether the change is a rename
func (c *Change) IsRenamed() bool {
	return c.ChangeType == Renamed
}

// HasSelectedHunks returns whether any hunks in the change are selected
func (c *Change) HasSelectedHunks() bool {
	for _, hunk := range c.Hunks {
		if hunk.IsSelected || hunk.HasSelectedLines() {
			return true
		}
	}
	
	return false
}

// SelectedHunks returns the hunks that are selected
func (c *Change) SelectedHunks() []Hunk {
	selected := []Hunk{}
	
	for _, hunk := range c.Hunks {
		if hunk.IsSelected || hunk.HasSelectedLines() {
			selected = append(selected, hunk)
		}
	}
	
	return selected
}

// AllHunksSelected returns whether all hunks in the change are selected
func (c *Change) AllHunksSelected() bool {
	if len(c.Hunks) == 0 {
		return false
	}
	
	for _, hunk := range c.Hunks {
		if !hunk.IsSelected {
			return false
		}
	}
	
	return true
}

// SelectAllHunks selects all hunks in the change
func (c *Change) SelectAllHunks() {
	for i := range c.Hunks {
		c.Hunks[i].Select()
	}
}

// DeselectAllHunks deselects all hunks in the change
func (c *Change) DeselectAllHunks() {
	for i := range c.Hunks {
		c.Hunks[i].Deselect()
	}
}