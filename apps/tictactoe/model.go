package tictactoe

import (
	"github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	One   key.Binding
	Two   key.Binding
	Three key.Binding
	Four  key.Binding
	Five  key.Binding
	Six   key.Binding
	Seven key.Binding
	Eight key.Binding
	Nine  key.Binding

	Save key.Binding
	Help key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.One,
			k.Two,
			k.Three,
			k.Four,
			k.Five,
			k.Six,
			k.Seven,
			k.Eight,
			k.Nine,
		},
		{k.Save, k.Help, k.Quit},
	}
}

type Model struct {
	keys   KeyMap
	help   help.Model
	turn   uint8
	board  [][]uint8
	winner int
}

func NewModel() Model {
	board := make([][]uint8, 3)
	for i := range board {
		board[i] = make([]uint8, 3)
	}

	return Model{
		keys: KeyMap{
			One:   key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "one")),
			Two:   key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "two")),
			Three: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "three")),
			Four:  key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "four")),
			Five:  key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "five")),
			Six:   key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "six")),
			Seven: key.NewBinding(key.WithKeys("7"), key.WithHelp("7", "seven")),
			Eight: key.NewBinding(key.WithKeys("8"), key.WithHelp("8", "eight")),
			Nine:  key.NewBinding(key.WithKeys("9"), key.WithHelp("9", "nine")),
			Save:  key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "save")),
			Help:  common.NewHelpBinding(),
			Quit:  common.NewBackBinding(),
		},
		help:  help.New(),
		board: board,
		turn:  1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
