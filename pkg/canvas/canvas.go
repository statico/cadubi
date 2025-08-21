package canvas

import (
	"fmt"
	"strings"
	"sync"

	"github.com/cadubi/cadubi/pkg/commands"
	"github.com/charmbracelet/lipgloss"
)

// Canvas represents the drawing surface
type Canvas struct {
	pixels map[string]commands.Pixel // key: "x,y"
	width  int
	height int
	mutex  sync.RWMutex
}

// NewCanvas creates a new canvas
func NewCanvas(width, height int) *Canvas {
	return &Canvas{
		pixels: make(map[string]commands.Pixel),
		width:  width,
		height: height,
	}
}

// GetPixel gets a pixel at the specified coordinates
func (c *Canvas) GetPixel(x, y int) (commands.Pixel, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	key := fmt.Sprintf("%d,%d", x, y)
	if pixel, exists := c.pixels[key]; exists {
		return pixel, nil
	}
	
	// Return empty pixel if not set
	return commands.Pixel{Char: ' '}, nil
}

// SetPixel sets a pixel at the specified coordinates
func (c *Canvas) SetPixel(x, y int, pixel commands.Pixel) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Expand canvas if needed
	if x >= c.width {
		c.width = x + 1
	}
	if y >= c.height {
		c.height = y + 1
	}
	
	key := fmt.Sprintf("%d,%d", x, y)
	c.pixels[key] = pixel
	return nil
}

// GetDimensions returns the canvas dimensions
func (c *Canvas) GetDimensions() (width, height int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.width, c.height
}

// Clear clears the entire canvas
func (c *Canvas) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.pixels = make(map[string]commands.Pixel)
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
			key := fmt.Sprintf("%d,%d", x, y)
			if pixel, exists := c.pixels[key]; exists && pixel.Char != ' ' {
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
				
				line = append(line, style.Render(string(pixel.Char)))
			} else {
				line = append(line, " ")
			}
		}
		
		lines = append(lines, strings.Join(line, ""))
	}
	
	return strings.Join(lines, "\n")
}

// Resize changes the canvas dimensions
func (c *Canvas) Resize(width, height int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.width = width
	c.height = height
}

// GetBounds returns the actual bounds of drawn content
func (c *Canvas) GetBounds() (minX, minY, maxX, maxY int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	if len(c.pixels) == 0 {
		return 0, 0, 0, 0
	}
	
	minX, minY = c.width, c.height
	maxX, maxY = 0, 0
	
	for key := range c.pixels {
		var x, y int
		fmt.Sscanf(key, "%d,%d", &x, &y)
		
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	
	return minX, minY, maxX, maxY
}

// GetPixelCount returns the number of non-empty pixels
func (c *Canvas) GetPixelCount() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	count := 0
	for _, pixel := range c.pixels {
		if pixel.Char != ' ' {
			count++
		}
	}
	return count
}