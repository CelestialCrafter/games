package chess

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"time"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/CelestialCrafter/games/saveManager"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

type gameSave struct {
	FEN   string
	Color chess.Color
}

func (m *Model) handleMove() tea.Msg {
	move := m.findMove()
	if move != nil {
		m.selectedPiece = nil
		return moveMsg(move)
	}

	return nil
}

func handleEngineMove(game *chess.Game) func() tea.Msg {
	return func() tea.Msg {
		cmdPos := uci.CmdPosition{Position: game.Position()}
		cmdGo := uci.CmdGo{MoveTime: time.Second / 100}

		err := engine.Run(cmdPos, cmdGo)
		if err != nil {
			return common.ErrorWithBack(fmt.Errorf("engine errored: %v", err))
		}

		move := engine.SearchResults().BestMove
		return moveMsg(move)
	}
}

// sorry for the mess (not sorry)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case moveMsg:
		_ = m.game.Move(msg)
	case multiplayer.DisconnectMsg:
		data, _ := m.multiplayer.Lobby.Data.(*lobbyData)
		loser := data.colors[string(msg)]

		m.game.Resign(loser)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Select):
			cmds = append(cmds, m.handleSelection())
		case key.Matches(msg, m.keys.Save):
			cmds = append(cmds, func() tea.Msg {
				bytes := bytes.Buffer{}
				encoder := gob.NewEncoder(&bytes)
				err := encoder.Encode(gameSave{
					FEN:   m.game.FEN(),
					Color: m.color,
				})
				if err != nil {
					return common.ErrorMsg{
						Err: err,
					}
				}

				return saveManager.SaveMsg{
					Data: bytes.Bytes(),
					ID:   common.Chess.ID,
				}
			})

		// make sure this is the last case
		case m.color == chess.White:
			if m.game.Outcome() != chess.NoOutcome {
				break
			}
			switch {

			case key.Matches(msg, m.keys.Up):
				m.upRank()
			case key.Matches(msg, m.keys.Down):
				m.downRank()
			case key.Matches(msg, m.keys.Left):
				m.backwardsFile()
			case key.Matches(msg, m.keys.Right):
				m.forwardFile()
			}
		case m.color == chess.Black:
			if m.game.Outcome() != chess.NoOutcome {
				break
			}
			switch {
			case key.Matches(msg, m.keys.Up):
				m.downRank()
			case key.Matches(msg, m.keys.Down):
				m.upRank()
			case key.Matches(msg, m.keys.Left):
				m.forwardFile()
			case key.Matches(msg, m.keys.Right):
				m.backwardsFile()
			}
		}

	case saveManager.LoadMsg:
		// disable loading if multiplayer is on
		if m.multiplayer.Lobby != nil {
			break
		}
		var saveData gameSave

		bytes := bytes.Buffer{}
		bytes.Write(msg.Data)

		decoder := gob.NewDecoder(&bytes)
		err := decoder.Decode(&saveData)
		if err != nil {
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}

		fen, err := chess.FEN(saveData.FEN)
		if err != nil {
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}

		m.game = chess.NewGame(fen)
		m.color = saveData.Color

	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	case multiplayer.InitialReadyMsg:
		data, _ := m.multiplayer.Lobby.Data.(*lobbyData)

		m.color = data.colors[m.multiplayer.Self.ID]

		if m.color == chess.White {
			m.selectedSquare = chess.A1
		} else {
			m.selectedSquare = chess.H8
		}

		m.ready = true
	}

	m.selectedSquare = chess.Square(math.Max(0, float64(m.selectedSquare)))
	m.selectedSquare = chess.Square(math.Min(boardSize*boardSize-1, float64(m.selectedSquare)))

	if m.multiplayer.Lobby != nil {
		var multiplayerCmd tea.Cmd
		m.multiplayer, multiplayerCmd = m.multiplayer.Update(msg)

		cmds = append(cmds, multiplayerCmd)

	}

	return m, tea.Batch(cmds...)
}
