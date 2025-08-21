package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// MainModel is the root model that handles the main app
type MainModel struct {
	app          *App
	currentModel tea.Model
	initialized  bool
	width        int
	height       int
}

// NewMainModel creates a new main model
func NewMainModel() *MainModel {
	// Go directly to main app - no first-run dialog needed
	app := NewApp()
	
	m := &MainModel{
		app:          app,
		currentModel: app,
		initialized:  true,
	}
	
	return m
}

// Init initializes the main model
func (m *MainModel) Init() tea.Cmd {
	return m.currentModel.Init()
}

// Update handles messages for the main model
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.initialized {
		return m, nil
	}
	
	// Track window size changes
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Set the app dimensions immediately
		if m.app != nil {
			m.app.width = m.width
			m.app.height = m.height
		}
	}
	
	// Delegate to the main app
	newModel, cmd := m.currentModel.Update(msg)
	m.app = newModel.(*App)
	m.currentModel = m.app
	
	return m, cmd
}

// View renders the current model
func (m *MainModel) View() string {
	if !m.initialized {
		return "Initializing..."
	}
	
	return m.currentModel.View()
}