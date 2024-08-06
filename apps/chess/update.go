package chess

import (
	"fmt"
	"math"
	"time"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func (m *Model) upRank() {
	if m.selectedSquare.Rank() == chess.Rank8 {
		return
	}
	m.selectedSquare += 8
}

func (m *Model) downRank() {
	if m.selectedSquare.Rank() == chess.Rank1 {
		return
	}
	m.selectedSquare -= 8
}

func (m *Model) forwardFile() {
	if m.selectedSquare.File() == chess.FileH {
		return
	}
	m.selectedSquare++
}

func (m *Model) backwardsFile() {
	if m.selectedSquare.File() == chess.FileA {
		return
	}
	m.selectedSquare--
}

func (m Model) findMove() *chess.Move {
	position := m.game.Position()

	validMoves := position.ValidMoves()
	var selectedMove *chess.Move

	for _, move := range validMoves {
		if move.S1() != m.selectedPiece.square || move.S2() != m.selectedSquare {
			continue
		}

		selectedMove = move
		break
	}

	return selectedMove
}

func (m Model) handleMove() tea.Msg {
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
			if !m.ready || m.game.Outcome() != chess.NoOutcome {
				break
			}

			position := m.game.Position()
			board := position.Board()
			selectedPiece := board.Piece(m.selectedSquare)

			pieceUnowned := selectedPiece.Color() != m.color
			pieceAlreadySelected := m.selectedPiece != nil && m.selectedPiece.square == m.selectedSquare

			if position.Turn() == m.color && m.selectedPiece != nil {
				newMoveMsg := m.handleMove()
				if newMoveMsg != nil {
					m.multiplayer.Lobby.Broadcast(newMoveMsg)
				}
			}

			if pieceUnowned || pieceAlreadySelected || selectedPiece == chess.NoPiece {
				m.selectedPiece = nil
				break
			}

			m.selectedPiece = &pieceSquare{piece: selectedPiece, square: m.selectedSquare}
		// PLEASE make sure this is the last case
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

	var mutliplayerCmd tea.Cmd
	m.multiplayer, mutliplayerCmd = m.multiplayer.Update(msg)

	return m, mutliplayerCmd
}
