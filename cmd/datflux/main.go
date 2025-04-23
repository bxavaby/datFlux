package main

import (
	"fmt"
	"os"
	"time"

	"datflux/internal/entropy"
	"datflux/internal/password"
	"datflux/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// are there args?
	if len(os.Args) > 1 {
		handleSubcommands(os.Args[1:])
		return
	}

	launchTUI()
}

func launchTUI() {
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

func handleSubcommands(args []string) {
	if len(args) == 0 {
		launchTUI()
		return
	}

	switch args[0] {
	case "now":
		generatePasswordNow(args[1:])
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func generatePasswordNow(args []string) {
	// flag parsing
	paranoiaMode := false
	for _, arg := range args {
		if arg == "--paranoia" || arg == "-p" {
			paranoiaMode = true
		} else if arg == "--help" || arg == "-h" {
			printHelp()
			return
		}
	}

	ui.InitializeStyles(ui.GetDefaultTheme())

	// entropy collector with shorter initialization time
	collector := entropy.NewCollector(time.Millisecond*50, 20)
	defer collector.Close()

	// run the noise generator briefly
	noiseGen := entropy.NewNoiseGenerator(collector)

	// run for a short period to gather entropy
	time.Sleep(200 * time.Millisecond)
	noiseGen.Stop()

	passGen := password.NewGenerator(collector)
	passGen.SetParanoiaMode(paranoiaMode, 5) // fewer samples for CLI

	pw := passGen.Generate()

	// output to stdout
	fmt.Println(pw)
}

func printHelp() {
	ui.InitializeStyles(ui.GetDefaultTheme())

	header := ui.Logo() + "\n"

	fmt.Println(header)
}
