package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Toast represents a temporary notification message
type Toast struct {
	message   string
	visible   bool
	duration  time.Duration
	startTime time.Time
}

// ToastMsg is sent when a toast should be shown
type ToastMsg struct {
	Message  string
	Duration time.Duration
}

// ToastTimeoutMsg is sent when a toast should be hidden
type ToastTimeoutMsg struct{}

// NewToast creates a new toast notification
func NewToast() *Toast {
	return &Toast{
		visible:  false,
		duration: 2 * time.Second, // Default 2 seconds
	}
}

// Show displays a toast message for the specified duration
func (t *Toast) Show(message string, duration time.Duration) tea.Cmd {
	t.message = message
	t.visible = true
	t.duration = duration
	t.startTime = time.Now()

	// Return a command that will send a timeout message after the duration
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return ToastTimeoutMsg{}
	})
}

// Hide hides the toast
func (t *Toast) Hide() {
	t.visible = false
	t.message = ""
}

// IsVisible returns whether the toast is currently visible
func (t *Toast) IsVisible() bool {
	return t.visible
}

// Update handles toast-related messages
func (t *Toast) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ToastMsg:
		return t.Show(msg.Message, msg.Duration)
	case ToastTimeoutMsg:
		// Only hide if enough time has passed (prevents race conditions)
		if time.Since(t.startTime) >= t.duration {
			t.Hide()
		}
		return nil
	}
	return nil
}

// View renders the toast notification
func (t *Toast) View() string {
	if !t.visible {
		return ""
	}

	// Create toast style with subtle background and border
	toastStyle := lipgloss.NewStyle().
		Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#333333"}).
		Foreground(lipgloss.AdaptiveColor{Light: "#333333", Dark: "#EEEEEE"}).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#CCCCCC", Dark: "#555555"}).
		Padding(0, 1).
		Margin(0, 1)

	return toastStyle.Render(t.message)
}