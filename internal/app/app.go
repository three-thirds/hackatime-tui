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
	Width          int            // Width of terminal screen
	Height         int            // Height of terminal screen
	CurrentTab     int            // Currently selected tab on TUI
	Viewport       viewport.Model // State of viewport
	Ready          bool           // Whether the app is ready to be rendered or not
	ActiveZone     ActiveZone     // Currently Focused interative area
	ActiveFilter   int            // Currently focused filter button
	FilterValues   [5]string      // Commited selection for each filter
	DropdownOpen   bool           // Whether the options dropdown is open
	DropdownCursor int            // Selected index inside the open dropdown
}

type ActiveZone int

const (
	FocusNav ActiveZone = iota
	FocusFilter
	FocusDeck
)

func NewApp() AppModel {
	return AppModel{
		CurrentTab:   0,
		ActiveZone:   FocusNav,
		ActiveFilter: 0,
		FilterValues: [5]string{
			"Last 7 Days",
			"All",
			"All",
			"All",
			"All",
		},
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

		headerContent := renderFixedHeader(
			m.Width,
			m.CurrentTab,
			m.ActiveZone,
			m.ActiveFilter,
			m.FilterValues,
			m.DropdownOpen,
			m.DropdownCursor,
		)

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
		if key.Matches(msg, Keys.Quit) {
			return m, tea.Quit
		}
		if key.Matches(msg, Keys.TabNext) {
			m.ActiveZone = (m.ActiveZone + 1) % 3
			return m, nil
		}
		switch m.ActiveZone {
		case FocusNav:
			return m.handleNavKeys(msg)
		case FocusFilter:
			return m.handleFilterKeys(msg)
		case FocusDeck:
			return m.handleDeckKeys(msg)
		}
	}
	return m, cmd
}

func (m AppModel) handleNavKeys(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, Keys.Down):
		m.CurrentTab = (m.CurrentTab + 1) % 3
		m.Viewport.SetContent(m.renderDeck())
		return m, nil

	case key.Matches(msg, Keys.Up):
		if m.CurrentTab == 0 {
			m.CurrentTab = 2
		} else {
			m.CurrentTab--
		}
		m.Viewport.SetContent(m.renderDeck())
		return m, nil
	}

	return m, nil
}

func (m AppModel) handleFilterKeys(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	options := getFilterOptions(m.ActiveFilter)

	if m.DropdownOpen {
		switch {
		case key.Matches(msg, Keys.Down):
			m.DropdownCursor = (m.DropdownCursor + 1) % len(options)
		case key.Matches(msg, Keys.Up):
			if m.DropdownCursor <= 0 {
				m.DropdownCursor = len(options) - 1
			} else {
				m.DropdownCursor--
			}
		case key.Matches(msg, Keys.Select):
			m.FilterValues[m.ActiveFilter] = options[m.DropdownCursor]
			m.DropdownOpen = false

		case key.Matches(msg, Keys.Cancel):
			m.DropdownOpen = false
		}
		return m, nil
	}
	switch {
	case key.Matches(msg, Keys.Right):
		m.ActiveFilter = (m.ActiveFilter + 1) % 5

	case key.Matches(msg, Keys.Left):
		if m.ActiveFilter <= 0 {
			m.ActiveFilter = 4
		} else {
			m.ActiveFilter--
		}

	case key.Matches(msg, Keys.Select):
		m.DropdownOpen = true
		m.DropdownCursor = 0
	}
	return m, nil
}

func (m AppModel) handleDeckKeys(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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

	header := renderFixedHeader(
		m.Width,
		m.CurrentTab,
		m.ActiveZone,
		m.ActiveFilter,
		m.FilterValues,
		m.DropdownOpen,
		m.DropdownCursor,
	)

	deck := m.Viewport.View()

	var modeBadge string
	switch m.ActiveZone {
	case FocusNav:
		modeBadge = ActiveTabStyle.Render("[NAV]") + DimText.Render(" (j/k switch tab • Tab next zone)")
	case FocusFilter:
		modeBadge = ActiveTabStyle.Render("[FILTERS]") + DimText.Render(" (h/l select • Enter open • Tab next zone)")
	case FocusDeck:
		modeBadge = ActiveTabStyle.Render("[DECK]") + DimText.Render(" (j/k scroll • Tab next zone)")
	}

	scrollPercent := int(m.Viewport.ScrollPercent() * 100)
	rightInfo := DimText.Render(fmt.Sprintf("Scroll: %d%% │ [q] Quit ", scrollPercent))

	gap := max(m.Width-lipgloss.Width(modeBadge)-lipgloss.Width(rightInfo), 1)
	footer := lipgloss.JoinHorizontal(lipgloss.Top, modeBadge, strings.Repeat(" ", gap), rightInfo)

	fullUI := lipgloss.JoinVertical(lipgloss.Left, header, deck, footer)

	view := tea.NewView(fullUI)
	view.AltScreen = true

	return view
}

func getFilterOptions(filterIdx int) []string {
	switch filterIdx {
	case 0:
		return []string{"Last 7 Days", "Last 30 Days", "Today", "All Time"}
	case 1:
		return []string{"All", "hackatime-tui", "kasumi", "skora-backend"}
	case 2:
		return []string{"All", "Rust", "Python", "Go", "Svelte"}
	case 3:
		return []string{"All", "Linux", "Mac", "Windows"}
	case 4:
		return []string{"All", "Neovim", "VSCode", "Zed"}
	default:
		return []string{"All"}
	}
}

// renderDeck builds the scrollable wireframe slot grid, calculating equal
// half-width and full-width card dimensions to match the terminal bounds.
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
