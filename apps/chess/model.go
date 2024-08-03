package chess

import (
	"fmt"
	"math/rand"

	"github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

var engine *uci.Engine

type pieceSquare struct {
	piece  chess.Piece
	square chess.Square
}

type nextMoveMsg struct{}

type KeyMap struct {
	common.ArrowsKeyMap
	Save   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Select key.Binding
	Resign key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Save,
		k.Help,
		k.Quit,
		k.Select,
		k.Resign,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Save, k.Help, k.Quit},
		{k.Select, k.Resign},
	}
}

// @TODO drawing
// @TODO mutliplayer

type Model struct {
	game           *chess.Game
	help           help.Model
	keys           KeyMap
	selectedSquare chess.Square
	selectedPiece  *pieceSquare
	color          chess.Color
	height         int
	width          int
}

func NewModel() Model {
	// @TODO support multiplayer
	color := chess.Color(rand.Intn(2) + 1)
	var selectedSquare chess.Square
	if color == chess.White {
		selectedSquare = chess.A1
	} else {
		selectedSquare = chess.H8
	}

	return Model{
		help: help.New(),
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Save:         common.NewSaveBinding(),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
			Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select piece")),
			// dont ask me why it's p
			Resign: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "resign game")),
		},
		selectedSquare: selectedSquare,
		color:          color,
		game:           chess.NewGame(),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		if engine == nil {
			var err error
			engine, err = uci.New("stockfish")
			if err != nil {
				return common.ErrorWithBack(fmt.Errorf("couldn't create stockfish engine: %v", err))
			}
		}

		err := engine.Run(uci.CmdUCINewGame, uci.CmdIsReady, uci.CmdUCINewGame)
		if err != nil {
			return common.ErrorWithBack(fmt.Errorf("couldn't initialize stockfish engine: %v", err))
		}

		return nextMoveMsg{}
	}
}
