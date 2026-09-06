// Package app manages the root state machine, viewport layout,
// and keyboard navigation for the dashboard interface
package app

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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

	scrollPercent := int(m.Viewport.ScrollPercent() * 100)
	footer := DimText.Render(lipgloss.PlaceHorizontal(m.Width, lipgloss.Right,
		fmt.Sprintf("Scroll: %d%% │ [j/k/↑/↓] Scroll │ [h/l/Tab] Nav │ [q] Quit ", scrollPercent)))

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
	halfW := m.Width / 2
	fullW := halfW * 2

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
