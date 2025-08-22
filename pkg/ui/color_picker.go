package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ColorPicker represents an interactive color picker modal
type ColorPicker struct {
	visible       bool
	selectedIndex int
	colors        []ColorInfo
	width         int
	height        int
	gridCols      int
	gridRows      int
}

// ColorInfo represents a color with its display information
type ColorInfo struct {
	Name  string
	Value string
	Style lipgloss.Style
}

// ColorSelectedMsg is sent when a color is selected
type ColorSelectedMsg struct {
	Color string
	Name  string
}

// NewColorPicker creates a new color picker with predefined colors
func NewColorPicker() *ColorPicker {
	colors := []ColorInfo{
		{"White", "white", lipgloss.NewStyle().Foreground(lipgloss.Color("white"))},
		{"Black", "black", lipgloss.NewStyle().Foreground(lipgloss.Color("black"))},
		{"Red", "red", lipgloss.NewStyle().Foreground(lipgloss.Color("red"))},
		{"Green", "green", lipgloss.NewStyle().Foreground(lipgloss.Color("green"))},
		{"Yellow", "yellow", lipgloss.NewStyle().Foreground(lipgloss.Color("yellow"))},
		{"Blue", "blue", lipgloss.NewStyle().Foreground(lipgloss.Color("blue"))},
		{"Magenta", "magenta", lipgloss.NewStyle().Foreground(lipgloss.Color("magenta"))},
		{"Cyan", "cyan", lipgloss.NewStyle().Foreground(lipgloss.Color("cyan"))},
		{"Bright Red", "9", lipgloss.NewStyle().Foreground(lipgloss.Color("9"))},
		{"Bright Green", "10", lipgloss.NewStyle().Foreground(lipgloss.Color("10"))},
		{"Bright Yellow", "11", lipgloss.NewStyle().Foreground(lipgloss.Color("11"))},
		{"Bright Blue", "12", lipgloss.NewStyle().Foreground(lipgloss.Color("12"))},
		{"Bright Magenta", "13", lipgloss.NewStyle().Foreground(lipgloss.Color("13"))},
		{"Bright Cyan", "14", lipgloss.NewStyle().Foreground(lipgloss.Color("14"))},
		{"Gray", "8", lipgloss.NewStyle().Foreground(lipgloss.Color("8"))},
		{"Bright White", "15", lipgloss.NewStyle().Foreground(lipgloss.Color("15"))},
	}

	return &ColorPicker{
		visible:       false,
		selectedIndex: 0,
		colors:        colors,
		width:         60,
		height:        20,
		gridCols:      4,
		gridRows:      4,
	}
}

// Show displays the color picker
func (cp *ColorPicker) Show() {
	cp.visible = true
	cp.selectedIndex = 0
}

// Hide hides the color picker
func (cp *ColorPicker) Hide() {
	cp.visible = false
}

// IsVisible returns true if the color picker is visible
func (cp *ColorPicker) IsVisible() bool {
	return cp.visible
}

// SetDimensions sets the terminal dimensions
func (cp *ColorPicker) SetDimensions(width, height int) {
	cp.width = width
	cp.height = height
}

// Update handles input and state updates
func (cp *ColorPicker) Update(msg tea.Msg) tea.Cmd {
	if !cp.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			cp.selectedIndex -= cp.gridCols
			if cp.selectedIndex < 0 {
				cp.selectedIndex = len(cp.colors) - 1
			}

		case "down", "j":
			cp.selectedIndex += cp.gridCols
			if cp.selectedIndex >= len(cp.colors) {
				cp.selectedIndex = cp.selectedIndex % cp.gridCols
			}

		case "left", "h":
			cp.selectedIndex--
			if cp.selectedIndex < 0 {
				cp.selectedIndex = len(cp.colors) - 1
			}

		case "right", "l":
			cp.selectedIndex++
			if cp.selectedIndex >= len(cp.colors) {
				cp.selectedIndex = 0
			}

		case "enter", " ":
			// Select the current color
			if cp.selectedIndex >= 0 && cp.selectedIndex < len(cp.colors) {
				color := cp.colors[cp.selectedIndex]
				cp.Hide()
				return func() tea.Msg {
					return ColorSelectedMsg{
						Color: color.Value,
						Name:  color.Name,
					}
				}
			}

		case "escape", "q":
			cp.Hide()
		}

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			// Calculate grid position from mouse coordinates
			// This is a simplified version - would need proper coordinate mapping
			gridX := (msg.X - 10) / 12 // Approximate based on color swatch width
			gridY := (msg.Y - 5) / 2   // Approximate based on color swatch height

			if gridX >= 0 && gridX < cp.gridCols && gridY >= 0 && gridY < cp.gridRows {
				index := gridY*cp.gridCols + gridX
				if index >= 0 && index < len(cp.colors) {
					cp.selectedIndex = index
					// Double-click to select
					color := cp.colors[cp.selectedIndex]
					cp.Hide()
					return func() tea.Msg {
						return ColorSelectedMsg{
							Color: color.Value,
							Name:  color.Name,
						}
					}
				}
			}
		}
	}

	return nil
}

// View renders the color picker modal
func (cp *ColorPicker) View() string {
	if !cp.visible {
		return ""
	}

	// Create styles
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Padding(0, 1)

	// Build color grid
	title := titleStyle.Render("Color Picker")
	
	var rows []string
	for row := 0; row < cp.gridRows; row++ {
		var colorSwatches []string
		for col := 0; col < cp.gridCols; col++ {
			index := row*cp.gridCols + col
			if index < len(cp.colors) {
				color := cp.colors[index]
				swatch := color.Style.Render("████")
				
				// Add selection border if this is the selected color
				if index == cp.selectedIndex {
					swatch = selectedStyle.Render(swatch)
				} else {
					swatch = normalStyle.Render(swatch)
				}
				
				colorSwatches = append(colorSwatches, swatch)
			}
		}
		if len(colorSwatches) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, colorSwatches...))
		}
	}

	colorGrid := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Show selected color info
	selectedColorInfo := ""
	if cp.selectedIndex >= 0 && cp.selectedIndex < len(cp.colors) {
		color := cp.colors[cp.selectedIndex]
		selectedColorInfo = fmt.Sprintf("Selected: %s (%s)", color.Name, color.Value)
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Align(lipgloss.Center).
		Render("Navigate: hjkl/arrows · Select: Enter/Space · Cancel: Esc")

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		colorGrid,
		"",
		selectedColorInfo,
		"",
		instructions,
	)

	return modalStyle.Render(content)
}