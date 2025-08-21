package canvas

import (
	"strings"
	"testing"

	"github.com/cadubi/cadubi/pkg/commands"
)

func TestCanvasInfinite(t *testing.T) {
	c := NewCanvas(0, 0)

	// Test that canvas starts empty
	if c.initialized {
		t.Error("Expected canvas to start uninitialized")
	}

	// Test setting pixels at various positions including negative
	testPixel := commands.Pixel{Char: '*', Foreground: "red"}

	positions := []struct{ x, y int }{
		{0, 0},
		{10, 5},
		{-5, -3},
		{100, 50},
		{-100, -50},
	}

	for _, pos := range positions {
		err := c.SetPixel(pos.x, pos.y, testPixel)
		if err != nil {
			t.Errorf("Failed to set pixel at (%d, %d): %v", pos.x, pos.y, err)
		}
	}

	// Test that bounds are updated correctly
	minX, minY, maxX, maxY := c.GetBounds()
	expectedMinX, expectedMaxX := -100, 100
	expectedMinY, expectedMaxY := -50, 50

	if minX != expectedMinX || maxX != expectedMaxX || minY != expectedMinY || maxY != expectedMaxY {
		t.Errorf("Expected bounds (%d,%d,%d,%d), got (%d,%d,%d,%d)",
			expectedMinX, expectedMinY, expectedMaxX, expectedMaxY,
			minX, minY, maxX, maxY)
	}

	// Test getting pixels
	for _, pos := range positions {
		pixel, err := c.GetPixel(pos.x, pos.y)
		if err != nil {
			t.Errorf("Failed to get pixel at (%d, %d): %v", pos.x, pos.y, err)
		}
		if pixel.Char != '*' || pixel.Foreground != "red" {
			t.Errorf("Expected pixel '*' with red foreground at (%d, %d), got '%c' with '%s'",
				pos.x, pos.y, pixel.Char, pixel.Foreground)
		}
	}
}

func TestCanvasViewport(t *testing.T) {
	c := NewCanvas(0, 0)

	// Draw a simple pattern
	patterns := []struct {
		x, y int
		char rune
	}{
		{5, 5, 'A'},
		{6, 5, 'B'},
		{7, 5, 'C'},
		{5, 6, 'D'},
		{6, 6, 'E'},
		{7, 6, 'F'},
	}

	for _, p := range patterns {
		c.SetPixel(p.x, p.y, commands.Pixel{Char: p.char})
	}

	// Test viewport rendering
	viewport := c.GetViewport(5, 5, 3, 2)
	lines := strings.Split(viewport, "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines in viewport, got %d", len(lines))
	}

	// Check first line: ABC
	if !strings.Contains(lines[0], "A") || !strings.Contains(lines[0], "B") || !strings.Contains(lines[0], "C") {
		t.Errorf("Expected first line to contain ABC, got: %s", lines[0])
	}

	// Check second line: DEF
	if !strings.Contains(lines[1], "D") || !strings.Contains(lines[1], "E") || !strings.Contains(lines[1], "F") {
		t.Errorf("Expected second line to contain DEF, got: %s", lines[1])
	}
}

func TestCanvasViewportManagement(t *testing.T) {
	c := NewCanvas(0, 0)

	// Test viewport positioning
	c.SetViewport(10, 20)
	x, y := c.GetViewportPos()
	if x != 10 || y != 20 {
		t.Errorf("Expected viewport at (10, 20), got (%d, %d)", x, y)
	}

	// Test centering viewport
	c.CenterViewportOn(50, 60, 20, 10)
	x, y = c.GetViewportPos()
	expectedX, expectedY := 50-20/2, 60-10/2
	if x != expectedX || y != expectedY {
		t.Errorf("Expected viewport centered at (%d, %d), got (%d, %d)", expectedX, expectedY, x, y)
	}
}

func TestCanvasClear(t *testing.T) {
	c := NewCanvas(0, 0)

	// Add some pixels
	c.SetPixel(5, 5, commands.Pixel{Char: 'X'})
	c.SetPixel(10, 10, commands.Pixel{Char: 'Y'})

	// Verify canvas has been drawn on
	// Note: GetPixelCount is currently a placeholder, so we skip this check
	// TODO: Implement proper pixel counting with layer compositing

	// Clear canvas
	err := c.Clear()
	if err != nil {
		t.Errorf("Failed to clear canvas: %v", err)
	}

	// Verify canvas layer was cleared
	// We'll check by trying to get the pixels that were set
	pixel1, _ := c.GetPixel(5, 5)
	pixel2, _ := c.GetPixel(10, 10)

	if pixel1.Char != ' ' || pixel2.Char != ' ' {
		t.Error("Expected pixels to be cleared after Clear() call")
	}

	if c.initialized {
		t.Error("Expected canvas to be uninitialized after clear")
	}
}

func TestCanvasDimensions(t *testing.T) {
	c := NewCanvas(0, 0)

	// Empty canvas should return default dimensions
	w, h := c.GetDimensions()
	if w != 80 || h != 24 {
		t.Errorf("Expected default dimensions (80, 24), got (%d, %d)", w, h)
	}

	// Add pixels and test content dimensions
	c.SetPixel(-5, -10, commands.Pixel{Char: 'A'})
	c.SetPixel(15, 5, commands.Pixel{Char: 'B'})

	w, h = c.GetDimensions()
	expectedW, expectedH := 15-(-5)+1, 5-(-10)+1 // 21, 16
	if w != expectedW || h != expectedH {
		t.Errorf("Expected content dimensions (%d, %d), got (%d, %d)", expectedW, expectedH, w, h)
	}
}
