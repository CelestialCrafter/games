package chess

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/notnil/chess"
)

// this file doesnt have much logic due to it being handled by the chess library

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

func (m *Model) handleSelection() tea.Cmd {
	if !m.ready || m.game.Outcome() != chess.NoOutcome {
		return nil
	}

	position := m.game.Position()
	board := position.Board()
	selectedPiece := board.Piece(m.selectedSquare)

	pieceUnowned := selectedPiece.Color() != m.color
	pieceAlreadySelected := m.selectedPiece != nil && m.selectedPiece.square == m.selectedSquare

	// handle move
	if position.Turn() == m.color && m.selectedPiece != nil {
		move := m.handleMove()
		if move != nil {
			if m.multiplayer.Lobby != nil {
				m.multiplayer.Lobby.Broadcast(move)
				return nil
			}

			var moveCmds []tea.Cmd
			moveCmds = append(moveCmds, func() tea.Msg {
				return move
			})

			if m.multiplayer.Lobby == nil {
				moveCmds = append(moveCmds, handleEngineMove(m.game))
			}

			return tea.Sequence(moveCmds...)
		}
	}

	if pieceUnowned || pieceAlreadySelected || selectedPiece == chess.NoPiece {
		m.selectedPiece = nil
		return nil
	}

	m.selectedPiece = &pieceSquare{piece: selectedPiece, square: m.selectedSquare}
	return nil
}
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
