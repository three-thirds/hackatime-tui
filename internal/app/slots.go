package app

import (
	"fmt"
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

	titleText := SlotTitleStyle.Render(fmt.Sprintf("─ %s ", title))
	dimText := DimText.Render(fmt.Sprintf("(w: %d, h: %d)", width, height))
	instructions := DimText.Render(fmt.Sprintf("\n[Widget Slot] %s", hint))

	content := fmt.Sprintf("%s%s\n%s", titleText, dimText, instructions)
	return SlotBorder.Width(boxW).Height(boxH).Render(content)
}
