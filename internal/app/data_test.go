package app

import (
	"testing"
	"time"

	"github.com/three-thirds/hackatime-tui/pkg/api"
)

func TestDashboardFromAPIMapsStatsAndToday(t *testing.T) {
	today := &api.StatusToday{}
	today.Data.GrandTotal.TotalSecs = 2610
	today.Data.GrandTotal.Text = "43 mins"
	today.Data.Goal.TargetSecs = 3600
	today.Data.Goal.TrackedSecs = 2610
	today.Data.Goal.CompletionPercent = 72.5

	stats := &api.Last7Days{}
	stats.Data.Username = "devaansh"
	stats.Data.TotalSecs = 7200
	stats.Data.Projects = []api.Days7Thing{
		{Name: "hackatime-tui", TotalSeconds: 3600, Percent: 50},
		{Name: "other", TotalSeconds: 3600, Percent: 50},
	}
	stats.Data.Languages = []api.Days7Thing{
		{Name: "Go", TotalSeconds: 4800, Percent: 66.7},
	}
	stats.Data.Editors = []api.Days7Thing{
		{Name: "Neovim", TotalSeconds: 7200, Percent: 100},
	}
	stats.Data.OperatingSystems = []api.Days7Thing{
		{Name: "Windows", TotalSeconds: 7200, Percent: 100},
	}

	data := dashboardFromAPI(today, stats)

	if data.Username != "devaansh" {
		t.Fatalf("username = %q", data.Username)
	}
	if data.TotalTime != 2*time.Hour {
		t.Fatalf("total time = %v", data.TotalTime)
	}
	if data.TodayLogged != 43*time.Minute+30*time.Second {
		t.Fatalf("today logged = %v", data.TodayLogged)
	}
	if data.GoalPercentage != 73 {
		t.Fatalf("goal percent = %d", data.GoalPercentage)
	}
	if data.TopProject != "hackatime-tui" {
		t.Fatalf("top project = %q", data.TopProject)
	}
	if data.TopLanguage != "Go" {
		t.Fatalf("top language = %q", data.TopLanguage)
	}
	if len(data.Projects) != 2 {
		t.Fatalf("projects len = %d", len(data.Projects))
	}
	if data.Projects[0].Duration != time.Hour {
		t.Fatalf("project duration = %v", data.Projects[0].Duration)
	}
}

func TestFilterBreakdown(t *testing.T) {
	items := breakdownFromDays7([]api.Days7Thing{
		{Name: "Go", TotalSeconds: 10, Percent: 50},
		{Name: "Rust", TotalSeconds: 10, Percent: 50},
	})

	filtered := filterBreakdown(items, "Rust")
	if len(filtered) != 1 || filtered[0].Name != "Rust" {
		t.Fatalf("filtered = %+v", filtered)
	}

	all := filterBreakdown(items, "All")
	if len(all) != 2 {
		t.Fatalf("all filter len = %d", len(all))
	}
}

func TestStatsRangeForFilter(t *testing.T) {
	cases := map[string]string{
		"Last 7 Days":  "last_7_days",
		"Last 30 Days": "last_30_days",
		"Today":        "today",
		"All Time":     "all_time",
		"something":    "last_7_days",
	}
	for in, want := range cases {
		if got := statsRangeForFilter(in); got != want {
			t.Errorf("statsRangeForFilter(%q) = %q, want %q", in, got, want)
		}
	}
}
