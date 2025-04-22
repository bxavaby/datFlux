package main

import (
	"fmt"
	"os"
	"time"

	"datflux/internal/entropy"
	"datflux/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	collector := entropy.NewCollector(time.Millisecond*100, 50)
	noiseGen := entropy.NewNoiseGenerator(collector)
	defer collector.Close()
	defer noiseGen.Stop()

	p := tea.NewProgram(
		ui.NewDashboardModel(collector),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
