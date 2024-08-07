package chess

import (
	"fmt"
	"math/rand"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"github.com/puzpuzpuz/xsync/v3"
)

var engine *uci.Engine

type pieceSquare struct {
	piece  chess.Piece
	square chess.Square
}

type moveMsg *chess.Move

type KeyMap struct {
	common.ArrowsKeyMap
	Save   key.Binding
	Help   key.Binding
	Quit   key.Binding
	Select key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Save,
		k.Help,
		k.Quit,
		k.Select,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Select, k.Up, k.Down, k.Left, k.Right},
		{k.Save, k.Help, k.Quit},
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
	multiplayer    multiplayer.Model
	ready          bool
}

type lobbyData struct {
	colors map[string]chess.Color
}

func NewModel() Model {
	m := Model{
		help: help.New(),
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Save:         common.NewSaveBinding(),
			Help:         common.NewHelpBinding(),
			Quit:         common.NewBackBinding(),
			Select:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select piece")),
		},
		game:           chess.NewGame(),
		selectedSquare: chess.A1,
	}

	if *common.EnableMultiplayer {
		m.multiplayer = multiplayer.NewModel(
			2,
			common.Chess.ID,
			func(players *xsync.MapOf[string, *multiplayer.Player]) interface{} {
				nextPlayer := 1
				colors := map[string]chess.Color{}

				players.Range(func(id string, _ *multiplayer.Player) bool {
					colors[id] = chess.Color(nextPlayer)
					nextPlayer++

					return true
				})

				return &lobbyData{
					colors,
				}
			},
		)
	} else {
		m.color = chess.Color(rand.Intn(2) + 1)
		if m.color == chess.White {
			m.selectedSquare = chess.A1
		} else {
			m.selectedSquare = chess.H8
		}
		m.ready = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if m.multiplayer.Lobby != nil {
		return m.multiplayer.Init()
	}

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

		if m.color == chess.Black {
			return handleEngineMove(m.game)()
		}

		return nil
	}
}
