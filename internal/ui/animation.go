package ui

import (
	"math/rand"
	"strings"
	"time"

	"datflux/internal/password"

	"github.com/charmbracelet/lipgloss"
)

type PasswordAnimation struct {
	Target          string
	Current         string
	ColorPhase      int
	IsAnimating     bool
	Progress        int
	FlickersPerChar int
	Delay           time.Duration
	Generator       *password.Generator
	lastUpdateTime  time.Time
}

func NewPasswordAnimation(generator *password.Generator) *PasswordAnimation {
	return &PasswordAnimation{
		Target:          "",
		Current:         "Press 'r' to generate a password",
		ColorPhase:      0,
		IsAnimating:     false,
		Progress:        0,
		FlickersPerChar: 3,
		Delay:           time.Millisecond * 40,
		Generator:       generator,
		lastUpdateTime:  time.Now(),
	}
}

// begins a new password reveal animation
func (pa *PasswordAnimation) StartAnimation(password string) {
	pa.Target = password
	pa.IsAnimating = true
	pa.Progress = 0
	pa.ColorPhase = 0
	pa.Current = strings.Repeat("?", len(password))
	pa.lastUpdateTime = time.Now()
}

func (pa *PasswordAnimation) Update() bool {
	if !pa.IsAnimating {
		return false
	}

	now := time.Now()
	if now.Sub(pa.lastUpdateTime) < pa.Delay {
		return false
	}
	pa.lastUpdateTime = now

	charPos := pa.Progress / pa.FlickersPerChar
	flickerPos := pa.Progress % pa.FlickersPerChar

	if charPos >= len(pa.Target) {
		pa.Current = pa.Target
		pa.IsAnimating = false
		return true
	}

	currentRunes := []rune(pa.Current)

	if flickerPos == pa.FlickersPerChar-1 {
		currentRunes[charPos] = []rune(pa.Target)[charPos]
	} else {
		currentRunes[charPos] = rune(pa.Generator.GenerateRandomChar())
	}

	for i := charPos + 1; i < len(currentRunes); i++ {
		currentRunes[i] = rune(pa.Generator.GenerateRandomChar())
	}

	pa.Current = string(currentRunes)
	pa.Progress++

	pa.ColorPhase = (pa.ColorPhase + 1) % 6

	return true
}

func (pa *PasswordAnimation) IsComplete() bool {
	return !pa.IsAnimating
}

func (pa *PasswordAnimation) CurrentPassword() string {
	return pa.Current
}

func (pa *PasswordAnimation) StyledPassword() string {
	if !pa.IsAnimating {
		return PasswordStyle.Render(pa.Current)
	}

	var result strings.Builder
	chars := []rune(pa.Current)

	charPos := pa.Progress / pa.FlickersPerChar

	for i, char := range chars {
		var styledChar string

		if i < charPos {
			styledChar = PasswordStyle.Render(string(char))
		} else if i == charPos {
			colorIndex := pa.ColorPhase % len(PasswordAnimationColors)
			animColor := PasswordAnimationColors[colorIndex]
			styledChar = lipgloss.NewStyle().Foreground(animColor).Bold(true).Render(string(char))
		} else {
			randColorIdx := rand.Intn(len(PasswordAnimationColors))
			animColor := PasswordAnimationColors[randColorIdx]
			styledChar = lipgloss.NewStyle().Foreground(animColor).Render(string(char))
		}

		result.WriteString(styledChar)
	}

	return result.String()
}
