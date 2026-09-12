// This file consists with the logic to render Header
// the header has user info, filtering options, navigation between
// different tabs
//

package app

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/three-thirds/hackatime-tui/internal/model"
)

// renderFixedHeader renders the Header with the info about width and currently selectedTab
// this info is supposed to be received from central state model `AppModel`
// It returns the widget string with user Info, filter tabs, navigation panels and high level summaries.
func renderFixedHeader(width int,
	currentTab int,
	activeZone ActiveZone,
	activeFilter int,
	filterValues [5]string,
	dropdownOpen bool,
	dropdownCursor int,
	data model.DashboardData,
) string {
	if width < 50 {
		return "Terminal width is too small"
	}

	// Basic layout info. The header shares the deck's grid so its border lines
	// up with the cards scrolling underneath it.
	headerWidth := fullCardWidth(width)
	contentWidth := headerWidth - 2 // the header border sits inside headerWidth
	leftColWidth := 22
	rightColWidth := max(contentWidth-leftColWidth-3, 20)

	username := data.Username
	if username == "" {
		username = "User"
	}
	profileContent := fmt.Sprintf("[●] %s\n%s", username, StreakStyle.Render(fmt.Sprintf("[%d day streak]", data.StreakDays)))
	profileBox := lipgloss.NewStyle().Width(leftColWidth).Render(profileContent)

	todayTime := MetricValueStyle.Render(formatDuration(data.TodayLogged))
	summary := data.TodaySummary
	if summary == "" {
		summary = "coding"
	}
	greetingContent := fmt.Sprintf("Keep Track of Your Coding Time\nToday: %s logged %s", todayTime, summary)
	greetingBox := lipgloss.NewStyle().Width(rightColWidth).Render(greetingContent)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, profileBox, " | ", greetingBox)

	divider := DimText.Render(strings.Repeat("─", max(0, contentWidth)))

	// The navigation builder, it renders the navigation list, only three, I don't think
	// we would ever need to add or remove any of these.
	tabNames := []string{"Home", "Project", "Settings"}
	var navLines strings.Builder

	if activeZone == FocusNav {
		navLines.WriteString(ActiveTabStyle.Render("NAV MENU"))
		navLines.WriteString("\n")
	} else {
		navLines.WriteString(DimText.Render("NAV MENU"))
		navLines.WriteString("\n")
	}
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

	filterLabels := []string{
		fmt.Sprintf("Date: %s ▾", filterValues[0]),
		fmt.Sprintf("Project: %s ▾", filterValues[1]),
		fmt.Sprintf("Lang: %s ▾", filterValues[2]),
		fmt.Sprintf("OS: %s ▾", filterValues[3]),
		fmt.Sprintf("Editor: %s ▾", filterValues[4]),
	}

	// clamps the width of card to atleast 10 to avoid panic
	cardW := max((rightColWidth-8)/5, 10)

	var renderedFilter []string
	for i, name := range filterLabels {
		if i == activeFilter && activeZone == FocusFilter {
			styled := ActiveTabStyle.Render(fmt.Sprintf("[%s]", name))
			renderedFilter = append(renderedFilter, styled)
		} else {
			styled := FilterBarStyle.Render(fmt.Sprintf("[%s]", name))
			renderedFilter = append(renderedFilter, styled)
		}
	}

	filterBar := strings.Join(renderedFilter, " ")

	var bottomBlock string

	if dropdownOpen && activeZone == FocusFilter {
		options := breakdownNames(nil)
		switch activeFilter {
		case 0:
			options = []string{"Last 7 Days", "Last 30 Days", "Today", "All Time"}
		case 1:
			options = breakdownNames(data.Projects)
		case 2:
			options = breakdownNames(data.Languages)
		case 3:
			options = breakdownNames(data.OSList)
		case 4:
			options = breakdownNames(data.Editors)
		}
		var optionLines []string
		for idx, opt := range options {
			if idx == dropdownCursor {
				optionLines = append(optionLines, ActiveTabStyle.Render(fmt.Sprintf("  ▸ %s", opt)))
			} else {
				optionLines = append(optionLines, DimText.Render(fmt.Sprintf("    %s", opt)))
			}
		}

		boxContent := fmt.Sprintf("Select %s:\n%s", filterLabels[activeFilter], strings.Join(optionLines, "\n"))
		bottomBlock = CardBorder.Width(rightColWidth - 4).Render(boxContent)
	} else {
		c1 := renderCard("TOTAL TIME", formatOrDash(data.TotalTime), cardW)
		c2 := renderCard("TOP PROJECT", stringOrDash(data.TopProject), cardW)
		c3 := renderCard("TOP LANGUAGE", stringOrDash(data.TopLanguage), cardW)
		c4 := renderCard("TOP OS", stringOrDash(data.TopOS), cardW)
		c5 := renderCard("TOP EDITOR", stringOrDash(data.TopEditor), cardW)
		bottomBlock = lipgloss.JoinHorizontal(lipgloss.Top, c1, c2, c3, c4, c5)
	}

	rightBlock := lipgloss.JoinVertical(lipgloss.Left, filterBar, bottomBlock)
	middleRow := lipgloss.JoinHorizontal(lipgloss.Top, navBox, " │ ", rightBlock)

	fullHeader := lipgloss.JoinVertical(lipgloss.Left, topRow, divider, middleRow)
	return HeaderBorder.Width(headerWidth).Render(fullHeader)
}

// renderCard draws simple cards with title and value
// returns a small card, which can be later added into a deck
func renderCard(title, val string, width int) string {
	t := MetricTitleStyle.Render(title)
	v := MetricValueStyle.Render(val)
	return CardBorder.Width(width).Render(fmt.Sprintf("%s\n%s", t, v))
}

func stringOrDash(value string) string {
	if value == "" {
		return "--"
	}
	return value
}

func formatOrDash(d time.Duration) string {
	if d <= 0 {
		return "--"
	}
	return formatDuration(d)
}
