package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cadubi/cadubi/pkg/trace"
	"github.com/cadubi/cadubi/pkg/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command line flags
	traceEnabled := flag.Bool("trace", false, "Enable comprehensive trace logging")
	traceFile := flag.String("trace-file", "trace.log", "Trace log file path")
	flag.Parse()

	// Initialize trace logging
	if err := trace.Init(*traceEnabled, *traceFile); err != nil {
		fmt.Printf("Failed to initialize trace logging: %v\n", err)
		os.Exit(1)
	}
	defer trace.Close()

	trace.Infof("main", "CADUBI v2 starting with trace=%v, trace-file=%s", *traceEnabled, *traceFile)

	// Create the main model
	mainModel := ui.NewMainModel()

	// Create Bubble Tea program with mouse support
	p := tea.NewProgram(
		mainModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(), // Enable mouse support with motion
	)

	trace.Infof("main", "Bubble Tea program created with AltScreen and mouse support")

	// Run the program
	trace.Infof("main", "Starting Bubble Tea program execution")
	if _, err := p.Run(); err != nil {
		trace.LogError("main", err, map[string]interface{}{
			"operation": "program_run",
		})
		fmt.Printf("Error running CADUBI v2: %v\n", err)
		os.Exit(1)
	}

	trace.Infof("main", "CADUBI v2 shutting down normally")
}
