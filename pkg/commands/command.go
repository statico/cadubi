package commands

import "time"

// Command represents an undoable action in the application
type Command interface {
	// Execute performs the command
	Execute(ctx *CommandContext) error

	// Undo reverses the command
	Undo(ctx *CommandContext) error

	// Description returns a human-readable description of the command
	Description() string

	// Timestamp returns when the command was created
	Timestamp() time.Time
}

// CommandContext provides the environment for command execution
type CommandContext struct {
	Canvas    CanvasManager
	Layers    LayerManager
	Cursor    CursorState
	Selection SelectionManager
	Clipboard ClipboardManager
}

// CanvasManager interface for canvas operations
type CanvasManager interface {
	GetPixel(x, y int) (Pixel, error)
	SetPixel(x, y int, pixel Pixel) error
	GetDimensions() (width, height int)
	Clear() error
}

// LayerManager interface for layer operations
type LayerManager interface {
	GetCurrentLayer() int
	SetCurrentLayer(layer int) error
	GetLayerCount() int
	AddLayer(name string) (int, error)
	RemoveLayer(layer int) error
	SetLayerVisible(layer int, visible bool) error
	IsLayerVisible(layer int) bool
}

// CursorState represents cursor position and state
type CursorState struct {
	X, Y int
	Tool string
	Pen  Pixel
}

// SelectionManager interface for selection operations
type SelectionManager interface {
	GetSelection() (x1, y1, x2, y2 int, hasSelection bool)
	SetSelection(x1, y1, x2, y2 int) error
	ClearSelection() error
}

// ClipboardManager interface for clipboard operations
type ClipboardManager interface {
	Copy(data [][]Pixel) error
	Paste() ([][]Pixel, error)
	HasData() bool
}

// Pixel represents a single character with formatting
type Pixel struct {
	Char       rune
	Foreground string // ANSI color
	Background string // ANSI color
	Bold       bool
	Italic     bool
	Underline  bool
}

// BaseCommand provides common functionality for commands
type BaseCommand struct {
	timestamp   time.Time
	description string
}

// NewBaseCommand creates a new base command
func NewBaseCommand(description string) BaseCommand {
	return BaseCommand{
		timestamp:   time.Now(),
		description: description,
	}
}

// Timestamp returns when the command was created
func (bc BaseCommand) Timestamp() time.Time {
	return bc.timestamp
}

// Description returns the command description
func (bc BaseCommand) Description() string {
	return bc.description
}
