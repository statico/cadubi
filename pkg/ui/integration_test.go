package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cadubi/cadubi/pkg/commands"
	tea "github.com/charmbracelet/bubbletea"
)

func TestHSLColorPickerCreation(t *testing.T) {
	picker := NewHSLColorPicker()
	
	if picker == nil {
		t.Fatal("Failed to create HSL color picker")
	}

	if picker.IsVisible() {
		t.Error("New color picker should not be visible")
	}

	// Test show/hide
	picker.Show()
	if !picker.IsVisible() {
		t.Error("Color picker should be visible after Show()")
	}

	picker.Hide()
	if picker.IsVisible() {
		t.Error("Color picker should not be visible after Hide()")
	}
}

func TestHSLColorPickerNavigation(t *testing.T) {
	picker := NewHSLColorPicker()
	picker.Show()

	// Test keyboard navigation
	initialIndex := picker.selectedIndex

	// Test right movement
	picker.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	if picker.selectedIndex != initialIndex+1 {
		t.Errorf("Expected selectedIndex to be %d, got %d", initialIndex+1, picker.selectedIndex)
	}

	// Test down movement (should move by 16 positions in wheel mode)
	picker.selectedIndex = 0
	picker.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if picker.selectedIndex != 16 {
		t.Errorf("Expected selectedIndex to be 16 after down movement, got %d", picker.selectedIndex)
	}

	// Test mode switching
	initialMode := picker.mode
	picker.Update(tea.KeyMsg{Type: tea.KeyTab})
	if picker.mode == initialMode {
		t.Error("Mode should change after Tab key")
	}
}

func TestAppFileOperations(t *testing.T) {
	app := NewApp()
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "app_test.cadubi")

	// Draw something on the canvas
	pixel := commands.Pixel{
		Char:       '@',
		Foreground: "red",
		Bold:       true,
	}
	
	// Set cursor position and draw
	app.cursor.X = 5
	app.cursor.Y = 3
	app.cursor.Pen = pixel
	app.ctx.Cursor = app.cursor
	
	// Draw at cursor position  
	app.drawAtCursor()

	// Verify canvas is modified
	if !app.canvas.IsModified() {
		t.Error("Canvas should be modified after drawing")
	}

	// Set filename and save
	app.canvas.SetFilename(testFile)
	app.saveFile()

	// Verify canvas is no longer modified after save
	if app.canvas.IsModified() {
		t.Error("Canvas should not be modified after save")
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Save file was not created: %s", testFile)
	}

	// Test loading
	app.loadFile(testFile)
	
	// Verify the pixel was loaded back
	loadedPixel, err := app.canvas.GetPixel(5, 3)
	if err != nil {
		t.Errorf("Failed to get loaded pixel: %v", err)
	}

	if loadedPixel.Char != '@' {
		t.Errorf("Expected loaded char '@', got '%c'", loadedPixel.Char)
	}

	if loadedPixel.Foreground != "red" {
		t.Errorf("Expected loaded color 'red', got '%s'", loadedPixel.Foreground)
	}
}

func TestAppColorSelection(t *testing.T) {
	app := NewApp()

	initialColor := app.cursor.Pen.Foreground

	// Test basic color change
	app.setPenColor("blue")
	if app.cursor.Pen.Foreground != "blue" {
		t.Errorf("Expected pen color 'blue', got '%s'", app.cursor.Pen.Foreground)
	}

	// Test color cycling
	app.cyclePenColor()
	if app.cursor.Pen.Foreground == initialColor {
		t.Error("Pen color should have changed after cycling")
	}

	// Test color picker integration
	colorMsg := ColorSelectedMsg{
		Color: "123",
		Name:  "Test Color",
	}
	
	app.Update(colorMsg)
	if app.cursor.Pen.Foreground != "123" {
		t.Errorf("Expected pen color '123' after ColorSelectedMsg, got '%s'", app.cursor.Pen.Foreground)
	}
}

func TestAppModifiedStateTracking(t *testing.T) {
	app := NewApp()

	// New app should not be modified
	if app.canvas.IsModified() {
		t.Error("New app canvas should not be modified")
	}

	// Drawing should mark as modified
	app.drawAtCursor()
	if !app.canvas.IsModified() {
		t.Error("Canvas should be modified after drawing")
	}

	// Test quit handling with unsaved changes
	_, cmd := app.handleQuit()
	if cmd != nil {
		// Should not quit immediately if modified
		t.Error("Should not quit immediately when canvas is modified")
	}

	// Dialog should be visible
	if !app.confirmDialog.IsVisible() {
		t.Error("Confirmation dialog should be visible when quitting with unsaved changes")
	}

	// Mark as not modified and test quit again
	app.canvas.SetModified(false)
	app.confirmDialog.Hide()
	
	_, cmd = app.handleQuit()
	if cmd == nil {
		t.Error("Should quit immediately when canvas is not modified")
	}
}