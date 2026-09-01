package app

import "charm.land/lipgloss/v2"

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
)
