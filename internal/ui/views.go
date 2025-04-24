package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"

	"datflux/internal/monitor"
	"datflux/internal/password"
)

// layout maintenance
const (
	MinPasswordPanelWidth = 80
	MinMetricsPanelWidth  = 80
	MinScreenWidth        = 84
	MinScreenHeight       = 40
)

func styledHeader(text string, width int) string {

	headerText := SectionTitleStyle.Render(text)

	totalHeaderWidth := lipgloss.Width(headerText)

	padding := max(0, (width-4-totalHeaderWidth)/2)

	header := strings.Repeat(" ", padding) + headerText + strings.Repeat(" ", padding)

	if (width-4-totalHeaderWidth)%2 != 0 {
		header += " "
	}

	return header
}

func renderHexQuality(quality float64) string {
	percentage := int(quality * 100)

	// percentage 2 hex
	hexValue := fmt.Sprintf("%02X", percentage)

	var statusStyle lipgloss.Style
	var statusText string
	var statusIndicator string

	switch {
	case percentage < 13:
		statusText = "VOID"
		statusIndicator = "⣀" // empty gauge
		statusStyle = DangerStyle
	case percentage < 25:
		statusText = "POOR"
		statusIndicator = "⣀" // 25% filled gauge
		statusStyle = DangerStyle
	case percentage < 38:
		statusText = "WEAK"
		statusIndicator = "⣄" // 37.5% filled gauge
		statusStyle = DangerStyle
	case percentage < 50:
		statusText = "FAIR"
		statusIndicator = "⣤" // 50% filled gauge
		statusStyle = WarningStyle
	case percentage < 63:
		statusText = "MILD"
		statusIndicator = "⣦" // 62.5% filled gauge
		statusStyle = WarningStyle
	case percentage < 75:
		statusText = "FIRM"
		statusIndicator = "⣶" // 75% filled gauge
		statusStyle = WarningStyle
	case percentage < 88:
		statusText = "GOOD"
		statusIndicator = "⣷" // 87.5% filled gauge
		statusStyle = StrongPwdStyle
	default:
		statusText = "PEAK"
		statusIndicator = "⣿" // 100% filled gauge
		statusStyle = VeryStrongPwdStyle
	}

	// hexadecimal timestamp
	// timestamp := time.Now().UnixNano() % 0xFFFFFF
	// timeHex := fmt.Sprintf("%06X", timestamp)

	indicatorStyle := statusStyle
	indicatorStyle = indicatorStyle.Faint(true)

	return fmt.Sprintf("%s%s%s %s %s %s",
		BracketStyle.Render("["),
		HexStyle.Render("0x"+hexValue+"%"),
		BracketStyle.Render("]"),
		LabelStyle.Render("ENTROPY"),
		indicatorStyle.Render(statusIndicator),
		statusStyle.Render(statusText))
}

func renderCPUView(cpuUsage float64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(styledHeader("CPU LOAD", width))
	builder.WriteString("\n\n")

	progressView := FormatProgressBar(progressBar, cpuUsage, width-10)
	progressWithPercentage := AddPercentage(progressView, cpuUsage, width)

	builder.WriteString(progressWithPercentage)

	return BorderStyle.Width(width).Render(builder.String())
}

func renderMemoryView(memUsage float64, memTotal uint64, memUsed uint64, progressBar progress.Model, width int) string {
	var builder strings.Builder

	builder.WriteString(styledHeader("MEMORY LOAD", width))
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

	title := fmt.Sprintf("NETWORK TRAFFIC (%s)", iface)
	builder.WriteString(styledHeader(title, width))
	builder.WriteString("\n\n")

	downloadText := fmt.Sprintf("↓ %s", monitor.FormatSpeed(rxSpeed))
	uploadText := fmt.Sprintf("↑ %s", monitor.FormatSpeed(txSpeed))

	downloadFormatted := ValueStyle.Render(downloadText)
	uploadFormatted := ValueStyle.Render(uploadText)

	columnWidth := (width - 4) / 2
	builder.WriteString(lipgloss.NewStyle().
		Width(columnWidth).
		Align(lipgloss.Center).
		Render(downloadFormatted))

	builder.WriteString(lipgloss.NewStyle().
		Width(columnWidth).
		Align(lipgloss.Center).
		Render(uploadFormatted))

	return BorderStyle.Width(width).Render(builder.String())
}

func renderPasswordView(animation *PasswordAnimation, quality float64, width int, passwordGen *password.Generator, attackModel password.AttackModelType) string {
	var builder strings.Builder

	// title := "CRYPTOGRAPHICALLY SECURED PASSWORD"
	// builder.WriteString(styledHeader(title, width))
	// builder.WriteString("\n")

	// if !animation.IsAnimating {
	qualityText := renderHexQuality(quality)
	qualityStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width - 4)
	builder.WriteString(qualityStyle.Render(qualityText))
	builder.WriteString("\n\n")
	// } else {
	// builder.WriteString("\n")
	// }

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

	if animation.IsAnimating || passwordText == "Press 'r' to generate" || len(passwordText) == 0 {
		// when animating or showing initial prompt, use bottom padding
		return PassBorderStyle.Width(width).Render(builder.String())
	} else {
		// when showing stats, no bottom padding
		return BorderStyle.Width(width).Render(builder.String())
	}
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
