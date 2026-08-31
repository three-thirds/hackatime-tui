// This file will store the shape of data required for app to function properly
// @Willgob will use this definitions to parse his JSON
package model

import "time"

type ActiveTab int

const (
	TabHome ActiveTab = iota
	TabProjects
	TabSettings
)

type FilterState struct {
	DateRange string //let's just keep it till today, yesterday, last week and all time
	Project   string //"All" or type in project name maybe...
	Language  string
	OS        string
	Editor    string
}

// BreakdownItem represents barcharts like Projects, Editor, Language, OS, Filter, etc.
type BreakdownItem struct {
	Name       string
	Duration   time.Duration
	Percentage float64
}

type TimelinePoint struct {
	Label    string
	Duration time.Duration
}

// HeatmapDay represents a single cell in the GitHub-style activity graph
type HeatmapDay struct {
	Date     time.Time
	Duration time.Duration
	Level    int // 0 to 4 (for color intensity: 0=dim/empty, 4=brightest)
}

// The main data structure for app state
type DashboardData struct {
	Username     string
	CountryFlag  string
	StreakDays   int
	TodayLogged  time.Duration
	TodaySummary string

	Filters FilterState

	TotalTime   time.Duration
	TopProject  string
	TopLanguage string
	TopOS       string
	TopEditor   string

	Projects  []BreakdownItem
	Languages []BreakdownItem
	Editors   []BreakdownItem
	OSList    []BreakdownItem

	Timeline []TimelinePoint
	Heatmap  []HeatmapDay

	// Goal & AI vs Human
	GoalPercentage  int
	GoalStatusText  string
	GoalDetailText  string
	AITime          time.Duration
	AIPercentage    int
	HumanTime       time.Duration
	HumanPercentage int
}
