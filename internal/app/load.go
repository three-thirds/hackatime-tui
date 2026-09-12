package app

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/three-thirds/hackatime-tui/internal/model"
	"github.com/three-thirds/hackatime-tui/pkg/api"
)

type dataLoadedMsg struct {
	data model.DashboardData
	err  error
}

func loadDashboardCmd(client *api.Client, dateRange string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return dataLoadedMsg{err: fmt.Errorf("missing API credentials (check ~/.wakatime.cfg)")}
		}

		today, todayErr := client.GetStatusToday()
		stats, statsErr := client.GetStats(statsRangeForFilter(dateRange))
		if statsErr != nil {
			if todayErr != nil {
				return dataLoadedMsg{err: fmt.Errorf("stats: %v; today: %v", statsErr, todayErr)}
			}
			return dataLoadedMsg{err: statsErr}
		}

		var todayPtr *api.StatusToday
		if todayErr == nil {
			todayPtr = today
		}

		data := dashboardFromAPI(todayPtr, stats)
		data.Filters.DateRange = dateRange
		return dataLoadedMsg{data: data}
	}
}
