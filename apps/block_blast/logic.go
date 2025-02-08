package blockblast

import (
	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/rand"
)

func (m *Model) integrate(board boardType, color lipgloss.Color) {
	for y, r := range m.pieces[m.piece] {
		for x, c := range r {
			if c {
				board[y+m.position[0]][x+m.position[1]] = &color
			}
		}
	}
}

func (m *Model) addPlacePoints() {
	for _, r := range m.pieces[m.piece] {
		for _, c := range r {
			if c {
				m.points++
			}
		}
	}
}

func (m Model) valid(position [2]int, checkOverlap bool) bool {
	row := len(m.board) < len(m.pieces[m.piece])+position[0]
	col := len(m.board[0]) < len(m.pieces[m.piece][0])+position[1]

	if row || col {
		return false
	}

	if checkOverlap {
		for y, r := range m.pieces[m.piece] {
			for x, c := range r {
				if c {
					if m.board[y+m.position[0]][x+m.position[1]] != nil {
						return false
					}
				}
			}
		}
	}

	return true
}

func (m *Model) processLines(first bool) {
	for y, r := range m.board {
		complete := true
		for _, c := range r {
			if c == nil {
				complete = false
				break
			}
		}

		if complete {
			m.board[y] = make([]*lipgloss.Color, len(m.board[y]))
			m.points += 10
		}
	}

	if first {
		common.TransposeBoard(m.board)
		m.processLines(false)
	} else {
		common.TransposeBoard(m.board)
	}
}

func (m *Model) process(msg tea.KeyMsg) {
	position := m.position
	switch {
	case key.Matches(msg, m.keys.Up):
		position[1]--
	case key.Matches(msg, m.keys.Down):
		position[1]++
	case key.Matches(msg, m.keys.Left):
		position[0]--
	case key.Matches(msg, m.keys.Right):
		position[0]++
	case key.Matches(msg, m.keys.Next):
		m.piece++
		if m.piece >= len(m.pieces) {
			m.piece = 0
		}
	case key.Matches(msg, m.keys.Previous):
		m.piece--
		if m.piece < 1 {
			m.piece = len(m.pieces) - 1
		}
	case key.Matches(msg, m.keys.Place):
		if m.valid(position, true) {
			m.integrate(m.board, styles.CellColors[rand.Intn(len(styles.CellColors))])
			m.addPlacePoints()
			m.processLines(true)
			m.pieces = m.pieces[1:]
			m.piece = 0
			if len(m.pieces) < 1 {
				m.pieces = randomPieces()
			}
		}
	}

	if position != m.position {
		row := position[0] < len(m.board) && position[0] >= 0
		col := position[1] < len(m.board[0]) && position[1] >= 0
		if row && col {
			m.position = position
		}
	}
}
