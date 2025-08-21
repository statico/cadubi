package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAppInitialization(t *testing.T) {
	app := NewApp()

	// Test initial state
	if app.canvas == nil {
		t.Error("Expected app to have a canvas")
	}

	if app.history == nil {
		t.Error("Expected app to have a command history")
	}

	if app.layers == nil {
		t.Error("Expected app to have a layer manager")
	}

	if app.currentTool != "pen" {
		t.Errorf("Expected initial tool to be 'pen', got '%s'", app.currentTool)
	}

	if app.cursor.Pen.Char != '*' {
		t.Errorf("Expected initial pen character to be '*', got '%c'", app.cursor.Pen.Char)
	}

	if app.cursor.Pen.Foreground != "white" && app.cursor.Pen.Foreground != "" {
		t.Errorf("Expected initial pen color to be empty or white, got '%s'", app.cursor.Pen.Foreground)
	}
}

func TestStatusBarFormatting(t *testing.T) {
	app := NewApp()

	// Set window size
	app.width = 80
	app.height = 24

	// Test view rendering
	view := app.View()

	// Check for Unicode symbols from our enhancement
	if !strings.Contains(view, "■ Tool:") {
		t.Error("Expected status bar to contain tool indicator with ■ symbol")
	}

	if !strings.Contains(view, "●●●") {
		t.Error("Expected status bar to contain color swatch (●●●)")
	}

	// Check for layer indicators - should show ⧉ for visible layers or ○ for hidden
	if !strings.Contains(view, "⧉") && !strings.Contains(view, "○") {
		t.Error("Expected status bar to contain layer indicator (⧉ for visible or ○ for hidden)")
	}
}

func TestColorPicker(t *testing.T) {
	app := NewApp()

	// Test initial color
	initialColor := app.cursor.Pen.Foreground

	// Test setPenColor
	app.setPenColor("red")
	if app.cursor.Pen.Foreground != "red" {
		t.Errorf("Expected pen color to be 'red', got '%s'", app.cursor.Pen.Foreground)
	}

	if app.ctx.Cursor.Pen.Foreground != "red" {
		t.Error("Expected context cursor to be updated with new color")
	}

	// Test cyclePenColor
	app.cyclePenColor()
	if app.cursor.Pen.Foreground == "red" {
		t.Error("Expected pen color to change after cycling")
	}

	// Test multiple cycles to ensure it loops
	colors := make([]string, 10)
	for i := 0; i < 10; i++ {
		colors[i] = app.cursor.Pen.Foreground
		app.cyclePenColor()
	}

	// Should have looped back to the first color
	if colors[0] != colors[8] {
		t.Error("Expected color cycling to loop back after 8 colors")
	}

	// Restore initial color
	if initialColor != "" {
		app.setPenColor(initialColor)
	}
}

func TestKeyHandling(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	tests := []struct {
		key      string
		testFunc func(*testing.T, *App)
	}{
		{"*", func(t *testing.T, a *App) {
			if a.cursor.Pen.Char != '*' {
				t.Errorf("Expected pen char to be '*', got '%c'", a.cursor.Pen.Char)
			}
		}},
		{"#", func(t *testing.T, a *App) {
			if a.cursor.Pen.Char != '#' {
				t.Errorf("Expected pen char to be '#', got '%c'", a.cursor.Pen.Char)
			}
		}},
		{"w", func(t *testing.T, a *App) {
			if a.cursor.Pen.Foreground != "white" {
				t.Errorf("Expected pen color to be 'white', got '%s'", a.cursor.Pen.Foreground)
			}
		}},
		{"g", func(t *testing.T, a *App) {
			if a.cursor.Pen.Foreground != "green" {
				t.Errorf("Expected pen color to be 'green', got '%s'", a.cursor.Pen.Foreground)
			}
		}},
	}

	for _, test := range tests {
		t.Run("key_"+test.key, func(t *testing.T) {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
			_, _ = app.Update(keyMsg)
			test.testFunc(t, app)
		})
	}
}

func TestCursorMovement(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	initialX, initialY := app.cursor.X, app.cursor.Y

	// Test vi-style movement
	movements := []struct {
		key                    string
		expectedDx, expectedDy int
	}{
		{"h", -1, 0}, // left
		{"j", 0, 1},  // down
		{"k", 0, -1}, // up
		{"l", 1, 0},  // right
	}

	for _, movement := range movements {
		t.Run("movement_"+movement.key, func(t *testing.T) {
			startX, startY := app.cursor.X, app.cursor.Y

			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(movement.key)}
			_, _ = app.Update(keyMsg)

			expectedX := startX + movement.expectedDx
			expectedY := startY + movement.expectedDy

			// Account for bounds checking
			if expectedX < 0 {
				expectedX = 0
			}
			if expectedY < 0 {
				expectedY = 0
			}
			if expectedX > 200 {
				expectedX = 200
			}
			if expectedY > 100 {
				expectedY = 100
			}

			if app.cursor.X != expectedX || app.cursor.Y != expectedY {
				t.Errorf("Expected cursor to move to (%d,%d), got (%d,%d)",
					expectedX, expectedY, app.cursor.X, app.cursor.Y)
			}
		})
	}

	// Reset cursor position
	app.cursor.X, app.cursor.Y = initialX, initialY
}

func TestHelpModal(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Initially help should be hidden
	if app.showHelp {
		t.Error("Expected help modal to be initially hidden")
	}

	// Toggle help with '?'
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	_, _ = app.Update(keyMsg)

	if !app.showHelp {
		t.Error("Expected help modal to be shown after pressing '?'")
	}

	// Test view contains help content when modal is shown
	view := app.View()
	if !strings.Contains(view, "hjkl") && !strings.Contains(view, "space") {
		t.Error("Expected help modal view to contain key binding help content")
	}

	// Test escape key closes help
	escapeKey := tea.KeyMsg{Type: tea.KeyEscape}
	_, _ = app.Update(escapeKey)

	if app.showHelp {
		t.Error("Expected help modal to be hidden after pressing escape")
	}
}

func TestLayerManagement(t *testing.T) {
	app := NewApp()

	// Test initial layer state
	initialLayerCount := app.canvas.GetLayerCount()
	if initialLayerCount == 0 {
		t.Error("Expected at least one initial layer")
	}

	// Test layer switching (key '1')
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	_, _ = app.Update(keyMsg)

	if app.canvas.GetCurrentLayer() != 0 {
		t.Errorf("Expected current layer to be 0 after pressing '1', got %d", app.canvas.GetCurrentLayer())
	}

	// Test adding new layer (key 'n')
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	_, _ = app.Update(keyMsg)

	newLayerCount := app.canvas.GetLayerCount()
	if newLayerCount != initialLayerCount+1 {
		t.Errorf("Expected layer count to increase by 1, got %d -> %d", initialLayerCount, newLayerCount)
	}

	// Test layer visibility toggle (key 'v')
	currentLayer := app.canvas.GetCurrentLayer()
	initialVisibility := app.canvas.IsLayerVisible(currentLayer)

	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
	_, _ = app.Update(keyMsg)

	newVisibility := app.canvas.IsLayerVisible(currentLayer)
	if newVisibility == initialVisibility {
		t.Error("Expected layer visibility to toggle after pressing 'v'")
	}
}

func TestMouseInteraction(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Test mouse click
	mouseMsg := tea.MouseMsg{
		X:    10,
		Y:    5,
		Type: tea.MouseLeft,
	}

	initialX, initialY := app.cursor.X, app.cursor.Y
	_, _ = app.Update(mouseMsg)

	// Cursor should move to approximate mouse position
	// (exact calculation depends on screenToCanvas implementation)
	if app.cursor.X == initialX && app.cursor.Y == initialY {
		// This might be expected if the click was outside the canvas area
		t.Log("Mouse click did not move cursor (possibly outside canvas area)")
	}
}

func TestDrawingAndErasing(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24

	// Set cursor to known position
	app.cursor.X = 5
	app.cursor.Y = 5

	// Test drawing with space key
	spaceKey := tea.KeyMsg{Type: tea.KeySpace}
	_, _ = app.Update(spaceKey)

	// Check that a drawing command was executed (we can't easily verify the pixel
	// without exposing more internals, but we can check history)
	total, position := app.history.GetStats()
	if total == 0 {
		t.Error("Expected drawing action to add command to history")
	}

	// Test erasing with 'x' key
	xKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	_, _ = app.Update(xKey)

	// Check that an erase command was executed
	newTotal, _ := app.history.GetStats()
	if newTotal != total+1 {
		t.Error("Expected erasing action to add command to history")
	}

	// Test undo with 'u' key
	uKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
	_, _ = app.Update(uKey)

	// Check that history position changed
	_, newPosition := app.history.GetStats()
	if newPosition == position+1 {
		t.Error("Expected undo to change history position")
	}
}

func TestWindowResize(t *testing.T) {
	app := NewApp()

	// Test window resize
	resizeMsg := tea.WindowSizeMsg{Width: 120, Height: 30}
	_, _ = app.Update(resizeMsg)

	if app.width != 120 || app.height != 30 {
		t.Errorf("Expected app dimensions to be (120, 30), got (%d, %d)", app.width, app.height)
	}

	// Test that view renders correctly with new size
	view := app.View()
	if view == "" {
		t.Error("Expected view to render after resize")
	}
}

func TestHeadlessSequence(t *testing.T) {
	// This test demonstrates the headless testing approach mentioned in the plan
	app := NewApp()

	// Sequence of actions: resize window, move cursor, change pen, draw
	actions := []tea.Msg{
		tea.WindowSizeMsg{Width: 80, Height: 24},
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'*'}}, // Set pen to *
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}, // Set color to green
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}, // Move right
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}, // Move down
		tea.KeyMsg{Type: tea.KeySpace},                     // Draw
		tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}, // Undo
	}

	var model tea.Model = app

	// Execute sequence
	for i, action := range actions {
		var cmd tea.Cmd
		model, cmd = model.Update(action)
		_ = cmd // Ignore commands for this test

		// Verify state at each step
		appModel := model.(*App)
		view := appModel.View()

		if view == "" {
			t.Errorf("Expected non-empty view at step %d", i)
		}

		// At step 0 (resize), verify dimensions
		if i == 0 {
			if appModel.width != 80 || appModel.height != 24 {
				t.Errorf("Expected dimensions (80,24) at step %d, got (%d,%d)",
					i, appModel.width, appModel.height)
			}
		}

		// At step 1 (*), verify pen character
		if i == 1 {
			if appModel.cursor.Pen.Char != '*' {
				t.Errorf("Expected pen char '*' at step %d, got '%c'", i, appModel.cursor.Pen.Char)
			}
		}

		// At step 2 (g), verify color
		if i == 2 {
			if appModel.cursor.Pen.Foreground != "green" {
				t.Errorf("Expected pen color 'green' at step %d, got '%s'",
					i, appModel.cursor.Pen.Foreground)
			}
		}
	}

	t.Log("Headless sequence test completed successfully")
}

func TestGetCharDisplayWidth(t *testing.T) {
	app := NewApp()
	
	tests := []struct {
		char     rune
		expected int
		desc     string
	}{
		{' ', 1, "space"},
		{'A', 1, "regular ASCII letter"},
		{'*', 1, "regular ASCII symbol"},
		{'🤖', 2, "robot emoji"},
		{'🎉', 2, "party emoji"},
		{'🌟', 2, "star emoji"},
		{'中', 2, "CJK character"},
		{'あ', 2, "hiragana"},
		{'カ', 2, "katakana"},
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result := app.getCharDisplayWidth(test.char)
			if result != test.expected {
				t.Errorf("getCharDisplayWidth('%c') = %d, expected %d", test.char, result, test.expected)
			}
		})
	}
}

func TestScreenToCanvasNarrowChars(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24
	
	// Place some narrow characters on the canvas
	app.cursor.X = 0
	app.cursor.Y = 0
	
	// Draw some regular characters: "ABC"
	chars := []rune{'A', 'B', 'C'}
	for i, char := range chars {
		app.cursor.X = i
		app.cursor.Pen.Char = char
		app.drawAtCursor()
	}
	
	tests := []struct {
		screenX, screenY   int
		expectedX, expectedY int
		desc               string
	}{
		{2, 1, 0, 0, "click on first character A"},
		{3, 1, 1, 0, "click on second character B"},  
		{4, 1, 2, 0, "click on third character C"},
		{5, 1, 3, 0, "click after all characters"},
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			canvasX, canvasY := app.screenToCanvas(test.screenX, test.screenY)
			if canvasX != test.expectedX || canvasY != test.expectedY {
				t.Errorf("screenToCanvas(%d, %d) = (%d, %d), expected (%d, %d)", 
					test.screenX, test.screenY, canvasX, canvasY, test.expectedX, test.expectedY)
			}
		})
	}
}

func TestScreenToCanvasWideChars(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24
	
	// Place some wide characters on the canvas
	app.cursor.X = 0
	app.cursor.Y = 0
	
	// Draw some emoji: "🤖🎉🌟"
	emojis := []rune{'🤖', '🎉', '🌟'}
	for i, emoji := range emojis {
		app.cursor.X = i
		app.cursor.Pen.Char = emoji
		app.drawAtCursor()
	}
	
	tests := []struct {
		screenX, screenY   int
		expectedX, expectedY int
		desc               string
	}{
		{2, 1, 0, 0, "click on first emoji 🤖 (left side)"},
		{3, 1, 0, 0, "click on first emoji 🤖 (right side)"},
		{4, 1, 1, 0, "click on second emoji 🎉 (left side)"},
		{5, 1, 1, 0, "click on second emoji 🎉 (right side)"},
		{6, 1, 2, 0, "click on third emoji 🌟 (left side)"},
		{7, 1, 2, 0, "click on third emoji 🌟 (right side)"},
		{8, 1, 3, 0, "click after all emoji"},
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			canvasX, canvasY := app.screenToCanvas(test.screenX, test.screenY)
			if canvasX != test.expectedX || canvasY != test.expectedY {
				t.Errorf("screenToCanvas(%d, %d) = (%d, %d), expected (%d, %d)", 
					test.screenX, test.screenY, canvasX, canvasY, test.expectedX, test.expectedY)
			}
		})
	}
}

func TestScreenToCanvasMixedChars(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24
	
	// Place mixed narrow and wide characters: "A🤖B🎉C"
	app.cursor.Y = 0
	mixedChars := []struct {
		x    int
		char rune
	}{
		{0, 'A'},   // narrow: screen cols 0
		{1, '🤖'},  // wide: screen cols 1-2  
		{2, 'B'},   // narrow: screen col 3
		{3, '🎉'},  // wide: screen cols 4-5
		{4, 'C'},   // narrow: screen col 6
	}
	
	for _, item := range mixedChars {
		app.cursor.X = item.x
		app.cursor.Pen.Char = item.char
		app.drawAtCursor()
	}
	
	tests := []struct {
		screenX, screenY   int
		expectedX, expectedY int
		desc               string
	}{
		{2, 1, 0, 0, "click on A"},
		{3, 1, 1, 0, "click on 🤖 (left side)"},
		{4, 1, 1, 0, "click on 🤖 (right side)"},
		{5, 1, 2, 0, "click on B"},
		{6, 1, 3, 0, "click on 🎉 (left side)"},
		{7, 1, 3, 0, "click on 🎉 (right side)"},
		{8, 1, 4, 0, "click on C"},
		{9, 1, 5, 0, "click after C"},
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			canvasX, canvasY := app.screenToCanvas(test.screenX, test.screenY)
			if canvasX != test.expectedX || canvasY != test.expectedY {
				t.Errorf("screenToCanvas(%d, %d) = (%d, %d), expected (%d, %d)", 
					test.screenX, test.screenY, canvasX, canvasY, test.expectedX, test.expectedY)
			}
		})
	}
}

func TestScreenXToCanvasX(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24
	
	// Test the internal helper function directly
	app.cursor.Y = 5
	
	// Place characters: "🤖ABC🎉"  
	chars := []struct {
		x    int
		char rune
	}{
		{10, '🤖'}, // wide: screen cols 0-1
		{11, 'A'},  // narrow: screen col 2  
		{12, 'B'},  // narrow: screen col 3
		{13, 'C'},  // narrow: screen col 4
		{14, '🎉'}, // wide: screen cols 5-6
	}
	
	for _, item := range chars {
		app.cursor.X = item.x
		app.cursor.Pen.Char = item.char
		app.drawAtCursor()
	}
	
	canvasY := 5
	viewportX := 10
	canvasWidth := 70
	
	tests := []struct {
		screenRelativeX int
		expectedCanvasX int
		desc           string
	}{
		{0, 10, "screen pos 0 -> 🤖 at canvas 10"},
		{1, 10, "screen pos 1 -> 🤖 at canvas 10 (right side)"},
		{2, 11, "screen pos 2 -> A at canvas 11"},
		{3, 12, "screen pos 3 -> B at canvas 12"},
		{4, 13, "screen pos 4 -> C at canvas 13"},
		{5, 14, "screen pos 5 -> 🎉 at canvas 14"},
		{6, 14, "screen pos 6 -> 🎉 at canvas 14 (right side)"},
		{7, 15, "screen pos 7 -> after 🎉 at canvas 15"},
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result := app.screenXToCanvasX(test.screenRelativeX, canvasY, viewportX, canvasWidth)
			if result != test.expectedCanvasX {
				t.Errorf("screenXToCanvasX(%d, %d, %d, %d) = %d, expected %d",
					test.screenRelativeX, canvasY, viewportX, canvasWidth, result, test.expectedCanvasX)
			}
		})
	}
}

func TestMouseClickWideCharacters(t *testing.T) {
	app := NewApp()
	app.width = 80
	app.height = 24
	
	// Place an emoji at position (5, 5)
	app.cursor.X = 5
	app.cursor.Y = 5
	app.cursor.Pen.Char = '🤖'
	app.drawAtCursor()
	
	// Test mouse clicks around the emoji
	tests := []struct {
		mouseX, mouseY int
		expectedX, expectedY int
		desc string
	}{
		{7, 6, 5, 5, "click on left side of emoji"},   // screen pos (7,6) -> canvas (5,5)  
		{8, 6, 5, 5, "click on right side of emoji"},  // screen pos (8,6) -> canvas (5,5)
		{9, 6, 6, 5, "click after emoji"},             // screen pos (9,6) -> canvas (6,5)
	}
	
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			initialX, initialY := app.cursor.X, app.cursor.Y
			
			// Send mouse click
			mouseMsg := tea.MouseMsg{
				X:    test.mouseX,
				Y:    test.mouseY,
				Type: tea.MouseLeft,
			}
			_, _ = app.Update(mouseMsg)
			
			if app.cursor.X != test.expectedX || app.cursor.Y != test.expectedY {
				t.Errorf("Mouse click at (%d, %d) moved cursor to (%d, %d), expected (%d, %d)",
					test.mouseX, test.mouseY, app.cursor.X, app.cursor.Y, test.expectedX, test.expectedY)
			}
			
			// Reset cursor for next test
			app.cursor.X, app.cursor.Y = initialX, initialY
		})
	}
}
