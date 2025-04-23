package ui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var currentTheme Theme

var (
	TitleStyle         lipgloss.Style
	SectionTitleStyle  lipgloss.Style
	ValueStyle         lipgloss.Style
	WarningStyle       lipgloss.Style
	DangerStyle        lipgloss.Style
	PasswordStyle      lipgloss.Style
	HelpStyle          lipgloss.Style
	BorderStyle        lipgloss.Style
	StrongPwdStyle     lipgloss.Style
	VeryStrongPwdStyle lipgloss.Style
	LogoStyle          lipgloss.Style

	CPUProgress    progress.Model
	MemoryProgress progress.Model
)

var CustomBorder = lipgloss.Border{
	Top:         "═",
	Bottom:      "═", // ━ or ╌
	Left:        "║", // ┃
	Right:       "║",
	TopLeft:     "┲", // ┭ or ╆
	TopRight:    "┱", // ┮ or ╅
	BottomLeft:  "┺", // ┵ or ╄
	BottomRight: "┹", // ┶ or ╃
}

var PasswordAnimationColors []lipgloss.Color

func InitializeStyles(theme Theme) {
	currentTheme = theme

	PasswordAnimationColors = []lipgloss.Color{
		theme.Accent,
		theme.Highlight,
		theme.StrongColor,
		theme.Warning,
		theme.Danger,
	}

	TitleStyle = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true).
		Padding(0, 0, 1, 0).
		Align(lipgloss.Center)

	SectionTitleStyle = lipgloss.NewStyle().
		Foreground(theme.Highlight).
		Bold(true)

	ValueStyle = lipgloss.NewStyle().
		Foreground(theme.StrongColor)

	WarningStyle = lipgloss.NewStyle().
		Foreground(theme.Warning)

	DangerStyle = lipgloss.NewStyle().
		Foreground(theme.Danger)

	PasswordStyle = lipgloss.NewStyle().
		Foreground(theme.PasswordColor).
		Bold(true)

	HelpStyle = lipgloss.NewStyle().
		Foreground(theme.Secondary).
		Italic(true)

	BorderStyle = lipgloss.NewStyle().
		Border(CustomBorder).
		BorderForeground(theme.Highlight).
		Padding(1, 0)

	StrongPwdStyle = lipgloss.NewStyle().
		Foreground(theme.StrongColor)

	VeryStrongPwdStyle = lipgloss.NewStyle().
		Foreground(theme.VeryStrongColor)

	LogoStyle = lipgloss.NewStyle().
		Foreground(theme.Accent).
		Bold(true)

	CPUProgress = progress.New(
		progress.WithGradient(theme.StrongColorStr, theme.DangerStr),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	MemoryProgress = progress.New(
		progress.WithGradient(theme.StrongColorStr, theme.HighlightStr),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
}

func FormatProgressBar(model progress.Model, value float64, width int) string {
	model.Width = width

	return model.ViewAs(value / 100)
}

func AddPercentage(progressBar string, percentage float64, width int) string {
	percentText := fmt.Sprintf(" %.1f%%", percentage)

	var style lipgloss.Style
	if percentage > 85 {
		style = DangerStyle
	} else if percentage > 60 {
		style = WarningStyle
	} else {
		style = ValueStyle
	}

	return progressBar + style.Render(percentText)
}

func Logo() string {
	logo := `
dF== == == == == == == == == == == == == == == == == == == == == == == == == == ==dF
||  	  8888888888888   ||   USAGE:                                             ||
||        888888          ||     datflux          Launch interactive mode (TUI)   ||
||        888888          ||     datflux now      Generate password immediately   ||
||    .d888888888888      ||     datflux now -p   Generate ultra-secure password  ||
||   d88" 888888       == 🌸 == == == == == == == == == == == == == == == == == ==
||   888  888888          ||   OPTIONS:                                           ||
||   Y88b 888888          ||     --paranoia, -p          Enable paranoia mode     ||
||    "Y88888888          ||     --help, -h, help        Show help                ||
dF== == == == == == ==| v1.0.0 |== == == == == == == == == == == == == == == == ==dF
`
	return LogoStyle.Render(strings.TrimRight(logo, "\n"))
}

func Wiper() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
