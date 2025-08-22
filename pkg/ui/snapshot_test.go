package ui

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestSnapshotBasicView tests visual regression using view snapshots
func TestSnapshotBasicView(t *testing.T) {
	app := NewApp()
	
	// Initialize with standard terminal size
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = model.(*App)
	
	// Get the initial view
	view := app.View()
	
	// Verify the snapshot
	verifySnapshot(t, "basic_view", view)
}

// TestSnapshotWithContent tests view after drawing operations
func TestSnapshotWithContent(t *testing.T) {
	app := NewApp()
	
	// Initialize with standard terminal size
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = model.(*App)
	
	// Draw some content
	app.drawAtCursor()
	app.moveCursor(2, 1)
	app.setPenChar('#')
	app.drawAtCursor()
	app.moveCursor(4, 2)
	app.setPenChar('+')
	app.drawAtCursor()
	
	// Get the view with content
	view := app.View()
	
	// Verify the snapshot
	verifySnapshot(t, "with_content", view)
}

// TestSnapshotHelpModal tests help modal view
func TestSnapshotHelpModal(t *testing.T) {
	app := NewApp()
	
	// Initialize with standard terminal size
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = model.(*App)
	
	// Open help modal
	model, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = model.(*App)
	
	// Get the view with help modal
	view := app.View()
	
	// Verify the snapshot
	verifySnapshot(t, "help_modal", view)
}

// TestSnapshotLayerInfo tests layer information display
func TestSnapshotLayerInfo(t *testing.T) {
	app := NewApp()
	
	// Initialize with standard terminal size
	model, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = model.(*App)
	
	// Add multiple layers
	app.canvas.AddLayer("Background")
	app.canvas.AddLayer("Foreground") 
	app.canvas.SetCurrentLayer(1)
	
	// Draw something
	app.drawAtCursor()
	
	// Toggle layer visibility
	app.canvas.SetLayerVisible(1, false)
	
	// Get the view with layer information
	view := app.View()
	
	// Verify the snapshot
	verifySnapshot(t, "layer_info", view)
}

// verifySnapshot compares the current view with a stored snapshot
func verifySnapshot(t *testing.T, testName, view string) {
	t.Helper()
	
	// Create snapshot directory if it doesn't exist
	snapshotDir := "testdata/snapshots"
	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		t.Fatalf("Failed to create snapshot directory: %v", err)
	}
	
	snapshotFile := filepath.Join(snapshotDir, fmt.Sprintf("%s.golden", testName))
	
	// Generate hash of current view for fast comparison
	currentHash := sha256.Sum256([]byte(view))
	
	// Check if running with update flag (go test -update)
	updateSnapshots := os.Getenv("UPDATE_SNAPSHOTS") == "1"
	
	if updateSnapshots {
		// Write new snapshot
		if err := os.WriteFile(snapshotFile, []byte(view), 0644); err != nil {
			t.Fatalf("Failed to write snapshot file %s: %v", snapshotFile, err)
		}
		t.Logf("Updated snapshot: %s", snapshotFile)
		return
	}
	
	// Read existing snapshot
	expected, err := os.ReadFile(snapshotFile)
	if err != nil {
		if os.IsNotExist(err) {
			// Create initial snapshot
			if err := os.WriteFile(snapshotFile, []byte(view), 0644); err != nil {
				t.Fatalf("Failed to create initial snapshot file %s: %v", snapshotFile, err)
			}
			t.Logf("Created initial snapshot: %s", snapshotFile)
			return
		}
		t.Fatalf("Failed to read snapshot file %s: %v", snapshotFile, err)
	}
	
	// Compare hashes first for performance
	expectedHash := sha256.Sum256(expected)
	if currentHash == expectedHash {
		return // Views match
	}
	
	// If hashes don't match, compare line by line for better error messages
	expectedLines := splitLines(string(expected))
	actualLines := splitLines(view)
	
	if len(expectedLines) != len(actualLines) {
		t.Errorf("Snapshot %s: line count mismatch. Expected %d lines, got %d lines",
			testName, len(expectedLines), len(actualLines))
	}
	
	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}
	
	differenceFound := false
	for i := 0; i < maxLines; i++ {
		var expectedLine, actualLine string
		
		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}
		
		if expectedLine != actualLine {
			if !differenceFound {
				t.Errorf("Snapshot %s differs from expected:", testName)
				differenceFound = true
			}
			t.Errorf("Line %d:", i+1)
			t.Errorf("  Expected: %q", expectedLine)
			t.Errorf("  Actual:   %q", actualLine)
		}
	}
	
	if differenceFound {
		t.Errorf("To update snapshots, run: UPDATE_SNAPSHOTS=1 go test")
	}
}

// splitLines splits a string into lines, preserving empty lines
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	
	lines := []string{}
	start := 0
	
	for i, char := range s {
		if char == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	
	// Add the last line if there's content after the last newline
	// or if the string doesn't end with a newline
	if start < len(s) {
		lines = append(lines, s[start:])
	} else if len(s) > 0 && s[len(s)-1] == '\n' {
		// If string ends with newline, add an empty line
		lines = append(lines, "")
	}
	
	return lines
}

// TestSnapshotHelpers tests the snapshot testing helpers
func TestSnapshotHelpers(t *testing.T) {
	// Test splitLines function
	testCases := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{"single line", []string{"single line"}},
		{"line1\nline2", []string{"line1", "line2"}},
		{"line1\nline2\n", []string{"line1", "line2", ""}},
		{"line1\n\nline3", []string{"line1", "", "line3"}},
	}
	
	for i, tc := range testCases {
		result := splitLines(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("Test case %d: expected %d lines, got %d", i, len(tc.expected), len(result))
			continue
		}
		
		for j, line := range result {
			if line != tc.expected[j] {
				t.Errorf("Test case %d, line %d: expected %q, got %q", i, j, tc.expected[j], line)
			}
		}
	}
}