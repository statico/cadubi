package main

import (
	"fmt"
	"os"

	"github.com/cadubi/cadubi/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Create the main model
	mainModel := ui.NewMainModel()

	// Create Bubble Tea program with mouse support
	p := tea.NewProgram(
		mainModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse support with motion
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running CADUBI v2: %v\n", err)
		os.Exit(1)
	}
}
