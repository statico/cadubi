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

	// Create Bubble Tea program with mouse support
	p := tea.NewProgram(
		app, 
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse support with motion
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running CADUBI v2: %v", err)
		os.Exit(1)
	}
}