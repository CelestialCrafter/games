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

type KeyMap struct {
	common.ArrowsKeyMap
	Save  key.Binding
	Help  key.Binding
	Quit  key.Binding
	Reset key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Save, k.Help, k.Quit, k.Reset}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Save, k.Help, k.Quit, k.Reset},
	}
}

type Model struct {
	keys     KeyMap
	help     help.Model
	board    [][]uint16
	finished bool
	height   int
	width    int
}

func NewModel() Model {
	board := createBoard(boardWidth, boardHeight)

	for i := range 2 {
		_ = i

		// this should never error
		board = addSquare(board)
	}

	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Save:         common.NewSaveBinding(),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
			Reset:        common.NewResetBinding(),
		},
		help:  help.New(),
		board: board,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
