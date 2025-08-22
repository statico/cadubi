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

// HexPosition represents a position in hexagonal coordinates
type HexPosition struct {
	Q int // column
	R int // row
	X int // screen x
	Y int // screen y
}

// generateHSLColorWheel creates a rectangular color grid arranged by hue for visual appeal
func generateHSLColorWheel() []HSLColor {
	var colors []HSLColor
	
	// Add special default colors first
	colors = append(colors, HSLColor{
		Name: "Default FG", Hue: 0, Saturation: 0, Lightness: 50,
		TerminalCode: "default", Style: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}),
	})
	colors = append(colors, HSLColor{
		Name: "Default BG", Hue: 0, Saturation: 0, Lightness: 50,
		TerminalCode: "default_bg", Style: lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "0"}),
	})
	
	// Create a visually pleasing rectangular arrangement
	// Target: 12 rows x 18 columns = 216 colors
	
	// Generate grid: X-axis = hue, Y-axis = brightness (light on top, dark on bottom)
	for lightnessStep := 0; lightnessStep < 12; lightnessStep++ {
		// Start with light colors (high lightness) and go down to dark colors
		baseLightness := 90 - (lightnessStep * 7) // 90, 83, 76, ..., 20
		
		for hueStep := 0; hueStep < 18; hueStep++ {
			baseHue := hueStep * 20 // 0, 20, 40, 60, ..., 340 degrees
			
			// Vary saturation based on lightness for better color distribution
			var saturation int
			if baseLightness > 70 {
				// Light colors: lower saturation (pastels)
				saturation = 30 + (hueStep % 4) * 15 // 30-75%
			} else if baseLightness > 40 {
				// Medium colors: high saturation (vibrant)
				saturation = 70 + (hueStep % 4) * 10 // 70-100%
			} else {
				// Dark colors: medium saturation
				saturation = 50 + (hueStep % 4) * 15 // 50-95%
			}
			
			// Find closest terminal color
			terminalCode := findClosestTerminalColor(baseHue, saturation, baseLightness)
			
			colors = append(colors, HSLColor{
				Name:         fmt.Sprintf("H%d-S%d-L%d", baseHue, saturation, baseLightness),
				Hue:          baseHue,
				Saturation:   saturation,
				Lightness:    baseLightness,
				TerminalCode: terminalCode,
				Style:        lipgloss.NewStyle().Foreground(lipgloss.Color(terminalCode)),
			})
		}
	}
	
	return colors
}

// findClosestTerminalColor finds the closest 256-color terminal code for given HSL
func findClosestTerminalColor(h, s, l int) string {
	// For grayscale colors, use the 24-step grayscale ramp (232-255)
	if s < 15 { 
		grayLevel := int(float64(l) / 100.0 * 23) // 0-23 range
		if grayLevel <= 0 { return "0" }  // black
		if grayLevel >= 23 { return "15" } // white
		return fmt.Sprintf("%d", 232 + grayLevel)
	}
	
	// Convert HSL to RGB for better color matching
	r, g, b := hslToRGB(h, s, l)
	
	// For the 216-color cube (16-231), find closest match
	// The cube is 6x6x6 with RGB values at positions: 0, 95, 135, 175, 215, 255
	cubeValues := []int{0, 95, 135, 175, 215, 255}
	
	bestR, bestG, bestB := 0, 0, 0
	minDistance := float64(999999)
	
	for ri, rv := range cubeValues {
		for gi, gv := range cubeValues {
			for bi, bv := range cubeValues {
				// Calculate distance in RGB space
				dr := float64(r - rv)
				dg := float64(g - gv) 
				db := float64(b - bv)
				distance := dr*dr + dg*dg + db*db
				
				if distance < minDistance {
					minDistance = distance
					bestR, bestG, bestB = ri, gi, bi
				}
			}
		}
	}
	
	// Convert back to color cube index (16-231)
	colorCode := 16 + bestR*36 + bestG*6 + bestB
	return fmt.Sprintf("%d", colorCode)
}

// hslToRGB converts HSL to RGB (0-255 range)
func hslToRGB(h, s, l int) (int, int, int) {
	hf := float64(h) / 360.0
	sf := float64(s) / 100.0
	lf := float64(l) / 100.0
	
	var r, g, b float64
	
	if s == 0 {
		r, g, b = lf, lf, lf // achromatic
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 { t += 1 }
			if t > 1 { t -= 1 }
			if t < 1.0/6.0 { return p + (q-p)*6*t }
			if t < 1.0/2.0 { return q }
			if t < 2.0/3.0 { return p + (q-p)*(2.0/3.0-t)*6 }
			return p
		}
		
		var q float64
		if lf < 0.5 {
			q = lf * (1 + sf)
		} else {
			q = lf + sf - lf*sf
		}
		p := 2*lf - q
		
		r = hue2rgb(p, q, hf+1.0/3.0)
		g = hue2rgb(p, q, hf)
		b = hue2rgb(p, q, hf-1.0/3.0)
	}
	
	return int(r * 255), int(g * 255), int(b * 255)
}

// calculateHexPositions calculates screen positions for hexagonal layout
func calculateHexPositions(colors []HSLColor) map[int]HexPosition {
	positions := make(map[int]HexPosition)
	
	// Center position (index 0)
	positions[0] = HexPosition{Q: 0, R: 0, X: 0, Y: 0}
	
	colorIndex := 1
	maxRings := 7 // Match the number of rings in generation
	
	// Calculate positions for each ring
	for ring := 1; ring <= maxRings && colorIndex < len(colors); ring++ {
		colorsInRing := ring * 6 // Match the generation algorithm
		
		for i := 0; i < colorsInRing && colorIndex < len(colors); i++ {
			// Calculate hexagonal coordinates with better spacing
			angle := float64(i) / float64(colorsInRing) * 2 * math.Pi
			
			// Use smaller spacing to fit more colors on screen
			ringRadius := float64(ring) * 1.5 // Reduced from 2 to 1.5
			
			// Convert to screen coordinates
			x := int(ringRadius * math.Cos(angle))
			y := int(ringRadius * math.Sin(angle) * 0.6) // Compress vertically for better fit
			
			positions[colorIndex] = HexPosition{
				Q: int(ringRadius * math.Cos(angle)),
				R: int(ringRadius * math.Sin(angle)),
				X: x,
				Y: y,
			}
			colorIndex++
		}
	}
	
	return positions
}

// findNearestColor finds the nearest color in the specified direction for hexagonal navigation
func (hcp *HSLColorPicker) findNearestColor(direction string) int {
	if hcp.selectedIndex < 0 || hcp.selectedIndex >= len(hcp.colors) {
		return 0
	}
	
	positions := calculateHexPositions(hcp.colors)
	currentPos, exists := positions[hcp.selectedIndex]
	if !exists {
		// If selected color is in grayscale section, use simple navigation
		switch direction {
		case "up":
			if hcp.selectedIndex > 0 { return hcp.selectedIndex - 1 }
		case "down": 
			if hcp.selectedIndex < len(hcp.colors)-1 { return hcp.selectedIndex + 1 }
		case "left":
			if hcp.selectedIndex > 0 { return hcp.selectedIndex - 1 }
		case "right":
			if hcp.selectedIndex < len(hcp.colors)-1 { return hcp.selectedIndex + 1 }
		}
		return hcp.selectedIndex
	}
	
	bestIndex := hcp.selectedIndex
	bestDistance := float64(1000) // large number
	
	// Find the nearest color in the specified direction
	for index, pos := range positions {
		if index == hcp.selectedIndex { continue }
		
		dx := float64(pos.X - currentPos.X)
		dy := float64(pos.Y - currentPos.Y)
		distance := math.Sqrt(dx*dx + dy*dy)
		
		// Check if this color is in the right direction
		validDirection := false
		switch direction {
		case "up":
			validDirection = dy < -0.5 // must be significantly above
		case "down":
			validDirection = dy > 0.5  // must be significantly below
		case "left":
			validDirection = dx < -0.5 // must be significantly to the left
		case "right":
			validDirection = dx > 0.5  // must be significantly to the right
		}
		
		if validDirection && distance < bestDistance {
			bestDistance = distance
			bestIndex = index
		}
	}
	
	return bestIndex
}

// renderRectangularGrid creates a rectangular color grid layout
func (hcp *HSLColorPicker) renderRectangularGrid() string {
	selectedStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"})
	normalStyle := lipgloss.NewStyle()
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"})
	
	var rows []string
	
	// First row: Default colors with labels
	var defaultRow []string
	for i := 0; i < 2 && i < len(hcp.colors); i++ {
		color := hcp.colors[i]
		var swatch string
		
		// Special rendering for default colors
		if color.TerminalCode == "default" {
			swatch = "FG"
		} else if color.TerminalCode == "default_bg" {
			swatch = "BG"
		} else {
			swatch = "██"
		}
		
		rendered := color.Style.Render(swatch)
		if i == hcp.selectedIndex {
			rendered = selectedStyle.Render(rendered)
		} else {
			rendered = normalStyle.Render(rendered)
		}
		defaultRow = append(defaultRow, rendered)
	}
	
	// Add spacing and label for default colors
	if len(defaultRow) > 0 {
		defaultRowStr := lipgloss.JoinHorizontal(lipgloss.Top, defaultRow...)
		labeledDefault := lipgloss.JoinHorizontal(lipgloss.Top, 
			labelStyle.Render("Default: "), defaultRowStr)
		rows = append(rows, labeledDefault)
		rows = append(rows, "") // spacer
	}
	
	// Main color grid: 12 rows x 18 columns
	colsPerRow := 18
	colorIndex := 2 // Start after default colors
	
	for row := 0; row < 12 && colorIndex < len(hcp.colors); row++ {
		var rowColors []string
		for col := 0; col < colsPerRow && colorIndex < len(hcp.colors); col++ {
			color := hcp.colors[colorIndex]
			swatch := color.Style.Render("██")
			
			if colorIndex == hcp.selectedIndex {
				swatch = selectedStyle.Render(swatch)
			} else {
				swatch = normalStyle.Render(swatch)
			}
			
			rowColors = append(rowColors, swatch)
			colorIndex++
		}
		
		if len(rowColors) > 0 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, rowColors...))
		}
	}
	
	// Add grayscale colors at the bottom if we have any remaining
	if colorIndex < len(hcp.colors) {
		rows = append(rows, "") // spacer
		
		// Add grayscale colors - use 24 colors from 232-255 plus basic black/white
		grayscaleColors := []HSLColor{
			{"Black", 0, 0, 0, "0", lipgloss.NewStyle().Foreground(lipgloss.Color("0"))},
		}
		
		// Add 24-step grayscale ramp (232-255)
		for i := 232; i <= 255; i++ {
			grayLevel := int(float64(i-232) / 23.0 * 100)
			grayscaleColors = append(grayscaleColors, HSLColor{
				Name: fmt.Sprintf("Gray %d", grayLevel),
				Hue: 0, Saturation: 0, Lightness: grayLevel,
				TerminalCode: fmt.Sprintf("%d", i),
				Style: lipgloss.NewStyle().Foreground(lipgloss.Color(fmt.Sprintf("%d", i))),
			})
		}
		
		grayscaleColors = append(grayscaleColors, HSLColor{
			Name: "White", Hue: 0, Saturation: 0, Lightness: 100,
			TerminalCode: "15", Style: lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		})
		
		var grayscaleRows []string
		colsPerGrayRow := 13 // 26 colors in 2 rows
		for grayRowIdx := 0; grayRowIdx < 2; grayRowIdx++ {
			var grayRow []string
			for col := 0; col < colsPerGrayRow; col++ {
				grayIndex := grayRowIdx*colsPerGrayRow + col
				if grayIndex < len(grayscaleColors) {
					color := grayscaleColors[grayIndex]
					swatch := color.Style.Render("██")
					
					// Check if this grayscale color is selected (need to map to main color index)
					isSelected := false
					if hcp.selectedIndex >= colorIndex {
						relativeIndex := hcp.selectedIndex - colorIndex
						isSelected = (relativeIndex == grayIndex)
					}
					
					if isSelected {
						swatch = selectedStyle.Render(swatch)
					} else {
						swatch = normalStyle.Render(swatch)
					}
					
					grayRow = append(grayRow, swatch)
				}
			}
			if len(grayRow) > 0 {
				grayscaleRows = append(grayscaleRows, lipgloss.JoinHorizontal(lipgloss.Top, grayRow...))
			}
		}
		
		if len(grayscaleRows) > 0 {
			grayLabel := labelStyle.Render("Grayscale:")
			rows = append(rows, grayLabel)
			rows = append(rows, grayscaleRows...)
		}
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
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
				// Grid navigation - handle special cases for default colors and main grid
				if hcp.selectedIndex < 2 {
					// In default colors section - don't move up
					return nil
				} else if hcp.selectedIndex < 2+18 {
					// First row of main grid - move to default colors
					hcp.selectedIndex = hcp.selectedIndex - 2 // map to default colors (0 or 1)
					if hcp.selectedIndex < 0 { hcp.selectedIndex = 0 }
				} else {
					// In main grid - move up one row (18 colors)
					hcp.selectedIndex -= 18
					if hcp.selectedIndex < 2 {
						hcp.selectedIndex = 2 + (hcp.selectedIndex % 18) // wrap to first row
					}
				}
			} else {
				hcp.selectedIndex -= 10 // 10 recent colors per row
				if hcp.selectedIndex < 0 {
					hcp.selectedIndex = len(hcp.recentColors) - 1
				}
			}

		case "down", "j":
			if hcp.mode == "wheel" {
				// Grid navigation - handle special cases
				if hcp.selectedIndex < 2 {
					// In default colors - move to first row of main grid
					hcp.selectedIndex = 2 + hcp.selectedIndex
				} else {
					// In main grid - move down one row (18 colors)
					hcp.selectedIndex += 18
					if hcp.selectedIndex >= len(hcp.colors) {
						// Wrap to top of main grid
						col := (hcp.selectedIndex - 2) % 18
						hcp.selectedIndex = 2 + col
					}
				}
			} else {
				hcp.selectedIndex += 10
				if hcp.selectedIndex >= len(hcp.recentColors) {
					hcp.selectedIndex = hcp.selectedIndex % 10
				}
			}

		case "left", "h":
			if hcp.mode == "wheel" {
				hcp.selectedIndex--
				if hcp.selectedIndex < 0 {
					hcp.selectedIndex = len(hcp.colors) - 1
				}
			} else {
				hcp.selectedIndex--
				if hcp.selectedIndex < 0 {
					hcp.selectedIndex = len(hcp.recentColors) - 1
				}
			}

		case "right", "l":
			if hcp.mode == "wheel" {
				hcp.selectedIndex++
				if hcp.selectedIndex >= len(hcp.colors) {
					hcp.selectedIndex = 0
				}
			} else {
				hcp.selectedIndex++
				if hcp.selectedIndex >= len(hcp.recentColors) {
					hcp.selectedIndex = 0
				}
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
		// Build rectangular color grid layout
		content = hcp.renderRectangularGrid()
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