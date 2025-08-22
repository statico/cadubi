package ui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestToastShowAndHide(t *testing.T) {
	toast := NewToast()

	// Initially should not be visible
	if toast.IsVisible() {
		t.Error("New toast should not be visible")
	}

	// Show toast
	cmd := toast.Show("Test message", 100*time.Millisecond)
	if !toast.IsVisible() {
		t.Error("Toast should be visible after Show()")
	}

	// Check message
	view := toast.View()
	if view == "" {
		t.Error("Toast view should not be empty when visible")
	}

	// Test timeout command
	if cmd == nil {
		t.Error("Show should return a timeout command")
	}

	// Test hide
	toast.Hide()
	if toast.IsVisible() {
		t.Error("Toast should not be visible after Hide()")
	}

	view = toast.View()
	if view != "" {
		t.Error("Toast view should be empty when hidden")
	}
}

func TestToastUpdate(t *testing.T) {
	toast := NewToast()

	// Test ToastMsg
	msg := ToastMsg{Message: "Update test", Duration: 50 * time.Millisecond}
	cmd := toast.Update(msg)

	if !toast.IsVisible() {
		t.Error("Toast should be visible after ToastMsg")
	}

	if cmd == nil {
		t.Error("Update with ToastMsg should return timeout command")
	}

	// Wait for duration to pass
	time.Sleep(60 * time.Millisecond)

	// Test ToastTimeoutMsg - now it should hide since time has passed
	timeoutMsg := ToastTimeoutMsg{}
	cmd = toast.Update(timeoutMsg)

	if toast.IsVisible() {
		t.Error("Toast should be hidden after ToastTimeoutMsg when duration has passed")
	}

	if cmd != nil {
		t.Error("Update with ToastTimeoutMsg should return nil")
	}
}

func TestToastInApp(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Make sure toast is not visible initially
	if app.toast.IsVisible() {
		t.Error("Toast should not be visible initially")
	}

	// Create a dirty canvas and set filename
	app.cursor.X = 5
	app.cursor.Y = 5
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)
	app.canvas.SetFilename("test_toast.cadubi")

	// Save file - should show toast
	ctrlSKey := tea.KeyMsg{Type: tea.KeyCtrlS}
	_, cmd := app.Update(ctrlSKey)

	// The save command should return a toast message command
	if cmd == nil {
		t.Error("Save should return a command")
	}

	// Execute the toast message command
	if cmd != nil {
		msg := cmd()
		if toastMsg, ok := msg.(ToastMsg); ok {
			if toastMsg.Message != "File saved!" {
				t.Errorf("Expected 'File saved!' toast message, got '%s'", toastMsg.Message)
			}
			
			// Update app with the toast message
			_, timeoutCmd := app.Update(toastMsg)
			
			// Toast should now be visible
			if !app.toast.IsVisible() {
				t.Error("Toast should be visible after ToastMsg")
			}
			
			// Check that timeout command is returned
			if timeoutCmd == nil {
				t.Error("Toast update should return timeout command")
			}
		} else {
			t.Error("Save command should return a ToastMsg")
		}
	}
}

func TestToastView(t *testing.T) {
	toast := NewToast()

	// Hidden toast should return empty view
	view := toast.View()
	if view != "" {
		t.Error("Hidden toast should return empty view")
	}

	// Visible toast should return styled message
	toast.Show("Test message", 1*time.Second)
	view = toast.View()
	if view == "" {
		t.Error("Visible toast should return non-empty view")
	}

	// View should contain the message
	if !contains(view, "Test message") {
		t.Error("Toast view should contain the message")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(substr)] == substr ||
		   len(s) > len(substr) && contains(s[1:], substr)
}