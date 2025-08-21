package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cadubi/cadubi/pkg/unicode"
)

// UnicodeCategory represents a category of Unicode characters
type UnicodeCategory struct {
	Name       string
	StartCode  int
	EndCode    int
	Characters []rune
}

// UnicodeBrowser manages Unicode character browsing state
type UnicodeBrowser struct {
	categories     []UnicodeCategory
	currentCat     int
	currentChar    int
	selectedChar   rune
	visible        bool
	width          int
	height         int
	charactersPerRow int
	rowsPerPage    int
	recentChars    []rune
	maxRecentChars int
	showRecent     bool
	lastClickTime  time.Time
	lastClickChar  int
}

// NewUnicodeBrowser creates a new Unicode character browser
func NewUnicodeBrowser() *UnicodeBrowser {
	categories := []UnicodeCategory{
		{
			Name:      "Dingbats",
			StartCode: 0x2700,
			EndCode:   0x27BF,
		},
		{
			Name:      "Mathematical Operators",
			StartCode: 0x2200,
			EndCode:   0x22FF,
		},
		{
			Name:      "Geometric Shapes",
			StartCode: 0x25A0,
			EndCode:   0x25FF,
		},
		{
			Name:      "Mathematical Symbols",
			StartCode: 0x2900,
			EndCode:   0x297F,
		},
		{
			Name:      "Block Elements",
			StartCode: 0x2580,
			EndCode:   0x259F,
		},
		{
			Name:      "Box Drawing",
			StartCode: 0x2500,
			EndCode:   0x257F,
		},
		{
			Name:      "Arrows",
			StartCode: 0x2190,
			EndCode:   0x21FF,
		},
		{
			Name:      "Miscellaneous Symbols",
			StartCode: 0x2600,
			EndCode:   0x26FF,
		},
		{
			Name:      "Emoji",
			StartCode: 0x1F600,
			EndCode:   0x1F64F,
		},
		{
			Name:      "Emoji Symbols",
			StartCode: 0x1F300,
			EndCode:   0x1F5FF,
		},
	}

	// Populate characters for each category
	for i := range categories {
		cat := &categories[i]
		for code := cat.StartCode; code <= cat.EndCode; code++ {
			char := rune(code)
			// Filter out control characters and private use areas
			if isPrintableUnicode(char) {
				cat.Characters = append(cat.Characters, char)
			}
		}
	}

	return &UnicodeBrowser{
		categories:       categories,
		currentCat:       0,
		currentChar:      0,
		selectedChar:     '●',
		visible:          false,
		charactersPerRow: 16,
		rowsPerPage:      8,
		recentChars:      []rune{},
		maxRecentChars:   20,
		showRecent:       false,
	}
}

// isPrintableUnicode checks if a Unicode character is printable and safe to display
func isPrintableUnicode(r rune) bool {
	// Basic checks for printable Unicode characters
	if r < 32 {
		return false // Control characters
	}
	if r >= 0xFFF0 && r <= 0xFFFF {
		return false // Specials block
	}
	if r >= 0xFDD0 && r <= 0xFDEF {
		return false // Arabic Presentation Forms-A noncharacters
	}
	// Check for surrogate pairs (invalid in UTF-8)
	if r >= 0xD800 && r <= 0xDFFF {
		return false
	}
	return true
}

// Toggle shows/hides the Unicode browser
func (ub *UnicodeBrowser) Toggle() {
	ub.visible = !ub.visible
}

// Show displays the Unicode browser
func (ub *UnicodeBrowser) Show() {
	ub.visible = true
}

// Hide conceals the Unicode browser
func (ub *UnicodeBrowser) Hide() {
	ub.visible = false
}

// IsVisible returns whether the browser is currently visible
func (ub *UnicodeBrowser) IsVisible() bool {
	return ub.visible
}

// Update handles input events for the Unicode browser
func (ub *UnicodeBrowser) Update(msg tea.Msg) tea.Cmd {
	if !ub.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ub.width = msg.Width
		ub.height = msg.Height
		// Calculate modal content width
		modalWidth := ub.width - 8  // Account for screen margins
		if modalWidth < 60 {
			modalWidth = 60
		}
		if modalWidth > 120 {
			modalWidth = 120
		}
		// Store available width for dynamic row calculation
		ub.width = msg.Width
		ub.height = msg.Height
		// Set reasonable defaults that will be dynamically adjusted
		ub.charactersPerRow = 16  // This will be recalculated per row
		ub.rowsPerPage = max(4, min(12, (msg.Height-12)/2))

	case tea.KeyMsg:
		return ub.handleKeyPress(msg)
		
	case tea.MouseMsg:
		return ub.handleMouseEvent(msg)
	}

	return nil
}

// handleKeyPress processes keyboard input for the Unicode browser
func (ub *UnicodeBrowser) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	var maxChars int
	
	// Determine max characters based on current view
	if ub.showRecent {
		maxChars = len(ub.recentChars)
		if maxChars == 0 {
			return nil
		}
	} else {
		if len(ub.categories) == 0 || len(ub.categories[ub.currentCat].Characters) == 0 {
			return nil
		}
		currentCategory := &ub.categories[ub.currentCat]
		maxChars = len(currentCategory.Characters)
	}

	switch msg.String() {
	case "left", "h":
		if ub.currentChar > 0 {
			ub.currentChar--
		}
	case "right", "l":
		if ub.currentChar < maxChars-1 {
			ub.currentChar++
		}
	case "up", "k":
		ub.navigateUp(maxChars)
	case "down", "j":
		ub.navigateDown(maxChars)
	case "pageup", "ctrl+b":
		newPos := ub.currentChar - (ub.charactersPerRow * ub.rowsPerPage)
		if newPos < 0 {
			newPos = 0
		}
		ub.currentChar = newPos
	case "pagedown", "ctrl+f":
		newPos := ub.currentChar + (ub.charactersPerRow * ub.rowsPerPage)
		if newPos >= maxChars {
			newPos = maxChars - 1
		}
		ub.currentChar = newPos
	case "tab", ">":
		// Next category (only in category view)
		if !ub.showRecent && ub.currentCat < len(ub.categories)-1 {
			ub.currentCat++
			ub.currentChar = 0
		}
	case "shift+tab", "<":
		// Previous category (only in category view)
		if !ub.showRecent && ub.currentCat > 0 {
			ub.currentCat--
			ub.currentChar = 0
		}
	case "r":
		// Toggle recent view
		ub.ToggleRecentView()
	case "c":
		// Clear recent characters (only in recent view)
		if ub.showRecent {
			ub.ClearRecentChars()
		}
	case "enter", " ":
		// Select current character
		if ub.showRecent {
			if ub.currentChar < len(ub.recentChars) {
				ub.selectedChar = ub.recentChars[ub.currentChar]
			}
		} else {
			currentCategory := &ub.categories[ub.currentCat]
			if ub.currentChar < len(currentCategory.Characters) {
				ub.selectedChar = currentCategory.Characters[ub.currentChar]
			}
		}
	case "escape", "esc", "q":
		ub.Hide()
	}

	return nil
}

// handleMouseEvent processes mouse input for the Unicode browser
func (ub *UnicodeBrowser) handleMouseEvent(msg tea.MouseMsg) tea.Cmd {
	if msg.Type != tea.MouseLeft {
		return nil
	}

	// Calculate modal bounds - we need to determine where the modal is positioned
	modalWidth := ub.width - 8
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 120 {
		modalWidth = 120
	}
	
	// Modal is centered on screen
	modalX := (ub.width - modalWidth) / 2
	modalY := 2 // Approximate modal top position with margin
	
	// Check if click is within modal bounds
	if msg.X < modalX || msg.X > modalX+modalWidth || msg.Y < modalY {
		return nil
	}
	
	// Convert mouse coordinates relative to modal content
	relativeX := msg.X - modalX - 2 // Account for border and padding
	relativeY := msg.Y - modalY - 1 // Account for border
	
	// Handle category tab clicks (around line 4-5 of modal content)
	if relativeY >= 3 && relativeY <= 4 && !ub.showRecent {
		return ub.handleCategoryTabClick(relativeX)
	}
	
	// Handle recent/category toggle area click (line 3-4 for recent view)
	if relativeY >= 3 && relativeY <= 4 && ub.showRecent {
		// Could add click to switch back to categories
		return nil
	}
	
	// Handle character grid clicks (around lines 6-14 of modal content)
	if relativeY >= 6 && relativeY <= 14 {
		return ub.handleCharacterGridClick(relativeX, relativeY-6)
	}
	
	return nil
}

// handleCategoryTabClick processes clicks on category tabs
func (ub *UnicodeBrowser) handleCategoryTabClick(x int) tea.Cmd {
	if ub.showRecent {
		return nil
	}
	
	// Estimate tab positions - each category name plus spacing
	tabPositions := []int{0} // Start positions for each tab
	currentPos := 0
	
	for i, cat := range ub.categories {
		tabWidth := len(cat.Name) + 2 // Name plus padding
		if x >= currentPos && x < currentPos+tabWidth {
			ub.currentCat = i
			ub.currentChar = 0
			break
		}
		currentPos += tabWidth + 1 // Add space between tabs
		tabPositions = append(tabPositions, currentPos)
	}
	
	return nil
}

// CharacterSelectedMsg is sent when a character is double-clicked for selection
type CharacterSelectedMsg struct {
	Character rune
}

// handleCharacterGridClick processes clicks on character grid
func (ub *UnicodeBrowser) handleCharacterGridClick(x, y int) tea.Cmd {
	var characters []rune
	var maxChars int
	
	if ub.showRecent {
		characters = ub.recentChars
		maxChars = len(ub.recentChars)
	} else {
		if ub.currentCat >= len(ub.categories) {
			return nil
		}
		characters = ub.categories[ub.currentCat].Characters
		maxChars = len(characters)
	}
	
	if maxChars == 0 {
		return nil
	}
	
	// Calculate character position from click coordinates
	// Each character takes approximately 4 spaces (to account for wide chars), so x position / 4 gives column
	col := x / 4
	row := y
	
	if col >= ub.charactersPerRow || row >= ub.rowsPerPage {
		return nil
	}
	
	// Calculate character index
	currentPage := ub.currentChar / (ub.charactersPerRow * ub.rowsPerPage)
	startIdx := currentPage * ub.charactersPerRow * ub.rowsPerPage
	charIdx := startIdx + (row * ub.charactersPerRow) + col
	
	if charIdx >= maxChars {
		return nil
	}
	
	// Check for double-click (within 500ms)
	now := time.Now()
	isDoubleClick := charIdx == ub.lastClickChar && 
		now.Sub(ub.lastClickTime) < 500*time.Millisecond
	
	// Update click tracking
	ub.lastClickTime = now
	ub.lastClickChar = charIdx
	
	// Select the clicked character
	ub.currentChar = charIdx
	ub.selectedChar = characters[charIdx]
	
	// If double-click, send selection message
	if isDoubleClick {
		return func() tea.Msg {
			return CharacterSelectedMsg{Character: ub.selectedChar}
		}
	}
	
	return nil
}

// GetSelectedChar returns the currently selected character
func (ub *UnicodeBrowser) GetSelectedChar() rune {
	return ub.selectedChar
}

// AddRecentChar adds a character to the recent characters list
func (ub *UnicodeBrowser) AddRecentChar(char rune) {
	// Remove the character if it already exists
	for i, existing := range ub.recentChars {
		if existing == char {
			ub.recentChars = append(ub.recentChars[:i], ub.recentChars[i+1:]...)
			break
		}
	}
	
	// Add to the front
	ub.recentChars = append([]rune{char}, ub.recentChars...)
	
	// Trim to max size
	if len(ub.recentChars) > ub.maxRecentChars {
		ub.recentChars = ub.recentChars[:ub.maxRecentChars]
	}
}

// GetRecentChars returns the recent characters list
func (ub *UnicodeBrowser) GetRecentChars() []rune {
	return ub.recentChars
}

// ClearRecentChars clears the recent characters list
func (ub *UnicodeBrowser) ClearRecentChars() {
	ub.recentChars = []rune{}
}

// ToggleRecentView toggles between recent characters and category view
func (ub *UnicodeBrowser) ToggleRecentView() {
	ub.showRecent = !ub.showRecent
	ub.currentChar = 0 // Reset position when switching views
}

// View renders the Unicode browser
func (ub *UnicodeBrowser) View() string {
	if !ub.visible {
		return ""
	}

	// Check if we have content to show
	if ub.showRecent {
		if len(ub.recentChars) == 0 {
			// Show empty recent message
			return ub.renderEmptyRecent()
		}
	} else {
		if len(ub.categories) == 0 || len(ub.categories[ub.currentCat].Characters) == 0 {
			return ""
		}
	}

	// Create modal styling
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(modalBorder).
		Background(modalBackground).
		Padding(1, 2).
		Margin(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(modalBorder).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Background(activeTabBg).
		Foreground(activeTabFg).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(inactiveTabFg)

	var title, navigation, charGrid, charInfo, instructions string

	if ub.showRecent {
		title = "Recent Characters"
		navigation = ub.renderRecentHeader()
		charGrid = ub.renderRecentGrid(selectedStyle, normalStyle)
		
		// Character info for recent view
		selectedChar := ub.selectedChar
		if ub.currentChar < len(ub.recentChars) {
			selectedChar = ub.recentChars[ub.currentChar]
		}
		charInfo = fmt.Sprintf("Selected: %c (U+%04X) · Position: %d/%d",
			selectedChar, int(selectedChar), ub.currentChar+1, len(ub.recentChars))
		
		// Instructions for recent view
		instructions = "Navigation: ← → ↑ ↓ · Page: PgUp/PgDn · Toggle: r · Clear: c · Select: Enter/Space · Click: select · Double-click: apply · Close: Esc/q"
	} else {
		title = "Unicode Character Browser"
		navigation = ub.renderCategoryTabs()
		charGrid = ub.renderCharacterGrid(selectedStyle, normalStyle)
		
		// Character info for category view
		currentCategory := &ub.categories[ub.currentCat]
		selectedChar := ub.selectedChar
		if ub.currentChar < len(currentCategory.Characters) {
			selectedChar = currentCategory.Characters[ub.currentChar]
		}
		charInfo = fmt.Sprintf("Selected: %c (U+%04X) · Position: %d/%d",
			selectedChar, int(selectedChar), ub.currentChar+1, len(currentCategory.Characters))
		
		// Instructions for category view
		instructions = "Navigation: ← → ↑ ↓ · Page: PgUp/PgDn · Category: < > · Recent: r · Select: Enter/Space · Click: select · Double-click: apply · Close: Esc/q"
	}

	// Build content
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(title),
		"",
		navigation,
		"",
		charGrid,
		"",
		charInfo,
		"",
		instructions,
	)

	// Calculate modal dimensions
	maxWidth := ub.width - 8
	if maxWidth < 60 {
		maxWidth = 60
	}
	if maxWidth > 120 {
		maxWidth = 120
	}

	modal := modalStyle.
		Width(maxWidth).
		MaxWidth(maxWidth).
		Render(content)

	return modal
}

// renderCategoryTabs renders category navigation tabs
func (ub *UnicodeBrowser) renderCategoryTabs() string {
	var tabs []string

	for i, cat := range ub.categories {
		if i == ub.currentCat {
			tab := lipgloss.NewStyle().
				Background(activeTabBg).
				Foreground(activeTabFg).
				Bold(true).
				Padding(0, 1).
				Render(cat.Name)
			tabs = append(tabs, tab)
		} else {
			tab := lipgloss.NewStyle().
				Foreground(inactiveTabFg).
				Padding(0, 1).
				Render(cat.Name)
			tabs = append(tabs, tab)
		}
	}

	return strings.Join(tabs, " ")
}

// renderCharacterGrid renders the character selection grid with proper width handling
func (ub *UnicodeBrowser) renderCharacterGrid(selectedStyle, normalStyle lipgloss.Style) string {
	currentCategory := &ub.categories[ub.currentCat]
	characters := currentCategory.Characters

	if len(characters) == 0 {
		return normalStyle.Render("No characters available in this category")
	}

	contentWidth := ub.getContentWidth()
	currentIdx := 0
	var lines []string
	
	// Find which page we should show based on current character
	targetRowStart := 0
	charsSoFar := 0
	
	for row := 0; row < 100 && charsSoFar <= ub.currentChar; row++ { // Limit iterations for safety
		rowChars := ub.calculateCharsPerRow(characters, charsSoFar, len(characters), contentWidth)
		if charsSoFar <= ub.currentChar && ub.currentChar < charsSoFar+rowChars {
			targetRowStart = row
			break
		}
		charsSoFar += rowChars
	}
	
	// Calculate page start (show a few rows before target if possible)
	pageStartRow := max(0, targetRowStart-ub.rowsPerPage/2)
	
	// Skip to page start
	currentIdx = 0
	for row := 0; row < pageStartRow && currentIdx < len(characters); row++ {
		charsInRow := ub.calculateCharsPerRow(characters, currentIdx, len(characters), contentWidth)
		currentIdx += charsInRow
	}

	// Render the visible rows
	for row := 0; row < ub.rowsPerPage && currentIdx < len(characters); row++ {
		line := ub.renderCharacterRow(characters, currentIdx, len(characters), contentWidth, selectedStyle, normalStyle)
		if line != "" {
			lines = append(lines, line)
		}
		charsInRow := ub.calculateCharsPerRow(characters, currentIdx, len(characters), contentWidth)
		currentIdx += charsInRow
	}

	if len(lines) == 0 {
		return normalStyle.Render("No characters to display")
	}

	return strings.Join(lines, "\n")
}

// calculateCharsPerRow determines how many characters fit in a row given the available width
func (ub *UnicodeBrowser) calculateCharsPerRow(characters []rune, startIdx, totalChars, availableWidth int) int {
	if startIdx >= totalChars {
		return 0
	}
	
	currentWidth := 0
	charCount := 0
	
	for i := startIdx; i < totalChars && charCount < 32; i++ { // Limit to reasonable number
		char := characters[i]
		charDisplayWidth := unicode.CharWidth(char)
		// Each character needs its display width + 2 spaces padding
		charTotalWidth := charDisplayWidth + 2
		
		if currentWidth+charTotalWidth > availableWidth {
			break
		}
		
		currentWidth += charTotalWidth
		charCount++
	}
	
	return max(1, charCount) // Always show at least 1 character per row
}

// renderCharacterRow renders a single row of characters
func (ub *UnicodeBrowser) renderCharacterRow(characters []rune, startIdx, totalChars, availableWidth int, selectedStyle, normalStyle lipgloss.Style) string {
	if startIdx >= totalChars {
		return ""
	}
	
	charsInRow := ub.calculateCharsPerRow(characters, startIdx, totalChars, availableWidth)
	var line strings.Builder
	
	for i := 0; i < charsInRow; i++ {
		charIdx := startIdx + i
		if charIdx >= totalChars {
			break
		}
		
		char := characters[charIdx]
		charDisplayWidth := unicode.CharWidth(char)
		
		// Create character string with appropriate padding
		var charStr string
		if charDisplayWidth == 1 {
			charStr = fmt.Sprintf(" %c ", char) // 1-width char + 2 spaces = 3 total
		} else {
			charStr = fmt.Sprintf(" %c", char) // 2-width char + 1 space = 3 total
		}
		
		if charIdx == ub.currentChar {
			line.WriteString(selectedStyle.Render(charStr))
		} else {
			line.WriteString(normalStyle.Render(charStr))
		}
	}
	
	return line.String()
}

// renderEmptyRecent renders the empty recent characters view
func (ub *UnicodeBrowser) renderEmptyRecent() string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(modalBorder).
		Background(modalBackground).
		Padding(1, 2).
		Margin(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(modalBorder).
		Align(lipgloss.Center)

	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Recent Characters"),
		"",
		"No recent characters yet.",
		"Use characters from the Unicode browser to populate this list.",
		"",
		"Press 'r' to return to category view · Press 'Esc' to close",
	)

	return modalStyle.Width(60).Render(content)
}

// renderRecentHeader renders the header for recent characters view
func (ub *UnicodeBrowser) renderRecentHeader() string {
	return fmt.Sprintf("Recent Characters (%d/%d used)", len(ub.recentChars), ub.maxRecentChars)
}

// renderRecentGrid renders the recent characters grid with proper width handling
func (ub *UnicodeBrowser) renderRecentGrid(selectedStyle, normalStyle lipgloss.Style) string {
	if len(ub.recentChars) == 0 {
		return normalStyle.Render("No recent characters")
	}

	contentWidth := ub.getContentWidth()
	currentIdx := 0
	var lines []string
	
	// Find which page we should show based on current character
	targetRowStart := 0
	charsSoFar := 0
	
	for row := 0; row < 100 && charsSoFar <= ub.currentChar; row++ { // Limit iterations for safety
		rowChars := ub.calculateCharsPerRowRecent(ub.recentChars, charsSoFar, len(ub.recentChars), contentWidth)
		if charsSoFar <= ub.currentChar && ub.currentChar < charsSoFar+rowChars {
			targetRowStart = row
			break
		}
		charsSoFar += rowChars
	}
	
	// Calculate page start (show a few rows before target if possible)
	pageStartRow := max(0, targetRowStart-ub.rowsPerPage/2)
	
	// Skip to page start
	currentIdx = 0
	for row := 0; row < pageStartRow && currentIdx < len(ub.recentChars); row++ {
		charsInRow := ub.calculateCharsPerRowRecent(ub.recentChars, currentIdx, len(ub.recentChars), contentWidth)
		currentIdx += charsInRow
	}

	// Render the visible rows
	for row := 0; row < ub.rowsPerPage && currentIdx < len(ub.recentChars); row++ {
		line := ub.renderCharacterRowRecent(ub.recentChars, currentIdx, len(ub.recentChars), contentWidth, selectedStyle, normalStyle)
		if line != "" {
			lines = append(lines, line)
		}
		charsInRow := ub.calculateCharsPerRowRecent(ub.recentChars, currentIdx, len(ub.recentChars), contentWidth)
		currentIdx += charsInRow
	}

	if len(lines) == 0 {
		return normalStyle.Render("No recent characters to display")
	}

	return strings.Join(lines, "\n")
}

// calculateCharsPerRowRecent determines how many recent characters fit in a row
func (ub *UnicodeBrowser) calculateCharsPerRowRecent(characters []rune, startIdx, totalChars, availableWidth int) int {
	if startIdx >= totalChars {
		return 0
	}
	
	currentWidth := 0
	charCount := 0
	
	for i := startIdx; i < totalChars && charCount < 32; i++ { // Limit to reasonable number
		char := characters[i]
		charDisplayWidth := unicode.CharWidth(char)
		// Each character needs its display width + 2 spaces padding
		charTotalWidth := charDisplayWidth + 2
		
		if currentWidth+charTotalWidth > availableWidth {
			break
		}
		
		currentWidth += charTotalWidth
		charCount++
	}
	
	return max(1, charCount) // Always show at least 1 character per row
}

// renderCharacterRowRecent renders a single row of recent characters
func (ub *UnicodeBrowser) renderCharacterRowRecent(characters []rune, startIdx, totalChars, availableWidth int, selectedStyle, normalStyle lipgloss.Style) string {
	if startIdx >= totalChars {
		return ""
	}
	
	charsInRow := ub.calculateCharsPerRowRecent(characters, startIdx, totalChars, availableWidth)
	var line strings.Builder
	
	for i := 0; i < charsInRow; i++ {
		charIdx := startIdx + i
		if charIdx >= totalChars {
			break
		}
		
		char := characters[charIdx]
		charDisplayWidth := unicode.CharWidth(char)
		
		// Create character string with appropriate padding
		var charStr string
		if charDisplayWidth == 1 {
			charStr = fmt.Sprintf(" %c ", char) // 1-width char + 2 spaces = 3 total
		} else {
			charStr = fmt.Sprintf(" %c", char) // 2-width char + 1 space = 3 total
		}
		
		if charIdx == ub.currentChar {
			line.WriteString(selectedStyle.Render(charStr))
		} else {
			line.WriteString(normalStyle.Render(charStr))
		}
	}
	
	return line.String()
}

// getContentWidth calculates the available width for character content
func (ub *UnicodeBrowser) getContentWidth() int {
	modalWidth := ub.width - 8  // Account for screen margins
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 120 {
		modalWidth = 120
	}
	// Account for modal border (2) + padding (4) = 6 chars
	return modalWidth - 6
}

// navigateUp moves selection up one row, handling variable width rows
func (ub *UnicodeBrowser) navigateUp(maxChars int) {
	if ub.currentChar <= 0 {
		return
	}
	
	var characters []rune
	if ub.showRecent {
		characters = ub.recentChars
	} else {
		characters = ub.categories[ub.currentCat].Characters
	}
	
	if len(characters) == 0 {
		return
	}
	
	contentWidth := ub.getContentWidth()
	
	// Find current row start
	currentRowStart := 0
	for i := 0; i < ub.currentChar; {
		var charsInRow int
		if ub.showRecent {
			charsInRow = ub.calculateCharsPerRowRecent(characters, i, len(characters), contentWidth)
		} else {
			charsInRow = ub.calculateCharsPerRow(characters, i, len(characters), contentWidth)
		}
		if i + charsInRow > ub.currentChar {
			currentRowStart = i
			break
		}
		i += charsInRow
	}
	
	// If already at first row, do nothing
	if currentRowStart == 0 {
		return
	}
	
	// Find previous row start
	prevRowStart := 0
	for i := 0; i < currentRowStart; {
		var charsInRow int
		if ub.showRecent {
			charsInRow = ub.calculateCharsPerRowRecent(characters, i, len(characters), contentWidth)
		} else {
			charsInRow = ub.calculateCharsPerRow(characters, i, len(characters), contentWidth)
		}
		if i + charsInRow >= currentRowStart {
			prevRowStart = i
			break
		}
		i += charsInRow
	}
	
	// Calculate position in current row
	positionInRow := ub.currentChar - currentRowStart
	
	// Find chars in previous row
	var prevRowChars int
	if ub.showRecent {
		prevRowChars = ub.calculateCharsPerRowRecent(characters, prevRowStart, len(characters), contentWidth)
	} else {
		prevRowChars = ub.calculateCharsPerRow(characters, prevRowStart, len(characters), contentWidth)
	}
	
	// Move to corresponding position in previous row (or last char if row is shorter)
	newPos := prevRowStart + min(positionInRow, prevRowChars-1)
	if newPos >= 0 && newPos < maxChars {
		ub.currentChar = newPos
	}
}

// navigateDown moves selection down one row, handling variable width rows  
func (ub *UnicodeBrowser) navigateDown(maxChars int) {
	if ub.currentChar >= maxChars-1 {
		return
	}
	
	var characters []rune
	if ub.showRecent {
		characters = ub.recentChars
	} else {
		characters = ub.categories[ub.currentCat].Characters
	}
	
	if len(characters) == 0 {
		return
	}
	
	contentWidth := ub.getContentWidth()
	
	// Find current row start and chars in current row
	currentRowStart := 0
	currentRowChars := 0
	for i := 0; i <= ub.currentChar; {
		var charsInRow int
		if ub.showRecent {
			charsInRow = ub.calculateCharsPerRowRecent(characters, i, len(characters), contentWidth)
		} else {
			charsInRow = ub.calculateCharsPerRow(characters, i, len(characters), contentWidth)
		}
		if i + charsInRow > ub.currentChar {
			currentRowStart = i
			currentRowChars = charsInRow
			break
		}
		i += charsInRow
	}
	
	// Calculate next row start
	nextRowStart := currentRowStart + currentRowChars
	if nextRowStart >= maxChars {
		return // Already at last row
	}
	
	// Calculate position in current row
	positionInRow := ub.currentChar - currentRowStart
	
	// Find chars in next row
	var nextRowChars int
	if ub.showRecent {
		nextRowChars = ub.calculateCharsPerRowRecent(characters, nextRowStart, len(characters), contentWidth)
	} else {
		nextRowChars = ub.calculateCharsPerRow(characters, nextRowStart, len(characters), contentWidth)
	}
	
	// Move to corresponding position in next row (or last char if row is shorter)
	newPos := nextRowStart + min(positionInRow, nextRowChars-1)
	if newPos >= 0 && newPos < maxChars {
		ub.currentChar = newPos
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}