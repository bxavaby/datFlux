package main

import (
	"fmt"
	"os"

	"datflux/internal/entropy"
	"datflux/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	noiseGen := entropy.NewNoiseGenerator()
	defer noiseGen.Stop()

	p := tea.NewProgram(
		ui.NewDashboardModel(noiseGen.Collector()),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
