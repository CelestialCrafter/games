package common

import "github.com/charmbracelet/bubbles/key"

type BackMsg struct{}

type ArrowsKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
}

func NewArrowsKeyMap() ArrowsKeyMap {
	return ArrowsKeyMap{
		Up:    key.NewBinding(key.WithKeys("k", "up", "w"), key.WithHelp("↑/k/w", "move up")),
		Down:  key.NewBinding(key.WithKeys("j", "down", "s"), key.WithHelp("↑/j/s", "move down")),
		Left:  key.NewBinding(key.WithKeys("h", "left", "a"), key.WithHelp("↑/h/a", "move left")),
		Right: key.NewBinding(key.WithKeys("l", "right", "d"), key.WithHelp("↑/l/d", "move right")),
	}
}

func NewBackBinding() key.Binding {
	return key.NewBinding(key.WithKeys("esc", "q"), key.WithHelp("esc/q", "back"))
}

func NewHelpBinding() key.Binding {
	return key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help"))
}

func NewResetBinding() key.Binding {
	return key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset"))
}
