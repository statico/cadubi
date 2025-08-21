package commands

import (
	"fmt"
	"testing"
)

// MockCanvas implements CanvasManager for testing
type MockCanvas struct {
	pixels map[string]Pixel // key: "x,y"
	width  int
	height int
}

func NewMockCanvas(width, height int) *MockCanvas {
	return &MockCanvas{
		pixels: make(map[string]Pixel),
		width:  width,
		height: height,
	}
}

func (mc *MockCanvas) GetPixel(x, y int) (Pixel, error) {
	key := fmt.Sprintf("%d,%d", x, y)
	if pixel, exists := mc.pixels[key]; exists {
		return pixel, nil
	}
	return Pixel{Char: ' '}, nil
}

func (mc *MockCanvas) SetPixel(x, y int, pixel Pixel) error {
	key := fmt.Sprintf("%d,%d", x, y)
	mc.pixels[key] = pixel
	return nil
}

func (mc *MockCanvas) GetDimensions() (width, height int) {
	return mc.width, mc.height
}

func (mc *MockCanvas) Clear() error {
	mc.pixels = make(map[string]Pixel)
	return nil
}

// MockLayerManager implements LayerManager for testing
type MockLayerManager struct {
	currentLayer int
	layers       []string
}

func NewMockLayerManager() *MockLayerManager {
	return &MockLayerManager{
		currentLayer: 0,
		layers:       []string{"Layer 1"},
	}
}

func (ml *MockLayerManager) GetCurrentLayer() int {
	return ml.currentLayer
}

func (ml *MockLayerManager) SetCurrentLayer(layer int) error {
	if layer >= 0 && layer < len(ml.layers) {
		ml.currentLayer = layer
	}
	return nil
}

func (ml *MockLayerManager) GetLayerCount() int {
	return len(ml.layers)
}

func (ml *MockLayerManager) AddLayer(name string) (int, error) {
	ml.layers = append(ml.layers, name)
	return len(ml.layers) - 1, nil
}

func (ml *MockLayerManager) RemoveLayer(layer int) error {
	if layer >= 0 && layer < len(ml.layers) {
		ml.layers = append(ml.layers[:layer], ml.layers[layer+1:]...)
		if ml.currentLayer >= len(ml.layers) {
			ml.currentLayer = len(ml.layers) - 1
		}
	}
	return nil
}

func (ml *MockLayerManager) SetLayerVisible(layer int, visible bool) error {
	return nil
}

func (ml *MockLayerManager) IsLayerVisible(layer int) bool {
	return true
}

func TestDrawCommand(t *testing.T) {
	canvas := NewMockCanvas(80, 24)
	ctx := &CommandContext{
		Canvas: canvas,
		Layers: NewMockLayerManager(),
		Cursor: CursorState{X: 0, Y: 0},
	}

	// Create a draw command
	pixel := Pixel{Char: '*', Foreground: "red"}
	cmd := NewDrawCommand(10, 5, pixel, 0)

	// Test command properties
	if cmd.Description() == "" {
		t.Error("Expected command to have a description")
	}

	if cmd.Timestamp().IsZero() {
		t.Error("Expected command to have a timestamp")
	}

	// Execute the command
	err := cmd.Execute(ctx)
	if err != nil {
		t.Fatalf("Failed to execute draw command: %v", err)
	}

	// Verify pixel was drawn
	drawnPixel, err := canvas.GetPixel(10, 5)
	if err != nil {
		t.Fatalf("Failed to get pixel: %v", err)
	}

	if drawnPixel.Char != '*' || drawnPixel.Foreground != "red" {
		t.Errorf("Expected pixel to be '*' with red foreground, got '%c' with '%s'",
			drawnPixel.Char, drawnPixel.Foreground)
	}

	// Undo the command
	err = cmd.Undo(ctx)
	if err != nil {
		t.Fatalf("Failed to undo draw command: %v", err)
	}

	// Verify pixel was restored
	restoredPixel, err := canvas.GetPixel(10, 5)
	if err != nil {
		t.Fatalf("Failed to get restored pixel: %v", err)
	}

	if restoredPixel.Char != ' ' {
		t.Errorf("Expected restored pixel to be ' ', got '%c'", restoredPixel.Char)
	}
}

func TestCommandHistory(t *testing.T) {
	canvas := NewMockCanvas(80, 24)
	ctx := &CommandContext{
		Canvas: canvas,
		Layers: NewMockLayerManager(),
		Cursor: CursorState{X: 0, Y: 0},
	}

	history := NewCommandHistory(100)

	// Test initial state
	if history.CanUndo() {
		t.Error("Expected new history to have no commands to undo")
	}

	if history.CanRedo() {
		t.Error("Expected new history to have no commands to redo")
	}

	// Execute a command
	cmd1 := NewDrawCommand(5, 5, Pixel{Char: 'A'}, 0)
	err := history.Execute(ctx, cmd1)
	if err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check pixel was drawn
	pixel, _ := canvas.GetPixel(5, 5)
	if pixel.Char != 'A' {
		t.Error("Expected pixel to be 'A' after command execution")
	}

	// Execute another command
	cmd2 := NewDrawCommand(6, 6, Pixel{Char: 'B'}, 0)
	err = history.Execute(ctx, cmd2)
	if err != nil {
		t.Fatalf("Failed to execute second command: %v", err)
	}

	// Test undo
	if !history.CanUndo() {
		t.Error("Expected history to have commands to undo")
	}

	err = history.Undo(ctx)
	if err != nil {
		t.Fatalf("Failed to undo: %v", err)
	}

	// Check second pixel was undone
	pixel, _ = canvas.GetPixel(6, 6)
	if pixel.Char != ' ' {
		t.Error("Expected second pixel to be undone")
	}

	// First pixel should still be there
	pixel, _ = canvas.GetPixel(5, 5)
	if pixel.Char != 'A' {
		t.Error("Expected first pixel to remain after undo")
	}

	// Test redo
	if !history.CanRedo() {
		t.Error("Expected history to have commands to redo")
	}

	err = history.Redo(ctx)
	if err != nil {
		t.Fatalf("Failed to redo: %v", err)
	}

	// Check second pixel was redrawn
	pixel, _ = canvas.GetPixel(6, 6)
	if pixel.Char != 'B' {
		t.Error("Expected second pixel to be redrawn after redo")
	}

	// Test history display
	historyStrings := history.GetHistory()
	if len(historyStrings) != 2 {
		t.Errorf("Expected 2 commands in history, got %d", len(historyStrings))
	}

	// Test stats
	total, position := history.GetStats()
	if total != 2 || position != 2 {
		t.Errorf("Expected stats (2, 2), got (%d, %d)", total, position)
	}
}
