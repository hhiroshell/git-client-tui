package models

// Hunk represents a contiguous section of changes within a file
type Hunk struct {
	OldStart   int     // Starting line in the old file
	OldLines   int     // Number of lines in the old file
	NewStart   int     // Starting line in the new file
	NewLines   int     // Number of lines in the new file
	Lines      []Line  // Lines in the hunk
	IsSelected bool    // Whether the hunk is selected for staging/unstaging
}

// NewHunk creates a new hunk
func NewHunk(oldStart, oldLines, newStart, newLines int) Hunk {
	return Hunk{
		OldStart:   oldStart,
		OldLines:   oldLines,
		NewStart:   newStart,
		NewLines:   newLines,
		Lines:      []Line{},
		IsSelected: false,
	}
}

// AddLine adds a line to the hunk
func (h *Hunk) AddLine(line Line) {
	h.Lines = append(h.Lines, line)
}

// ToggleSelection switches the selection state of the hunk and all its lines
func (h *Hunk) ToggleSelection() {
	h.IsSelected = !h.IsSelected
	
	for i := range h.Lines {
		h.Lines[i].IsSelected = h.IsSelected
	}
}

// Select marks the hunk and all its lines as selected
func (h *Hunk) Select() {
	h.IsSelected = true
	
	for i := range h.Lines {
		h.Lines[i].IsSelected = true
	}
}

// Deselect marks the hunk and all its lines as not selected
func (h *Hunk) Deselect() {
	h.IsSelected = false
	
	for i := range h.Lines {
		h.Lines[i].IsSelected = false
	}
}

// SelectedLines returns the lines that are selected
func (h *Hunk) SelectedLines() []Line {
	selected := []Line{}
	
	for _, line := range h.Lines {
		if line.IsSelected {
			selected = append(selected, line)
		}
	}
	
	return selected
}

// HasSelectedLines returns whether any lines in the hunk are selected
func (h *Hunk) HasSelectedLines() bool {
	for _, line := range h.Lines {
		if line.IsSelected {
			return true
		}
	}
	
	return false
}

// UpdateSelectionState updates the hunk selection state based on its lines
func (h *Hunk) UpdateSelectionState() {
	allSelected := true
	
	for _, line := range h.Lines {
		if !line.IsSelected {
			allSelected = false
			break
		}
	}
	
	h.IsSelected = allSelected
}