package ui

import (
	"fmt"
	"strings"

	"github.com/cadubi/cadubi/pkg/canvas"
	"github.com/cadubi/cadubi/pkg/commands"
	"github.com/cadubi/cadubi/pkg/io"
	"github.com/cadubi/cadubi/pkg/trace"
	"github.com/cadubi/cadubi/pkg/unicode"
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
	// Minimal styling - no custom colors for modals
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
	CycleColors   key.Binding
	ColorPicker   key.Binding
	ColorWhite    key.Binding
	ColorBlack    key.Binding
	ColorGreen    key.Binding
	ColorYellow   key.Binding
	ColorMagenta  key.Binding

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

	// File operations
	Save key.Binding
	Load key.Binding
	New  key.Binding

	// General
	Help           key.Binding
	UnicodeBrowser key.Binding
	Quit           key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.UnicodeBrowser, k.Save, k.Load, k.Quit, k.Draw, k.Erase, k.Undo, k.Redo}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // Movement
		{k.Draw, k.Erase},               // Drawing
		{k.PenAsterisk, k.PenHash, k.PenPlus, k.PenMinus, k.PenPipe},           // Pen tools
		{k.CycleColors, k.ColorPicker, k.ColorWhite, k.ColorGreen, k.ColorYellow}, // Colors
		{k.Layer1, k.Layer2, k.Layer3, k.NewLayer, k.ToggleVisibility}, // Layers
		{k.Undo, k.Redo},                       // History
		{k.Save, k.Load, k.New},                // File operations
		{k.Help, k.UnicodeBrowser, k.Quit},     // General
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
		ColorPicker: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "color picker"),
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

		// File operations
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		Load: key.NewBinding(
			key.WithKeys("ctrl+o"),
			key.WithHelp("ctrl+o", "load"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "new"),
		),

		// General
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		UnicodeBrowser: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "unicode browser"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

// ClickableRegion represents a clickable area in the UI
type ClickableRegion struct {
	X, Y, Width, Height int
	Action              string
	Data                interface{}
}

// App represents the main application state
type App struct {
	canvas           *canvas.Canvas
	history          *commands.CommandHistory
	ctx              *commands.CommandContext
	cursor           commands.CursorState
	layers           *LayerManager
	width            int
	height           int
	viewportX        int
	viewportY        int
	currentTool      string
	showHelp         bool
	unicodeBrowser   *UnicodeBrowser
	confirmDialog    *ConfirmationDialog
	colorPicker      *HSLColorPicker
	keyMap           KeyMap
	help             help.Model
	clickableRegions []ClickableRegion
	quitting         bool
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
	trace.Infof("app", "Creating new CADUBI application")
	
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

	app := &App{
		canvas:         c,
		history:        commands.NewCommandHistory(1000),
		ctx:            ctx,
		cursor:         cursor,
		layers:         layers,
		currentTool:    "pen",
		unicodeBrowser: NewUnicodeBrowser(),
		confirmDialog:  NewSaveConfirmationDialog(),
		colorPicker:    NewHSLColorPicker(),
		keyMap:         DefaultKeyMap(),
		help:           help.New(),
		quitting:       false,
	}
	
	trace.Infof("app", "Application created successfully with canvas size 80x24, initial tool: pen")
	return app
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		trace.LogTerminalInfo(msg.Width, msg.Height)
		a.width = msg.Width
		a.height = msg.Height
		// Forward to Unicode browser
		cmd := a.unicodeBrowser.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		trace.LogUIEvent("key_press", map[string]interface{}{
			"key": msg.String(),
			"unicode_browser_visible": a.unicodeBrowser.IsVisible(),
			"help_visible": a.showHelp,
		})
		
		// If confirmation dialog is visible, handle its keys first
		if a.confirmDialog.IsVisible() {
			cmd := a.confirmDialog.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}
		
		// If color picker is visible, handle its keys first
		if a.colorPicker.IsVisible() {
			cmd := a.colorPicker.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}
		
		// If Unicode browser is visible, handle its keys first
		if a.unicodeBrowser.IsVisible() {
			cmd := a.unicodeBrowser.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			
			// Check if character was selected (Enter/Space was pressed)
			if msg.String() == "enter" || msg.String() == " " {
				newChar := a.unicodeBrowser.GetSelectedChar()
				a.setPenChar(newChar)
				a.unicodeBrowser.AddRecentChar(newChar) // Track in recent characters
				a.unicodeBrowser.Hide() // Close browser after selection
			}
			
			// Don't process main app keys when browser is open, except for close
			if msg.String() != "escape" && msg.String() != "esc" && msg.String() != "q" {
				return a, tea.Batch(cmds...)
			}
		}
		
		model, cmd := a.handleKeyPress(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		return model, tea.Batch(cmds...)

	case tea.MouseMsg:
		trace.LogUIEvent("mouse_event", map[string]interface{}{
			"type": int(msg.Type),
			"x": msg.X,
			"y": msg.Y,
			"button": int(msg.Button),
			"unicode_browser_visible": a.unicodeBrowser.IsVisible(),
			"confirm_dialog_visible": a.confirmDialog.IsVisible(),
		})
		
		// If confirmation dialog is visible, let it handle mouse events first
		if a.confirmDialog.IsVisible() {
			cmd := a.confirmDialog.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}
		
		// If color picker is visible, let it handle mouse events first
		if a.colorPicker.IsVisible() {
			cmd := a.colorPicker.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return a, tea.Batch(cmds...)
		}
		
		// If Unicode browser is visible, let it handle mouse events first
		if a.unicodeBrowser.IsVisible() {
			cmd := a.unicodeBrowser.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			// Don't pass mouse events to main app when browser is open
			return a, tea.Batch(cmds...)
		} else {
			return a.handleMouseEvent(msg)
		}
		
	case CharacterSelectedMsg:
		// Handle character selection from Unicode browser (double-click)
		a.setPenChar(msg.Character)
		a.unicodeBrowser.AddRecentChar(msg.Character)
		a.unicodeBrowser.Hide()
		return a, nil

	case ConfirmationResultMsg:
		// Handle confirmation dialog results
		switch msg.Action {
		case "yes": // Save & Exit
			a.saveFile()
			return a, tea.Quit
		case "no": // Exit Without Saving  
			return a, tea.Quit
		case "cancel": // Cancel
			a.quitting = false
		}
		return a, nil

	case ColorSelectedMsg:
		// Handle color selection from color picker
		a.setPenColor(msg.Color)
		trace.Debugf("ui", "Color selected from picker: %s (%s)", msg.Name, msg.Color)
		return a, nil
	}

	return a, nil
}

// handleMouseEvent processes mouse input
func (a *App) handleMouseEvent(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		// First check if click is on any clickable UI regions
		if region := a.getClickedRegion(msg.X, msg.Y); region != nil {
			return a.handleUIClick(*region)
		}
		
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

// getClickedRegion returns the first clickable region that contains the given coordinates
func (a *App) getClickedRegion(x, y int) *ClickableRegion {
	for _, region := range a.clickableRegions {
		if x >= region.X && x < region.X+region.Width &&
		   y >= region.Y && y < region.Y+region.Height {
			return &region
		}
	}
	return nil
}

// addClickableRegion adds a clickable region to the tracking list
func (a *App) addClickableRegion(x, y, width, height int, action string, data interface{}) {
	a.clickableRegions = append(a.clickableRegions, ClickableRegion{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Action: action,
		Data:   data,
	})
}

// cycleTool cycles through available tools
func (a *App) cycleTool() {
	tools := []string{"pen", "brush", "fill"} // Available tools
	currentIndex := 0
	
	// Find current tool index
	for i, tool := range tools {
		if tool == a.currentTool {
			currentIndex = i
			break
		}
	}
	
	// Move to next tool
	nextIndex := (currentIndex + 1) % len(tools)
	oldTool := a.currentTool
	a.currentTool = tools[nextIndex]
	a.cursor.Tool = a.currentTool
	a.ctx.Cursor.Tool = a.currentTool
	
	trace.Debugf("ui", "Tool cycled from '%s' to '%s'", oldTool, a.currentTool)
}

// handleUIClick processes clicks on UI elements
func (a *App) handleUIClick(region ClickableRegion) (tea.Model, tea.Cmd) {
	trace.LogUIEvent("ui_click", map[string]interface{}{
		"action": region.Action,
		"x": region.X,
		"y": region.Y,
		"width": region.Width,
		"height": region.Height,
	})
	
	switch region.Action {
	case "cycle_colors":
		a.cyclePenColor()
	case "cycle_tool":
		a.cycleTool()
	case "toggle_layer_visibility":
		currentLayer := a.canvas.GetCurrentLayer()
		visible := a.canvas.IsLayerVisible(currentLayer)
		a.canvas.SetLayerVisible(currentLayer, !visible)
		trace.Debugf("ui", "Layer %d visibility toggled to %v", currentLayer, !visible)
	case "canvas_info":
		// For now, canvas click doesn't do anything (placeholder for future resize dialog)
		// TODO: Implement canvas resize/info dialog
		trace.Debugf("ui", "Canvas info clicked (not implemented)")
	case "undo":
		a.undo()
	case "redo":
		a.redo()
	// TODO: Add more actions as needed (color picker dialog, resize dialog, etc.)
	}
	return a, nil
}

// screenToCanvas converts screen coordinates to canvas coordinates accounting for wide characters
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

	// Account for border (1 char) and padding (1 char) = 2 chars offset
	screenRelativeX := screenX - 2
	screenRelativeY := screenY - 1 // Top border

	// Ensure we're within canvas bounds
	if screenRelativeY >= 0 && screenRelativeY < canvasHeight && screenRelativeX >= 0 {
		// Calculate the canvas Y coordinate (straightforward)
		canvasY = viewportY + screenRelativeY
		
		// For X coordinate, we need to account for wide characters
		canvasX = a.screenXToCanvasX(screenRelativeX, canvasY, viewportX, canvasWidth)
	} else {
		// Default to current cursor position if click is outside canvas
		canvasX = a.cursor.X
		canvasY = a.cursor.Y
	}

	return canvasX, canvasY
}

// screenXToCanvasX converts screen X coordinate to canvas X coordinate, accounting for wide characters
func (a *App) screenXToCanvasX(screenRelativeX, canvasY, viewportX, canvasWidth int) int {
	// Get the content for this row from the canvas
	currentScreenX := 0
	canvasX := viewportX
	
	// Iterate through canvas positions and accumulate screen width
	for currentScreenX < screenRelativeX && (canvasX - viewportX) < canvasWidth {
		// Get the character at this canvas position
		pixel, err := a.canvas.GetPixel(canvasX, canvasY)
		var char rune = ' '
		if err == nil {
			char = pixel.Char
		}
		
		// Calculate how much screen space this character takes
		charWidth := a.getCharDisplayWidth(char)
		
		// If adding this character would exceed our target screen position,
		// we've found the right canvas position
		if currentScreenX + charWidth > screenRelativeX {
			return canvasX
		}
		
		currentScreenX += charWidth
		canvasX++
	}
	
	return canvasX
}

// getCharDisplayWidth returns the display width of a character using proper Unicode width calculation
func (a *App) getCharDisplayWidth(char rune) int {
	return unicode.CharWidth(char)
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
		return a.handleQuit()
	case "q":
		// Don't quit if help modal is open
		if !a.showHelp {
			return a.handleQuit()
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
	case "p":
		// Open color picker modal
		if !a.showHelp && !a.unicodeBrowser.IsVisible() {
			a.colorPicker.Show()
		}
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

	case "s":
		// Toggle Unicode browser
		if !a.showHelp { // Don't open browser if help is open
			a.unicodeBrowser.Toggle()
			// If a character was selected, update pen character
			if a.unicodeBrowser.IsVisible() {
				selectedChar := a.unicodeBrowser.GetSelectedChar()
				a.setPenChar(selectedChar)
			}
		}

	// File operations
	case "ctrl+s":
		// Save current file
		a.saveFile()

	case "ctrl+o":
		// Load file (for now, load a default test file)
		// TODO: Implement file picker dialog
		a.loadFile("test.cadubi")

	case "ctrl+n":
		// New file
		a.newFile()

	case "escape", "esc":
		// Close modals if open
		if a.showHelp {
			a.showHelp = false
		}
		if a.colorPicker.IsVisible() {
			a.colorPicker.Hide()
		}
		if a.unicodeBrowser.IsVisible() {
			a.unicodeBrowser.Hide()
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
	trace.LogCommand("DrawCommand", map[string]interface{}{
		"x": a.cursor.X,
		"y": a.cursor.Y,
		"char": string(a.cursor.Pen.Char),
		"color": a.cursor.Pen.Foreground,
		"layer": a.canvas.GetCurrentLayer(),
	})
	cmd := commands.NewDrawCommand(a.cursor.X, a.cursor.Y, a.cursor.Pen, 0)
	a.history.Execute(a.ctx, cmd)
}

// eraseAtCursor erases at the current cursor position
func (a *App) eraseAtCursor() {
	trace.LogCommand("EraseCommand", map[string]interface{}{
		"x": a.cursor.X,
		"y": a.cursor.Y,
		"layer": a.canvas.GetCurrentLayer(),
	})
	cmd := commands.NewEraseCommand(a.cursor.X, a.cursor.Y, 0)
	a.history.Execute(a.ctx, cmd)
}

// undo undoes the last command
func (a *App) undo() {
	trace.LogCommand("UndoCommand", map[string]interface{}{})
	a.history.Undo(a.ctx)
}

// redo redoes the next command
func (a *App) redo() {
	trace.LogCommand("RedoCommand", map[string]interface{}{})
	a.history.Redo(a.ctx)
}

// setPenChar changes the pen character
func (a *App) setPenChar(char rune) {
	trace.Debugf("ui", "Changing pen character from '%c' to '%c'", a.cursor.Pen.Char, char)
	a.cursor.Pen.Char = char
	a.ctx.Cursor.Pen.Char = char
	// Add to recent characters if it's not a basic ASCII character
	if char > 127 {
		a.unicodeBrowser.AddRecentChar(char)
	}
}

// setPenColor changes the pen color
func (a *App) setPenColor(color string) {
	trace.Debugf("ui", "Changing pen color from '%s' to '%s'", a.cursor.Pen.Foreground, color)
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

	// Clear previous clickable regions  
	a.clickableRegions = []ClickableRegion{}

	// Create status bar with clickable regions
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

	// Build status bar with clickable tracking
	statusY := canvasHeight + 2 // Position after canvas and border
	currentX := 1 // Start with left margin
	
	// Build status parts with position tracking
	toolPart := fmt.Sprintf("■ Tool: %s", a.currentTool)
	penPart := fmt.Sprintf("Pen: 〈%c〉", a.cursor.Pen.Char)
	colorPart := fmt.Sprintf("%s Colors", colorSwatch)
	layerPart := fmt.Sprintf("Layer: %s%d/%d", layerIcon, currentLayer+1, layerCount)
	canvasPart := fmt.Sprintf("Canvas: %dx%d", canvasW, canvasH)
	cursorPart := fmt.Sprintf("Cursor: (%d,%d)", a.cursor.X, a.cursor.Y)
	historyPart := fmt.Sprintf("History: %d/%d", position, total)
	
	// Add file status indicator
	fileStatus := "○" // Saved
	if a.canvas.IsModified() {
		fileStatus = "●" // Modified
	}
	filename := a.canvas.GetFilename()
	if filename == "" {
		filename = "untitled"
	}
	filePart := fmt.Sprintf("File: %s%s", fileStatus, filename)
	
	// Add clickable regions for each part
	// Tool (not clickable for now, but could cycle tools)
	a.addClickableRegion(currentX, statusY, len(toolPart), 1, "cycle_tool", nil)
	currentX += len(toolPart) + 3 // Add separator length " · "
	
	// Skip pen part (not clickable)
	currentX += len(penPart) + 3
	
	// Color swatch (clickable to cycle colors)  
	a.addClickableRegion(currentX, statusY, len(colorPart), 1, "cycle_colors", nil)
	currentX += len(colorPart) + 3
	
	// Layer info (clickable to toggle visibility)
	a.addClickableRegion(currentX, statusY, len(layerPart), 1, "toggle_layer_visibility", nil)
	currentX += len(layerPart) + 3
	
	// Canvas dimensions (clickable to show info or reset)
	a.addClickableRegion(currentX, statusY, len(canvasPart), 1, "canvas_info", nil)
	currentX += len(canvasPart) + 3
	
	// Skip cursor part (not clickable)
	currentX += len(cursorPart) + 3
	
	// History (clickable for undo/redo - left side undo, right side redo)
	historyWidth := len(historyPart)
	halfWidth := historyWidth / 2
	
	// Left half for undo
	a.addClickableRegion(currentX, statusY, halfWidth, 1, "undo", nil)
	// Right half for redo
	a.addClickableRegion(currentX + halfWidth, statusY, historyWidth - halfWidth, 1, "redo", nil)

	status := statusStyle.Render(fmt.Sprintf(
		"%s · %s · %s · %s · %s · %s · %s · %s",
		toolPart, penPart, colorPart, layerPart, canvasPart, cursorPart, historyPart, filePart,
	))

	// Instructions
	instructions := statusStyle.Render(
		"Move: hjkl · Draw: space · Erase: x · Undo: u · Redo: r · Pen: *#+|- · Colors: p,c,w,g,y,m · File: Ctrl+s,Ctrl+o,Ctrl+n · Layer: 1-5,n,v · Unicode: s · Help: ?",
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

	// If confirmation dialog is visible, overlay it on the main view  
	if a.confirmDialog.IsVisible() {
		a.confirmDialog.SetDimensions(a.width, a.height)
		dialogView := a.confirmDialog.View()
		// Center the dialog on the screen
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, dialogView, lipgloss.WithWhitespaceChars(""))
	}

	// If color picker is visible, overlay it on the main view
	if a.colorPicker.IsVisible() {
		a.colorPicker.SetDimensions(a.width, a.height)
		pickerView := a.colorPicker.View()
		// Center the color picker on the screen
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, pickerView, lipgloss.WithWhitespaceChars(""))
	}

	// If Unicode browser is visible, overlay it on the main view
	if a.unicodeBrowser.IsVisible() {
		browserView := a.unicodeBrowser.View()
		// Center the browser on the screen
		return lipgloss.Place(a.width, a.height, lipgloss.Center, lipgloss.Center, browserView, lipgloss.WithWhitespaceChars(""))
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

// File operation methods

// handleQuit handles the quit process, checking for unsaved changes
func (a *App) handleQuit() (tea.Model, tea.Cmd) {
	if a.canvas.IsModified() {
		// Show save confirmation dialog
		a.confirmDialog.Show()
		a.quitting = true
		return a, nil
	} else {
		// No unsaved changes, quit directly
		return a, tea.Quit
	}
}

// saveFile saves the current canvas to a file
func (a *App) saveFile() {
	filename := a.canvas.GetFilename()
	if filename == "" {
		// Default filename if none set
		filename = "untitled.cadubi"
		a.canvas.SetFilename(filename)
	}

	trace.Infof("file", "Saving canvas to file: %s", filename)
	
	// Save using the io package
	err := io.SaveToFile(filename, a.canvas)
	if err != nil {
		trace.Errorf("file", "Failed to save file %s: %v", filename, err)
		// TODO: Show error dialog to user
		return
	}

	// Mark as not modified after successful save
	a.canvas.SetModified(false)
	trace.Infof("file", "Successfully saved canvas to %s", filename)
}

// loadFile loads a canvas from a file
func (a *App) loadFile(filename string) {
	trace.Infof("file", "Loading canvas from file: %s", filename)
	
	// Load using the io package
	newCanvas, layerInfos, err := io.LoadFromFile(filename)
	if err != nil {
		trace.Errorf("file", "Failed to load file %s: %v", filename, err)
		// TODO: Show error dialog to user
		return
	}

	// Replace current canvas
	a.canvas = newCanvas
	a.canvas.SetFilename(filename)
	a.canvas.SetModified(false)

	// Update context
	a.ctx.Canvas = newCanvas

	// Reset cursor to origin
	a.cursor.X = 0
	a.cursor.Y = 0
	a.ctx.Cursor = a.cursor

	trace.Infof("file", "Successfully loaded canvas from %s with %d layers", filename, len(layerInfos))
}

// newFile creates a new blank canvas
func (a *App) newFile() {
	trace.Infof("file", "Creating new blank canvas")
	
	// Create new canvas
	a.canvas = canvas.NewCanvas(80, 24)
	a.canvas.SetModified(false)
	a.canvas.SetFilename("")

	// Update context
	a.ctx.Canvas = a.canvas

	// Reset cursor
	a.cursor.X = 0
	a.cursor.Y = 0
	a.ctx.Cursor = a.cursor

	// Reset history
	a.history = commands.NewCommandHistory(1000)

	trace.Infof("file", "New blank canvas created")
}
