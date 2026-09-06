package app

import (
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/three-thirds/hackatime-tui/internal/model"
)

func TestRenderBarChartRendersPercentValues(t *testing.T) {
	chart := RenderBarChart(BarChartConfig{
		Title: "Languages",
		Width: 40,
		Items: []model.BreakdownItem{
			{Name: "Rust", Percentage: 70},
			{Name: "Go", Percentage: 30},
		},
		Format: ValuePercent,
	})

	for _, want := range []string{"Languages", "Rust", "70%", "Go", "30%", "█"} {
		if !strings.Contains(chart, want) {
			t.Errorf("chart did not contain %q:\n%s", want, chart)
		}
	}
}

func TestRenderColumnChartRendersBarsAndAxis(t *testing.T) {
	chart := RenderColumnChart(ColumnChartConfig{
		Title:      "Timeline",
		Width:      40,
		PlotHeight: 5,
		Points: []model.TimelinePoint{
			{Label: "Jun15", Duration: 10 * time.Hour},
			{Label: "Jun22", Duration: 2 * time.Hour},
		},
	})

	for _, want := range []string{"Timeline", "Jun15", "Jun22", "█"} {
		if !strings.Contains(chart, want) {
			t.Errorf("chart did not contain %q:\n%s", want, chart)
		}
	}
}

func TestRenderHeatmapClipsToAvailableWidth(t *testing.T) {
	start := time.Date(2026, time.January, 4, 0, 0, 0, 0, time.UTC) // a Sunday

	var days []model.HeatmapDay
	for i := range 365 {
		days = append(days, model.HeatmapDay{Date: start.AddDate(0, 0, i), Level: i % 5})
	}

	const width = 60
	chart := RenderHeatmap(HeatmapConfig{Title: "Rhythm", Width: width, Days: days})

	for _, line := range strings.Split(chart, "\n") {
		if got := lipgloss.Width(line); got > width {
			t.Fatalf("heatmap line overflowed width %d (got %d):\n%s", width, got, chart)
		}
	}
}

func TestRenderSplitBarFillsTheWholeTrack(t *testing.T) {
	const width = 40

	chart := RenderSplitBar(SplitBarConfig{
		Title: "AI vs Human",
		Width: width,
		Segments: []SplitSegment{
			{Label: "AI", Percent: 7, Detail: "3h 00m"},
			{Label: "Human", Percent: 93, Detail: "40h 00m"},
		},
	})

	// The bar is the line right below the title; rounding leftovers must land
	// somewhere rather than leaving the track short.
	bar := strings.Split(chart, "\n")[2]
	if got, want := strings.Count(bar, "█"), width-chartChrome; got != want {
		t.Errorf("split bar filled %d cells, want %d:\n%s", got, want, chart)
	}
}

func TestChartsRenderEmptyStates(t *testing.T) {
	charts := map[string]string{
		"bar":    RenderBarChart(BarChartConfig{Title: "Bar", Width: 30, Empty: "nothing here"}),
		"column": RenderColumnChart(ColumnChartConfig{Title: "Column", Width: 30, Empty: "nothing here"}),
		"heat":   RenderHeatmap(HeatmapConfig{Title: "Heat", Width: 30, Empty: "nothing here"}),
		"split":  RenderSplitBar(SplitBarConfig{Title: "Split", Width: 30, Empty: "nothing here"}),
	}

	for name, chart := range charts {
		if !strings.Contains(chart, "nothing here") {
			t.Errorf("%s chart did not render its empty state:\n%s", name, chart)
		}
	}
}

// TestRenderDeckFitsTerminalWidth guards the deck geometry: Lipgloss counts a
// card's border and padding inside its width, so an off-by-two here silently
// wraps every card.
func TestRenderDeckFitsTerminalWidth(t *testing.T) {
	for _, width := range []int{60, 80, 100, 120, 160, 200} {
		app := AppModel{Width: width, Data: model.GetMockDashboardData()}

		for _, line := range strings.Split(app.renderDeck(), "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Errorf("deck line overflowed terminal width %d (got %d): %q", width, got, line)
				break
			}
		}
	}
}
