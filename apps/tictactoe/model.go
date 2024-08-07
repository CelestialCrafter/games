package tictactoe

import (
	"math/rand"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puzpuzpuz/xsync/v3"
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
		},
		{
			k.Six,
			k.Seven,
			k.Eight,
			k.Nine,
		},
		{k.Save, k.Help, k.Quit},
	}
}

type Model struct {
	keys        KeyMap
	help        help.Model
	multiplayer multiplayer.Model
	turn        uint8
	player      uint8
	board       [][]uint8
	ready       bool
	winner      uint8
	height      int
	width       int
}

type lobbyData struct {
	colors       map[string]uint8
	startingTurn uint8
}

func randomPlayer() uint8 {
	return uint8(rand.Intn(2)) + 1
}

func NewModel() Model {
	board := make([][]uint8, 3)
	for i := range board {
		board[i] = make([]uint8, 3)
	}

	m := Model{
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
	}

	if *common.EnableMultiplayer {
		m.multiplayer = multiplayer.NewModel(
			2,
			common.TicTacToe.ID,
			func(players *xsync.MapOf[string, *multiplayer.Player]) interface{} {
				var nextPlayer uint8 = 1
				colors := map[string]uint8{}

				players.Range(func(id string, _ *multiplayer.Player) bool {
					colors[id] = nextPlayer
					nextPlayer++

					return true
				})

				return &lobbyData{
					colors:       colors,
					startingTurn: randomPlayer(),
				}
			},
		)
	} else {
		m.player = randomPlayer()
		m.turn = randomPlayer()
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return m.multiplayer.Init()
}
