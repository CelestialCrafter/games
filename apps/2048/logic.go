package twenty48

import (
	"math/rand"

	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// @TODO score

func getEmpty(board [][]uint16) []*uint16 {
	empty := make([]*uint16, 0)

	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[0]); y++ {
			c := &board[x][y]
			if *c == 0 {
				empty = append(empty, c)
			}
		}
	}

	return empty
}

func addSquare(board [][]uint16) [][]uint16 {
	empty := getEmpty(board)

	if len(empty) < 1 {
		return board
	}

	*empty[rand.Intn(len(empty))] = uint16((rand.Intn(2) + 1) * 2)

	return board
}

func rotate90(matrix [][]uint16) {
	common.TransposeBoard(matrix)
	common.ReverseBoard(matrix)
}

func rotateN90(matrix [][]uint16) {
	common.ReverseBoard(matrix)
	common.TransposeBoard(matrix)
}

func findDownPair(board [][]uint16) bool {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0])-1; j++ {
			if board[i][j] == board[i][j+1] {
				return true
			}
		}
	}

	return false
}

func checkLost(originalBoard [][]uint16) bool {
	board := common.CreateBoard[uint16](
		len(originalBoard),
		len(originalBoard[0]),
	)
	copy(board, originalBoard)

	if len(getEmpty(board)) > 0 {
		return false
	}

	if findDownPair(board) {
		return false
	}

	rotate90(board)
	return !findDownPair(board)
}

func push(board [][]uint16) {
	for i := 0; i < len(board); i++ {
		position := 0
		for j := 0; j < len(board[0]); j++ {
			current := &board[i][j]
			next := &board[i][position]

			if *current != 0 {
				tmp := *next
				*next = *current
				*current = tmp
				position++
			}
		}
	}
}

func merge(board [][]uint16) (points int) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board)-1; j++ {
			current := &board[i][j]
			next := &board[i][j+1]
			if *current == *next && *current != 0 {
				*current = *current * 2
				points += int(*current)
				*next = 0
			}
		}
	}

	return
}

func up(board boardType) int {
	push(board)
	p := merge(board)
	push(board)
	return p
}

func right(board boardType) int {
	rotateN90(board)
	p := up(board)
	rotate90(board)
	return p
}

func left(board boardType) int {
	rotate90(board)
	p := up(board)
	rotateN90(board)
	return p
}

func down(board boardType) int {
	for range 2 {
		rotate90(board)
	}
	p := up(board)
	for range 2 {
		rotateN90(board)
	}
	return p
}

func (m *Model) process(msg tea.Msg) {
	if m.finished {
		return
	}

	before := common.CreateBoard[uint16](boardWidth, boardHeight)
	for i, row := range m.board {
		copy(before[i], row)
	}

	var p int
	switch {
	case key.Matches(msg.(tea.KeyMsg), m.keys.Up):
		p = up(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Down):
		p = down(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Left):
		p = left(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Right):
		p = right(m.board)
	}
	m.points += p

	if !common.CompareBoards(before, m.board) {
		addSquare(m.board)
		m.finished = checkLost(m.board)
	}
}
