package canvas

import (
	"strings"
	"sync"

	"github.com/cadubi/cadubi/pkg/commands"
	"github.com/cadubi/cadubi/pkg/unicode"
	"github.com/charmbracelet/lipgloss"
)

// Canvas represents the drawing surface with infinite dimensions and layers
type Canvas struct {
	layers      *LayerStack // Layer management
	minX, minY  int         // Actual content bounds
	maxX, maxY  int         // Actual content bounds
	viewportX   int         // Current viewport position
	viewportY   int         // Current viewport position
	mutex       sync.RWMutex
	initialized bool   // Track if bounds are initialized
	modified    bool   // Track if canvas has been modified since last save
	filename    string // Current filename (empty if unsaved)
}

// NewCanvas creates a new infinite canvas with layers
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		layers:      NewLayerStack(),
		minX:        0,
		minY:        0,
		maxX:        0,
		maxY:        0,
		viewportX:   0,
		viewportY:   0,
		initialized: false,
		modified:    false,
		filename:    "",
	}
}

// GetPixel gets the composited pixel at the specified coordinates
func (c *Canvas) GetPixel(x, y int) (commands.Pixel, error) {
	// Return the composited pixel from all layers
	return c.layers.GetCompositePixel(x, y), nil
}

// SetPixel sets a pixel at the specified coordinates on the current layer
func (c *Canvas) SetPixel(x, y int, pixel commands.Pixel) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Update bounds for infinite canvas
	if !c.initialized {
		c.minX, c.minY = x, y
		c.maxX, c.maxY = x, y
		c.initialized = true
	} else {
		if x < c.minX {
			c.minX = x
		}
		if x > c.maxX {
			c.maxX = x
		}
		if y < c.minY {
			c.minY = y
		}
		if y > c.maxY {
			c.maxY = y
		}
	}

	// Set pixel on current layer
	currentLayer := c.layers.GetCurrentLayer()
	err := c.layers.SetPixel(currentLayer, x, y, pixel)
	if err == nil {
		c.modified = true
	}
	return err
}

// GetDimensions returns the actual content dimensions
func (c *Canvas) GetDimensions() (width, height int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.initialized {
		return 80, 24 // Default dimensions
	}

	return c.maxX - c.minX + 1, c.maxY - c.minY + 1
}

// Clear clears the current layer
func (c *Canvas) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	currentLayer := c.layers.GetCurrentLayer()
	err := c.layers.ClearLayer(currentLayer)
	if err != nil {
		return err
	}

	// Reset bounds
	c.initialized = false
	c.minX, c.minY = 0, 0
	c.maxX, c.maxY = 0, 0
	return nil
}

// GetViewport returns a string representation of a viewport region
func (c *Canvas) GetViewport(startX, startY, width, height int) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var lines []string

	for y := startY; y < startY+height; y++ {
		var line []string

		for x := startX; x < startX+width; x++ {
			// Get composited pixel from all layers
			pixel := c.layers.GetCompositePixel(x, y)

			if pixel.Char != ' ' {
				// Create styled character using Lipgloss
				style := lipgloss.NewStyle()

				if pixel.Foreground != "" {
					style = style.Foreground(lipgloss.Color(pixel.Foreground))
				}
				if pixel.Background != "" {
					style = style.Background(lipgloss.Color(pixel.Background))
				}
				if pixel.Bold {
					style = style.Bold(true)
				}
				if pixel.Italic {
					style = style.Italic(true)
				}
				if pixel.Underline {
					style = style.Underline(true)
				}

				styledChar := style.Render(string(pixel.Char))
				line = append(line, styledChar)
				
				// For wide characters (like emoji), skip the next column to avoid double-rendering
				if unicode.IsWideChar(pixel.Char) {
					x++ // Skip next column since wide char occupies 2 cells
				}
			} else {
				line = append(line, " ")
			}
		}

		lines = append(lines, strings.Join(line, ""))
	}

	return strings.Join(lines, "\n")
}

// SetViewport sets the viewport position
func (c *Canvas) SetViewport(x, y int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.viewportX = x
	c.viewportY = y
}

// GetViewportPos returns the current viewport position
func (c *Canvas) GetViewportPos() (x, y int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.viewportX, c.viewportY
}

// CenterViewportOn centers the viewport on the given coordinates
func (c *Canvas) CenterViewportOn(x, y, viewWidth, viewHeight int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.viewportX = x - viewWidth/2
	c.viewportY = y - viewHeight/2
}

// GetViewportCentered returns viewport content centered on current viewport position
func (c *Canvas) GetViewportCentered(width, height int) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// The viewportX/Y from CenterViewportOn represents the top-left corner of the view
	startX := c.viewportX
	startY := c.viewportY

	return c.getViewportInternal(startX, startY, width, height)
}

// getViewportInternal is the internal viewport rendering method (assumes lock held)
func (c *Canvas) getViewportInternal(startX, startY, width, height int) string {
	var lines []string

	for y := startY; y < startY+height; y++ {
		var line []string

		for x := startX; x < startX+width; x++ {
			// Get composited pixel from all layers
			pixel := c.layers.GetCompositePixel(x, y)

			if pixel.Char != ' ' {
				// Create styled character using Lipgloss
				style := lipgloss.NewStyle()

				if pixel.Foreground != "" {
					style = style.Foreground(lipgloss.Color(pixel.Foreground))
				}
				if pixel.Background != "" {
					style = style.Background(lipgloss.Color(pixel.Background))
				}
				if pixel.Bold {
					style = style.Bold(true)
				}
				if pixel.Italic {
					style = style.Italic(true)
				}
				if pixel.Underline {
					style = style.Underline(true)
				}

				styledChar := style.Render(string(pixel.Char))
				line = append(line, styledChar)
				
				// For wide characters (like emoji), skip the next column to avoid double-rendering
				if unicode.IsWideChar(pixel.Char) {
					x++ // Skip next column since wide char occupies 2 cells
				}
			} else {
				line = append(line, " ")
			}
		}

		lines = append(lines, strings.Join(line, ""))
	}

	return strings.Join(lines, "\n")
}

// GetBounds returns the actual bounds of drawn content
func (c *Canvas) GetBounds() (minX, minY, maxX, maxY int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.initialized {
		return 0, 0, 0, 0
	}

	return c.minX, c.minY, c.maxX, c.maxY
}

// GetPixelCount returns the number of non-empty pixels across all layers
func (c *Canvas) GetPixelCount() int {
	// This is an approximation - we'd need to iterate through all coordinates
	// to get the exact count after compositing. For now, return a simple count.

	// TODO: Implement proper pixel counting with layer compositing
	return 42 // Placeholder - should iterate through actual drawn area
}

// Layer management methods

// GetCurrentLayer returns the current layer index
func (c *Canvas) GetCurrentLayer() int {
	return c.layers.GetCurrentLayer()
}

// SetCurrentLayer sets the current layer
func (c *Canvas) SetCurrentLayer(layer int) error {
	return c.layers.SetCurrentLayer(layer)
}

// GetLayerCount returns the number of layers
func (c *Canvas) GetLayerCount() int {
	return c.layers.GetLayerCount()
}

// AddLayer adds a new layer
func (c *Canvas) AddLayer(name string) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	layerIndex, err := c.layers.AddLayer(name)
	if err == nil {
		c.modified = true
	}
	return layerIndex, err
}

// RemoveLayer removes a layer
func (c *Canvas) RemoveLayer(layer int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	err := c.layers.RemoveLayer(layer)
	if err == nil {
		c.modified = true
	}
	return err
}

// SetLayerVisible sets layer visibility
func (c *Canvas) SetLayerVisible(layer int, visible bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	err := c.layers.SetLayerVisible(layer, visible)
	if err == nil {
		c.modified = true
	}
	return err
}

// IsLayerVisible checks if a layer is visible
func (c *Canvas) IsLayerVisible(layer int) bool {
	return c.layers.IsLayerVisible(layer)
}

// GetLayerNames returns formatted layer names with status indicators
func (c *Canvas) GetLayerNames() []string {
	return c.layers.GetLayerNames()
}

// File state management methods

// IsModified returns true if the canvas has been modified since last save
func (c *Canvas) IsModified() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.modified
}

// SetModified sets the modified state of the canvas
func (c *Canvas) SetModified(modified bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.modified = modified
}

// GetFilename returns the current filename
func (c *Canvas) GetFilename() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.filename
}

// SetFilename sets the current filename
func (c *Canvas) SetFilename(filename string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.filename = filename
}

// HasFilename returns true if the canvas has a filename
func (c *Canvas) HasFilename() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.filename != ""
}
