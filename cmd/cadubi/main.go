package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cadubi/cadubi/pkg/ui"
)

func main() {
	// Create the application
	app := ui.NewApp()

	// Create Bubble Tea program
	p := tea.NewProgram(app, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running CADUBI v2: %v", err)
		os.Exit(1)
	}
}