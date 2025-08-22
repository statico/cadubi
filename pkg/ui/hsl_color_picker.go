package ui

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HSLColorPicker represents an advanced HSL color wheel with 256 colors
type HSLColorPicker struct {
	visible       bool
	selectedIndex int
	colors        []HSLColor
	width         int
	height        int
	mode          string // "wheel", "palette", "recent"
	recentColors  []HSLColor
}

// HSLColor represents a color in HSL space with terminal color code
type HSLColor struct {
	Name        string
	Hue         int   // 0-360
	Saturation  int   // 0-100
	Lightness   int   // 0-100
	TerminalCode string
	Style       lipgloss.Style
}

// NewHSLColorPicker creates a new HSL color picker with 256 colors
func NewHSLColorPicker() *HSLColorPicker {
	colors := generateHSLColorWheel()
	
	return &HSLColorPicker{
		visible:       false,
		selectedIndex: 0,
		colors:        colors,
		width:         70,
		height:        25,
		mode:          "wheel",
		recentColors:  []HSLColor{},
	}
}

// generateHSLColorWheel creates a comprehensive 256-color palette
func generateHSLColorWheel() []HSLColor {
	var colors []HSLColor
	
	// Standard 16 colors (0-15)
	standardColors := []HSLColor{
		{"Black", 0, 0, 0, "0", lipgloss.NewStyle().Foreground(lipgloss.Color("0"))},
		{"Dark Red", 0, 100, 25, "1", lipgloss.NewStyle().Foreground(lipgloss.Color("1"))},
		{"Dark Green", 120, 100, 25, "2", lipgloss.NewStyle().Foreground(lipgloss.Color("2"))},
		{"Dark Yellow", 60, 100, 25, "3", lipgloss.NewStyle().Foreground(lipgloss.Color("3"))},
		{"Dark Blue", 240, 100, 25, "4", lipgloss.NewStyle().Foreground(lipgloss.Color("4"))},
		{"Dark Magenta", 300, 100, 25, "5", lipgloss.NewStyle().Foreground(lipgloss.Color("5"))},
		{"Dark Cyan", 180, 100, 25, "6", lipgloss.NewStyle().Foreground(lipgloss.Color("6"))},
		{"Light Gray", 0, 0, 75, "7", lipgloss.NewStyle().Foreground(lipgloss.Color("7"))},
		{"Dark Gray", 0, 0, 50, "8", lipgloss.NewStyle().Foreground(lipgloss.Color("8"))},
		{"Red", 0, 100, 50, "9", lipgloss.NewStyle().Foreground(lipgloss.Color("9"))},
		{"Green", 120, 100, 50, "10", lipgloss.NewStyle().Foreground(lipgloss.Color("10"))},
		{"Yellow", 60, 100, 50, "11", lipgloss.NewStyle().Foreground(lipgloss.Color("11"))},
		{"Blue", 240, 100, 50, "12", lipgloss.NewStyle().Foreground(lipgloss.Color("12"))},
		{"Magenta", 300, 100, 50, "13", lipgloss.NewStyle().Foreground(lipgloss.Color("13"))},
		{"Cyan", 180, 100, 50, "14", lipgloss.NewStyle().Foreground(lipgloss.Color("14"))},
		{"White", 0, 0, 100, "15", lipgloss.NewStyle().Foreground(lipgloss.Color("15"))},
	}
	colors = append(colors, standardColors...)
	
	// 216 color cube (16-231) - 6x6x6 RGB cube
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				// Convert RGB cube to RGB values (0-255)
				rVal := 0
				if r > 0 {
					rVal = 55 + r*40
				}
				gVal := 0
				if g > 0 {
					gVal = 55 + g*40
				}
				bVal := 0
				if b > 0 {
					bVal = 55 + b*40
				}
				
				// Convert RGB to HSL
				h, s, l := rgbToHSL(rVal, gVal, bVal)
				
				colorCode := fmt.Sprintf("%d", 16+r*36+g*6+b)
				colors = append(colors, HSLColor{
					Name:         fmt.Sprintf("RGB(%d,%d,%d)", rVal, gVal, bVal),
					Hue:          h,
					Saturation:   s,
					Lightness:    l,
					TerminalCode: colorCode,
					Style:        lipgloss.NewStyle().Foreground(lipgloss.Color(colorCode)),
				})
			}
		}
	}
	
	// 24 grayscale colors (232-255)
	for i := 0; i < 24; i++ {
		gray := 8 + i*10 // 8 to 238
		lightness := int(float64(gray) / 255.0 * 100)
		colorCode := fmt.Sprintf("%d", 232+i)
		colors = append(colors, HSLColor{
			Name:         fmt.Sprintf("Gray %d", gray),
			Hue:          0,
			Saturation:   0,
			Lightness:    lightness,
			TerminalCode: colorCode,
			Style:        lipgloss.NewStyle().Foreground(lipgloss.Color(colorCode)),
		})
	}
	
	return colors
}

// rgbToHSL converts RGB (0-255) to HSL (H: 0-360, S: 0-100, L: 0-100)
func rgbToHSL(r, g, b int) (h, s, l int) {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0
	
	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))
	
	lf := (max + min) / 2.0
	
	if max == min {
		// Achromatic
		h, s = 0, 0
	} else {
		d := max - min
		if lf > 0.5 {
			s = int(d / (2.0 - max - min) * 100)
		} else {
			s = int(d / (max + min) * 100)
		}
		
		switch max {
		case rf:
			hf := (gf - bf) / d
			if gf < bf {
				hf += 6
			}
			h = int(hf / 6.0 * 360)
		case gf:
			h = int(((bf-rf)/d + 2) / 6.0 * 360)
		case bf:
			h = int(((rf-gf)/d + 4) / 6.0 * 360)
		}
	}
	
	l = int(lf * 100)
	return
}

// Show displays the HSL color picker
func (hcp *HSLColorPicker) Show() {
	hcp.visible = true
	hcp.selectedIndex = 0
}

// Hide hides the HSL color picker
func (hcp *HSLColorPicker) Hide() {
	hcp.visible = false
}

// IsVisible returns true if the color picker is visible
func (hcp *HSLColorPicker) IsVisible() bool {
	return hcp.visible
}

// SetDimensions sets the terminal dimensions
func (hcp *HSLColorPicker) SetDimensions(width, height int) {
	hcp.width = width
	hcp.height = height
}

// AddRecentColor adds a color to the recent colors list
func (hcp *HSLColorPicker) AddRecentColor(color HSLColor) {
	// Remove if already exists
	for i, c := range hcp.recentColors {
		if c.TerminalCode == color.TerminalCode {
			hcp.recentColors = append(hcp.recentColors[:i], hcp.recentColors[i+1:]...)
			break
		}
	}
	
	// Add to front
	hcp.recentColors = append([]HSLColor{color}, hcp.recentColors...)
	
	// Keep only last 20
	if len(hcp.recentColors) > 20 {
		hcp.recentColors = hcp.recentColors[:20]
	}
}

// Update handles input and state updates
func (hcp *HSLColorPicker) Update(msg tea.Msg) tea.Cmd {
	if !hcp.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Switch between modes
			switch hcp.mode {
			case "wheel":
				hcp.mode = "recent"
			case "recent":
				hcp.mode = "wheel"
			}
			hcp.selectedIndex = 0

		case "up", "k":
			if hcp.mode == "wheel" {
				hcp.selectedIndex -= 16 // 16 colors per row
			} else {
				hcp.selectedIndex -= 10 // 10 recent colors per row
			}
			if hcp.selectedIndex < 0 {
				if hcp.mode == "wheel" {
					hcp.selectedIndex = len(hcp.colors) - 1
				} else {
					hcp.selectedIndex = len(hcp.recentColors) - 1
				}
			}

		case "down", "j":
			if hcp.mode == "wheel" {
				hcp.selectedIndex += 16
				if hcp.selectedIndex >= len(hcp.colors) {
					hcp.selectedIndex = hcp.selectedIndex % 16
				}
			} else {
				hcp.selectedIndex += 10
				if hcp.selectedIndex >= len(hcp.recentColors) {
					hcp.selectedIndex = hcp.selectedIndex % 10
				}
			}

		case "left", "h":
			hcp.selectedIndex--
			if hcp.selectedIndex < 0 {
				if hcp.mode == "wheel" {
					hcp.selectedIndex = len(hcp.colors) - 1
				} else {
					hcp.selectedIndex = len(hcp.recentColors) - 1
				}
			}

		case "right", "l":
			hcp.selectedIndex++
			maxIndex := len(hcp.colors)
			if hcp.mode == "recent" {
				maxIndex = len(hcp.recentColors)
			}
			if hcp.selectedIndex >= maxIndex {
				hcp.selectedIndex = 0
			}

		case "enter", " ":
			// Select the current color
			var selectedColor HSLColor
			if hcp.mode == "wheel" && hcp.selectedIndex >= 0 && hcp.selectedIndex < len(hcp.colors) {
				selectedColor = hcp.colors[hcp.selectedIndex]
			} else if hcp.mode == "recent" && hcp.selectedIndex >= 0 && hcp.selectedIndex < len(hcp.recentColors) {
				selectedColor = hcp.recentColors[hcp.selectedIndex]
			} else {
				return nil
			}
			
			// Add to recent colors
			hcp.AddRecentColor(selectedColor)
			hcp.Hide()
			
			return func() tea.Msg {
				return ColorSelectedMsg{
					Color: selectedColor.TerminalCode,
					Name:  selectedColor.Name,
				}
			}

		case "escape", "q":
			hcp.Hide()
		}
	}

	return nil
}

// View renders the HSL color picker modal
func (hcp *HSLColorPicker) View() string {
	if !hcp.visible {
		return ""
	}

	// Create styles
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Padding(1, 2).
		Width(66)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	selectedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}).
		Padding(0)

	normalStyle := lipgloss.NewStyle().
		Padding(0)

	// Build title with mode indicator
	title := titleStyle.Render(fmt.Sprintf("Color Picker (%s)", hcp.mode))

	var content string
	
	if hcp.mode == "wheel" {
		// Build color wheel grid (16 colors per row)
		var rows []string
		colsPerRow := 16
		for i := 0; i < len(hcp.colors); i += colsPerRow {
			var colorSwatches []string
			for j := 0; j < colsPerRow && i+j < len(hcp.colors); j++ {
				index := i + j
				color := hcp.colors[index]
				swatch := color.Style.Render("██")
				
				if index == hcp.selectedIndex {
					swatch = selectedStyle.Render(swatch)
				} else {
					swatch = normalStyle.Render(swatch)
				}
				
				colorSwatches = append(colorSwatches, swatch)
			}
			if len(colorSwatches) > 0 {
				rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, colorSwatches...))
			}
		}
		content = lipgloss.JoinVertical(lipgloss.Left, rows...)
	} else {
		// Recent colors mode
		if len(hcp.recentColors) == 0 {
			content = "No recent colors"
		} else {
			var rows []string
			colsPerRow := 10
			for i := 0; i < len(hcp.recentColors); i += colsPerRow {
				var colorSwatches []string
				for j := 0; j < colsPerRow && i+j < len(hcp.recentColors); j++ {
					index := i + j
					color := hcp.recentColors[index]
					swatch := color.Style.Render("███")
					
					if index == hcp.selectedIndex {
						swatch = selectedStyle.Render(swatch)
					} else {
						swatch = normalStyle.Render(swatch)
					}
					
					colorSwatches = append(colorSwatches, swatch)
				}
				if len(colorSwatches) > 0 {
					rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, colorSwatches...))
				}
			}
			content = lipgloss.JoinVertical(lipgloss.Left, rows...)
		}
	}

	// Show selected color info
	selectedColorInfo := ""
	if hcp.mode == "wheel" && hcp.selectedIndex >= 0 && hcp.selectedIndex < len(hcp.colors) {
		color := hcp.colors[hcp.selectedIndex]
		selectedColorInfo = fmt.Sprintf("Selected: %s (H:%d S:%d L:%d) [%s]", 
			color.Name, color.Hue, color.Saturation, color.Lightness, color.TerminalCode)
	} else if hcp.mode == "recent" && hcp.selectedIndex >= 0 && hcp.selectedIndex < len(hcp.recentColors) {
		color := hcp.recentColors[hcp.selectedIndex]
		selectedColorInfo = fmt.Sprintf("Selected: %s [%s]", color.Name, color.TerminalCode)
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Align(lipgloss.Center).
		Render("Navigate: hjkl/arrows · Select: Enter/Space · Mode: Tab · Cancel: Esc")

	// Combine all elements
	allContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		content,
		"",
		selectedColorInfo,
		"",
		instructions,
	)

	return modalStyle.Render(allContent)
}