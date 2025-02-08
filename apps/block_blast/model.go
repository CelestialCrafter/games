package blockblast

import (
	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/rand"
)

const boardSize = 12
const pieceBufferSize = 3

type boardType = [][]*lipgloss.Color
type piece = [][]bool

var (
	Square1 piece = piece{
		{true, true, true},
		{true, true, true},
		{true, true, true},
	}

	Square2 piece = piece{
		{true, true},
		{true, true},
	}

	Rectangle1 piece = piece{
		{true, true, true},
		{true, true, true},
	}

	Rectangle2 piece = piece{
		{true, true},
		{true, true},
		{true, true},
	}

	Line1 piece = piece{
		{true},
		{true},
		{true},
	}

	Line2 piece = piece{
		{true, true, true},
	}

	L1 piece = piece{
		{true, false},
		{true, false},
		{true, true},
	}

	L2 piece = piece{
		{false, true},
		{false, true},
		{true, true},
	}

	L3 piece = piece{
		{true, true},
		{true, false},
		{true, false},
	}

	L4 piece = piece{
		{true, true},
		{false, true},
		{false, true},
	}

	L5 piece = piece{
		{true, false, false},
		{true, true, true},
	}

	L6 piece = piece{
		{false, false, true},
		{true, true, true},
	}

	L7 piece = piece{
		{true, true, true},
		{true, false, false},
	}

	L8 piece = piece{
		{true, true, true},
		{false, false, true},
	}

	T1 piece = piece{
		{true, true, true},
		{false, true, false},
	}

	T2 piece = piece{
		{false, true, false},
		{true, true, true},
	}

	T3 piece = piece{
		{false, true},
		{true, true},
		{false, true},
	}

	T4 piece = piece{
		{true, false},
		{true, true},
		{true, false},
	}

	Z1 piece = piece{
		{true, true, false},
		{false, true, true},
	}

	Z2 piece = piece{
		{false, true, true},
		{true, true, false},
	}

	Z3 piece = piece{
		{false, true},
		{true, true},
		{true, false},
	}

	Z4 piece = piece{
		{true, false},
		{true, true},
		{false, true},
	}

	Pieces = []piece{
		Square1,
		Square2,
		Rectangle1,
		Rectangle2,
		Line1,
		Line2,
		L1,
		L2,
		L3,
		L4,
		L5,
		L6,
		L7,
		L8,
		T1,
		T2,
		T3,
		T4,
		Z1,
		Z2,
		Z3,
		Z4,
	}
)

func randomPieces() []piece {
	pieces := make([]piece, pieceBufferSize)
	for i := range pieceBufferSize {
		pieces[i] = Pieces[rand.Intn(len(Pieces))]
	}

	return pieces
}

type KeyMap struct {
	common.ArrowsKeyMap
	Next     key.Binding
	Previous key.Binding
	Place    key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Place, k.Next, k.Previous, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Place, k.Next, k.Previous, k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

type Model struct {
	keys     KeyMap
	help     help.Model
	board    [][]*lipgloss.Color
	preview  [][]*lipgloss.Color
	points   int
	piece    int
	pieces   []piece
	position [2]int
	width    int
	height   int
}

func NewModel() Model {
	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Place:        key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter/space", "place piece")),
			Next:         key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "next piece")),
			Previous:     key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "prev piece")),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
		},
		help:    help.New(),
		board:   common.CreateBoard[*lipgloss.Color](boardSize, boardSize),
		preview: common.CreateBoard[*lipgloss.Color](boardSize, boardSize),
		pieces:  randomPieces(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
