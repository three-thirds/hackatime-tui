package app

import (
	"strings"
	"testing"
	"time"

	"github.com/three-thirds/hackatime-tui/internal/model"
)

func TestRenderProjectDurationsWidgetRendersSeedEntries(t *testing.T) {
	widget := RenderProjectDurationsWidget(60, []model.BreakdownItem{
		{Name: "atlas", Duration: 2*time.Hour + 15*time.Minute, Percentage: 75},
		{Name: "notes", Duration: 45 * time.Minute, Percentage: 25},
	})

	for _, want := range []string{"Project Durations", "atlas", "2h 15m", "notes", "45m", "█"} {
		if !strings.Contains(widget, want) {
			t.Errorf("widget did not contain %q:\n%s", want, widget)
		}
	}
}

func TestRenderProjectDurationsWidgetRendersEmptyState(t *testing.T) {
	widget := RenderProjectDurationsWidget(60, nil)

	if !strings.Contains(widget, "No project activity yet") {
		t.Errorf("widget did not render an empty state:\n%s", widget)
	}
}
