package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSaveDialogPrePopulationWithExistingFile(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Set a filename on the canvas
	testFilename := "/path/to/existing/myfile.cadubi"
	app.canvas.SetFilename(testFilename)

	// Make canvas dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Clear the filename to force the save dialog to appear
	app.canvas.SetFilename("")

	// Now try to save - should show dialog with pre-populated filename from getSaveDialogInitialValue
	ctrlSKey := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, _ = app.Update(ctrlSKey)

	if !app.fileInputDialog.IsVisible() {
		t.Error("File input dialog should be visible when saving without filename")
	}

	// For a new file (no filename), should be pre-populated with untitled.cadubi
	expectedValue := "untitled.cadubi"
	if app.fileInputDialog.textInput.Value() != expectedValue {
		t.Errorf("Expected dialog to be pre-populated with '%s', got '%s'", 
			expectedValue, app.fileInputDialog.textInput.Value())
	}

	// Check that cursor is at the beginning (position 0)
	if app.fileInputDialog.textInput.Position() != 0 {
		t.Errorf("Expected cursor at beginning (position 0), got position %d",
			app.fileInputDialog.textInput.Position())
	}
}

func TestSaveDialogPrePopulationWithNewFile(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// New file - no filename set
	if app.canvas.GetFilename() != "" {
		t.Error("New canvas should not have a filename")
	}

	// Make canvas dirty to trigger save dialog
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Now try to save - should show dialog with .cadubi extension
	ctrlSKey := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, _ = app.Update(ctrlSKey)

	if !app.fileInputDialog.IsVisible() {
		t.Error("File input dialog should be visible when saving new file")
	}

	// Check that the dialog is pre-populated with untitled.cadubi
	expectedValue := "untitled.cadubi"
	if app.fileInputDialog.textInput.Value() != expectedValue {
		t.Errorf("Expected dialog to be pre-populated with '%s', got '%s'", 
			expectedValue, app.fileInputDialog.textInput.Value())
	}

	// Check that cursor is at the beginning (position 0)
	if app.fileInputDialog.textInput.Position() != 0 {
		t.Errorf("Expected cursor at beginning (position 0), got position %d",
			app.fileInputDialog.textInput.Position())
	}
}

func TestQuitSaveDialogPrePopulationWithNewFile(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// New file - no filename, make canvas dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Try to quit - should show confirmation dialog
	_, _ = app.handleQuit()
	if !app.confirmDialog.IsVisible() {
		t.Error("Confirmation dialog should be visible when quitting with unsaved changes")
	}

	// Choose "Save & Exit" - should trigger save dialog since no filename
	confirmMsg := ConfirmationResultMsg{Action: "yes"}
	_, _ = app.Update(confirmMsg)

	if !app.fileInputDialog.IsVisible() {
		t.Error("File input dialog should be visible when choosing to save without filename")
	}

	// Should be pre-populated with untitled.cadubi for new file
	expectedValue := "untitled.cadubi"
	if app.fileInputDialog.textInput.Value() != expectedValue {
		t.Errorf("Expected dialog to be pre-populated with '%s', got '%s'", 
			expectedValue, app.fileInputDialog.textInput.Value())
	}
}

func TestQuitSaveDialogPrePopulationWithExistingFile(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Set a filename and make canvas dirty
	testFilename := "/some/path/myproject.cadubi"
	app.canvas.SetFilename(testFilename)
	
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Temporarily clear filename to force the save dialog to appear
	app.canvas.SetFilename("")

	// Try to quit - should show confirmation dialog
	_, _ = app.handleQuit()
	if !app.confirmDialog.IsVisible() {
		t.Error("Confirmation dialog should be visible when quitting with unsaved changes")
	}

	// Choose "Save & Exit" - should trigger save dialog
	confirmMsg := ConfirmationResultMsg{Action: "yes"}
	_, _ = app.Update(confirmMsg)

	if !app.fileInputDialog.IsVisible() {
		t.Error("File input dialog should be visible when choosing to save without filename")
	}

	// Should be pre-populated with untitled.cadubi since canvas has no filename
	expectedValue := "untitled.cadubi"
	if app.fileInputDialog.textInput.Value() != expectedValue {
		t.Errorf("Expected dialog to be pre-populated with '%s', got '%s'", 
			expectedValue, app.fileInputDialog.textInput.Value())
	}
}

func TestGetSaveDialogInitialValue(t *testing.T) {
	app := NewApp()

	// Test with no filename (new file)
	initialValue := app.getSaveDialogInitialValue()
	if initialValue != "untitled.cadubi" {
		t.Errorf("Expected 'untitled.cadubi' for new file, got '%s'", initialValue)
	}

	// Test with existing filename
	app.canvas.SetFilename("/path/to/mywork.cadubi")
	initialValue = app.getSaveDialogInitialValue()
	if initialValue != "mywork.cadubi" {
		t.Errorf("Expected 'mywork.cadubi' for existing file, got '%s'", initialValue)
	}

	// Test with complex path
	app.canvas.SetFilename("/very/deep/nested/path/with spaces/my project.cadubi")
	initialValue = app.getSaveDialogInitialValue()
	if initialValue != "my project.cadubi" {
		t.Errorf("Expected 'my project.cadubi' for complex path, got '%s'", initialValue)
	}
}

func TestFileInputDialogShowWithValue(t *testing.T) {
	dialog := NewFileInputDialog()

	// Test showing with empty value
	dialog.ShowWithValue("save", "")
	if dialog.textInput.Value() != "" {
		t.Error("Expected empty value when showing with empty string")
	}

	// Test showing with pre-populated value
	testValue := "test.cadubi"
	dialog.ShowWithValue("save", testValue)
	
	if dialog.textInput.Value() != testValue {
		t.Errorf("Expected '%s', got '%s'", testValue, dialog.textInput.Value())
	}

	// Test cursor positioning at start
	if dialog.textInput.Position() != 0 {
		t.Errorf("Expected cursor at beginning (position 0), got position %d",
			dialog.textInput.Position())
	}

	// Test that Show() still works without value
	dialog.Show("load")
	if dialog.textInput.Value() != "" {
		t.Error("Expected empty value when using Show() without value")
	}
}