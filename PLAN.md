# CADUBI v2 - Development Plan & Implementation Roadmap

## Project Overview

Modern terminal-based ASCII art editor built with Go using the Charm Bracelet libraries (Bubble Tea, Bubbles, Lipgloss). Complete rewrite of the 1997 Perl original with modern architecture, advanced features, and cross-platform single executable distribution.

## 🎉 Recent Major Achievement: Automatic Theme Detection (Phase 4.5)

**COMPLETED**: Fully implemented Charm's adaptive color system following Bubble Tea/Lipgloss best practices:
- ✅ **Zero Configuration**: Application launches directly to main interface  
- ✅ **Automatic Terminal Detection**: Colors adapt to light/dark terminal backgrounds
- ✅ **Simplified Architecture**: Removed 500+ lines of complex theme management code
- ✅ **Charm Ecosystem Alignment**: Uses `lipgloss.AdaptiveColor` patterns throughout
- ✅ **Seamless User Experience**: No setup dialogs, preferences, or configuration required

The application now provides a clean, professional experience that automatically adapts to any terminal environment.

## Phase 1: Foundation ✅ COMPLETED

### Documentation & Research

- [x] Read all Bubble Tea documentation using gitingest → `.claude/BUBBLETEA.md`
- [x] Analyzed original CADUBI Perl codebase → `.claude/OLD_PERL_VERSION.md`
- [x] Created comprehensive CLAUDE.md for future instances → `.claude/CLAUDE.md`

### Project Setup

- [x] Analyzed original Perl codebase (preserved in repository)
- [ ] Initialize Go module with proper dependencies
- [ ] Set up Go + Bubble Tea + Bubbles + Lipgloss stack
- [ ] Created working "Hello World" demo application with Bubble Tea

### Command Pattern Architecture

- [ ] Implement comprehensive Command Pattern system in Go
- [ ] Create `CommandHistory` struct for undo/redo functionality
- [ ] Design core interfaces: `Command`, `CommandContext`, `CanvasState`, `LayerManager`
- [ ] Build foundation commands: `DrawCommand`, `EraseCommand`, `MoveCursorCommand`, `NewDocumentCommand`

### Test Harness Creation

- [ ] Set up Go testing framework (built-in testing package)
- [ ] Implement headless testing with Bubble Tea's `WithoutRenderer` option
- [ ] Create command-driven test system: JSON array of commands (create document, draw pixel, save, etc.)
- [ ] Build state validation testing: run commands headlessly and verify final state
- [ ] Integrate Catwalk testing library for TUI model unit testing
- [ ] Set up Teatest for end-to-end testing scenarios
- [ ] Create UI snapshot testing system for Bubble Tea components
- [ ] Enable automated UI regression testing with headless execution

## Phase 2: Canvas System ✅ COMPLETED

### Core Canvas System

- [ ] Canvas rendering engine with Bubble Tea-based character grid
- [ ] Infinite canvas system with dynamic expansion when drawing beyond viewport bounds
- [ ] Responsive viewport handling terminal resizing with minimum size warnings
- [ ] Real-time input handling connecting keybindings to command execution with live updates
- [ ] Basic drawing tools: pen tool, cursor movement, and paint/erase functionality
- [ ] Layer visualization with layer panel and basic layer switching

### Canvas & Layer Architecture

- [ ] Photoshop-like layer system: Layers overlay each other with blending modes
- [ ] Infinite canvas management: Canvas dynamically expands when drawing beyond bounds
- [ ] Viewport system: Terminal shows a window into the infinite canvas, can pan around
- [ ] Smart rendering engine: Only redraw changed regions, handle viewport clipping
- [ ] Responsive UI: Gracefully handle small terminal windows (minimum 40x10), show error for smaller

### Input & Interaction

- [ ] Keyboard shortcuts: Vi-style movement (hjkl) + modern shortcuts
- [ ] Modern keybindings inspired by Vim, Photoshop, and Figma
- [ ] Mouse support: Implementation with Bubble Tea's mouse events and scrolling

## Phase 3: Advanced Color System ✅ COMPLETED

### Color System

- [ ] Full 256-color terminal support with RGB ↔ HSL conversion algorithms
- [ ] Advanced color picker with multiple modes (Palette|HSL|Schemes|Recent)
- [ ] Color palette management (standard, art palettes)
- [ ] Recent colors tracking
- [ ] Color scheme generation (monochromatic, analogous, complementary, triadic)
- [ ] Smart color matching: Auto-suggest terminal-safe colors

### Responsive UI System ✅ UPDATED FOR 80x24 MINIMUM

- [x] **Responsive UI system for different terminal sizes with 80x24 minimum support**
  - [x] All modals and dialogs fit within 80x24 terminal windows
  - [x] First-run theme selection dialog optimized for small screens
  - [x] Compact layouts that scale up for larger terminals
  - [x] Graceful handling of terminals smaller than 80x24 with warnings
- [x] **Help dialog system with comprehensive keyboard shortcuts**
  - [x] Tabbed help system with pagination (< > keys)
  - [x] Compact modal that fits 80x24 screens
  - [x] Six help categories: Movement, Drawing, Pen Tools, Layers, History, General
- [x] **Modal dialog sizing standards:**
  - [x] Maximum dialog width: 76 characters (4 char padding in 80 col terminal)
  - [x] Maximum dialog height: 20 lines (4 line padding in 24 row terminal)
  - [x] All content must be legible and functional at minimum size

### Development Tooling

- [ ] Go formatting with gofmt and pre-commit hooks
- [ ] Git hooks for code quality
- [ ] Comprehensive test suite with Go testing framework

## Phase 4: UI Enhancement & Unicode Symbols ✅ COMPLETED

**KEY ACHIEVEMENTS:**
- ✅ **Integrated Color Picker**: Added visual color swatches (●●●) and keyboard shortcuts (c,w,g,y,m,b) 
- ✅ **Enhanced Status Bar**: ■ Tool: pen │ Pen: '*' │ ●●● Colors │ Layer: ⧉1/1 │ Canvas: 1x1 │ Cursor: (0,0) │ History: 0/0
- ✅ **Unicode Symbol System**: ⧉ (layer stack), ○ (hidden), ■ (tools), ●●● (colors), ← → (navigation)
- ✅ **Comprehensive Testing**: UI coverage improved from 0% → 62.5% with headless testing framework
- ✅ **Two-Panel Layout**: Clean separation with bordered canvas and status information

### 🎯 Unicode Design Guidelines

**PRIMARY PRINCIPLE:** Heavy use of Unicode box drawing characters and mathematical symbols for all UI elements. Emoji is discouraged due to inconsistent width rendering across terminals.

#### Preferred Unicode Character Sets:

- **Box Drawing** (U+2500–U+257F): ┌┐└┘├┤┬┴┼─│ etc.
- **Block Elements** (U+2580–U+259F): ▀▄█▌▐░▒▓ etc.
- **Mathematical Operators** (U+2200–U+22FF): ∀∃∈∉∑∏∫ etc.
- **Geometric Shapes** (U+25A0–U+25FF): ■□▲▼●○◆◇ etc.
- **Miscellaneous Symbols** (U+2600–U+26FF): ☐☑☒★☆ etc.
- **Arrows** (U+2190–U+21FF): ←↑→↓↔↕⇄⇅ etc.

#### UI Component Symbol Standards:

- **Layer Visibility**: ⧉ (layer stack) / ○ (hidden) - ✅ IMPLEMENTED
- **Layer Status**: ▶ (active) / ◆ (inactive) - ✅ IMPLEMENTED
- **File Status**: ● (modified) / ○ (saved)
- **Tools**: ■ (brush) - ✅ IMPLEMENTED
- **Colors**: ●●● (color swatches) - ✅ IMPLEMENTED
- **Alerts**: ⚠ (warning) / ✗ (error) / ✓ (success)
- **Navigation**: ◀▶▲▼ (arrows) / ┌┘ (corners)
- **Modal Borders**: ╭─── (rounded corners) instead of ┌─── (square)

### Terminal Enhancement Features:

- [ ] Support advanced terminal text modes:
  - [ ] Dim text rendering
  - [ ] Italic text support
  - [ ] Strikethrough text support
- [ ] Rounded borders for all UI elements (╭─╮ ╰─╯ style)
- [ ] Preferences system with ~/.config/cadubi storage
- [ ] Dual UI theme system:
  - [ ] Dark mode theme (default terminal colors)
  - [ ] Light mode theme (inverted/bright colors)
  - [ ] Theme switching in preferences dialog
- [ ] First-run experience:
  - [ ] Welcome dialog on initial startup
  - [ ] Theme selection prompt for new users
  - [ ] Quick tutorial/onboarding flow

### Two-Panel Layout Design:

**PRIMARY LAYOUT:**
```
╭─ Main Canvas (Drawing Area) ────────────────────────╮
│ ╭─ Drawable Area ─────────────────────────────────╮ │
│ │                                                 │ │
│ │          █ Cursor here                          │ │
│ │                                                 │ │
│ │                 ASCII Art                       │ │
│ │                                                 │ │
│ ╰─────────────────────────────────────────────────╯ │
├─ Status Bar (Bottom) ───────────────────────────────┤
│ ■ Tool: Pen │ ●●● Colors │ Layer 1 │ 25x12 │ 100%   │
╰─────────────────────────────────────────────────────╯
```

**Main Canvas Features (Top):**
- Rounded border in dim grey (same color as status bar) to clearly define drawable area
- Infinite scrolling viewport with pan indicators
- Cursor visualization and position tracking
- Grid overlay (toggle with Ctrl+G)
- Real-time drawing feedback
- Hero element - the primary focus area

**Status Bar Features (Bottom):**
- Current tool and mode indicators
- Color picker access (quick color swatches)
- Layer information (current layer, visibility)
- Canvas dimensions and cursor position
- Zoom level and viewport info
- File status (saved/modified)
- Always visible reference information

### 🎯 Bubble Tea Terminal Handling

**ADVANTAGE:** Bubble Tea provides excellent terminal handling with proper fullscreen support and accurate dimension calculation.

**IMPLEMENTATION:**
```go
// Bubble Tea automatically handles terminal dimensions correctly
func (m model) Init() tea.Cmd {
    return tea.EnterAltScreen  // Enter fullscreen mode
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        // Bubble Tea provides accurate dimensions
    }
    return m, nil
}
```

**BENEFIT:** No workarounds needed - Bubble Tea handles terminal interaction robustly.

### Current Tasks:

- [x] Update LayerPanel emoji (👁🚫) to Unicode symbols (⧉○) - ✅ COMPLETED
- [x] Redesign UI for two-panel layout - ✅ COMPLETED
- [x] Implement status bar component with tool/color/layer info - ✅ COMPLETED  
- [x] Create bordered canvas component with rounded corners - ✅ COMPLETED
- [x] Remove sidebar - integrate color picker into status bar - ✅ COMPLETED
- [x] Audit all UI components for emoji usage - ✅ COMPLETED
- [x] Standardize symbol usage across components - ✅ COMPLETED
- [ ] Create UI symbol reference guide

## Phase 4.5: Modern UI & Automatic Theme Detection ✅ COMPLETED

### Charm's Adaptive Color System ✅ COMPLETED
- [x] **Automatic terminal background detection using Lipgloss AdaptiveColor**
- [x] **Simplified theme system following Charm's best practices:**
  - [x] `borderColor = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}`
  - [x] `statusColor = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}`  
  - [x] `cursorBg = lipgloss.AdaptiveColor{Light: "0", Dark: "7"}`
  - [x] `cursorFg = lipgloss.AdaptiveColor{Light: "15", Dark: "0"}`
- [x] **No manual theme selection needed - automatic adaptation**
- [x] **Removed complex preferences system in favor of simplicity**
- [x] **Clean, minimal codebase following Charm ecosystem patterns**
- [x] **Complete UI component coverage:**
  - [x] Main canvas border and status bar styles
  - [x] Cursor styling with adaptive colors
  - [x] Help modal dialog styling with adaptive colors
  - [x] Modal borders, backgrounds, and tab styling
  - [x] Comprehensive audit completed - no hardcoded colors remaining
  - [x] All UI components use AdaptiveColor patterns throughout

### Advanced Text Formatting
- [ ] Extended ANSI support beyond basic colors:
  - [ ] Dim text (SGR 2) for subtle UI elements
  - [ ] Italic text (SGR 3) for emphasis and labels
  - [ ] Strikethrough (SGR 9) for deleted/disabled items
  - [ ] Double underline (SGR 21) for special highlighting
- [ ] Terminal capability detection
- [ ] Graceful fallbacks for limited terminals

### Streamlined User Experience ✅ COMPLETED  
- [x] **Direct launch to main application - no setup dialogs**
- [x] **Zero configuration required - works out of the box**
- [x] **Automatic terminal adaptation - seamless experience**
- [x] **Consistent behavior across all terminal environments**

### Modern Help System ✅ COMPLETED
- [x] **Bubbles Help Component Integration**: Replaced manual ASCII help boxes with proper Bubbles components
  - [x] Removed 200+ lines of manual help rendering functions (getHelpContent, renderHelpModal, renderHelpModalOverlay)
  - [x] Implemented KeyMap struct with comprehensive key bindings and help text
  - [x] Added ShortHelp() and FullHelp() methods for proper help interface
  - [x] Automatic help text generation from key bindings using help.Model
  - [x] Fixed Unicode character width issues and border misalignment problems
  - [x] Clean, maintainable help system following Charm ecosystem patterns
  - [x] Help display integrated into main View method with full help mode

### Two-Panel UI Architecture
- [ ] StatusBar component (top panel):
  - [ ] Tool indicator with Unicode symbols (■ Pen, ◆ Select, etc.)
  - [ ] Inline color picker with swatches and quick access
  - [ ] Layer info display (name, visibility, opacity)
  - [ ] Canvas stats (dimensions, cursor pos, zoom)
  - [ ] File status indicator (● modified / ○ saved)
- [ ] BorderedCanvas component (main panel):
  - [ ] Rounded border frame in dim grey (#6B7280) around drawable area
  - [ ] Infinite canvas viewport with smooth scrolling
  - [ ] Grid overlay system (Ctrl+G toggle)
  - [ ] Cursor visualization and tracking
  - [ ] Pan indicators when scrolling
- [ ] Modal dialogs (overlays on main canvas):
  - [ ] Help dialog with rounded corners and pagination (< > keys for multiple pages)
  - [ ] Advanced color picker (popup from status bar) with tabbed interface (< > keys)
  - [ ] Preferences dialog with multiple sections navigable via < > keys
  - [ ] Layer management dialog with pagination for many layers
  - [ ] Status messages with bubble styling

## Phase 5: Advanced Drawing Tools 📋 PLANNED

- [ ] Advanced brush system with custom characters
- [ ] Selection tools (rectangular, freehand)
- [ ] Copy/paste operations
- [ ] Flood fill tool
- [ ] Line/shape drawing tools
- [ ] Text tool with font effects

## Phase 6: File & Export System 📋 PLANNED

### Native Format

- [ ] `.cadubi` format: JSON-based with layer data, metadata, version tracking
- [ ] Project files: Save entire workspace including layers, settings, history
- [ ] Template system: Reusable project templates and components

### Import/Export

- [ ] Multiple file format support (.ans, .txt, .cadubi)
- [ ] ANSI color export with high-fidelity terminal output
- [ ] PNG/HTML export for web sharing
- [ ] Import from various ASCII art formats
- [ ] File recovery and autosave

### Advanced I/O

- [ ] Version control integration: Git-friendly diff/merge for .cadubi files
- [ ] Batch processing: Command-line tools for bulk operations
- [ ] Web export: HTML/CSS output for web display

## Phase 7: Professional Features 📋 PLANNED

### Advanced Features

- [ ] Custom color palettes and themes
- [ ] Plugin/extension system
- [ ] Macro recording and playbook
- [ ] Multiple document tabs
- [ ] Advanced layer operations (merge, blend modes)
- [ ] Performance optimizations for large canvases

### Comprehensive ANSI Support

- [ ] Modern terminal parsing: Full SGR sequence support, 24-bit color
- [ ] Advanced text formatting modes:
  - [ ] Dim text (SGR 2)
  - [ ] Italic text (SGR 3)
  - [ ] Strikethrough text (SGR 9)
  - [ ] Double underline (SGR 21)
- [ ] Terminal compatibility: Auto-detect and adapt to terminal capabilities
- [ ] ANSI library integration: Leverage existing parsing libraries

### Resizing & Transformation

- [ ] Smart resize tools: Scale layers without quality loss where possible
- [ ] Canvas resizing: Add/remove canvas space with content preservation
- [ ] Layer transformations: Rotate, skew, flip individual layers

### Developer Tools

- [ ] Library mode: Use CADUBI command system programmatically
- [ ] Scripting API: `create document → write text in red → draw box → get output`
- [ ] Plugin system: Custom tools, effects, and file formats
- [ ] CLI utilities: Headless rendering, batch processing tools

## Phase 8: Distribution & Production 📋 PLANNED

### Single Executable Applications

- [ ] Go cross-compilation: `go build` for Windows, macOS, Linux native executables
- [ ] Zero-dependency deployment: Self-contained static binaries
- [ ] Package manager integration: Homebrew, apt, Go modules, GitHub releases

### Advanced Brush System

- [ ] Procedural brushes: Star fields, snowflakes, grass, clouds with randomization
- [ ] Seed control: `set-seed <value>` command for reproducible randomness
- [ ] Interactive brushes: Brushes that interact with existing canvas content
- [ ] Texture brushes: Pattern-based brushes with rotation and scaling
- [ ] Animation brushes: Brushes that create animated ASCII sequences

### Performance & Optimization

- [ ] Memory efficiency: Lazy loading for large canvases
- [ ] Render optimization: Incremental updates, viewport culling
- [ ] Terminal optimization: Minimize escape sequence output

### User Experience & Preferences

- [ ] Preferences system:
  - [ ] ~/.config/cadubi configuration storage
  - [ ] JSON-based preference files
  - [ ] Theme selection persistence
  - [ ] Keybinding customization
  - [ ] Default tool settings
- [ ] Multi-theme system:
  - [ ] Dark mode (default terminal colors with dim grey borders)
  - [ ] Light mode (bright/inverted colors with adapted borders)
  - [ ] High contrast mode for accessibility
  - [ ] Custom theme support with user-defined color schemes
  - [ ] Dynamic theme switching (Ctrl+T)
  - [ ] Theme-aware color palettes
- [ ] First-run experience:
  - [ ] Welcome screen with mandatory theme selection (light/dark mode)
  - [ ] Theme preference saved to ~/.config/cadubi/settings.json on first run
  - [ ] Quick start tutorial
  - [ ] Preference migration from v1
- [ ] Interactive tutorials: Built-in learning mode
- [ ] Accessibility: Screen reader support, keyboard navigation
- [ ] Documentation: Comprehensive help system, video tutorials

## Technical Architecture

### Core Components

- **CanvasState**: Central state management with infinite canvas (Go struct with methods)
- **CommandHistory**: Undo/redo with command pattern (interface-based design)
- **ColorSystem**: 256-color support with RGB/HSL conversion (Lipgloss integration)
- **StatusBar**: Top panel with tool/color/layer/file information (Bubble Tea component)
- **BorderedCanvas**: Main drawing area with rounded border frame (Lipgloss styling)
- **ResponsiveUI**: Two-panel layout adapting to terminal sizes (Bubble Tea model)
- **KeybindingSystem**: Vim/Photoshop-inspired shortcuts (Bubble Tea key handling)
- **PreferencesManager**: User configuration and theme management (JSON + file I/O)
- **ThemeEngine**: Dynamic light/dark mode switching with rounded UI (Lipgloss themes)
- **ANSIRenderer**: Advanced text formatting leveraging Lipgloss capabilities

### Design Principles

1. **Two-Panel Layout**: Clean separation of controls (status bar) and content (canvas)
2. **Unicode-First UI**: Use box drawing and mathematical symbols
3. **Terminal Native**: Optimized for terminal environments
4. **Rounded Aesthetics**: Modern rounded borders for all UI elements (Lipgloss)
5. **Responsive Design**: Adapts to any terminal size (80x24 to 200x50+)
6. **Professional Workflow**: Modern keybindings and tools
7. **Performance**: Efficient rendering and memory usage with Go's concurrency
8. **Accessibility**: Clear visual hierarchy with Unicode symbols
9. **Code Quality**: Simple, well-structured, maintainable codebase
10. **Code Simplicity**: Easy to understand architecture and clear separation of concerns

### Bubble Tea Architecture

**The Elm Architecture Pattern:**
- **Model**: Application state (canvas, layers, tools, preferences)
- **Update**: State transitions based on messages (key presses, mouse events)
- **View**: Rendering the current state to terminal output

**Key Benefits:**
- **Functional**: Immutable state updates and predictable behavior
- **Concurrent**: Go's goroutines for background tasks (file I/O, auto-save)
- **Composable**: Small, reusable components that combine naturally
- **Testable**: Pure functions for state transitions enable easy testing

### Testing Architecture

**Headless Testing Framework:**
```go
// Example headless test with command sequence
func TestDrawingSequence(t *testing.T) {
    commands := []Command{
        {"type": "new_document", "width": 80, "height": 24},
        {"type": "draw_pixel", "x": 10, "y": 5, "char": "*", "color": "red"},
        {"type": "draw_line", "x1": 0, "y1": 0, "x2": 10, "y2": 10},
        {"type": "save", "filename": "/tmp/test.cadubi"},
    }
    
    model := NewModel().WithoutRenderer()
    for _, cmd := range commands {
        model, _ = model.Update(cmd)
    }
    
    // Validate final state
    assert.Equal(t, expected, model.Canvas.GetPixel(10, 5))
}
```

**Testing Tools:**
- **Built-in Go Testing**: Standard `go test` framework
- **Catwalk Library**: Unit testing for Bubble Tea models and state transitions
- **Teatest Framework**: End-to-end testing scenarios
- **Headless Execution**: `WithoutRenderer()` option for automated testing
- **Command-driven Tests**: JSON arrays of user commands for reproducible testing

### Technology Stack

- **Runtime**: Go 1.21+
- **Language**: Go (with strict linting and formatting)
- **UI Framework**: Bubble Tea for terminal interfaces
- **Styling**: Lipgloss for terminal styling and layouts
- **Components**: Bubbles for common UI components
- **Architecture**: Command pattern, concurrent programming
- **Testing**: Go testing package with comprehensive coverage
- **Package Manager**: Go modules
- **Build System**: Go compiler with cross-compilation support

### Design Patterns

- **Command Pattern**: All actions as undoable commands
- **Layer System**: Compositional rendering with blend modes
- **Plugin Architecture**: Extensible tool and format system
- **Event System**: Decoupled component communication

### File Structure

```
cmd/
└── cadubi/            # Main application entry point
pkg/
├── commands/          # Command pattern implementation
├── canvas/            # Canvas and layer management
├── ui/                # Bubble Tea UI components
├── tools/             # Drawing tools and brushes
├── io/                # File import/export
├── color/             # Color system and picker
├── text/              # Typography and fonts
└── plugins/           # Plugin system and built-in plugins
internal/
└── config/            # Internal configuration and utilities
.claude/
├── CLAUDE.md          # Claude Code instructions
├── OLD_PERL_VERSION.md # Original Perl codebase documentation
└── BUBBLETEA.md       # Bubble Tea library documentation
```

## Quality Standards

### Code Quality & Simplicity
- **Simple Architecture**: Clear separation of concerns, easy to understand codebase
- **Maintainable Code**: Well-structured packages with logical organization
- **Comprehensive Go Types**: Strong type safety and well-defined interfaces
- **Go Best Practices**: gofmt, go vet, golint compliance
- **Documentation**: Clear comments and comprehensive README

### Testing & Reliability
- **100% Test Coverage**: Core logic fully tested with headless execution
- **Command-driven Testing**: JSON test scenarios for reproducible results
- **State Validation**: Automated verification of application state
- **Regression Testing**: UI snapshot testing to prevent visual regressions

### Performance & Design
- **Responsive Design**: Adapts to all screen sizes (80x24 to 200x50+)
- **Consistent Unicode**: Standardized symbol usage throughout UI
- **Professional Shortcuts**: Vim/Photoshop-inspired keyboard shortcuts
- **Memory Efficiency**: Optimal performance with Go's concurrency model

## Success Metrics

- **Performance**: Handle 1000x1000 canvas smoothly
- **Compatibility**: Work across major terminals (iTerm, Windows Terminal, etc.)
- **Usability**: Artists can create complex ASCII art efficiently
- **Extensibility**: Third-party plugins and tools
- **Distribution**: Easy installation across platforms
