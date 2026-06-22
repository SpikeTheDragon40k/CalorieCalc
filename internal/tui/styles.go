package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	colorGreen  = lipgloss.Color("#4ade80")
	colorYellow = lipgloss.Color("#fbbf24")
	colorRed    = lipgloss.Color("#f87171")
	colorBlue   = lipgloss.Color("#60a5fa")
	colorPurple = lipgloss.Color("#a78bfa")
	colorGray   = lipgloss.Color("#6b7280")
	colorWhite  = lipgloss.Color("#e5e7eb")
	colorBlack  = lipgloss.Color("#1f2937")
	colorOrange = lipgloss.Color("#fb923c")

	// Base styles
	appStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Background(lipgloss.Color("#0f172a"))

	titleStyle = lipgloss.NewStyle().
			Foreground(colorBlue).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 1)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray).
			Padding(1, 2)

	activeBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBlue).
				Padding(1, 2)

	itemStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			PaddingLeft(1).
			Width(40)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorBlue).
				Background(lipgloss.Color("#1e3a5f")).
				PaddingLeft(1).
				Width(40)

	dimmedStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	greenText = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	yellowText = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	redText = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	blueText = lipgloss.NewStyle().
			Foreground(colorBlue)

	purpleText = lipgloss.NewStyle().
			Foreground(colorPurple)

	errorText = lipgloss.NewStyle().
			Foreground(colorRed).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			PaddingTop(1).
			PaddingLeft(1)

	keyStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	highlightStyle = lipgloss.NewStyle().
			Foreground(colorOrange).
			Bold(true)

	totalBarStyle = lipgloss.NewStyle().
			Height(1)

	separator = lipgloss.NewStyle().
			Foreground(colorGray).
			Padding(0, 1).
			Render("│")

	checkboxEmpty = lipgloss.NewStyle().
			Foreground(colorGray).
			Render("○")

	checkboxFull = lipgloss.NewStyle().
			Foreground(colorGreen).
			Render("●")
)

func KcalColor(kcal, target float64) lipgloss.Style {
	if target <= 0 {
		return greenText
	}
	ratio := kcal / target
	switch {
	case ratio <= 0.8:
		return greenText
	case ratio <= 1.0:
		return yellowText
	default:
		return redText
	}
}

func KcalBar(kcal, target float64, width int) string {
	if target <= 0 {
		return ""
	}
	if kcal > target {
		kcal = target
	}
	filled := int((kcal / target) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
