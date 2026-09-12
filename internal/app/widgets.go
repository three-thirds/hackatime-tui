// This file binds dashboard data to the chart primitives in charts.go.
//
// Every widget here is a thin, stateless function over `model` values: it takes
// the slice or scalars it needs plus a width, and returns a rendered card. When
// the API lands, only the values passed in change, never these functions.

package app

import (
	"fmt"
	"image/color"
	"time"

	"github.com/three-thirds/hackatime-tui/internal/model"
)

// RenderLanguagesWidget renders the language share breakdown.
func RenderLanguagesWidget(width int, languages []model.BreakdownItem) string {
	return RenderBarChart(BarChartConfig{
		Title:     "Languages",
		Width:     width,
		MinHeight: 8,
		Items:     languages,
		Format:    ValuePercent,
		Empty:     "No language activity yet",
	})
}

// RenderEditorsWidget renders the editor share breakdown.
func RenderEditorsWidget(width int, editors []model.BreakdownItem) string {
	return RenderBarChart(BarChartConfig{
		Title:     "Editors",
		Width:     width,
		MinHeight: 7,
		Items:     editors,
		Palette:   []color.Color{DraculaPink, DraculaPurple, DraculaCyan, DraculaOrange},
		Format:    ValuePercent,
		Empty:     "No editor activity yet",
	})
}

// RenderOperatingSystemsWidget renders the OS share breakdown.
func RenderOperatingSystemsWidget(width int, systems []model.BreakdownItem) string {
	return RenderBarChart(BarChartConfig{
		Title:     "Operating Systems",
		Width:     width,
		MinHeight: 7,
		Items:     systems,
		Palette:   []color.Color{DraculaGreen, DraculaYellow, DraculaOrange},
		Format:    ValuePercent,
		Empty:     "No OS activity yet",
	})
}

// RenderTimelineWidget renders weekly coding activity as vertical columns.
func RenderTimelineWidget(width, plotHeight int, points []model.TimelinePoint) string {
	return RenderColumnChart(ColumnChartConfig{
		Title:      "Project Timeline (Weekly Activity)",
		Width:      width,
		PlotHeight: plotHeight,
		Points:     points,
		Color:      DraculaPurple,
		Empty:      "No timeline activity yet",
	})
}

// RenderCodingRhythmWidget renders the day-by-day activity heatmap.
func RenderCodingRhythmWidget(width int, days []model.HeatmapDay) string {
	return RenderHeatmap(HeatmapConfig{
		Title: "Coding Rhythm (Activity Heatmap)",
		Width: width,
		Days:  days,
		Empty: "No activity recorded yet",
	})
}

// RenderGoalWidget renders progress toward today's coding goal.
func RenderGoalWidget(width, percent int, status, detail string) string {
	fillColor := DraculaGreen
	switch {
	case percent < 40:
		fillColor = DraculaRed
	case percent < 80:
		fillColor = DraculaOrange
	}

	return RenderMeter(MeterConfig{
		Title:     "Today's Goal",
		Width:     width,
		Percent:   percent,
		Status:    status,
		Detail:    detail,
		Color:     fillColor,
		MinHeight: 6,
	})
}

// RenderAIvsHumanWidget renders the split between AI-assisted and hand-written
// coding time.
func RenderAIvsHumanWidget(width int, aiPercent int, aiTime time.Duration, humanPercent int, humanTime time.Duration) string {
	segments := []SplitSegment{}
	if aiPercent > 0 || humanPercent > 0 || aiTime > 0 || humanTime > 0 {
		segments = []SplitSegment{
			{Label: "AI", Percent: aiPercent, Detail: formatDuration(aiTime), Color: DraculaPink},
			{Label: "Human", Percent: humanPercent, Detail: formatDuration(humanTime), Color: DraculaCyan},
		}
	}

	return RenderSplitBar(SplitBarConfig{
		Title:     "AI vs Human Coding",
		Width:     width,
		MinHeight: 6,
		Segments:  segments,
		Empty:     "No coding time recorded yet",
	})
}

// formatDuration prints a duration the way the dashboard reads it: minutes only
// under an hour, hours and zero-padded minutes above.
func formatDuration(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}
	return fmt.Sprintf("%dh %02dm", int(duration.Hours()), int(duration.Minutes())%60)
}
