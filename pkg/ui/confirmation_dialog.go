package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ButtonOpts represents options for a selectable button
type ButtonOpts struct {
	Text           string
	UnderlineIndex int  // Index of character to underline (0-based)
	Selected       bool // Whether this button is selected
	Key            string // Keyboard shortcut key
}

// ConfirmationDialog represents a confirmation dialog with custom message and actions
type ConfirmationDialog struct {
	visible      bool
	title        string
	message      string
	yesText      string
	noText       string
	cancelText   string
	selected     int // 0=Yes, 1=No, 2=Cancel
	showCancel   bool
	width        int
	height       int
	windowWidth  int
	windowHeight int
}

// ConfirmationResult represents the result of the confirmation dialog
type ConfirmationResult struct {
	Action string // "yes", "no", "cancel"
}

// ConfirmationResultMsg is sent when the dialog is closed
type ConfirmationResultMsg ConfirmationResult

// NewConfirmationDialog creates a new confirmation dialog
func NewConfirmationDialog(title, message, yesText, noText string) *ConfirmationDialog {
	return &ConfirmationDialog{
		visible:    false,
		title:      title,
		message:    message,
		yesText:    yesText,
		noText:     noText,
		cancelText: "Cancel",
		selected:   2, // Default to "Cancel" for safety
		showCancel: false,
		width:      50,
		height:     10,
	}
}

// NewSaveConfirmationDialog creates a save confirmation dialog
func NewSaveConfirmationDialog() *ConfirmationDialog {
	return &ConfirmationDialog{
		visible:    false,
		title:      "Unsaved Changes",
		message:    "You have unsaved changes. Do you want to save before exiting?",
		yesText:    "Save & Exit",
		noText:     "Exit Without Saving",
		cancelText: "Cancel",
		selected:   2, // Default to "Cancel" for safety
		showCancel: true,
		width:      50,
		height:     8,
	}
}

// Show displays the confirmation dialog
func (d *ConfirmationDialog) Show() {
	d.visible = true
	d.selected = 0 // Reset to first option
}

// Hide hides the confirmation dialog
func (d *ConfirmationDialog) Hide() {
	d.visible = false
}

// IsVisible returns true if the dialog is visible
func (d *ConfirmationDialog) IsVisible() bool {
	return d.visible
}

// SetDimensions sets the terminal dimensions for proper centering
func (d *ConfirmationDialog) SetDimensions(width, height int) {
	d.windowWidth = width
	d.windowHeight = height
	
	// Calculate dialog size as percentage of window
	d.width = int(float64(width) * 0.6)
	if d.width > 60 {
		d.width = 60
	}
	if d.width < 40 {
		d.width = 40
	}
	
	d.height = 10
}

// Update handles input and returns commands
func (d *ConfirmationDialog) Update(msg tea.Msg) tea.Cmd {
	if !d.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h", "shift+tab":
			d.selected--
			maxOption := 1
			if d.showCancel {
				maxOption = 2
			}
			if d.selected < 0 {
				d.selected = maxOption
			}

		case "right", "l", "tab":
			d.selected++
			maxOption := 1
			if d.showCancel {
				maxOption = 2
			}
			if d.selected > maxOption {
				d.selected = 0
			}

		case "s", "S":
			// Save & Exit shortcut
			d.selected = 0
			d.Hide()
			return func() tea.Msg {
				return ConfirmationResultMsg{Action: "yes"}
			}

		case "e", "E":
			// Exit without saving shortcut
			d.selected = 1
			d.Hide()
			return func() tea.Msg {
				return ConfirmationResultMsg{Action: "no"}
			}

		case "c", "C":
			// Cancel shortcut
			if d.showCancel {
				d.selected = 2
				d.Hide()
				return func() tea.Msg {
					return ConfirmationResultMsg{Action: "cancel"}
				}
			}

		case "enter", " ":
			// Execute the selected action
			var action string
			switch d.selected {
			case 0:
				action = "yes"
			case 1:
				action = "no"
			case 2:
				action = "cancel"
			}
			
			d.Hide()
			return func() tea.Msg {
				return ConfirmationResultMsg{Action: action}
			}

		case "escape":
			// Cancel
			d.Hide()
			return func() tea.Msg {
				return ConfirmationResultMsg{Action: "cancel"}
			}
		}
	}

	return nil
}

// SelectableButton creates a button with underlined shortcut key following CRUSH patterns
func (d *ConfirmationDialog) SelectableButton(opts ButtonOpts) string {
	baseStyle := lipgloss.NewStyle().
		Padding(0, 2)
		
	if opts.Selected {
		baseStyle = baseStyle.
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
			Background(lipgloss.AdaptiveColor{Light: "#3B82F6", Dark: "#3B82F6"})
	} else {
		baseStyle = baseStyle.
			Background(lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#374151"})
	}
	
	text := opts.Text
	if opts.UnderlineIndex >= 0 && opts.UnderlineIndex < len(text) {
		// Split text and underline the shortcut character
		before := text[:opts.UnderlineIndex]
		char := string(text[opts.UnderlineIndex])
		after := text[opts.UnderlineIndex+1:]
		
		underlineStyle := baseStyle.Copy().Underline(true)
		
		return before + underlineStyle.Render(char) + baseStyle.Render(after)
	}
	
	return baseStyle.Render(text)
}

// SelectableButtons creates a horizontal row of buttons with spacing
func (d *ConfirmationDialog) SelectableButtons(buttons []ButtonOpts, spacing string) string {
	var parts []string
	for i, button := range buttons {
		parts = append(parts, d.SelectableButton(button))
		if i < len(buttons)-1 {
			parts = append(parts, spacing)
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

// SelectableButtonsVertical creates a vertical column of buttons
func (d *ConfirmationDialog) SelectableButtonsVertical(buttons []ButtonOpts, spacing int) string {
	var parts []string
	for i, button := range buttons {
		parts = append(parts, d.SelectableButton(button))
		if i < len(buttons)-1 {
			for j := 0; j < spacing; j++ {
				parts = append(parts, "")
			}
		}
	}
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

// View renders the confirmation dialog
func (d *ConfirmationDialog) View() string {
	if !d.visible {
		return ""
	}

	// Theme colors following CRUSH patterns
	borderColor := lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}
	textColor := lipgloss.AdaptiveColor{Light: "#374151", Dark: "#E5E7EB"}
	mutedColor := lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}

	// Base modal style
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(d.width)

	// Title style
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(textColor).
		Align(lipgloss.Center).
		Width(d.width - 4)

	// Message style
	messageStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		Width(d.width - 4)

	// Build dialog content
	title := titleStyle.Render(d.title)
	message := messageStyle.Render(d.message)

	// Create buttons using CRUSH patterns
	var buttons []ButtonOpts
	
	// Save & Exit button (S shortcut)
	buttons = append(buttons, ButtonOpts{
		Text:           d.yesText,
		UnderlineIndex: 0, // Underline 'S' in "Save & Exit"
		Selected:       d.selected == 0,
		Key:            "S",
	})
	
	// Exit Without Saving button (E shortcut)  
	buttons = append(buttons, ButtonOpts{
		Text:           d.noText,
		UnderlineIndex: 0, // Underline 'E' in "Exit Without Saving"
		Selected:       d.selected == 1,
		Key:            "E",
	})
	
	// Cancel button (C shortcut)
	if d.showCancel {
		buttons = append(buttons, ButtonOpts{
			Text:           d.cancelText,
			UnderlineIndex: 0, // Underline 'C' in "Cancel"
			Selected:       d.selected == 2,
			Key:            "C",
		})
	}

	// Try horizontal layout first
	buttonRow := d.SelectableButtons(buttons, "  ")
	
	// If horizontal layout is too wide, use vertical layout
	if lipgloss.Width(buttonRow) > d.width-4 {
		buttonRow = d.SelectableButtonsVertical(buttons, 1)
	}

	// Center the buttons
	buttonContainer := lipgloss.NewStyle().
		Width(d.width - 4).
		Align(lipgloss.Center).
		Render(buttonRow)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Width(d.width - 4).
		Render("Use Tab/← → to navigate, Enter to confirm, Esc to cancel")

	// Combine all elements with proper spacing
	content := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		"",
		message,
		"",
		buttonContainer,
		"",
		instructions,
	)

	return modalStyle.Render(content)
}