package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		// Adjust layout based on terminal size
		ub.charactersPerRow = max(8, min(32, (msg.Width-4)/3))
		ub.rowsPerPage = max(4, min(16, (msg.Height-8)/2))

	case tea.KeyMsg:
		return ub.handleKeyPress(msg)
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
		newPos := ub.currentChar - ub.charactersPerRow
		if newPos >= 0 {
			ub.currentChar = newPos
		}
	case "down", "j":
		newPos := ub.currentChar + ub.charactersPerRow
		if newPos < maxChars {
			ub.currentChar = newPos
		}
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
		charInfo = fmt.Sprintf("Selected: %c (U+%04X) │ Position: %d/%d",
			selectedChar, int(selectedChar), ub.currentChar+1, len(ub.recentChars))
		
		// Instructions for recent view
		instructions = "Navigation: ← → ↑ ↓ │ Page: PgUp/PgDn │ Toggle: r │ Clear: c │ Select: Enter/Space │ Close: Esc/q"
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
		charInfo = fmt.Sprintf("Selected: %c (U+%04X) │ Position: %d/%d",
			selectedChar, int(selectedChar), ub.currentChar+1, len(currentCategory.Characters))
		
		// Instructions for category view
		instructions = "Navigation: ← → ↑ ↓ │ Page: PgUp/PgDn │ Category: < > │ Recent: r │ Select: Enter/Space │ Close: Esc/q"
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

// renderCharacterGrid renders the character selection grid
func (ub *UnicodeBrowser) renderCharacterGrid(selectedStyle, normalStyle lipgloss.Style) string {
	currentCategory := &ub.categories[ub.currentCat]
	characters := currentCategory.Characters

	if len(characters) == 0 {
		return normalStyle.Render("No characters available in this category")
	}

	// Calculate visible range for current page
	currentPage := ub.currentChar / (ub.charactersPerRow * ub.rowsPerPage)
	startIdx := currentPage * ub.charactersPerRow * ub.rowsPerPage
	endIdx := min(startIdx+(ub.charactersPerRow*ub.rowsPerPage), len(characters))

	var lines []string
	for row := 0; row < ub.rowsPerPage; row++ {
		var line strings.Builder
		
		for col := 0; col < ub.charactersPerRow; col++ {
			charIdx := startIdx + (row * ub.charactersPerRow) + col
			
			if charIdx >= endIdx {
				line.WriteString("   ") // Empty space for alignment
				continue
			}

			char := characters[charIdx]
			charStr := fmt.Sprintf(" %c ", char)

			if charIdx == ub.currentChar {
				line.WriteString(selectedStyle.Render(charStr))
			} else {
				line.WriteString(normalStyle.Render(charStr))
			}
		}
		
		if line.Len() > 0 {
			lines = append(lines, line.String())
		}
	}

	return strings.Join(lines, "\n")
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
		"Press 'r' to return to category view │ Press 'Esc' to close",
	)

	return modalStyle.Width(60).Render(content)
}

// renderRecentHeader renders the header for recent characters view
func (ub *UnicodeBrowser) renderRecentHeader() string {
	return fmt.Sprintf("Recent Characters (%d/%d used)", len(ub.recentChars), ub.maxRecentChars)
}

// renderRecentGrid renders the recent characters grid
func (ub *UnicodeBrowser) renderRecentGrid(selectedStyle, normalStyle lipgloss.Style) string {
	if len(ub.recentChars) == 0 {
		return normalStyle.Render("No recent characters")
	}

	// Calculate visible range for current page
	currentPage := ub.currentChar / (ub.charactersPerRow * ub.rowsPerPage)
	startIdx := currentPage * ub.charactersPerRow * ub.rowsPerPage
	endIdx := min(startIdx+(ub.charactersPerRow*ub.rowsPerPage), len(ub.recentChars))

	var lines []string
	for row := 0; row < ub.rowsPerPage; row++ {
		var line strings.Builder
		
		for col := 0; col < ub.charactersPerRow; col++ {
			charIdx := startIdx + (row * ub.charactersPerRow) + col
			
			if charIdx >= endIdx {
				line.WriteString("   ") // Empty space for alignment
				continue
			}

			char := ub.recentChars[charIdx]
			charStr := fmt.Sprintf(" %c ", char)

			if charIdx == ub.currentChar {
				line.WriteString(selectedStyle.Render(charStr))
			} else {
				line.WriteString(normalStyle.Render(charStr))
			}
		}
		
		if line.Len() > 0 {
			lines = append(lines, line.String())
		}
	}

	return strings.Join(lines, "\n")
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