package app

import (
	"github.com/three-thirds/hackatime-tui/internal/model"
)

// RenderProjectDurationsWidget renders a horizontal bar chart of coding time
// per project. Like every widget it is data-source agnostic, so an API adapter
// can replace the seed data later without touching this code.
func RenderProjectDurationsWidget(width int, projects []model.BreakdownItem) string {
	return RenderBarChart(BarChartConfig{
		Title:     "Project Durations",
		Width:     width,
		MinHeight: 8,
		Items:     projects,
		Format:    ValueDuration,
		Empty:     "No project activity yet",
	})
}
