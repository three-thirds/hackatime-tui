package app

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	headerBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("240"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)

	accentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242"))
)

func renderFixedHeader(width int, currentTab int) string {
	if width < 50 {
		return "Terminal width is too small"
	}

	contentWidth := width - 4
	leftColWidth := 22
	rightColWidth := max(contentWidth-leftColWidth-3, 20)

	profileBox := lipgloss.NewStyle().Width(leftColWidth).Render("User \n[0 day streak]")
	greetingBox := lipgloss.NewStyle().Width(rightColWidth).Render("Keep Track of Your Coding Time\nToday: 0h 00m logged using Neovim & VSCode")
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, profileBox, " | ", greetingBox)

	divider := dimStyle.Render(strings.Repeat("-", max(0, contentWidth)))

	tabNames := []string{"Home", "Project", "Settings"}
	var navLines strings.Builder
	for i, name := range tabNames {
		line := fmt.Sprintf("   %s", name)
		if i == currentTab {
			line = accentStyle.Render(fmt.Sprintf("> [%s]", name))
		}
		if i > 0 {
			navLines.WriteString("\n")
		}
		navLines.WriteString(line)

	}
	navBox := lipgloss.NewStyle().Width(leftColWidth).Render(navLines.String())

	filterBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("75")).
		Render("[Date: All Time ▾] [Project: All ▾] [Lang: All ▾] [OS: All ▾] [Editor: All ▾]")

	cardW := max((rightColWidth-8)/5, 10)

	c1 := renderCard("TOTAL TIME", "--h --m", cardW)
	c2 := renderCard("TOP PROJECT", "--", cardW)
	c3 := renderCard("TOP LANGUAGE", "--", cardW)
	c4 := renderCard("TOP OS", "--", cardW)
	c5 := renderCard("TOP EDITOR", "--", cardW)
	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, c1, c2, c3, c4, c5)

	rightBlock := lipgloss.JoinVertical(lipgloss.Left, filterBar, cardsRow)
	middleRow := lipgloss.JoinHorizontal(lipgloss.Top, navBox, " │ ", rightBlock)

	fullHeader := lipgloss.JoinVertical(lipgloss.Left, topRow, divider, middleRow)
	return headerBorderStyle.Width(width - 2).Render(fullHeader)
}

func renderCard(title, val string, width int) string {
	t := dimStyle.Render(title)
	v := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255")).Render(val)
	return cardStyle.Width(width).Render(fmt.Sprintf("%s\n%s", t, v))
}
