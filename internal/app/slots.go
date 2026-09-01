package app

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

var (
	slotBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)

	slotTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
)

func RenderSlot(title string, width, height int, hint string) string {
	boxW := width - 4
	boxH := height - 2

	if boxW < 10 {
		boxW = 10
	}
	if boxH < 3 {
		boxH = 3
	}

	titleText := slotTitle.Render(fmt.Sprintf("- %s", title))
	dimText := dimStyle.Render(fmt.Sprintf("(w: %d, h: %d)", width, height))
	instructions := dimStyle.Render(fmt.Sprintf("\n[Widget Slot] %s", hint))

	content := fmt.Sprintf("%s%s\n%s", titleText, dimText, instructions)
	return slotBorder.Width(boxW).Height(boxH).Render(content)

}
