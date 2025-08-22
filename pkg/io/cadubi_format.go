package io

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cadubi/cadubi/pkg/canvas"
	"github.com/cadubi/cadubi/pkg/commands"
)

// CadubiDocument represents the complete file format for .cadubi files
type CadubiDocument struct {
	Version string        `json:"version"`
	Canvas  CadubiCanvas  `json:"canvas"`
	Layers  []CadubiLayer `json:"layers"`
}

// CadubiCanvas represents the canvas state
type CadubiCanvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	MinX   int `json:"min_x"`
	MinY   int `json:"min_y"`
	MaxX   int `json:"max_x"`
	MaxY   int `json:"max_y"`
}

// CadubiLayer represents a single layer
type CadubiLayer struct {
	Name    string               `json:"name"`
	Visible bool                 `json:"visible"`
	Pixels  map[string]CadubiPixel `json:"pixels"` // "x,y" -> pixel
}

// CadubiPixel represents a single pixel/character
type CadubiPixel struct {
	Char       string `json:"char"`
	Foreground string `json:"foreground"`
	Background string `json:"background"`
	Bold       bool   `json:"bold"`
	Italic     bool   `json:"italic"`
	Underline  bool   `json:"underline"`
}

// SaveToFile saves a canvas and layer data to a .cadubi file
func SaveToFile(filename string, c *canvas.Canvas) error {
	doc := CadubiDocument{
		Version: "2.0",
	}

	// Get canvas bounds
	minX, minY, maxX, maxY := c.GetBounds()
	// Get canvas dimensions
	width, height := c.GetDimensions()
	doc.Canvas = CadubiCanvas{
		Width:  width,
		Height: height,
		MinX:   minX,
		MinY:   minY,
		MaxX:   maxX,
		MaxY:   maxY,
	}

	// Export all layers
	layerCount := c.GetLayerCount()
	doc.Layers = make([]CadubiLayer, layerCount)

	for layerIndex := 0; layerIndex < layerCount; layerIndex++ {
		layer := CadubiLayer{
			Name:    getLayerName(c, layerIndex),
			Visible: c.IsLayerVisible(layerIndex),
			Pixels:  make(map[string]CadubiPixel),
		}

		// Export all pixels in the layer's bounds
		for y := minY; y <= maxY; y++ {
			for x := minX; x <= maxX; x++ {
				// Temporarily switch to this layer to get its pixel
				originalLayer := c.GetCurrentLayer()
				c.SetCurrentLayer(layerIndex)

				pixel, err := c.GetPixel(x, y)
				if err == nil && pixel.Char != ' ' {
					key := pixelKey(x, y)
					layer.Pixels[key] = CadubiPixel{
						Char:       string(pixel.Char),
						Foreground: pixel.Foreground,
						Background: pixel.Background,
						Bold:       pixel.Bold,
						Italic:     pixel.Italic,
						Underline:  pixel.Underline,
					}
				}

				// Restore original layer
				c.SetCurrentLayer(originalLayer)
			}
		}

		doc.Layers[layerIndex] = layer
	}

	// Marshal to JSON with proper indentation
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}

// LoadFromFile loads a .cadubi file into a canvas
func LoadFromFile(filename string) (*canvas.Canvas, []LayerInfo, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}

	var doc CadubiDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, nil, err
	}

	// Create canvas with original dimensions
	c := canvas.NewCanvas(doc.Canvas.Width, doc.Canvas.Height)

	// Load all layers
	var layerInfos []LayerInfo
	for i, layerData := range doc.Layers {
		// Add layer if it doesn't exist
		if i >= c.GetLayerCount() {
			c.AddLayer(layerData.Name)
		}

		// Set layer visibility
		c.SetLayerVisible(i, layerData.Visible)

		// Load pixels for this layer
		c.SetCurrentLayer(i)
		for posKey, pixelData := range layerData.Pixels {
			x, y := parsePixelKey(posKey)
			
			pixel := commands.Pixel{
				Char:       rune(pixelData.Char[0]), // Take first character
				Foreground: pixelData.Foreground,
				Background: pixelData.Background,
				Bold:       pixelData.Bold,
				Italic:     pixelData.Italic,
				Underline:  pixelData.Underline,
			}

			// Set pixel directly on canvas
			c.SetPixel(x, y, pixel)
		}

		layerInfos = append(layerInfos, LayerInfo{
			Name:    layerData.Name,
			Visible: layerData.Visible,
		})
	}

	return c, layerInfos, nil
}

// LayerInfo contains basic layer information
type LayerInfo struct {
	Name    string
	Visible bool
}

// Helper functions
func pixelKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

func parsePixelKey(key string) (x, y int) {
	fmt.Sscanf(key, "%d,%d", &x, &y)
	return x, y
}

func getLayerName(c *canvas.Canvas, layerIndex int) string {
	// This is a placeholder - canvas doesn't currently store layer names
	// TODO: Add layer names to canvas implementation
	return fmt.Sprintf("Layer %d", layerIndex+1)
}