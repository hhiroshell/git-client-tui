package models

// LineType represents the type of a line in a diff
type LineType int

const (
	// Context represents an unchanged line in a diff
	Context LineType = iota
	// Addition represents a line that was added
	Addition
	// Deletion represents a line that was removed
	Deletion
)

// Line represents a single line of change with selection state
type Line struct {
	Type       LineType // Context, Addition, Deletion
	Content    string   // The content of the line
	OldLineNo  int      // Line number in the old file (0 for additions)
	NewLineNo  int      // Line number in the new file (0 for deletions)
	IsSelected bool     // Whether the line is selected for staging/unstaging
}

// NewContextLine creates a new context line
func NewContextLine(content string, oldLineNo, newLineNo int) Line {
	return Line{
		Type:       Context,
		Content:    content,
		OldLineNo:  oldLineNo,
		NewLineNo:  newLineNo,
		IsSelected: false,
	}
}

// NewAdditionLine creates a new addition line
func NewAdditionLine(content string, newLineNo int) Line {
	return Line{
		Type:       Addition,
		Content:    content,
		OldLineNo:  0, // No line in old file
		NewLineNo:  newLineNo,
		IsSelected: false,
	}
}

// NewDeletionLine creates a new deletion line
func NewDeletionLine(content string, oldLineNo int) Line {
	return Line{
		Type:       Deletion,
		Content:    content,
		OldLineNo:  oldLineNo,
		NewLineNo:  0, // No line in new file
		IsSelected: false,
	}
}

// Toggle switches the selection state of the line
func (l *Line) Toggle() {
	l.IsSelected = !l.IsSelected
}

// Select marks the line as selected
func (l *Line) Select() {
	l.IsSelected = true
}

// Deselect marks the line as not selected
func (l *Line) Deselect() {
	l.IsSelected = false
}

// IsAddition returns whether the line is an addition
func (l *Line) IsAddition() bool {
	return l.Type == Addition
}

// IsDeletion returns whether the line is a deletion
func (l *Line) IsDeletion() bool {
	return l.Type == Deletion
}

// IsContext returns whether the line is a context line
func (l *Line) IsContext() bool {
	return l.Type == Context
}