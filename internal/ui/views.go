package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	"datflux/internal/monitor"
	"datflux/internal/password"
)

func renderCPUView(cpuUsage float64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(SectionTitleStyle.Render("CPU USAGE"))
	builder.WriteString("\n\n")

	progressView := FormatProgressBar(progressBar, cpuUsage, width-10)
	progressWithPercentage := AddPercentage(progressView, cpuUsage, width)

	builder.WriteString(progressWithPercentage)

	return BorderStyle.Width(width).Render(builder.String())
}

func renderMemoryView(memUsage float64, memTotal uint64, memUsed uint64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(SectionTitleStyle.Render("MEMORY USAGE"))
	builder.WriteString("\n\n")

	progressView := FormatProgressBar(progressBar, memUsage, width-10)
	progressWithPercentage := AddPercentage(progressView, memUsage, width)

	usedGB := float64(memUsed) / 1024 / 1024 / 1024
	totalGB := float64(memTotal) / 1024 / 1024 / 1024
	detailsText := fmt.Sprintf(" (%.1f/%.1f GB)", usedGB, totalGB)

	builder.WriteString(progressWithPercentage)
	builder.WriteString(ValueStyle.Render(detailsText))

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

func renderPasswordView(animation *PasswordAnimation, quality float64, width int, passwordGen *password.Generator, attackModel password.AttackModelType) string {
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

		if passwordText != "Press 'r' to generate" && len(passwordText) > 0 {
			strength := passwordGen.AnalyzeStrength(passwordText)

			builder.WriteString("\n\n")

			// strength meter
			builder.WriteString(renderStrengthMeter(strength.Score, width-10))

			// length and entropy stats
			statsText := fmt.Sprintf("Length: %d | Entropy: %.1f bits", len(passwordText), strength.EntropyBits)
			builder.WriteString("\n" + VeryStrongPwdStyle.Render(statsText))

			// special note for paranoia mode, due to the extreme security level
			paranoiaMode, _ := passwordGen.GetParanoiaMode()
			if paranoiaMode {
				builder.WriteString("\n" + ValueStyle.Render("Password exceeds quantum-resistant security threshold"))
				builder.WriteString("\n" + ValueStyle.Render("Time to crack: Beyond any feasible computation"))
			} else {
				// crack time w current model for standard mode
				crackTimeDesc := passwordGen.GetCrackTimeForModel(passwordText, attackModel)
				crackTimeText := fmt.Sprintf("Time to crack: %s", crackTimeDesc)
				builder.WriteString("\n" + renderStrengthText(crackTimeText, strength.Score))

				// attack model info
				models := password.GetAttackModels()
				currentModel := models[attackModel]
				modelText := fmt.Sprintf("Attack model: %s", currentModel.Name)
				builder.WriteString("\n" + ValueStyle.Render(modelText))
			}

			// feedback, if any
			if strength.Feedback != "" && width > 60 {
				builder.WriteString("\n" + WarningStyle.Render(strength.Feedback))
			}
		}
	}

	return BorderStyle.Width(width).Render(builder.String())
}

func renderStrengthMeter(score int, width int) string {
	colors := []lipgloss.Style{
		DangerStyle,        // 0 - Very Weak
		DangerStyle,        // 1 - Weak
		WarningStyle,       // 2 - Medium
		StrongPwdStyle,     // 3 - Strong
		VeryStrongPwdStyle, // 4 - Very Strong
	}

	labels := []string{
		"Very Weak", "Weak", "Reasonable", "Strong", "Very Strong",
	}

	barWidth := width - len(labels[score]) - 2
	filledWidth := int(float64(barWidth) * float64(score+1) / 5.0)
	emptyWidth := barWidth - filledWidth

	bar := colors[score].Render(strings.Repeat("█", filledWidth))
	bar += strings.Repeat("░", emptyWidth)

	return bar + " " + colors[score].Render(labels[score])
}

func renderStrengthText(text string, score int) string {
	switch score {
	case 0, 1:
		return DangerStyle.Render(text)
	case 2:
		return WarningStyle.Render(text)
	case 3:
		return StrongPwdStyle.Render(text)
	case 4:
		return VeryStrongPwdStyle.Render(text)
	default:
		return text
	}
}
