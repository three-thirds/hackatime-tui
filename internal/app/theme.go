// This file stores all the color related information
//
//

package app

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Official Dracula Palette
var (
	DraculaBg        = lipgloss.Color("#282A36")
	DraculaSelection = lipgloss.Color("#44475A")
	DraculaFg        = lipgloss.Color("#F8F8F2")
	DraculaComment   = lipgloss.Color("#6272A4")
	DraculaCyan      = lipgloss.Color("#8BE9FD")
	DraculaGreen     = lipgloss.Color("#50FA7B")
	DraculaOrange    = lipgloss.Color("#FFB86C")
	DraculaPink      = lipgloss.Color("#FF79C6")
	DraculaPurple    = lipgloss.Color("#BD93F9")
	DraculaRed       = lipgloss.Color("#FF5555")
	DraculaYellow    = lipgloss.Color("#F1FA8C")
)

// Semantic UI Styles
var (
	// Borders
	HeaderBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DraculaPurple)

	CardBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DraculaSelection).
			Padding(0, 1)

	SlotBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(DraculaSelection).
			Padding(0, 1)

	// Typography & Accents
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(DraculaPink)

	SlotTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(DraculaPurple)

	FilterBarStyle = lipgloss.NewStyle().
			Foreground(DraculaCyan)

	MetricTitleStyle = lipgloss.NewStyle().
				Foreground(DraculaComment)

	MetricValueStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(DraculaGreen)

	DimText = lipgloss.NewStyle().
		Foreground(DraculaComment)

	StreakStyle = lipgloss.NewStyle().
			Foreground(DraculaOrange)

	FilterInactiveStyle = lipgloss.NewStyle().
				Foreground(DraculaCyan).
				Padding(0, 1)

	// Active/Focused filter button (Highlighted background + Pink text)
	FilterActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(DraculaPink).
				Background(DraculaSelection).
				Padding(0, 1)

	// Focused Zone indicator badge for Header
	FocusedBadgeStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(DraculaBg).
				Background(DraculaPurple).
				Padding(0, 1)
)

// Chart palettes. ChartPalette is cycled across the bars of a breakdown chart,
// HeatmapLevels maps a HeatmapDay.Level (0-4) onto its intensity.
var (
	ChartPalette = []color.Color{
		DraculaCyan,
		DraculaPurple,
		DraculaPink,
		DraculaGreen,
		DraculaOrange,
		DraculaYellow,
	}

	HeatmapLevels = []color.Color{
		DraculaSelection,
		lipgloss.Color("#1F5B34"),
		lipgloss.Color("#2E8B4F"),
		lipgloss.Color("#3FCB68"),
		DraculaGreen,
	}
)
