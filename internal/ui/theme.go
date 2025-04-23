// Tested:
// - Tokyo Night (Perfect default)
// - Ozone 10 (Good)
// - Hydrangea 11 (Lovely)
// - Leopold's Dreams (Good)
// - Nyx8 (Does not fit)
// - Citrink (Does not fit)
// - Midnight Ablaze (Paranoia Mode)
// ----------------------------------
// To try:
// - Moonlight GB
// - 2bit demichrome
// - Oil 6
// - Berry Nebula
// - Cryptic Ocean
// - Kirokaze Gameboy
// - Dream Haze 8
// - EPHEMERA
// - Gothic Bit
// - bluem0ld
// - BloodMoon21
// - CUSTODIAN-8
// - slimy 05
// - ABYSS-9
// - vividmemory8

package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type ThemeType string

const (
	ThemeTokyoNight     ThemeType = "tokyo-night"
	ThemeOzone10        ThemeType = "ozone-10"
	ThemeHydrangea      ThemeType = "hydrangea 11"
	ThemeLeopoldsDreams ThemeType = "leopold's dreams"
	ThemeMidnightAblaze ThemeType = "midnight ablaze" // paranoia theme
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
	tm.RegisterTheme(ThemeOzone10, createOzone10Theme())
	tm.RegisterTheme(ThemeHydrangea, createHydrangeaTheme())
	tm.RegisterTheme(ThemeLeopoldsDreams, createLeopoldsDreamsTheme())

	return tm
}

func (tm *ThemeManager) RegisterTheme(themeType ThemeType, theme Theme) {
	tm.themes[themeType] = theme
}

func (tm *ThemeManager) GetCurrentTheme() Theme {
	return tm.themes[tm.currentTheme]
}

func GetDefaultTheme() Theme {
	return createTokyoNightTheme()
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
		VeryStrongColor: lipgloss.Color("#b4f9f8"), // PaleCyan

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
func createOzone10Theme() Theme {
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

// Leopold's Dreams
func createLeopoldsDreamsTheme() Theme {
	return Theme{
		Name:        "Leopold's Dreams",
		Description: "A blueish melancholic palette, good for water and/or sci-fi scenarios.",

		Primary:   lipgloss.Color("#372134"),
		Secondary: lipgloss.Color("#474476"),
		Accent:    lipgloss.Color("#6dbcb9"),
		Highlight: lipgloss.Color("#4888b7"),
		Warning:   lipgloss.Color("#6dbcb9"),
		Danger:    lipgloss.Color("#52c33f"),

		PasswordColor:   lipgloss.Color("#4888b7"),
		StrongColor:     lipgloss.Color("#6dbcb9"),
		VeryStrongColor: lipgloss.Color("#8cefb6"),

		PrimaryStr:         "#372134",
		SecondaryStr:       "#474476",
		AccentStr:          "#6dbcb9",
		HighlightStr:       "#4888b7",
		WarningStr:         "#6dbcb9",
		DangerStr:          "#52c33f",
		StrongColorStr:     "#6dbcb9",
		VeryStrongColorStr: "#8cefb6",
	}
}

// Midnight Ablaze
func createMidnightAblazeTheme() Theme {
	return Theme{
		Name:        "Midnight Ablaze",
		Description: "Made for a very ominous night sky.",

		Primary:   lipgloss.Color("#1f0510"),
		Secondary: lipgloss.Color("#460e2b"),
		Accent:    lipgloss.Color("#7c183c"),
		Highlight: lipgloss.Color("#d53c6a"),
		Warning:   lipgloss.Color("#7c183c"),
		Danger:    lipgloss.Color("#460e2b"),

		PasswordColor:   lipgloss.Color("#ff8274"),
		StrongColor:     lipgloss.Color("#d53c6a"),
		VeryStrongColor: lipgloss.Color("#ff8274"),

		PrimaryStr:         "#1f0510",
		SecondaryStr:       "#460e2b",
		AccentStr:          "#7c183c",
		HighlightStr:       "#d53c6a",
		WarningStr:         "#7c183c",
		DangerStr:          "#460e2b",
		StrongColorStr:     "#d53c6a",
		VeryStrongColorStr: "#ff8274",
	}
}
