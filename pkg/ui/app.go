package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cadubi/cadubi/pkg/canvas"
	"github.com/cadubi/cadubi/pkg/commands"
)

// App represents the main application state
type App struct {
	canvas      *canvas.Canvas
	history     *commands.CommandHistory
	ctx         *commands.CommandContext
	cursor      commands.CursorState
	layers      *LayerManager
	width       int
	height      int
	viewportX   int
	viewportY   int
	currentTool string
}

// LayerManager is a simple implementation of the layer manager interface
type LayerManager struct {
	layers       []Layer
	currentLayer int
}

// Layer represents a drawing layer
type Layer struct {
	Name    string
	Visible bool
}

// NewLayerManager creates a new layer manager
func NewLayerManager() *LayerManager {
	return &LayerManager{
		layers: []Layer{
			{Name: "Layer 1", Visible: true},
		},
		currentLayer: 0,
	}
}

// Implement LayerManager interface methods
func (lm *LayerManager) GetCurrentLayer() int {
	return lm.currentLayer
}

func (lm *LayerManager) SetCurrentLayer(layer int) error {
	if layer >= 0 && layer < len(lm.layers) {
		lm.currentLayer = layer
	}
	return nil
}

func (lm *LayerManager) GetLayerCount() int {
	return len(lm.layers)
}

func (lm *LayerManager) AddLayer(name string) (int, error) {
	lm.layers = append(lm.layers, Layer{Name: name, Visible: true})
	return len(lm.layers) - 1, nil
}

func (lm *LayerManager) RemoveLayer(layer int) error {
	if layer >= 0 && layer < len(lm.layers) {
		lm.layers = append(lm.layers[:layer], lm.layers[layer+1:]...)
		if lm.currentLayer >= len(lm.layers) {
			lm.currentLayer = len(lm.layers) - 1
		}
	}
	return nil
}

func (lm *LayerManager) SetLayerVisible(layer int, visible bool) error {
	if layer >= 0 && layer < len(lm.layers) {
		lm.layers[layer].Visible = visible
	}
	return nil
}

func (lm *LayerManager) IsLayerVisible(layer int) bool {
	if layer >= 0 && layer < len(lm.layers) {
		return lm.layers[layer].Visible
	}
	return false
}

// NewApp creates a new application
func NewApp() *App {
	c := canvas.NewCanvas(80, 24)
	layers := NewLayerManager()
	
	cursor := commands.CursorState{
		X:    0,
		Y:    0,
		Tool: "pen",
		Pen: commands.Pixel{
			Char:       '*',
			Foreground: "white",
			Bold:       false,
		},
	}
	
	ctx := &commands.CommandContext{
		Canvas: c,
		Layers: layers,
		Cursor: cursor,
	}
	
	return &App{
		canvas:      c,
		history:     commands.NewCommandHistory(1000),
		ctx:         ctx,
		cursor:      cursor,
		layers:      layers,
		currentTool: "pen",
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		
	case tea.KeyMsg:
		return a.handleKeyPress(msg)
	}
	
	return a, nil
}

// handleKeyPress processes keyboard input
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, tea.Quit
		
	// Cursor movement (vi-style)
	case "h", "left":
		a.moveCursor(-1, 0)
	case "l", "right":
		a.moveCursor(1, 0)
	case "k", "up":
		a.moveCursor(0, -1)
	case "j", "down":
		a.moveCursor(0, 1)
		
	// Drawing
	case " ":
		a.drawAtCursor()
		
	// Erase
	case "x", "backspace":
		a.eraseAtCursor()
		
	// Undo/Redo
	case "u":
		a.undo()
	case "r":
		a.redo()
		
	// Change pen character
	case "*":
		a.setPenChar('*')
	case "#":
		a.setPenChar('#')
	case "+":
		a.setPenChar('+')
	case "-":
		a.setPenChar('-')
	case "|":
		a.setPenChar('|')
	}
	
	return a, nil
}

// moveCursor moves the cursor by the specified offset
func (a *App) moveCursor(dx, dy int) {
	newX := a.cursor.X + dx
	newY := a.cursor.Y + dy
	
	// Keep cursor in bounds (with some buffer for infinite canvas)
	if newX < 0 {
		newX = 0
	}
	if newY < 0 {
		newY = 0
	}
	if newX > 200 { // Reasonable upper limit
		newX = 200
	}
	if newY > 100 {
		newY = 100
	}
	
	cmd := commands.NewMoveCursorCommand(newX, newY)
	a.history.Execute(a.ctx, cmd)
	a.cursor.X = newX
	a.cursor.Y = newY
}

// drawAtCursor draws at the current cursor position
func (a *App) drawAtCursor() {
	cmd := commands.NewDrawCommand(a.cursor.X, a.cursor.Y, a.cursor.Pen, 0)
	a.history.Execute(a.ctx, cmd)
}

// eraseAtCursor erases at the current cursor position
func (a *App) eraseAtCursor() {
	cmd := commands.NewEraseCommand(a.cursor.X, a.cursor.Y, 0)
	a.history.Execute(a.ctx, cmd)
}

// undo undoes the last command
func (a *App) undo() {
	a.history.Undo(a.ctx)
}

// redo redoes the next command
func (a *App) redo() {
	a.history.Redo(a.ctx)
}

// setPenChar changes the pen character
func (a *App) setPenChar(char rune) {
	a.cursor.Pen.Char = char
	a.ctx.Cursor.Pen.Char = char
}

// View renders the current application state
func (a *App) View() string {
	// Create styles
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)
	
	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Margin(0, 1)
	
	// Calculate canvas viewport dimensions
	canvasWidth := a.width - 4  // Account for border padding
	canvasHeight := a.height - 6 // Account for border, status bar, and margins
	
	if canvasWidth < 10 {
		canvasWidth = 10
	}
	if canvasHeight < 5 {
		canvasHeight = 5
	}
	
	// Get canvas content with cursor
	canvasContent := a.renderCanvasWithCursor(canvasWidth, canvasHeight)
	
	// Create bordered canvas
	canvas := borderStyle.
		Width(canvasWidth).
		Height(canvasHeight).
		Render(canvasContent)
	
	// Create status bar
	total, position := a.history.GetStats()
	pixelCount := a.canvas.GetPixelCount()
	
	status := statusStyle.Render(fmt.Sprintf(
		"Tool: %s │ Pen: '%c' │ Cursor: (%d,%d) │ History: %d/%d │ Pixels: %d │ Press 'q' to quit",
		a.currentTool, a.cursor.Pen.Char, a.cursor.X, a.cursor.Y, position, total, pixelCount,
	))
	
	// Instructions
	instructions := statusStyle.Render(
		"Movement: hjkl/arrows │ Draw: space │ Erase: x │ Undo: u │ Redo: r │ Pen: *#+|-",
	)
	
	return lipgloss.JoinVertical(lipgloss.Left, canvas, status, instructions)
}

// renderCanvasWithCursor renders the canvas content with cursor overlay
func (a *App) renderCanvasWithCursor(width, height int) string {
	// Calculate viewport to center cursor
	startX := a.cursor.X - width/2
	startY := a.cursor.Y - height/2
	
	if startX < 0 {
		startX = 0
	}
	if startY < 0 {
		startY = 0
	}
	
	// Get canvas content
	content := a.canvas.GetViewport(startX, startY, width, height)
	
	// Overlay cursor
	lines := strings.Split(content, "\n")
	
	// Ensure we have enough lines
	for len(lines) < height {
		lines = append(lines, strings.Repeat(" ", width))
	}
	
	// Calculate cursor position within viewport
	cursorX := a.cursor.X - startX
	cursorY := a.cursor.Y - startY
	
	if cursorY >= 0 && cursorY < len(lines) && cursorX >= 0 {
		line := lines[cursorY]
		
		// Ensure line is long enough
		for len(line) <= cursorX {
			line += " "
		}
		
		// Convert to runes for proper character handling
		runes := []rune(line)
		
		// Ensure we have enough runes
		for len(runes) <= cursorX {
			runes = append(runes, ' ')
		}
		
		// Add cursor styling
		if cursorX < len(runes) {
			cursorStyle := lipgloss.NewStyle().
				Background(lipgloss.Color("#7C3AED")).
				Foreground(lipgloss.Color("#FFFFFF"))
			
			cursorChar := runes[cursorX]
			if cursorChar == ' ' {
				cursorChar = '█' // Block cursor for empty spaces
			}
			
			// Replace character at cursor position
			before := string(runes[:cursorX])
			after := string(runes[cursorX+1:])
			lines[cursorY] = before + cursorStyle.Render(string(cursorChar)) + after
		}
	}
	
	return strings.Join(lines, "\n")
}