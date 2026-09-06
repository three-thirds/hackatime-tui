// This file owns the deck geometry: how much room each card gets and how much
// air sits between them. Widgets never compute their own position, they are
// handed a content width and render into it.

package app

const (
	// columnGap is the horizontal breathing room between two side-by-side cards.
	columnGap = 2
	// rowGap is the number of blank lines between two stacked rows.
	rowGap = 1
)

// halfCardWidth is the total width of one card in a two-column row, sized so
// both cards plus the gap fit exactly inside the terminal. Lipgloss counts a
// card's border and padding inside this width.
func halfCardWidth(terminalWidth int) int {
	return max((terminalWidth-columnGap)/2, minChartWidth)
}

// fullCardWidth is the width of a card spanning the whole row. It is derived
// from halfCardWidth so a full-width card lines up flush with the two
// half-width cards above it.
func fullCardWidth(terminalWidth int) int {
	return halfCardWidth(terminalWidth)*2 + columnGap
}
