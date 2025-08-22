package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHSLColorPickerRecentColors(t *testing.T) {
	picker := NewHSLColorPicker()

	// Initially, there should be no recent colors
	if len(picker.GetRecentColors()) != 0 {
		t.Errorf("Expected 0 recent colors initially, got %d", len(picker.GetRecentColors()))
	}

	// Test adding recent colors by terminal code
	picker.AddRecentColorByTerminalCode("red")
	recentColors := picker.GetRecentColors()
	if len(recentColors) != 1 {
		t.Errorf("Expected 1 recent color after adding red, got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "red" {
		t.Errorf("Expected first recent color to be 'red', got '%s'", recentColors[0].TerminalCode)
	}

	// Test adding another color
	picker.AddRecentColorByTerminalCode("blue")
	recentColors = picker.GetRecentColors()
	if len(recentColors) != 2 {
		t.Errorf("Expected 2 recent colors after adding blue, got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "blue" {
		t.Errorf("Expected first recent color to be 'blue' (most recent), got '%s'", recentColors[0].TerminalCode)
	}
	if recentColors[1].TerminalCode != "red" {
		t.Errorf("Expected second recent color to be 'red', got '%s'", recentColors[1].TerminalCode)
	}
}

func TestHSLColorPickerDeduplication(t *testing.T) {
	picker := NewHSLColorPicker()

	// Add same color twice - should not duplicate
	picker.AddRecentColorByTerminalCode("red")
	picker.AddRecentColorByTerminalCode("blue")
	picker.AddRecentColorByTerminalCode("red") // This should move red to front, not add duplicate

	recentColors := picker.GetRecentColors()
	if len(recentColors) != 2 {
		t.Errorf("Expected 2 recent colors (no duplicates), got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "red" {
		t.Errorf("Expected first recent color to be 'red' (moved to front), got '%s'", recentColors[0].TerminalCode)
	}
	if recentColors[1].TerminalCode != "blue" {
		t.Errorf("Expected second recent color to be 'blue', got '%s'", recentColors[1].TerminalCode)
	}
}

func TestHSLColorPickerMaxRecentColors(t *testing.T) {
	picker := NewHSLColorPicker()

	// Add more than 20 colors to test limit
	colors := []string{
		"red", "blue", "green", "yellow", "magenta", "cyan", "white", "black",
		"232", "233", "234", "235", "236", "237", "238", "239", "240", "241",
		"242", "243", "244", "245", "246", // 23 colors total
	}

	for _, color := range colors {
		picker.AddRecentColorByTerminalCode(color)
	}

	recentColors := picker.GetRecentColors()
	if len(recentColors) != 20 {
		t.Errorf("Expected maximum 20 recent colors, got %d", len(recentColors))
	}

	// Check that the most recent color is first
	if recentColors[0].TerminalCode != "246" {
		t.Errorf("Expected first recent color to be '246' (last added), got '%s'", recentColors[0].TerminalCode)
	}
}

func TestHSLColorPickerFindColorByTerminalCode(t *testing.T) {
	picker := NewHSLColorPicker()

	// Test finding existing colors
	redColor := picker.FindColorByTerminalCode("red")
	if redColor == nil {
		t.Error("Expected to find red color")
	} else {
		if redColor.Name != "Red" {
			t.Errorf("Expected red color name to be 'Red', got '%s'", redColor.Name)
		}
		if redColor.TerminalCode != "red" {
			t.Errorf("Expected red color terminal code to be 'red', got '%s'", redColor.TerminalCode)
		}
	}

	// Test finding non-existent color
	nonExistentColor := picker.FindColorByTerminalCode("nonexistent")
	if nonExistentColor != nil {
		t.Error("Expected not to find non-existent color")
	}
}

func TestHSLColorPickerModeSwitch(t *testing.T) {
	picker := NewHSLColorPicker()

	// Should start in wheel mode
	if picker.mode != "wheel" {
		t.Errorf("Expected initial mode to be 'wheel', got '%s'", picker.mode)
	}

	// Add some recent colors first
	picker.AddRecentColorByTerminalCode("red")
	picker.AddRecentColorByTerminalCode("blue")

	// Show the picker
	picker.Show()

	// Simulate Tab key press to switch modes
	msg := tea.KeyMsg{Type: tea.KeyTab}
	picker.Update(msg)
	if picker.mode != "recent" {
		t.Errorf("Expected mode to be 'recent' after tab key, got '%s'", picker.mode)
	}

	// Press Tab again to switch back
	picker.Update(msg)
	if picker.mode != "wheel" {
		t.Errorf("Expected mode to be 'wheel' after second tab key, got '%s'", picker.mode)
	}
}

func TestHSLColorPickerRecentModeNavigation(t *testing.T) {
	picker := NewHSLColorPicker()

	// Add several colors for navigation testing
	colors := []string{"red", "blue", "green", "yellow", "magenta"}
	for _, color := range colors {
		picker.AddRecentColorByTerminalCode(color)
	}

	picker.Show()

	// Switch to recent mode
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	picker.Update(tabMsg)

	// Test navigation
	if picker.selectedIndex != 0 {
		t.Errorf("Expected initial selected index to be 0, got %d", picker.selectedIndex)
	}

	// Test right movement
	rightMsg := tea.KeyMsg{Type: tea.KeyRight}
	picker.Update(rightMsg)
	if picker.selectedIndex != 1 {
		t.Errorf("Expected selected index to be 1 after right key, got %d", picker.selectedIndex)
	}

	// Test left movement
	leftMsg := tea.KeyMsg{Type: tea.KeyLeft}
	picker.Update(leftMsg)
	if picker.selectedIndex != 0 {
		t.Errorf("Expected selected index to be 0 after left key, got %d", picker.selectedIndex)
	}
}

func TestAppRecentColorsIntegration(t *testing.T) {
	app := NewApp()

	// Initially, there should be no recent colors
	if len(app.GetColorPicker().GetRecentColors()) != 0 {
		t.Errorf("Expected 0 recent colors initially, got %d", len(app.GetColorPicker().GetRecentColors()))
	}

	// Test that setPenColor adds to recent colors
	app.SetPenColor("red")
	recentColors := app.GetColorPicker().GetRecentColors()
	if len(recentColors) != 1 {
		t.Errorf("Expected 1 recent color after setting red, got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "red" {
		t.Errorf("Expected recent color to be 'red', got '%s'", recentColors[0].TerminalCode)
	}

	// Test multiple color changes
	app.SetPenColor("blue")
	app.SetPenColor("green")
	recentColors = app.GetColorPicker().GetRecentColors()
	if len(recentColors) != 3 {
		t.Errorf("Expected 3 recent colors, got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "green" {
		t.Errorf("Expected first recent color to be 'green' (most recent), got '%s'", recentColors[0].TerminalCode)
	}
}

func TestAppColorCyclingAddsToRecent(t *testing.T) {
	app := NewApp()

	// Test that color cycling adds colors to recent
	app.cyclePenColor()
	recentColors := app.GetColorPicker().GetRecentColors()
	if len(recentColors) != 1 {
		t.Errorf("Expected 1 recent color after cycling, got %d", len(recentColors))
	}

	// Cycle a few more times
	app.cyclePenColor()
	app.cyclePenColor()
	recentColors = app.GetColorPicker().GetRecentColors()
	if len(recentColors) != 3 {
		t.Errorf("Expected 3 recent colors after cycling 3 times, got %d", len(recentColors))
	}
}

func TestColorSelectedMsgAddsToRecent(t *testing.T) {
	app := NewApp()

	// Simulate a color being selected from the color picker
	colorMsg := ColorSelectedMsg{Color: "yellow", Name: "Yellow"}
	app.Update(colorMsg)

	recentColors := app.GetColorPicker().GetRecentColors()
	if len(recentColors) != 1 {
		t.Errorf("Expected 1 recent color after ColorSelectedMsg, got %d", len(recentColors))
	}
	if recentColors[0].TerminalCode != "yellow" {
		t.Errorf("Expected recent color to be 'yellow', got '%s'", recentColors[0].TerminalCode)
	}

	// Test that the pen color was also updated
	if app.cursor.Pen.Foreground != "yellow" {
		t.Errorf("Expected pen color to be 'yellow', got '%s'", app.cursor.Pen.Foreground)
	}
}

func TestHSLColorPickerRecentColorSelection(t *testing.T) {
	picker := NewHSLColorPicker()

	// Add some recent colors
	picker.AddRecentColorByTerminalCode("red")
	picker.AddRecentColorByTerminalCode("blue")

	// Switch to recent mode and select a color
	picker.Show()
	tabMsg := tea.KeyMsg{Type: tea.KeyTab} // Switch to recent mode
	picker.Update(tabMsg)

	// Select the first recent color (blue, since it's most recent)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	cmd := picker.Update(enterMsg)

	// The picker should be hidden after selection
	if picker.IsVisible() {
		t.Error("Expected picker to be hidden after selection")
	}

	// The command should return a ColorSelectedMsg
	if cmd == nil {
		t.Error("Expected command to be returned after color selection")
	} else {
		// Execute the command to get the message
		msg := cmd()
		if colorMsg, ok := msg.(ColorSelectedMsg); ok {
			if colorMsg.Color != "blue" {
				t.Errorf("Expected selected color to be 'blue', got '%s'", colorMsg.Color)
			}
		} else {
			t.Error("Expected ColorSelectedMsg from color selection command")
		}
	}
}