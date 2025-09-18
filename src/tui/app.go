package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI application state
type Model struct {
	ready  bool
	width  int
	height int
}

// NewModel creates a new TUI model
func NewModel() Model {
	return Model{}
}

// Init initializes the TUI
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	return fmt.Sprintf("Git Client TUI\nPress q to quit\nTerminal size: %dx%d", m.width, m.height)
}

// Start initializes and starts the TUI application
func Start(debug bool) error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
