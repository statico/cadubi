package main

import (
	"strings"
	"testing"

	"github.com/cadubi/cadubi/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// TestHeadlessExecution demonstrates running Bubble Tea headless for testing
func TestHeadlessExecution(t *testing.T) {
	app := ui.NewApp()

	// Simulate a sequence of events
	events := []tea.Msg{
		tea.WindowSizeMsg{Width: 80, Height: 24},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}, // Draw
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, // Move right
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}, // Draw
	}

	var finalModel tea.Model = app

	// Process events without rendering
	for _, event := range events {
		var cmd tea.Cmd
		finalModel, cmd = finalModel.Update(event)

		// For this test, we ignore commands but in a real test we'd handle them
		_ = cmd
	}

	// Test that view renders correctly
	view := finalModel.View()
	if !strings.Contains(view, "Tool:") {
		t.Error("Expected view to show tool information")
	}

	if !strings.Contains(view, "Cursor:") {
		t.Error("Expected view to show cursor information")
	}

	t.Log("Headless execution successful - UI components render properly")
}
