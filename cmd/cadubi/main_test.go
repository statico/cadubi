package main

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cadubi/cadubi/pkg/ui"
)

func TestAppInitialization(t *testing.T) {
	// Test that our app initializes correctly
	app := ui.NewApp()
	
	// Test Init
	cmd := app.Init()
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
	
	// Test that view renders without error
	view := app.View()
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}
}

func TestAppUpdate(t *testing.T) {
	app := ui.NewApp()
	
	// Test window resize
	msg := tea.WindowSizeMsg{Width: 80, Height: 24}
	updated, cmd := app.Update(msg)
	
	if cmd != nil {
		t.Error("Expected window size update to return nil command")
	}
	
	// Verify the app was updated
	if updated == nil {
		t.Error("Expected Update to return updated app")
	}
	
	// Test quit command
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd = app.Update(quitMsg)
	
	if cmd == nil {
		t.Error("Expected quit key to return quit command")
	}
}

func TestAppView(t *testing.T) {
	app := ui.NewApp()
	
	// Simulate window size
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	
	view := app.View()
	
	if view == "" {
		t.Error("Expected View to return non-empty string")
	}
	
	// Check that view contains expected elements
	expectedStrings := []string{
		"Tool:",
		"Pen:",
		"Cursor:",
		"Layer:",
		"Move:",
	}
	
	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Expected view to contain '%s'", expected)
		}
	}
}