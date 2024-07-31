package chess

import (
	"fmt"
	"math"
	"time"

	"github.com/CelestialCrafter/games/common"
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

// sorry for the mess (not sorry)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nextMoveMsg:
		if m.color == m.game.Position().Turn() {
			break
		}

		return m, func() tea.Msg {

			cmdPos := uci.CmdPosition{Position: m.game.Position()}
			cmdGo := uci.CmdGo{MoveTime: time.Second / 100}

			err := engine.Run(cmdPos, cmdGo)
			if err != nil {
				return common.ErrorWithBack(fmt.Errorf("engine errored: %v", err))
			}

			move := engine.SearchResults().BestMove
			err = m.game.Move(move)
			if err != nil {
				return common.ErrorWithBack(fmt.Errorf("couldn't execute move: %v", err))
			}

			return nextMoveMsg{}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Reset):
			newModel := NewModel()
			newModel.width = m.width
			newModel.height = m.height
			m = newModel
			return newModel, newModel.Init()
		case key.Matches(msg, m.keys.Select):
			position := m.game.Position()
			board := position.Board()
			selectedPiece := board.Piece(m.selectedSquare)

			pieceUnowned := selectedPiece.Color() != m.color
			pieceAlreadySelected := m.selectedPiece != nil && m.selectedPiece.square == m.selectedSquare

			if position.Turn() == m.color && m.selectedPiece != nil {
				move := m.findMove()
				if move != nil {
					m.selectedPiece = nil
					return m, func() tea.Msg {
						err := m.game.Move(move)
						// this should never happen
						if err != nil {
							return common.ErrorWithBack(fmt.Errorf("couldn't execute move: %v", err))
						}

						return nextMoveMsg{}
					}
				}

			}

			if pieceUnowned || pieceAlreadySelected || selectedPiece == chess.NoPiece {
				m.selectedPiece = nil
				break
			}

			m.selectedPiece = &pieceSquare{piece: selectedPiece, square: m.selectedSquare}
		// PLEASE make sure this is the last case
		case m.color == chess.White:
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
	}

	m.selectedSquare = chess.Square(math.Max(0, float64(m.selectedSquare)))
	m.selectedSquare = chess.Square(math.Min(boardSize*boardSize-1, float64(m.selectedSquare)))

	return m, nil
}
