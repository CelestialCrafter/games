package snake

import (
	"time"

	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	boardWidth  = 25
	boardHeight = 20
)

const tickRate = time.Millisecond * 100
const snakeSpeed = 10

const (
	empty = iota
	snake
	apple
)

type boardType [][]uint8

type KeyMap struct {
	common.ArrowsKeyMap
	Save key.Binding
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
	Board    boardType
	Finished bool
	height   int
	width    int
	progress float64
	lastTick time.Time
	snake []Point
	direction Point
	inputBuffer []tea.KeyMsg
	score int
}

func NewModel() Model {
	board := common.CreateBoard[uint8](boardWidth, boardHeight)
	board = addApple(board)

	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Save:         common.NewSaveBinding(),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
		},
		help:  help.New(),
		Board: board,
		snake: []Point{{X: 10, Y: 10}},
		inputBuffer: make([]tea.KeyMsg, 0),
		lastTick: time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tick()
}
