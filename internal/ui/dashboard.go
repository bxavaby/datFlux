package ui

import (
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"datflux/internal/entropy"
	"datflux/internal/monitor"
	"datflux/internal/password"
)

type tickMsg time.Time
type clipboardResultMsg struct {
	success bool
	message string
}

type clipboardClearMsg struct{}

type Dashboard struct {
	systemMonitor    *monitor.SystemMonitor
	passwordGen      *password.Generator
	entropyCollector *entropy.Collector
	animation        *PasswordAnimation
	width            int
	height           int
	ready            bool
	lastPassword     string
	clipboardStatus  string
	cpuProgress      progress.Model
	memProgress      progress.Model
}

func NewDashboardModel(collector *entropy.Collector) *Dashboard {
	sysMonitor := monitor.NewSystemMonitor()

	passGen := password.NewGenerator(collector)

	anim := NewPasswordAnimation(passGen)

	cpuBar := CPUProgress

	memBar := MemoryProgress

	return &Dashboard{
		systemMonitor:    sysMonitor,
		passwordGen:      passGen,
		entropyCollector: collector,
		animation:        anim,
		ready:            false,
		cpuProgress:      cpuBar,
		memProgress:      memBar,
	}
}

func (d *Dashboard) Init() tea.Cmd {
	return tickCmd()
}

func copyToClipboardCmd(text string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		var success bool
		var message string

		if _, err := exec.LookPath("xclip"); err == nil {
			// linux with X11
			cmd = exec.Command("xclip", "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(text)
			err = cmd.Run()
			success = (err == nil)
			if success {
				message = "Password copied to clipboard (xclip)"
			} else {
				message = "Failed to copy: " + err.Error()
			}
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			// linux with Wayland
			cmd = exec.Command("wl-copy")
			cmd.Stdin = strings.NewReader(text)
			err = cmd.Run()
			success = (err == nil)
			if success {
				message = "Password copied to clipboard (wl-copy)"
			} else {
				message = "Failed to copy: " + err.Error()
			}
		} else if _, err := exec.LookPath("pbcopy"); err == nil {
			// macOS
			cmd = exec.Command("pbcopy")
			cmd.Stdin = strings.NewReader(text)
			err = cmd.Run()
			success = (err == nil)
			if success {
				message = "Password copied to clipboard (pbcopy)"
			} else {
				message = "Failed to copy: " + err.Error()
			}
		} else {
			success = false
			message = "No clipboard command found"
		}

		return clipboardResultMsg{
			success: success,
			message: message,
		}
	}
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
		d.ready = true

		availableWidth := max(d.width-20, 20)

		d.cpuProgress.Width = availableWidth
		d.memProgress.Width = availableWidth

		return d, nil

	case tickMsg:
		d.systemMonitor.Update()

		entropySource := d.systemMonitor.GetEntropySource()
		d.entropyCollector.AddSample(entropySource)

		var cmds []tea.Cmd

		cmds = append(cmds, tickCmd())

		if d.animation.Update() {
			cmds = append(cmds, tea.Sequence())
		}

		prog, cmd := d.cpuProgress.Update(msg)
		d.cpuProgress = prog.(progress.Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		prog, cmd = d.memProgress.Update(msg)
		d.memProgress = prog.(progress.Model)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		return d, tea.Batch(cmds...)

	case clipboardResultMsg:
		d.clipboardStatus = msg.message
		return d, tea.Sequence(
			tea.Tick(3*time.Second, func(time.Time) tea.Msg {
				return clipboardClearMsg{}
			}),
		)

	case clipboardClearMsg:
		d.clipboardStatus = ""
		return d, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return d, tea.Quit

		case "r":
			if !d.animation.IsAnimating {
				newPassword := d.passwordGen.Generate()
				d.lastPassword = newPassword
				d.animation.StartAnimation(newPassword)
			}
			return d, nil

		case "c":
			if d.lastPassword != "" {
				return d, copyToClipboardCmd(d.lastPassword)
			}
		}
	}

	return d, nil
}

func (d *Dashboard) View() string {
	if !d.ready || d.width < 40 || d.height < 15 {
		return "Terminal too small. Please resize."
	}

	contentWidth := d.width - 4

	docStyle := lipgloss.NewStyle().Padding(0, 2)

	titleView := TitleStyle.Width(contentWidth).Render("🌸 [datFlux] Entropy-Borne Password Generator")

	// vertical layout always
	passwordView := renderPasswordView(
		d.animation,
		d.entropyCollector.GetEntropyQuality(),
		contentWidth-2,
	)

	cpuView := renderCPUView(
		d.systemMonitor.CPUUsage,
		d.cpuProgress,
		contentWidth-2,
	)

	memoryView := renderMemoryView(
		d.systemMonitor.MemoryUsage,
		d.systemMonitor.MemoryTotal,
		d.systemMonitor.MemoryUsed,
		d.memProgress,
		contentWidth-2,
	)

	networkView := renderNetworkView(
		d.systemMonitor.NetworkRxSpeed,
		d.systemMonitor.NetworkTxSpeed,
		d.systemMonitor.ActiveInterface,
		contentWidth-2,
	)

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		passwordView,
		"",
		cpuView,
		memoryView,
		networkView,
	)

	var helpText string
	if d.clipboardStatus != "" {
		helpText = ValueStyle.Render(d.clipboardStatus)
	} else {
		helpText = HelpStyle.Render("Press r to generate • Press c to copy • Press q to quit")
	}

	return docStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			titleView,
			"",
			mainView,
			"",
			helpText,
		),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
