package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type ThemeType string

const (
	ThemeTokyoNight ThemeType = "tokyo-night"
	ThemeOzoneIO    ThemeType = "ozone-10"
	ThemeHydrangea  ThemeType = "hydrangea 11"
)

type Theme struct {
	Name        string
	Description string

	// base colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Highlight lipgloss.Color
	Warning   lipgloss.Color
	Danger    lipgloss.Color

	// specials
	PasswordColor   lipgloss.Color
	StrongColor     lipgloss.Color
	VeryStrongColor lipgloss.Color

	// string versions for progress bars
	PrimaryStr         string
	SecondaryStr       string
	AccentStr          string
	HighlightStr       string
	WarningStr         string
	DangerStr          string
	StrongColorStr     string
	VeryStrongColorStr string
}

type ThemeManager struct {
	currentTheme ThemeType
	themes       map[ThemeType]Theme
}

func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		currentTheme: ThemeTokyoNight,
		themes:       make(map[ThemeType]Theme),
	}

	tm.RegisterTheme(ThemeTokyoNight, createTokyoNightTheme())
	tm.RegisterTheme(ThemeOzoneIO, createOzoneIOTheme())
	tm.RegisterTheme(ThemeHydrangea, createHydrangeaTheme())

	return tm
}

func (tm *ThemeManager) RegisterTheme(themeType ThemeType, theme Theme) {
	tm.themes[themeType] = theme
}

func (tm *ThemeManager) GetCurrentTheme() Theme {
	return tm.themes[tm.currentTheme]
}

func (tm *ThemeManager) SwitchTheme(themeType ThemeType) bool {
	if _, exists := tm.themes[themeType]; exists {
		tm.currentTheme = themeType
		return true
	}
	return false
}

func (tm *ThemeManager) CycleTheme() ThemeType {
	var themeTypes []ThemeType
	for tt := range tm.themes {
		themeTypes = append(themeTypes, tt)
	}

	currentIndex := 0
	for i, tt := range themeTypes {
		if tt == tm.currentTheme {
			currentIndex = i
			break
		}
	}

	nextIndex := (currentIndex + 1) % len(themeTypes)
	tm.currentTheme = themeTypes[nextIndex]

	return tm.currentTheme
}

func (tm *ThemeManager) GetAvailableThemes() map[ThemeType]Theme {
	return tm.themes
}

// Tokyo Night
func createTokyoNightTheme() Theme {
	return Theme{
		Name:        "Tokyo Night",
		Description: "A combo of dark blues and vivid neons, inspired by the beautiful city of Tokyo at night.",

		Primary:   lipgloss.Color("#24283b"), // StormBlue
		Secondary: lipgloss.Color("#565f89"), // DarkBlue
		Accent:    lipgloss.Color("#bb9af7"), // Lavender
		Highlight: lipgloss.Color("#7dcfff"), // SkyBlue
		Warning:   lipgloss.Color("#ffc777"), // OrangeYellow
		Danger:    lipgloss.Color("#ff757f"), // PeachRed

		PasswordColor:   lipgloss.Color("#bb9af7"), // Lavender
		StrongColor:     lipgloss.Color("#c3e88d"), // PaleGreen
		VeryStrongColor: lipgloss.Color("#b4f9f8"), // LightBlue

		PrimaryStr:         "#24283b",
		SecondaryStr:       "#565f89",
		AccentStr:          "#bb9af7",
		HighlightStr:       "#7dcfff",
		WarningStr:         "#ffc777",
		DangerStr:          "#ff757f",
		StrongColorStr:     "#c3e88d",
		VeryStrongColorStr: "#b4f9f8",
	}
}

// Ozone-10
func createOzoneIOTheme() Theme {
	return Theme{
		Name:        "Ozone-10",
		Description: "A soft palette inspired by polluted city sky colors.",

		Primary:   lipgloss.Color("#746576"),
		Secondary: lipgloss.Color("#888189"),
		Accent:    lipgloss.Color("#c2d199"),
		Highlight: lipgloss.Color("#c4b865"),
		Warning:   lipgloss.Color("#bb9135"),
		Danger:    lipgloss.Color("#e4728f"),

		PasswordColor:   lipgloss.Color("#fffbc1"),
		StrongColor:     lipgloss.Color("#94c481"),
		VeryStrongColor: lipgloss.Color("#809f8c"),

		PrimaryStr:         "#746576",
		SecondaryStr:       "#888189",
		AccentStr:          "#c2d199",
		HighlightStr:       "#c4b865",
		WarningStr:         "#bb9135",
		DangerStr:          "#e4728f",
		StrongColorStr:     "#94c481",
		VeryStrongColorStr: "#809f8c",
	}
}

// Hydrangea 11
func createHydrangeaTheme() Theme {
	return Theme{
		Name:        "Hydrangea 11",
		Description: "A pastel palette inspired by hydrangea flowers.",

		Primary:   lipgloss.Color("#413652"),
		Secondary: lipgloss.Color("#6f577e"),
		Accent:    lipgloss.Color("#eae4dd"),
		Highlight: lipgloss.Color("#6f919c"),
		Warning:   lipgloss.Color("#986f9c"),
		Danger:    lipgloss.Color("#c090a7"),

		PasswordColor:   lipgloss.Color("#eae4dd"),
		StrongColor:     lipgloss.Color("#c9d4b8"),
		VeryStrongColor: lipgloss.Color("#90c0a0"),

		PrimaryStr:         "#413652",
		SecondaryStr:       "#6f577e",
		AccentStr:          "#eae4dd",
		HighlightStr:       "#6f919c",
		WarningStr:         "#986f9c",
		DangerStr:          "#c090a7",
		StrongColorStr:     "#c9d4b8",
		VeryStrongColorStr: "#90c0a0",
	}
}
