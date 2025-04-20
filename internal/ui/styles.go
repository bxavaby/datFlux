package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

var (
	StormBlue    = lipgloss.Color("#24283b") // Primary background
	DarkBlue     = lipgloss.Color("#565f89") // Secondary background
	Lavender     = lipgloss.Color("#bb9af7") // Purple accent
	SkyBlue      = lipgloss.Color("#7dcfff") // Blue accent
	PaleGreen    = lipgloss.Color("#9ece6a") // Success/positive
	PeachRed     = lipgloss.Color("#f7768e") // Warning/accent
	OrangeYellow = lipgloss.Color("#e0af68") // Highlight

	StormBlueStr    = "#24283b"
	DarkBlueStr     = "#565f89"
	LavenderStr     = "#bb9af7"
	SkyBlueStr      = "#7dcfff"
	PaleGreenStr    = "#9ece6a"
	PeachRedStr     = "#f7768e"
	OrangeYellowStr = "#e0af68"
)

var PasswordAnimationColors = []lipgloss.Color{
	Lavender,
	SkyBlue,
	PaleGreen,
	OrangeYellow,
	PeachRed,
}

var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(Lavender). // purple 4 title
			Bold(true).
			Padding(0, 0, 1, 0).
			Align(lipgloss.Center)

	SectionTitleStyle = lipgloss.NewStyle().
				Foreground(SkyBlue). // blue 4 section titles
				Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(PaleGreen) // green 4 values

	WarningStyle = lipgloss.NewStyle().
			Foreground(OrangeYellow) // yellow-orange 4 warnings

	DangerStyle = lipgloss.NewStyle().
			Foreground(PeachRed) // red 4 danger/errors

	PasswordStyle = lipgloss.NewStyle().
			Foreground(Lavender). // purple 4 password
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(DarkBlue). // muted blue 4 help
			Italic(true)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SkyBlue). // blue borders
			Padding(1, 0)
)

var (
	CPUProgress = progress.New(
		progress.WithGradient(PaleGreenStr, PeachRedStr),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	MemoryProgress = progress.New(
		progress.WithGradient(PaleGreenStr, SkyBlueStr),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
)

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
