package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDirtyStateTracking(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// New canvas should not be dirty
	if app.canvas.IsModified() {
		t.Error("New canvas should not be dirty")
	}

	// Drawing should mark as dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after drawing")
	}

	// Saving should mark as clean
	app.canvas.SetFilename("test_dirty.cadubi")
	app.saveFile()

	if app.canvas.IsModified() {
		t.Error("Canvas should be clean after saving")
	}

	// Erasing should mark as dirty
	xKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	_, _ = app.Update(xKey)

	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after erasing")
	}

	// Undo should still be dirty (it's a modification)
	uKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	_, _ = app.Update(uKey)

	if !app.canvas.IsModified() {
		t.Error("Canvas should still be dirty after undo")
	}

	// New file should be clean
	ctrlNKey := tea.KeyMsg{Type: tea.KeyCtrlN}
	_, _ = app.Update(ctrlNKey)

	if app.canvas.IsModified() {
		t.Error("New canvas should be clean")
	}
}

func TestSaveDialogOnlyShowsWhenDirty(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Ctrl+S on clean canvas should not show dialog
	ctrlSKey := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, _ = app.Update(ctrlSKey)

	if app.fileInputDialog.IsVisible() {
		t.Error("Save dialog should not appear for clean canvas")
	}

	// Make canvas dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Now Ctrl+S should show dialog (since no filename is set)
	_, _ = app.Update(ctrlSKey)

	if !app.fileInputDialog.IsVisible() {
		t.Error("Save dialog should appear for dirty canvas without filename")
	}
}

func TestLayerOperationsMarkDirty(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// New canvas should be clean
	if app.canvas.IsModified() {
		t.Error("New canvas should not be dirty")
	}

	// Adding a layer should mark as dirty
	app.canvas.AddLayer("Test Layer")
	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after adding layer")
	}

	// Clean canvas
	app.canvas.SetModified(false)

	// Toggling layer visibility should mark as dirty
	currentLayer := app.canvas.GetCurrentLayer()
	visible := app.canvas.IsLayerVisible(currentLayer)
	app.canvas.SetLayerVisible(currentLayer, !visible)

	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after changing layer visibility")
	}
}

func TestLoadFileMarksCleam(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Make canvas dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after drawing")
	}

	// Create a test file to load
	testFilename := "test_load_clean.cadubi"
	app.canvas.SetFilename(testFilename)
	app.saveFile()

	// Make dirty again
	_, _ = app.Update(spaceKey)
	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after drawing again")
	}

	// Load file should mark as clean
	app.loadFile(testFilename)
	if app.canvas.IsModified() {
		t.Error("Canvas should be clean after loading file")
	}
}

func TestQuitWithSaveShowsFileDialog(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Make canvas dirty
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	if !app.canvas.IsModified() {
		t.Error("Canvas should be dirty after drawing")
	}

	// Try to quit - should show confirmation dialog
	_, _ = app.handleQuit()
	if !app.confirmDialog.IsVisible() {
		t.Error("Confirmation dialog should be visible when quitting with unsaved changes")
	}

	// Choose "Save & Exit" - should show file dialog since no filename is set
	confirmMsg := ConfirmationResultMsg{Action: "yes"}
	_, _ = app.Update(confirmMsg)

	if !app.fileInputDialog.IsVisible() {
		t.Error("File input dialog should be visible when choosing to save without filename")
	}

	if !app.quitting {
		t.Error("App should be in quitting state waiting for file selection")
	}

	// Select a file - should save and quit
	fileMsg := FileSelectedMsg{Path: "/tmp/test_quit.cadubi", Mode: "save"}
	_, cmd := app.Update(fileMsg)

	// Should return tea.Quit command
	if cmd == nil {
		t.Error("Should quit after saving file during quit process")
	}
}