package app

import (
	"fmt"
)

// RenderSlot draws a bordered wireframe container with sizing labels.
// It clamps minimum width and height to prevent Lipgloss layout panics
// on small terminal screens.
func RenderSlot(title string, width, height int, hint string) string {
	boxW := width
	boxH := height
	if boxW < 10 {
		boxW = 10
	}
	if boxH < 3 {
		boxH = 3
	}

	// TODO: Replace placeholder text later
	titleText := SlotTitleStyle.Render(fmt.Sprintf("─ %s ", title))
	dimText := DimText.Render(fmt.Sprintf("(w: %d, h: %d)", width, height))
	instructions := DimText.Render(fmt.Sprintf("\n[Widget Slot] %s", hint))

	content := fmt.Sprintf("%s%s\n%s", titleText, dimText, instructions)
	return SlotBorder.Width(boxW).Height(boxH).Render(content)
}
