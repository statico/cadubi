package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmationDialog represents a confirmation dialog with custom message and actions
type ConfirmationDialog struct {
	visible     bool
	title       string
	message     string
	yesText     string
	noText      string
	cancelText  string
	selected    int // 0=Yes, 1=No, 2=Cancel
	showCancel  bool
	width       int
	height      int
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
		selected:   0, // Default to "Yes"
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
		selected:   0, // Default to "Save & Exit"
		showCancel: true,
		width:      58,
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
	d.width = width
	d.height = height
}

// Update handles input and returns commands
func (d *ConfirmationDialog) Update(msg tea.Msg) tea.Cmd {
	if !d.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			d.selected--
			maxOption := 1
			if d.showCancel {
				maxOption = 2
			}
			if d.selected < 0 {
				d.selected = maxOption
			}

		case "right", "l":
			d.selected++
			maxOption := 1
			if d.showCancel {
				maxOption = 2
			}
			if d.selected > maxOption {
				d.selected = 0
			}

		case "enter":
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

		case "escape", "q":
			// Cancel
			d.Hide()
			return func() tea.Msg {
				return ConfirmationResultMsg{Action: "cancel"}
			}
		}

	case tea.MouseMsg:
		// Handle mouse clicks on buttons
		if msg.Type == tea.MouseLeft {
			// Calculate button positions and check if clicked
			// This is a simplified version - could be more precise
			buttonY := d.height/2 + 3 // Approximate button row
			if msg.Y == buttonY {
				buttonWidth := 15
				spacing := 2
				
				// Calculate button positions
				totalWidth := buttonWidth * 2
				if d.showCancel {
					totalWidth += buttonWidth + spacing
				}
				totalWidth += spacing
				
				startX := (d.width - totalWidth) / 2
				
				// Check which button was clicked
				if msg.X >= startX && msg.X < startX+buttonWidth {
					// Yes button
					d.selected = 0
					d.Hide()
					return func() tea.Msg {
						return ConfirmationResultMsg{Action: "yes"}
					}
				} else if msg.X >= startX+buttonWidth+spacing && msg.X < startX+buttonWidth*2+spacing {
					// No button
					d.selected = 1
					d.Hide()
					return func() tea.Msg {
						return ConfirmationResultMsg{Action: "no"}
					}
				} else if d.showCancel && msg.X >= startX+buttonWidth*2+spacing*2 && msg.X < startX+buttonWidth*3+spacing*2 {
					// Cancel button
					d.selected = 2
					d.Hide()
					return func() tea.Msg {
						return ConfirmationResultMsg{Action: "cancel"}
					}
				}
			}
		}
	}

	return nil
}

// View renders the confirmation dialog
func (d *ConfirmationDialog) View() string {
	if !d.visible {
		return ""
	}

	// Create styles using adaptive colors
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center)

	messageStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Margin(1, 0)

	buttonStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Margin(0, 1)

	selectedButtonStyle := buttonStyle.Copy().
		Bold(true).
		Underline(true)

	// Build dialog content
	title := titleStyle.Render(d.title)
	message := messageStyle.Render(d.message)

	// Create buttons
	var buttons []string
	
	// Yes button
	yesStyle := buttonStyle
	if d.selected == 0 {
		yesStyle = selectedButtonStyle
	}
	buttons = append(buttons, yesStyle.Render(d.yesText))
	
	// No button
	noStyle := buttonStyle
	if d.selected == 1 {
		noStyle = selectedButtonStyle
	}
	buttons = append(buttons, noStyle.Render(d.noText))
	
	// Cancel button (if enabled)
	if d.showCancel {
		cancelStyle := buttonStyle
		if d.selected == 2 {
			cancelStyle = selectedButtonStyle
		}
		buttons = append(buttons, cancelStyle.Render(d.cancelText))
	}

	buttonRow := lipgloss.JoinHorizontal(lipgloss.Top, buttons...)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"}).
		Align(lipgloss.Center).
		Render("Use ← → or h l to select, Enter to confirm, Esc to cancel")

	// Combine all elements
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		message,
		buttonRow,
		"",
		instructions,
	)

	return modalStyle.Render(content)
}