package selector

import (
	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	common.ArrowsKeyMap
	Play key.Binding
	Help key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Play, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Play, k.Help, k.Quit},
	}
}

type PlayMsg struct {
	ID uint
}

type Model struct {
	keys         KeyMap
	help         help.Model
	selectedGame int
	rowLength    int
	width        int
	height       int
}

func NewModel() Model {
	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
			Play:         key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "play game")),
			Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		},
		help: help.New(),
		// initial value till tea.WindowSizeMsg gets emitted
		rowLength: 5,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
