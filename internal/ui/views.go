package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	"datflux/internal/monitor"
)

func renderCPUView(cpuUsage float64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(SectionTitleStyle.Render("CPU USAGE"))
	builder.WriteString("\n\n")

	progressView := progressBar.ViewAs(cpuUsage / 100)

	percentage := fmt.Sprintf(" %.1f%%", cpuUsage)
	if cpuUsage > 85 {
		percentage = DangerStyle.Render(percentage)
	} else if cpuUsage > 60 {
		percentage = WarningStyle.Render(percentage)
	} else {
		percentage = ValueStyle.Render(percentage)
	}

	builder.WriteString(progressView)
	builder.WriteString(percentage)

	return BorderStyle.Width(width).Render(builder.String())
}

func renderMemoryView(memUsage float64, memTotal uint64, memUsed uint64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(SectionTitleStyle.Render("MEMORY USAGE"))
	builder.WriteString("\n\n")

	progressView := progressBar.ViewAs(memUsage / 100)

	usedGB := float64(memUsed) / 1024 / 1024 / 1024
	totalGB := float64(memTotal) / 1024 / 1024 / 1024

	usageText := fmt.Sprintf(" %.1f%%", memUsage)
	if memUsage > 85 {
		usageText = DangerStyle.Render(usageText)
	} else if memUsage > 60 {
		usageText = WarningStyle.Render(usageText)
	} else {
		usageText = ValueStyle.Render(usageText)
	}

	detailsText := fmt.Sprintf(" (%.1f/%.1f GB)", usedGB, totalGB)
	detailsText = ValueStyle.Render(detailsText)

	builder.WriteString(progressView)
	builder.WriteString(usageText)
	builder.WriteString(detailsText)

	return BorderStyle.Width(width).Render(builder.String())
}

func renderNetworkView(rxSpeed float64, txSpeed float64, iface string, width int) string {
	var builder strings.Builder

	title := fmt.Sprintf("NETWORK (%s)", iface)
	builder.WriteString(SectionTitleStyle.Render(title))
	builder.WriteString("\n\n")

	downloadText := fmt.Sprintf("↓ %s", monitor.FormatSpeed(rxSpeed))
	uploadText := fmt.Sprintf("↑ %s", monitor.FormatSpeed(txSpeed))

	downloadFormatted := ValueStyle.Render(downloadText)
	uploadFormatted := ValueStyle.Render(uploadText)

	builder.WriteString(lipgloss.NewStyle().
		Width(width / 2).
		Align(lipgloss.Center).
		Render(downloadFormatted))

	builder.WriteString(lipgloss.NewStyle().
		Width(width / 2).
		Align(lipgloss.Center).
		Render(uploadFormatted))

	return BorderStyle.Width(width).Render(builder.String())
}

func renderPasswordView(animation *PasswordAnimation, quality float64, width int) string {
	var builder strings.Builder

	qualityIndicator := ""
	if !animation.IsAnimating {
		qualityText := fmt.Sprintf(" [Entropy: %.0f%%]", quality*100)
		qualityIndicator = WarningStyle.Render(qualityText)
	}

	title := "GENERATED PASSWORD" + qualityIndicator
	builder.WriteString(SectionTitleStyle.Render(title))
	builder.WriteString("\n\n")

	passwordText := animation.CurrentPassword()
	passLen := len(passwordText)

	padding := max((width-passLen-4)/2, 0)

	builder.WriteString(strings.Repeat(" ", padding))

	if animation.IsAnimating {
		builder.WriteString(animation.StyledPassword())
	} else {
		builder.WriteString(PasswordStyle.Render(passwordText))
	}

	return BorderStyle.Width(width).Render(builder.String())
}
