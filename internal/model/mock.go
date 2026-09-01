package model

import "time"

// GetMockDashboardData provides rich dummy data for UI testing
func GetMockDashboardData() DashboardData {
	now := time.Now()

	// Generate 30 days of mock heatmap activity
	var mockHeatmap []HeatmapDay
	for i := 29; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		level := (i * 7) % 5 // pseudo-random level between 0 and 4
		mockHeatmap = append(mockHeatmap, HeatmapDay{
			Date:     day,
			Duration: time.Duration(level*2) * time.Hour,
			Level:    level,
		})
	}

	return DashboardData{
		Username:     "Chish",
		CountryFlag:  "🇮🇳",
		StreakDays:   3,
		TodayLogged:  43*time.Minute + 27*time.Second,
		TodaySummary: "(Python, JS, Rust, etc.) using Neovim & VSCode",

		Filters: FilterState{
			DateRange: "All Time",
			Project:   "All",
			Language:  "All",
			OS:        "All",
			Editor:    "All",
		},

		TotalTime:   520*time.Hour + 34*time.Minute,
		TopProject:  "kasumi",
		TopLanguage: "Rust",
		TopOS:       "Linux",
		TopEditor:   "VSCode",

		Projects: []BreakdownItem{
			{Name: "kasumi", Duration: 31*time.Hour + 28*time.Minute, Percentage: 30},
			{Name: "skora-backend", Duration: 31*time.Hour + 27*time.Minute, Percentage: 30},
			{Name: "stacksense", Duration: 22*time.Hour + 25*time.Minute, Percentage: 20},
			{Name: "the-verse-core", Duration: 19*time.Hour + 54*time.Minute, Percentage: 15},
			{Name: "desktop-editor", Duration: 16*time.Hour + 44*time.Minute, Percentage: 10},
		},

		Languages: []BreakdownItem{
			{Name: "Rust", Percentage: 38},
			{Name: "Python", Percentage: 22},
			{Name: "Svelte", Percentage: 12},
			{Name: "C / Go", Percentage: 8},
			{Name: "Other", Percentage: 20},
		},

		Editors: []BreakdownItem{
			{Name: "VSCode", Percentage: 48},
			{Name: "Neovim", Percentage: 28},
			{Name: "Zed / CLion", Percentage: 14},
			{Name: "Godot / Other", Percentage: 10},
		},

		OSList: []BreakdownItem{
			{Name: "Linux", Percentage: 100},
		},

		Timeline: []TimelinePoint{
			{Label: "Jun15", Duration: 22 * time.Hour},
			{Label: "Jun22", Duration: 11 * time.Hour},
			{Label: "Jun29", Duration: 5 * time.Hour},
			{Label: "Jul6", Duration: 11 * time.Hour},
			{Label: "Jul13", Duration: 5 * time.Hour},
			{Label: "Jul20", Duration: 5 * time.Hour},
			{Label: "Jul27", Duration: 5 * time.Hour},
			{Label: "Aug31", Duration: 5 * time.Hour},
		},

		Heatmap: mockHeatmap,

		GoalPercentage:  73,
		GoalStatusText:  "27% below usual",
		GoalDetailText:  "Today: 43m / Usual: 59m",
		AITime:          32*time.Hour + 15*time.Minute,
		AIPercentage:    7,
		HumanTime:       450*time.Hour + 58*time.Minute,
		HumanPercentage: 93,
	}
}
