package canvas

import (
	"fmt"
	"sync"

	"github.com/cadubi/cadubi/pkg/commands"
)

// Layer represents a single drawing layer
type Layer struct {
	ID        int
	Name      string
	Visible   bool
	Opacity   float64 // 0.0 to 1.0
	BlendMode string  // "normal", "multiply", "overlay", etc.
	pixels    map[string]commands.Pixel
	mutex     sync.RWMutex
}

// LayerStack manages multiple layers
type LayerStack struct {
	layers       []*Layer
	currentLayer int
	nextLayerID  int
	mutex        sync.RWMutex
}

// NewLayer creates a new layer
func NewLayer(id int, name string) *Layer {
	return &Layer{
		ID:        id,
		Name:      name,
		Visible:   true,
		Opacity:   1.0,
		BlendMode: "normal",
		pixels:    make(map[string]commands.Pixel),
	}
}

// NewLayerStack creates a new layer stack
func NewLayerStack() *LayerStack {
	// Create default layer
	defaultLayer := NewLayer(0, "Layer 1")

	return &LayerStack{
		layers:       []*Layer{defaultLayer},
		currentLayer: 0,
		nextLayerID:  1,
	}
}

// GetCurrentLayer returns the current layer index
func (ls *LayerStack) GetCurrentLayer() int {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()
	return ls.currentLayer
}

// SetCurrentLayer sets the current layer
func (ls *LayerStack) SetCurrentLayer(layer int) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if layer < 0 || layer >= len(ls.layers) {
		return fmt.Errorf("layer index %d out of range", layer)
	}

	ls.currentLayer = layer
	return nil
}

// GetLayerCount returns the number of layers
func (ls *LayerStack) GetLayerCount() int {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()
	return len(ls.layers)
}

// AddLayer adds a new layer
func (ls *LayerStack) AddLayer(name string) (int, error) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	newLayer := NewLayer(ls.nextLayerID, name)
	ls.layers = append(ls.layers, newLayer)
	layerIndex := len(ls.layers) - 1
	ls.nextLayerID++

	return layerIndex, nil
}

// RemoveLayer removes a layer
func (ls *LayerStack) RemoveLayer(layer int) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if layer < 0 || layer >= len(ls.layers) {
		return fmt.Errorf("layer index %d out of range", layer)
	}

	if len(ls.layers) <= 1 {
		return fmt.Errorf("cannot remove the last layer")
	}

	// Remove the layer
	ls.layers = append(ls.layers[:layer], ls.layers[layer+1:]...)

	// Adjust current layer if necessary
	if ls.currentLayer >= len(ls.layers) {
		ls.currentLayer = len(ls.layers) - 1
	} else if ls.currentLayer > layer {
		ls.currentLayer--
	}

	return nil
}

// SetLayerVisible sets layer visibility
func (ls *LayerStack) SetLayerVisible(layer int, visible bool) error {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	if layer < 0 || layer >= len(ls.layers) {
		return fmt.Errorf("layer index %d out of range", layer)
	}

	ls.layers[layer].Visible = visible
	return nil
}

// IsLayerVisible checks if a layer is visible
func (ls *LayerStack) IsLayerVisible(layer int) bool {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	if layer < 0 || layer >= len(ls.layers) {
		return false
	}

	return ls.layers[layer].Visible
}

// GetLayer returns a layer by index
func (ls *LayerStack) GetLayer(index int) (*Layer, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	if index < 0 || index >= len(ls.layers) {
		return nil, fmt.Errorf("layer index %d out of range", index)
	}

	return ls.layers[index], nil
}

// GetLayerByID returns a layer by ID
func (ls *LayerStack) GetLayerByID(id int) (*Layer, error) {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	for _, layer := range ls.layers {
		if layer.ID == id {
			return layer, nil
		}
	}

	return nil, fmt.Errorf("layer with ID %d not found", id)
}

// GetLayerNames returns names of all layers
func (ls *LayerStack) GetLayerNames() []string {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	names := make([]string, len(ls.layers))
	for i, layer := range ls.layers {
		visibilityMarker := "●"
		if !layer.Visible {
			visibilityMarker = "○"
		}
		currentMarker := " "
		if i == ls.currentLayer {
			currentMarker = "▶"
		}
		names[i] = fmt.Sprintf("%s%s %s", currentMarker, visibilityMarker, layer.Name)
	}

	return names
}

// SetPixel sets a pixel on a specific layer
func (ls *LayerStack) SetPixel(layerIndex int, x, y int, pixel commands.Pixel) error {
	ls.mutex.RLock()
	layer, err := ls.getLayerUnsafe(layerIndex)
	ls.mutex.RUnlock()

	if err != nil {
		return err
	}

	layer.mutex.Lock()
	defer layer.mutex.Unlock()

	key := fmt.Sprintf("%d,%d", x, y)
	layer.pixels[key] = pixel
	return nil
}

// GetPixel gets a pixel from a specific layer
func (ls *LayerStack) GetPixel(layerIndex int, x, y int) (commands.Pixel, error) {
	ls.mutex.RLock()
	layer, err := ls.getLayerUnsafe(layerIndex)
	ls.mutex.RUnlock()

	if err != nil {
		return commands.Pixel{}, err
	}

	layer.mutex.RLock()
	defer layer.mutex.RUnlock()

	key := fmt.Sprintf("%d,%d", x, y)
	if pixel, exists := layer.pixels[key]; exists {
		return pixel, nil
	}

	return commands.Pixel{Char: ' '}, nil
}

// GetCompositePixel gets the final composited pixel at coordinates
func (ls *LayerStack) GetCompositePixel(x, y int) commands.Pixel {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	// Start with transparent pixel
	result := commands.Pixel{Char: ' '}

	// Composite layers from bottom to top
	for _, layer := range ls.layers {
		if !layer.Visible {
			continue
		}

		layer.mutex.RLock()
		key := fmt.Sprintf("%d,%d", x, y)
		if pixel, exists := layer.pixels[key]; exists && pixel.Char != ' ' {
			// Simple alpha blending (replace for now)
			result = pixel
		}
		layer.mutex.RUnlock()
	}

	return result
}

// getLayerUnsafe returns layer without locking (caller must hold read lock)
func (ls *LayerStack) getLayerUnsafe(index int) (*Layer, error) {
	if index < 0 || index >= len(ls.layers) {
		return nil, fmt.Errorf("layer index %d out of range", index)
	}

	return ls.layers[index], nil
}

// ClearLayer clears all pixels from a layer
func (ls *LayerStack) ClearLayer(layerIndex int) error {
	ls.mutex.RLock()
	layer, err := ls.getLayerUnsafe(layerIndex)
	ls.mutex.RUnlock()

	if err != nil {
		return err
	}

	layer.mutex.Lock()
	defer layer.mutex.Unlock()

	layer.pixels = make(map[string]commands.Pixel)
	return nil
}
