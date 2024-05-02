package selector

import (
	common "github.com/CelestialCrafter/games/common"
	twenty48 "github.com/CelestialCrafter/games/games/2048"
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
	GameID uint
}

type Model struct {
	keys          KeyMap
	help          help.Model
	gamesMetadata []common.Metadata
	selectedGame  int
	rowLength     int
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
		gamesMetadata: []common.Metadata{
			twenty48.GetMetadata(),
		},
		rowLength: 5,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}