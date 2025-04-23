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
	ParanoiaMode    bool
}

func NewPasswordAnimation(generator *password.Generator) *PasswordAnimation {
	return &PasswordAnimation{
		Target:          "",
		Current:         "Press 'r' to generate",
		ColorPhase:      0,
		IsAnimating:     false,
		Progress:        0,
		FlickersPerChar: 3,
		Delay:           time.Millisecond * 40,
		Generator:       generator,
		lastUpdateTime:  time.Now(),
	}
}

// renders the revealed portion
func (pa *PasswordAnimation) RenderRevealedPart(revealLen int) string {
	if revealLen <= 0 {
		return ""
	}
	chars := []rune(pa.Current)
	return PasswordStyle.Render(string(chars[:revealLen]))
}

// renders the binary entropy stream
func (pa *PasswordAnimation) RenderBinaryStream(length int) string {
	var stream strings.Builder

	for i := 0; i < length; i++ {
		// binary digits based on animation phase and position
		bitValue := rand.Intn(2)
		if bitValue == 0 {
			stream.WriteString(DangerStyle.Render("0"))
		} else {
			stream.WriteString(StrongPwdStyle.Render("1"))
		}
	}
	return stream.String()
}

// renders the next character to be revealed
func (pa *PasswordAnimation) RenderNextChar(position int) string {
	chars := []rune(pa.Current)

	if position >= len(chars) {
		return ""
	}

	nextChar := string(chars[position])
	if pa.Progress%pa.FlickersPerChar == 0 {
		nextChar = string(pa.Generator.GenerateRandomChar())
	}

	return VeryStrongPwdStyle.Bold(true).Render(nextChar)
}

// special animation for paranoia mode
func (pa *PasswordAnimation) ParanoiaModeAnimation() string {
	if !pa.IsAnimating {
		return PasswordStyle.Render(pa.Current)
	}

	var result strings.Builder
	charPos := pa.Progress / pa.FlickersPerChar
	chars := []rune(pa.Current)
	revealLen := min(charPos, len(chars))

	// show what's revealed so far
	result.WriteString(pa.RenderRevealedPart(revealLen))

	// show the binary stream animation, if not complete
	if revealLen < len(chars) {
		// add any spacing needed between revealed part and animation
		result.WriteString(" ")

		// binary entropy stream (varies with progress)
		streamLength := min(20, len(chars)-revealLen) // up to 20 binary digits
		binaryStream := pa.RenderBinaryStream(streamLength)
		result.WriteString(" " +
			VeryStrongPwdStyle.Render("[") +
			binaryStream +
			VeryStrongPwdStyle.Render("]") +
			" ")

		// small preview of next character
		result.WriteString(pa.RenderNextChar(revealLen))
	}

	return result.String()
}

// begins a new password reveal animation
func (pa *PasswordAnimation) StartAnimation(password string) {
	pa.Target = password
	pa.IsAnimating = true
	pa.Progress = 0
	pa.ColorPhase = 0
	pa.Current = strings.Repeat("?", len(password))
	pa.lastUpdateTime = time.Now()

	// adjust for paranoia mode
	if pa.ParanoiaMode && len(password) > 40 {
		// Faster animation
		pa.FlickersPerChar = 2           // fewer flickers per char
		pa.Delay = time.Millisecond * 20 // halve delay
	} else {
		// standard speed
		pa.FlickersPerChar = 3
		pa.Delay = time.Millisecond * 40
	}
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
	if pa.ParanoiaMode {
		return pa.ParanoiaModeAnimation()
	}

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
