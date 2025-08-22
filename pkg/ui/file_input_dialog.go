package ui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileInputDialog struct {
	textInput     textinput.Model
	visible       bool
	mode          string // "load" or "save"
	completions   []string
	currentDir    string
	showingCompletion bool
}

type FileSelectedMsg struct {
	Path string
	Mode string
}

func NewFileInputDialog() *FileInputDialog {
	ti := textinput.New()
	ti.Placeholder = "Enter filename..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	currentDir, _ := os.Getwd()

	return &FileInputDialog{
		textInput:  ti,
		visible:    false,
		currentDir: currentDir,
	}
}

func (d *FileInputDialog) Show(mode string) {
	d.ShowWithValue(mode, "")
}

func (d *FileInputDialog) ShowWithValue(mode, initialValue string) {
	d.mode = mode
	d.visible = true
	d.textInput.Reset()
	
	if initialValue != "" {
		d.textInput.SetValue(initialValue)
		// Position cursor at the beginning for easy editing
		d.textInput.CursorStart()
	}
	
	if mode == "load" {
		d.textInput.Placeholder = "Load file..."
	} else {
		d.textInput.Placeholder = "Save as..."
	}
	d.textInput.Focus()
	d.showingCompletion = false
	d.updateCompletions()
}

func (d *FileInputDialog) Hide() {
	d.visible = false
	d.textInput.Blur()
	d.completions = nil
	d.showingCompletion = false
}

func (d *FileInputDialog) IsVisible() bool {
	return d.visible
}

func (d *FileInputDialog) updateCompletions() {
	input := d.textInput.Value()
	if input == "" {
		d.completions = nil
		return
	}

	dir := d.currentDir
	filename := input

	// If input contains a path separator, split it
	if strings.Contains(input, string(os.PathSeparator)) {
		dir = filepath.Join(d.currentDir, filepath.Dir(input))
		filename = filepath.Base(input)
	}

	// Read directory contents
	entries, err := os.ReadDir(dir)
	if err != nil {
		d.completions = nil
		return
	}

	var matches []string
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files
		if strings.HasPrefix(name, ".") {
			continue
		}

		// For load mode, only show .cadubi files and directories
		if d.mode == "load" && !entry.IsDir() && !strings.HasSuffix(name, ".cadubi") {
			continue
		}

		// Check if name starts with our partial filename
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(filename)) {
			if entry.IsDir() {
				name += string(os.PathSeparator)
			}
			if dir != d.currentDir {
				relDir, _ := filepath.Rel(d.currentDir, dir)
				name = filepath.Join(relDir, name)
			}
			matches = append(matches, name)
		}
	}

	sort.Strings(matches)
	d.completions = matches
}

func (d *FileInputDialog) handleTab() {
	if len(d.completions) == 0 {
		return
	}

	if len(d.completions) == 1 {
		// Single completion - use it
		d.textInput.SetValue(d.completions[0])
		d.textInput.CursorEnd()
		d.updateCompletions()
	} else {
		// Multiple completions - show them
		d.showingCompletion = true
	}
}

func (d *FileInputDialog) Update(msg tea.Msg) tea.Cmd {
	if !d.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			filename := strings.TrimSpace(d.textInput.Value())
			if filename != "" {
				d.Hide()
				// Convert relative path to absolute
				if !filepath.IsAbs(filename) {
					filename = filepath.Join(d.currentDir, filename)
				}
				// Add .cadubi extension if not present
				if !strings.HasSuffix(filename, ".cadubi") {
					filename += ".cadubi"
				}
				return func() tea.Msg {
					return FileSelectedMsg{Path: filename, Mode: d.mode}
				}
			}
		case "esc":
			d.Hide()
			return nil
		case "tab":
			d.updateCompletions()
			d.handleTab()
			return nil
		default:
			d.showingCompletion = false
			var cmd tea.Cmd
			d.textInput, cmd = d.textInput.Update(msg)
			d.updateCompletions()
			return cmd
		}
	}

	var cmd tea.Cmd
	d.textInput, cmd = d.textInput.Update(msg)
	return cmd
}

func (d *FileInputDialog) View() string {
	if !d.visible {
		return ""
	}

	title := "Load File"
	if d.mode == "save" {
		title = "Save File"
	}

	content := []string{
		lipgloss.NewStyle().Bold(true).Render(title),
		"",
		d.textInput.View(),
	}

	// Show completions if available
	if d.showingCompletion && len(d.completions) > 1 {
		content = append(content, "", "Completions:")
		maxShow := 8
		for i, comp := range d.completions {
			if i >= maxShow {
				content = append(content, "...")
				break
			}
			content = append(content, "  "+comp)
		}
	}

	content = append(content, "", "Enter: Select • Tab: Complete • Esc: Cancel")

	dialogContent := strings.Join(content, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(60).
		Render(dialogContent)
}