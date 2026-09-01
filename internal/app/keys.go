// This file contains all keybinds required for TUI
package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
	Top      key.Binding
	Bottom   key.Binding
	TabNext  key.Binding
	TabPrev  key.Binding
	Quit     key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "scroll down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("ctrl+u", "u", "pgup"),
		key.WithHelp("ctrl+u/u/pgup", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("ctrl+d", "d", "pgdown"),
		key.WithHelp("ctrl+d/d/pgdown", "page down"),
	),
	Top: key.NewBinding(
		key.WithKeys("g", "home"),
		key.WithHelp("g/home", "top"),
	),
	Bottom: key.NewBinding(
		key.WithKeys("G", "end"),
		key.WithHelp("G/end", "bottom"),
	),
	TabNext: key.NewBinding(
		key.WithKeys("tab", "l", "right"),
		key.WithHelp("tab/l", "next tab"),
	),
	TabPrev: key.NewBinding(
		key.WithKeys("shift+tab", "h", "left"),
		key.WithHelp("shift+tab/h", "prev tab"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
