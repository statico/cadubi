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

### Documentation & Research ✅ COMPLETED

- [x] Read all Bubble Tea documentation using gitingest → `.claude/BUBBLETEA.md`
- [x] Analyzed original CADUBI Perl codebase → `.claude/OLD_PERL_VERSION.md`
- [x] Created comprehensive CLAUDE.md for future instances → `.claude/CLAUDE.md`

### Project Setup ✅ COMPLETED

- [x] Analyzed original Perl codebase (preserved in repository)
- [x] Initialize Go module with proper dependencies
- [x] Set up Go + Bubble Tea + Bubbles + Lipgloss stack
- [x] Created working "Hello World" demo application with Bubble Tea

### Command Pattern Architecture ✅ COMPLETED

- [x] Implement comprehensive Command Pattern system in Go
- [x] Create `CommandHistory` struct for undo/redo functionality
- [x] Design core interfaces: `Command`, `CommandContext`, `CanvasState`, `LayerManager`
- [x] Build foundation commands: `DrawCommand`, `EraseCommand`, `MoveCursorCommand`, `NewDocumentCommand`
- [x] Advanced drawing tools: `LineCommand`, `RectangleCommand` with Bresenham's algorithm

### Test Harness Creation ✅ COMPLETED

- [x] Set up Go testing framework (built-in testing package)
- [x] Implement headless testing with Bubble Tea's `WithoutRenderer` option
- [x] Create command-driven test system: JSON array of commands (create document, draw pixel, save, etc.)
- [x] Build state validation testing: run commands headlessly and verify final state
- [x] Integrate Catwalk testing library for TUI model unit testing
- [x] Set up Teatest for end-to-end testing scenarios (experimental API, disabled for stability)
- [x] Create UI snapshot testing system for Bubble Tea components
- [x] Enable automated UI regression testing with headless execution

**Phase 1 Achievement Summary:**
- **Complete Command Pattern Implementation**: Fully functional undo/redo system with 100% test coverage
- **JSON-Driven Test System**: Advanced testing framework supporting complex scenarios and state validation
- **Headless Testing Framework**: Comprehensive TUI testing without requiring terminal rendering
- **UI Snapshot Testing**: Visual regression testing system with golden file comparison
- **Integration Testing**: TUI model testing with proper state transition validation

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

### 🎨 Minimal Color Styling Guidelines

**Philosophy**: Build functionality first, styling second.
- **Minimize background colors**: Only use when absolutely necessary for clarity (e.g., active tabs)
- **Minimal foreground colors**: Rely on terminal defaults, use color sparingly for highlights
- **Focus on structure over decoration**: Use borders, spacing, and typography for visual hierarchy
- **Adaptive colors only**: When colors are needed, use `lipgloss.AdaptiveColor` patterns
- **Examples of good minimal styling**:
  - Selected items: underline + bold (no background)
  - Navigation tabs: background color only for active state
  - Status indicators: symbols first, color second
  - Error states: bold text before adding color
- **Priority order**: Layout → Typography → Symbols → Colors

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

## Phase 4.6: Comprehensive Mouse Support ✅ COMPLETED

### Mouse Support Design Principles
**REQUIREMENT:** All interactive UI widgets must support mouse interaction when it makes logical sense.

### Core Mouse Interactions ✅ IMPLEMENTED
- [x] **Main Canvas**: Left-click to draw, right-click to erase, wheel scroll for cursor movement
- [x] **Unicode Browser**: Complete mouse support with tabs, character grid, and double-click selection
- [x] **Mouse Event Routing**: Proper event delegation between modals and main app

### Mouse Support Targets ✅ COMPLETED
- [x] **Status Bar Components**:
  - [x] Click on color swatches to cycle colors
  - [x] Click on tool indicator to cycle tools (pen/brush/fill)
  - [x] Click on layer info to toggle layer visibility
  - [x] Click on canvas dimensions (placeholder for future resize dialog)
  - [x] Click on history info for undo/redo actions (left side undo, right side redo)
- [ ] **Help Modal**:
  - [ ] Click on tab headers for navigation
  - [ ] Click outside modal to close
  - [ ] Scroll wheel for content navigation
- [ ] **Color Picker (Future)**:
  - [ ] Click on color swatches for selection
  - [ ] Click on HSL sliders for adjustment
  - [ ] Click on palette tabs for category switching
- [ ] **Layer Management (Future)**:
  - [ ] Click on layer items to select
  - [ ] Click on visibility icons to toggle
  - [ ] Drag and drop for layer reordering
- [ ] **File Operations (Future)**:
  - [ ] Click on menu items
  - [ ] Click on file browser items
  - [ ] Click and drag for selection

### Mouse UX Guidelines
- **Single Click**: Select/navigate/focus elements
- **Double Click**: Apply/confirm actions (like character selection)
- **Right Click**: Context actions or alternative operations
- **Wheel Scroll**: Navigate through lists, zoom, or move cursor
- **Click Outside**: Close modals/dialogs
- **Hover Effects**: Visual feedback where supported by terminal

## Phase 4.7: Comprehensive Trace Logging System ✅ COMPLETED

### Debug & Development Support
**PURPOSE:** Enable comprehensive debugging and development workflow with detailed trace logging.

### Core Logging System ✅ COMPLETED
- [x] **Command Line Flag Support**: `--trace` flag to enable comprehensive trace logging
- [x] **File-based Logging**: Save trace logs to configurable file (defaults to `trace.log`)
- [x] **Structured Log Format**: Timestamp-based structured text format for easy parsing
- [x] **Log Levels**: TRACE, DEBUG, INFO, WARN, ERROR with filtering capabilities
- [x] **Performance Timing**: Timer system for logging execution times
- [x] **Terminal Environment Capture**: Log terminal size, color support, capabilities

### Comprehensive Logging Coverage ✅ COMPLETED
- [x] **Command System**: Log all drawing commands, undo/redo, cursor movements
- [x] **UI Interactions**: Log key presses, mouse events, modal state changes
- [x] **Canvas Operations**: Log pixel operations, layer changes, tool switching  
- [x] **Application Lifecycle**: Log startup, shutdown, and error conditions
- [x] **Error Handling**: Detailed error context and stack traces
- [x] **User Actions**: Track all UI clicks, tool changes, color cycling

### Development Workflow Integration ✅ COMPLETED
- [x] **Build System**: Add `make run/trace` and `make dev/trace` targets for trace-enabled execution
- [x] **Git Integration**: Add `trace.log` to `.gitignore` to prevent commits
- [x] **Claude Code Integration**: Update `CLAUDE.local.md` with trace debugging workflow
- [x] **Trace Package**: Complete logging infrastructure in `pkg/trace/`

### Trace Log Format Design
```
[2025-01-21 10:15:30.123] [TRACE] [command] DrawCommand executed: x=15, y=8, char='*', color=white, layer=0, duration=1.2ms
[2025-01-21 10:15:30.125] [DEBUG] [ui] Mouse click detected: x=45, y=12, button=left, region=canvas
[2025-01-21 10:15:30.126] [INFO] [terminal] Terminal size: 120x40, color_support=true, mouse_support=true
[2025-01-21 10:15:30.128] [TRACE] [canvas] Viewport updated: x=0, y=0, width=100, height=30
```

### Development Benefits
- **Issue Reproduction**: Exact sequence of events leading to bugs
- **Performance Analysis**: Identify bottlenecks in rendering and command execution  
- **User Behavior**: Understand how artists use different features
- **Terminal Compatibility**: Debug issues across different terminal environments
- **Automated Testing**: Generate test scenarios from real usage patterns

## Phase 5: Unicode Explorer & Character Palette System ✅ IN PROGRESS

### Unicode Character Explorer ✅ COMPLETED
- [x] **Unicode Character Browser**: Interactive explorer for non-emoji Unicode characters
  - [x] **Dingbats Plane** (U+2700–U+27BF): ✂✉☎✏✈⚡⚾⛄⭐ etc.
  - [x] **Mathematical Operators** (U+2200–U+22FF): ∀∃∈∉∑∏∫∞√± etc.
  - [x] **Mathematical Symbols** (U+2900–U+297F): ⤀⤁⤂⤃⤄⤅⤆⤇ etc.
  - [x] **Geometric Shapes** (U+25A0–U+25FF): Box drawing, mathematical shapes etc.
  - [x] **Block Elements** (U+2580–U+259F): Terminal block characters
  - [x] **Box Drawing** (U+2500–U+257F): Border and line characters
  - [x] **Arrows** (U+2190–U+21FF): Directional arrows
  - [x] **Miscellaneous Symbols** (U+2600–U+26FF): Various utility symbols
  - [x] Navigation with keyboard shortcuts (arrow keys, hjkl, page up/down)
  - [x] Character preview with Unicode codepoint display
  - [x] Category switching with tab navigation (< >)
  - [x] Modal overlay interface with proper styling
  - [x] Keyboard shortcut 's' to open Unicode browser
  - [x] Enter/Space to select character and auto-close
  - [x] **Mouse Support**: Full mouse interaction capabilities
    - [x] Click on category tabs to switch categories
    - [x] Click on characters to select them
    - [x] Double-click on characters to select and auto-apply
    - [x] Proper modal bounds detection and coordinate mapping

### Quick Palette System ✅ COMPLETED
- [x] **Recent Characters Palette**: Automatic tracking of last 20 used characters
  - [x] Toggle between category and recent view ('r' key)
  - [x] Visual character grid with navigation
  - [x] Most recently used at top, chronological ordering
  - [x] Clear recent history option ('c' key in recent view)
  - [x] Automatic tracking when characters are selected from browser
  - [x] Smart filtering (only non-ASCII characters tracked)
  - [x] Full navigation support within recent characters
- [ ] **Recent Colors Palette**: Track last 20 used colors
  - [ ] Color swatch grid with visual preview
  - [ ] Integration with existing color picker system
  - [ ] Quick access via keyboard shortcuts
- [ ] **Palette Modes**: Different swatch display options
  - [ ] Grid mode (4x5 grid layout)
  - [ ] Linear mode (horizontal strip)
  - [ ] Compact mode (minimal space usage)

### Swatch Management System
- [ ] **Save/Load Character Swatches**: 
  - [ ] Export character palettes to JSON files
  - [ ] Import custom character sets
  - [ ] Named palette system (e.g., "Arrows", "Math", "Decorative")
- [ ] **Save/Load Color Swatches**:
  - [ ] Color palette export/import
  - [ ] Named color schemes (e.g., "Retro", "Neon", "Monochrome")
  - [ ] Palette sharing between documents
- [ ] **Document-Derived Palettes**: Auto-generate palettes from current document
  - [ ] Extract all unique characters used in canvas
  - [ ] Extract all unique colors used in canvas
  - [ ] "Use Document Palette" mode for consistent styling
  - [ ] Palette analysis and statistics

### Palette Integration
- [ ] **Seamless UI Integration**:
  - [ ] Toggle palette sidebar (keeps main two-panel layout)
  - [ ] Palette overlay mode (temporary popup)
  - [ ] Status bar integration for quick access
- [ ] **Palette Persistence**:
  - [ ] Save palette preferences in ~/.config/cadubi/
  - [ ] Per-project palette storage in .cadubi files
  - [ ] Global vs project-specific palette modes

## Phase 6: Advanced Drawing Tools 📋 PLANNED

- [ ] Advanced brush system with custom characters
- [ ] Selection tools (rectangular, freehand)
- [ ] Copy/paste operations
- [ ] Flood fill tool
- [ ] Line/shape drawing tools
- [ ] Text tool with font effects

## Phase 7: File & Export System 📋 IN PROGRESS

### Native Format ✅ PARTIALLY IMPLEMENTED

- [x] **Basic `.cadubi` format**: JSON-based with layer data and version tracking (v2.0)
  - [x] Complete layer preservation (visibility, pixel data with coordinates)  
  - [x] Canvas bounds and dimensions storage
  - [x] Full pixel attributes (char, foreground, background, bold, italic, underline)
  - [x] UTF-8 character support in JSON format
  - [x] Bi-directional save/load functionality with comprehensive test coverage

**Remaining Native Format Tasks:**
- [ ] Enhanced metadata storage (creation date, author, viewport position)
- [ ] Layer names storage (currently generates "Layer N" names)
- [ ] Compressed JSON for large canvases (gzip compression)
- [ ] Command history preservation for session restoration
- [ ] Project files: Save entire workspace including settings and history
- [ ] Template system: Reusable project templates and components

### Core File Format Support

**PRIORITY**: Implement dual file format support for maximum compatibility:

- [x] **Native .cadubi Format (JSON-based)** ✅ BASIC IMPLEMENTATION:
  - [x] Complete layer information preservation (visibility, pixel data)
  - [x] Canvas dimensions and bounds storage  
  - [x] Version tracking (currently v2.0)
  - [x] Full UTF-8 character support
  - [ ] Layer names storage (enhancement needed)
  - [ ] Enhanced metadata (creation date, author, viewport position)
  - [ ] Compressed JSON for large canvases
  - [ ] Command history preservation for session restoration

- [ ] **ANSI UTF-8 Text Format (.ans, .txt)**:
  - [ ] ANSI escape sequence export with full color fidelity
  - [ ] UTF-8 encoding support for Unicode characters and wide chars
  - [ ] Layer flattening for single-layer ANSI output
  - [ ] Import parser for existing ANSI art files
  - [ ] Terminal-ready output (can be cat'd directly to display)
  - [ ] Proper line ending handling (Unix/Windows/Mac compatibility)

### Import/Export

- [ ] **Bi-directional Format Support**:
  - [ ] .cadubi ↔ ANSI conversion with layer flattening options
  - [ ] Import ANSI files as single-layer .cadubi projects
  - [ ] Export .cadubi projects as flattened ANSI with layer merge options
  - [ ] Batch conversion utilities (cadubi-convert tool)

- [ ] **Advanced Export Options**:
  - [ ] PNG/HTML export for web sharing
  - [ ] SVG export for vector graphics
  - [ ] Multiple ANSI variants (16-color, 256-color, 24-bit)
  - [ ] File recovery and autosave in both formats

### Advanced I/O

- [ ] Version control integration: Git-friendly diff/merge for .cadubi files
- [ ] Batch processing: Command-line tools for bulk operations
- [ ] Web export: HTML/CSS output for web display
- [ ] Format detection and automatic handling
- [ ] Migration tools for upgrading file format versions

## Phase 8: Professional Features 📋 PLANNED

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

## Phase 9: Distribution & Production 📋 PLANNED

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
