package ui

import (
	"fmt"
	"strings"

	"github.com/cadubi/cadubi/pkg/canvas"
	"github.com/cadubi/cadubi/pkg/commands"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define adaptive colors following Charm's patterns
var (
	borderColor     = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"} // Dim grey for both
	statusColor     = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"} // Dim grey for both
	cursorBg        = lipgloss.AdaptiveColor{Light: "0", Dark: "7"}             // Dark on light, light on dark
	cursorFg        = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}            // Light on light, dark on dark
	modalBorder     = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"} // Purple highlight
	modalBackground = lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"} // Terminal background
	activeTabBg     = lipgloss.AdaptiveColor{Light: "#E5E7EB", Dark: "#374151"} // Light grey / dark grey
	activeTabFg     = lipgloss.AdaptiveColor{Light: "#111827", Dark: "#F9FAFB"} // Dark on light / light on dark
	inactiveTabFg   = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"} // Muted text
)

// KeyMap defines keybindings for the application
type KeyMap struct {
	// Movement
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Drawing
	Draw  key.Binding
	Erase key.Binding

	// Pen tools
	PenAsterisk key.Binding
	PenHash     key.Binding
	PenPlus     key.Binding
	PenMinus    key.Binding
	PenPipe     key.Binding

	// Colors
	CycleColors  key.Binding
	ColorWhite   key.Binding
	ColorBlack   key.Binding
	ColorGreen   key.Binding
	ColorYellow  key.Binding
	ColorMagenta key.Binding

	// Layers
	Layer1           key.Binding
	Layer2           key.Binding
	Layer3           key.Binding
	Layer4           key.Binding
	Layer5           key.Binding
	NewLayer         key.Binding
	ToggleVisibility key.Binding

	// History
	Undo key.Binding
	Redo key.Binding

	// General
	Help key.Binding
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Draw, k.Erase, k.Undo, k.Redo}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // Movement
		{k.Draw, k.Erase},               // Drawing
		{k.PenAsterisk, k.PenHash, k.PenPlus, k.PenMinus, k.PenPipe},   // Pen tools
		{k.CycleColors, k.ColorWhite, k.ColorGreen, k.ColorYellow},     // Colors
		{k.Layer1, k.Layer2, k.Layer3, k.NewLayer, k.ToggleVisibility}, // Layers
		{k.Undo, k.Redo}, // History
		{k.Help, k.Quit}, // General
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Movement
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),

		// Drawing
		Draw: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "draw"),
		),
		Erase: key.NewBinding(
			key.WithKeys("x", "backspace"),
			key.WithHelp("x/bksp", "erase"),
		),

		// Pen tools
		PenAsterisk: key.NewBinding(
			key.WithKeys("*"),
			key.WithHelp("*", "asterisk pen"),
		),
		PenHash: key.NewBinding(
			key.WithKeys("#"),
			key.WithHelp("#", "hash pen"),
		),
		PenPlus: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "plus pen"),
		),
		PenMinus: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "minus pen"),
		),
		PenPipe: key.NewBinding(
			key.WithKeys("|"),
			key.WithHelp("|", "pipe pen"),
		),

		// Colors
		CycleColors: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "cycle colors"),
		),
		ColorWhite: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "white"),
		),
		ColorBlack: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "black"),
		),
		ColorGreen: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "green"),
		),
		ColorYellow: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "yellow"),
		),
		ColorMagenta: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "magenta"),
		),

		// Layers
		Layer1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "layer 1"),
		),
		Layer2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "layer 2"),
		),
		Layer3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "layer 3"),
		),
		Layer4: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "layer 4"),
		),
		Layer5: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "layer 5"),
		),
		NewLayer: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new layer"),
		),
		ToggleVisibility: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "toggle visibility"),
		),

		// History
		Undo: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "undo"),
		),
		Redo: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "redo"),
		),

		// General
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

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
	showHelp    bool
	keyMap      KeyMap
	help        help.Model
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
		keyMap:      DefaultKeyMap(),
		help:        help.New(),
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

	case tea.MouseMsg:
		return a.handleMouseEvent(msg)
	}

	return a, nil
}

// handleMouseEvent processes mouse input
func (a *App) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		// Calculate canvas coordinates from mouse position
		canvasX, canvasY := a.screenToCanvas(msg.X, msg.Y)

		// Move cursor and draw
		a.moveCursorTo(canvasX, canvasY)
		a.drawAtCursor()

	case tea.MouseRight:
		// Calculate canvas coordinates from mouse position
		canvasX, canvasY := a.screenToCanvas(msg.X, msg.Y)

		// Move cursor and erase
		a.moveCursorTo(canvasX, canvasY)
		a.eraseAtCursor()

	case tea.MouseWheelUp:
		// Scroll up (move viewport up)
		a.cursor.Y -= 3
		if a.cursor.Y < 0 {
			a.cursor.Y = 0
		}
		cmd := commands.NewMoveCursorCommand(a.cursor.X, a.cursor.Y)
		a.history.Execute(a.ctx, cmd)

	case tea.MouseWheelDown:
		// Scroll down (move viewport down)
		a.cursor.Y += 3
		if a.cursor.Y > 1000 { // Reasonable limit
			a.cursor.Y = 1000
		}
		cmd := commands.NewMoveCursorCommand(a.cursor.X, a.cursor.Y)
		a.history.Execute(a.ctx, cmd)
	}

	return a, nil
}

// screenToCanvas converts screen coordinates to canvas coordinates
func (a *App) screenToCanvas(screenX, screenY int) (canvasX, canvasY int) {
	// Calculate canvas area bounds (account for border)
	canvasWidth := a.width - 4   // Account for border padding
	canvasHeight := a.height - 6 // Account for border, status bar, and margins

	if canvasWidth < 10 {
		canvasWidth = 10
	}
	if canvasHeight < 5 {
		canvasHeight = 5
	}

	// Get current viewport position (top-left corner of the view)
	viewportX, viewportY := a.canvas.GetViewportPos()

	// Convert screen coordinates to canvas coordinates
	// Account for border (1 char) and padding (1 char) = 2 chars offset
	canvasRelativeX := screenX - 2
	canvasRelativeY := screenY - 1 // Top border

	// Ensure we're within canvas bounds
	if canvasRelativeX >= 0 && canvasRelativeX < canvasWidth &&
		canvasRelativeY >= 0 && canvasRelativeY < canvasHeight {
		// Add viewport offset to get absolute canvas coordinates
		canvasX = viewportX + canvasRelativeX
		canvasY = viewportY + canvasRelativeY
	} else {
		// Default to current cursor position if click is outside canvas
		canvasX = a.cursor.X
		canvasY = a.cursor.Y
	}

	return canvasX, canvasY
}

// moveCursorTo moves cursor to specific coordinates
func (a *App) moveCursorTo(x, y int) {
	// Keep coordinates within reasonable bounds
	if x < -1000 {
		x = -1000
	}
	if x > 1000 {
		x = 1000
	}
	if y < -1000 {
		y = -1000
	}
	if y > 1000 {
		y = 1000
	}

	cmd := commands.NewMoveCursorCommand(x, y)
	a.history.Execute(a.ctx, cmd)
	a.cursor.X = x
	a.cursor.Y = y
}

// handleKeyPress processes keyboard input
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return a, tea.Quit
	case "q":
		// Don't quit if help modal is open
		if !a.showHelp {
			return a, tea.Quit
		}

	// Cursor movement (vi-style) - but not when help is open
	case "h":
		a.moveCursor(-1, 0)
	case "left":
		if !a.showHelp {
			a.moveCursor(-1, 0)
		}
	case "l":
		a.moveCursor(1, 0)
	case "right":
		if !a.showHelp {
			a.moveCursor(1, 0)
		}
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

	// Change pen color (basic colors)
	case "c":
		// Cycle through common colors
		a.cyclePenColor()
	case "w":
		a.setPenColor("white")
	case "b":
		a.setPenColor("black")
	case "g":
		a.setPenColor("green")
	case "y":
		a.setPenColor("yellow")
	case "m":
		a.setPenColor("magenta")

	// Layer controls
	case "1", "2", "3", "4", "5":
		// Switch to layer by number
		layerNum := int(msg.String()[0] - '1') // Convert '1'-'5' to 0-4
		if layerNum < a.canvas.GetLayerCount() {
			a.canvas.SetCurrentLayer(layerNum)
		}

	case "n":
		// Add new layer
		layerCount := a.canvas.GetLayerCount()
		layerName := fmt.Sprintf("Layer %d", layerCount+1)
		a.canvas.AddLayer(layerName)

	case "v":
		// Toggle current layer visibility
		currentLayer := a.canvas.GetCurrentLayer()
		visible := a.canvas.IsLayerVisible(currentLayer)
		a.canvas.SetLayerVisible(currentLayer, !visible)

	case "?":
		// Toggle help modal
		a.showHelp = !a.showHelp

	case "escape", "esc":
		// Close help modal if open
		if a.showHelp {
			a.showHelp = false
		}

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

// setPenColor changes the pen color
func (a *App) setPenColor(color string) {
	a.cursor.Pen.Foreground = color
	a.ctx.Cursor.Pen.Foreground = color
}

// cyclePenColor cycles through common colors
func (a *App) cyclePenColor() {
	colors := []string{"white", "red", "green", "yellow", "blue", "magenta", "cyan", "black"}
	currentColor := a.cursor.Pen.Foreground
	if currentColor == "" {
		currentColor = "white"
	}

	// Find current color index
	currentIndex := 0
	for i, color := range colors {
		if color == currentColor {
			currentIndex = i
			break
		}
	}

	// Move to next color
	nextIndex := (currentIndex + 1) % len(colors)
	a.setPenColor(colors[nextIndex])
}

// View renders the current application state
func (a *App) View() string {

	// Create styles using adaptive colors
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1)

	statusStyle := lipgloss.NewStyle().
		Foreground(statusColor).
		Margin(0, 1)

	// Calculate canvas viewport dimensions
	canvasWidth := a.width - 4   // Account for border padding
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
	currentLayer := a.canvas.GetCurrentLayer()
	layerCount := a.canvas.GetLayerCount()
	layerVisible := a.canvas.IsLayerVisible(currentLayer)
	layerIcon := "⧉"
	if !layerVisible {
		layerIcon = "○"
	}

	// Get canvas dimensions for display
	minX, minY, maxX, maxY := a.canvas.GetBounds()
	canvasW := maxX - minX + 1
	canvasH := maxY - minY + 1

	// Create color indicator
	penColor := a.cursor.Pen.Foreground
	if penColor == "" {
		penColor = "white"
	}
	colorSwatch := lipgloss.NewStyle().Foreground(lipgloss.Color(penColor)).Render("●●●")

	status := statusStyle.Render(fmt.Sprintf(
		"■ Tool: %s │ Pen: '%c' │ %s Colors │ Layer: %s%d/%d │ Canvas: %dx%d │ Cursor: (%d,%d) │ History: %d/%d",
		a.currentTool, a.cursor.Pen.Char, colorSwatch, layerIcon, currentLayer+1, layerCount, canvasW, canvasH, a.cursor.X, a.cursor.Y, position, total,
	))

	// Instructions
	instructions := statusStyle.Render(
		"Move: hjkl │ Draw: space │ Erase: x │ Undo: u │ Redo: r │ Pen: *#+|- │ Colors: c,w,g,y,m │ Layer: 1-5,n,v │ Help: ?",
	)

	// Create the main view
	mainView := lipgloss.JoinVertical(lipgloss.Left, canvas, status, instructions)

	// If help is shown, display it at the bottom
	if a.showHelp {
		// Set help width and show full help
		a.help.Width = a.width
		a.help.ShowAll = true
		helpView := a.help.View(a.keyMap)

		return lipgloss.JoinVertical(lipgloss.Left, mainView, "", helpView)
	}

	return mainView
}

// renderCanvasWithCursor renders the canvas content with cursor overlay
func (a *App) renderCanvasWithCursor(width, height int) string {
	// Position viewport so cursor is visible, but keep (0,0) near top-left when possible
	a.positionViewportForCursor(width, height)

	// Get canvas content centered on cursor
	content := a.canvas.GetViewportCentered(width, height)

	// Calculate viewport position for cursor rendering
	viewportX, viewportY := a.canvas.GetViewportPos()
	startX := viewportX // viewportX/Y is already the top-left corner
	startY := viewportY

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
				Background(cursorBg).
				Foreground(cursorFg)

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

// positionViewportForCursor positions the viewport to show the cursor intelligently
func (a *App) positionViewportForCursor(width, height int) {
	// Simple approach: keep (0,0) at top-left unless cursor moves too far right or down
	// Start scrolling when cursor approaches edges

	marginX := width / 6  // Start scrolling when cursor is 1/6 from edge
	marginY := height / 6 // Start scrolling when cursor is 1/6 from edge

	// Calculate viewport position to keep cursor visible with margin
	viewportX := 0
	viewportY := 0

	// Scroll horizontally if cursor is too far right
	if a.cursor.X >= width-marginX {
		viewportX = a.cursor.X - width + marginX + 1
	}

	// Scroll vertically if cursor is too far down
	if a.cursor.Y >= height-marginY {
		viewportY = a.cursor.Y - height + marginY + 1
	}

	// Scroll left if cursor goes negative and we're not at origin
	if a.cursor.X < marginX && a.cursor.X < 0 {
		viewportX = a.cursor.X - marginX
	}

	// Scroll up if cursor goes negative and we're not at origin
	if a.cursor.Y < marginY && a.cursor.Y < 0 {
		viewportY = a.cursor.Y - marginY
	}

	a.canvas.SetViewport(viewportX, viewportY)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
