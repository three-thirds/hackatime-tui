// Package app manages the root state machine, viewport layout,
// and keyboard navigation for the dashboard interface
package app

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/three-thirds/hackatime-tui/internal/model"
)

// AppModel holds important central state of entire applications
// it includes dimensions of terminal, currently active Tab, state of viewport, etc.
type AppModel struct {
	Width      int                 // Width of terminal screen
	Height     int                 // Height of terminal screen
	CurrentTab int                 // Currently selected tab on TUI
	Viewport   viewport.Model      // State of viewport
	Ready      bool                // Whether the app is ready to be rendered or not
	Data       model.DashboardData // Local seed data used by widgets during UI development
}

func NewApp() AppModel {
	return AppModel{
		CurrentTab: 0,
		Data:       model.GetMockDashboardData(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

// Update handles incoming Bubble Tea messages, including terminal window resizing,
// tab navigation keypresses, and viewport scrolling events.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

		headerContent := renderFixedHeader(m.Width, m.CurrentTab)
		headerHeight := lipgloss.Height(headerContent)

		footerHeight := 1
		viewportHeight := max(m.Height-headerHeight-footerHeight, 5)

		if !m.Ready {
			m.Viewport = viewport.New(viewport.WithWidth(m.Width), viewport.WithHeight(viewportHeight))
			m.Ready = true
		} else {
			m.Viewport.SetWidth(m.Width)
			m.Viewport.SetHeight(viewportHeight)
		}

		m.Viewport.SetContent(m.renderDeck())
		return m, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, Keys.TabNext):
			m.CurrentTab = (m.CurrentTab + 1) % 3
			m.Viewport.SetContent(m.renderDeck())
			return m, nil

		case key.Matches(msg, Keys.TabPrev):
			if m.CurrentTab == 0 {
				m.CurrentTab = 2
			} else {
				m.CurrentTab--
			}
			m.Viewport.SetContent(m.renderDeck())
			return m, nil

		case key.Matches(msg, Keys.Up), key.Matches(msg, Keys.Down),
			key.Matches(msg, Keys.PageUp), key.Matches(msg, Keys.PageDown),
			key.Matches(msg, Keys.Top), key.Matches(msg, Keys.Bottom):

			m.Viewport, cmd = m.Viewport.Update(msg)
			return m, cmd
		}
	}
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

// View renders the complete terminal user interface, assembling the pinned
// top header, the scrollable widget deck, and the status footer.
func (m AppModel) View() tea.View {
	if !m.Ready {
		v := tea.NewView("Initializing Wireframe...")
		v.AltScreen = true
		return v
	}

	header := renderFixedHeader(m.Width, m.CurrentTab)
	deck := m.Viewport.View()

	scrollPercent := int(m.Viewport.ScrollPercent() * 100)
	footer := DimText.Render(lipgloss.PlaceHorizontal(m.Width, lipgloss.Right,
		fmt.Sprintf("Scroll: %d%% │ [j/k/↑/↓] Scroll │ [h/l/Tab] Nav │ [q] Quit ", scrollPercent)))

	fullUI := lipgloss.JoinVertical(lipgloss.Left, header, deck, footer)

	view := tea.NewView(fullUI)
	view.AltScreen = true

	return view
}

// renderDeck builds the scrollable widget grid, sizing every card off the
// shared deck geometry so columns line up and gaps stay even.
func (m AppModel) renderDeck() string {
	halfW := halfCardWidth(m.Width)
	fullW := fullCardWidth(m.Width)
	gap := strings.Repeat(" ", columnGap)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		RenderProjectDurationsWidget(halfW, m.Data.Projects), gap,
		RenderLanguagesWidget(halfW, m.Data.Languages))

	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		RenderEditorsWidget(halfW, m.Data.Editors), gap,
		RenderOperatingSystemsWidget(halfW, m.Data.OSList))

	row3 := RenderTimelineWidget(fullW, 8, m.Data.Timeline)
	row4 := RenderCodingRhythmWidget(fullW, m.Data.Heatmap)

	row5 := lipgloss.JoinHorizontal(lipgloss.Top,
		RenderGoalWidget(halfW, m.Data.GoalPercentage, m.Data.GoalStatusText, m.Data.GoalDetailText), gap,
		RenderAIvsHumanWidget(halfW, m.Data.AIPercentage, m.Data.AITime, m.Data.HumanPercentage, m.Data.HumanTime))

	return joinRows(row1, row2, row3, row4, row5)
}

// joinRows stacks deck rows with rowGap blank lines between them so the cards
// are not visually glued together.
func joinRows(rows ...string) string {
	spaced := make([]string, 0, len(rows)*2)
	for i, row := range rows {
		if i > 0 {
			spaced = append(spaced, strings.Repeat("\n", rowGap-1))
		}
		spaced = append(spaced, row)
	}
	return lipgloss.JoinVertical(lipgloss.Left, spaced...)
}
