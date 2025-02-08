package twenty48

import (
	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	boardWidth  = 4
	boardHeight = 4
)

type boardType [][]uint16

type KeyMap struct {
	common.ArrowsKeyMap
	Help key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

type Model struct {
	keys     KeyMap
	help     help.Model
	board    boardType
	points   int
	finished bool
	height   int
	width    int
}

func NewModel() Model {
	board := common.CreateBoard[uint16](boardWidth, boardHeight)

	for range 2 {
		// this should never fail
		board = addSquare(board)
	}

	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
		},
		help:  help.New(),
		board: board,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
