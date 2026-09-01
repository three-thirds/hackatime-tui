package app

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type AppModel struct {
	Width      int
	Height     int
	CurrentTab int
	Viewport   viewport.Model
	Ready      bool
}

func NewApp() AppModel {
	return AppModel{
		CurrentTab: 0,
	}
}

func (m AppModel) Init() tea.Cmd {
	return nil
}

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

func (m AppModel) renderDeck() string {
	fullW := m.Width - 2
	halfW := fullW / 2

	r1Left := RenderSlot("Project Durations", halfW, 9, "Horizontal bar chart of Projects")
	r1Right := RenderSlot("Languages", halfW, 9, "Language % breakdown bar chart")
	row1 := lipgloss.JoinHorizontal(lipgloss.Top, r1Left, r1Right)

	r2Left := RenderSlot("Editors", halfW, 8, "Editor % breakdown progress bars")
	r2Right := RenderSlot("Operating Systems", halfW, 8, "OS % breakdown progress bars")
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, r2Left, r2Right)

	row3 := RenderSlot("Project Timeline (Stacked Weekly Activity)", fullW, 10, "Weekly activity bar chart")
	row4 := RenderSlot("Coding Rhythm (Activity Heatmap)", fullW, 11, "GitHub-style / Rhythm activity heatmap")

	r5Left := RenderSlot("Today's Goal", halfW, 8, "Goal progress ring / meter")
	r5Right := RenderSlot("AI vs Human Coding", halfW, 8, "AI vs Human ratio split bar")
	row5 := lipgloss.JoinHorizontal(lipgloss.Top, r5Left, r5Right)

	return lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3, row4, row5)
}
