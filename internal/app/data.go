package app

import (
	"fmt"
	"time"

	"github.com/three-thirds/hackatime-tui/internal/model"
	"github.com/three-thirds/hackatime-tui/pkg/api"
)

// dashboardFromAPI maps Hackatime API payloads into the dashboard shape the
// widgets already consume. Timeline, heatmap, streak, and AI/human split are
// left empty when the current endpoints do not provide them.
func dashboardFromAPI(today *api.StatusToday, stats *api.Last7Days) model.DashboardData {
	data := model.DashboardData{
		Filters: model.FilterState{
			DateRange: "Last 7 Days",
			Project:   "All",
			Language:  "All",
			OS:        "All",
			Editor:    "All",
		},
	}

	if stats != nil {
		data.Username = stats.Data.Username
		data.TotalTime = secondsToDuration(stats.Data.TotalSecs)
		data.Projects = breakdownFromDays7(stats.Data.Projects)
		data.Languages = breakdownFromDays7(stats.Data.Languages)
		data.Editors = breakdownFromDays7(stats.Data.Editors)
		data.OSList = breakdownFromDays7(stats.Data.OperatingSystems)
		data.TopProject = topName(data.Projects)
		data.TopLanguage = topName(data.Languages)
		data.TopOS = topName(data.OSList)
		data.TopEditor = topName(data.Editors)
		data.TodaySummary = buildTodaySummary(data.Languages, data.Editors)
		if stats.Data.HumanReadableRange != "" {
			data.Filters.DateRange = stats.Data.HumanReadableRange
		}
	}

	if today != nil {
		data.TodayLogged = secondsToDuration(today.Data.GrandTotal.TotalSecs)
		data.GoalPercentage = int(today.Data.Goal.CompletionPercent + 0.5)
		tracked := secondsToDuration(today.Data.Goal.TrackedSecs)
		target := secondsToDuration(today.Data.Goal.TargetSecs)

		switch {
		case today.Data.Goal.Complete:
			data.GoalStatusText = "Goal complete"
		case data.GoalPercentage >= 100:
			data.GoalStatusText = "Goal met"
		default:
			data.GoalStatusText = fmt.Sprintf("%d%% to go", max(100-data.GoalPercentage, 0))
		}

		if target > 0 {
			data.GoalDetailText = fmt.Sprintf("Today: %s / Goal: %s", formatDuration(tracked), formatDuration(target))
		} else if today.Data.GrandTotal.Text != "" {
			data.GoalDetailText = fmt.Sprintf("Today: %s", today.Data.GrandTotal.Text)
		} else {
			data.GoalDetailText = fmt.Sprintf("Today: %s", formatDuration(data.TodayLogged))
		}
	}

	return data
}

func breakdownFromDays7(items []api.Days7Thing) []model.BreakdownItem {
	out := make([]model.BreakdownItem, 0, len(items))
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		out = append(out, model.BreakdownItem{
			Name:       item.Name,
			Duration:   secondsToDuration(item.TotalSeconds),
			Percentage: item.Percent,
		})
	}
	return out
}

func secondsToDuration(secs float64) time.Duration {
	if secs <= 0 {
		return 0
	}
	return time.Duration(secs * float64(time.Second))
}

func topName(items []model.BreakdownItem) string {
	if len(items) == 0 {
		return ""
	}
	return items[0].Name
}

func buildTodaySummary(languages, editors []model.BreakdownItem) string {
	var parts []string
	if names := joinNames(languages, 3); names != "" {
		parts = append(parts, "("+names+")")
	}
	if names := joinNames(editors, 2); names != "" {
		parts = append(parts, "using "+names)
	}
	if len(parts) == 0 {
		return ""
	}
	return joinWithSpace(parts)
}

func joinNames(items []model.BreakdownItem, limit int) string {
	if limit <= 0 || len(items) == 0 {
		return ""
	}
	n := min(limit, len(items))
	names := make([]string, 0, n)
	for i := 0; i < n; i++ {
		names = append(names, items[i].Name)
	}
	if len(items) > n {
		names = append(names, "etc.")
	}
	return joinComma(names)
}

func joinComma(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " & " + parts[1]
	default:
		out := parts[0]
		for i := 1; i < len(parts)-1; i++ {
			out += ", " + parts[i]
		}
		return out + ", " + parts[len(parts)-1]
	}
}

func joinWithSpace(parts []string) string {
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += " " + parts[i]
	}
	return out
}

func filterBreakdown(items []model.BreakdownItem, selected string) []model.BreakdownItem {
	if selected == "" || selected == "All" {
		return items
	}
	for _, item := range items {
		if item.Name == selected {
			return []model.BreakdownItem{item}
		}
	}
	return nil
}

func breakdownNames(items []model.BreakdownItem) []string {
	names := make([]string, 0, len(items)+1)
	names = append(names, "All")
	for _, item := range items {
		names = append(names, item.Name)
	}
	return names
}

func statsRangeForFilter(label string) string {
	switch label {
	case "Last 30 Days":
		return "last_30_days"
	case "Today":
		return "today"
	case "All Time":
		return "all_time"
	default:
		return "last_7_days"
	}
}
