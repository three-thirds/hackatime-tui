package app

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	headerBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#EB6F92"))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1)

	accentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#44475A"))

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

	profileContent := fmt.Sprintf("[●] User\n%s", StreakStyle.Render("[0 day streak]"))
	profileBox := lipgloss.NewStyle().Width(leftColWidth).Render(profileContent)

	todayTime := MetricValueStyle.Render("0h 00m")
	greetingContent := fmt.Sprintf("Keep Track of Your Coding Time\nToday: %s logged using Neovim & VSCode", todayTime)
	greetingBox := lipgloss.NewStyle().Width(rightColWidth).Render(greetingContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, profileBox, " | ", greetingBox)

	divider := DimText.Render(strings.Repeat("─", max(0, contentWidth)))

	tabNames := []string{"Home", "Project", "Settings"}
	var navLines strings.Builder
	for i, name := range tabNames {
		line := DimText.Render(fmt.Sprintf("  %s", name))
		if i == currentTab {
			line = ActiveTabStyle.Render(fmt.Sprintf("> [%s]", name))
		}
		if i > 0 {
			navLines.WriteString("\n")
		}
		navLines.WriteString(line)

	}
	navBox := lipgloss.NewStyle().Width(leftColWidth).Render(navLines.String())

	filterBar := FilterBarStyle.Render("[Date: All Time ▾] [Project: All ▾] [Lang: All ▾] [OS: All ▾] [Editor: All ▾]")

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
	t := MetricTitleStyle.Render(title)
	v := MetricValueStyle.Render(val)
	return CardBorder.Width(width).Render(fmt.Sprintf("%s\n%s", t, v))
}
