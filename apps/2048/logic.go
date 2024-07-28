package twenty48

import (
	"math/rand"
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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

func reverse(matrix [][]uint16) {
	for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
		matrix[i], matrix[j] = matrix[j], matrix[i]
	}
}

func transpose(matrix [][]uint16) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < i; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}
}

func rotate90(matrix [][]uint16) {
	transpose(matrix)
	reverse(matrix)
}

func rotateN90(matrix [][]uint16) {
	reverse(matrix)
	transpose(matrix)
}

func createBoard(w int, h int) [][]uint16 {
	board := make([][]uint16, w)
	for i := range board {
		board[i] = make([]uint16, h)
	}

	return board
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
	board := createBoard(len(originalBoard), len(originalBoard[0]))
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

func merge(board [][]uint16) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board)-1; j++ {
			current := &board[i][j]
			next := &board[i][j+1]
			if *current == *next && *current != 0 {
				*current = *current * 2
				*next = 0
			}
		}
	}
}

func up(board boardType) {
	push(board)
	merge(board)
	push(board)
}

func right(board boardType) {
	rotateN90(board)
	up(board)
	rotate90(board)
}

func left(board boardType) {
	rotate90(board)
	up(board)
	rotateN90(board)
}

func down(board boardType) {
	for range 2 {
		rotate90(board)
	}
	up(board)
	for range 2 {
		rotateN90(board)
	}
}

func (m *Model) process(msg tea.Msg) {
	if m.finished {
		return
	}

	before := createBoard(boardWidth, boardHeight)
	copy(before, m.board)

	switch {
	case key.Matches(msg.(tea.KeyMsg), m.keys.Up):
		up(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Down):
		down(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Left):
		left(m.board)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Right):
		right(m.board)
	}

	if !reflect.DeepEqual(before, m.board) {
		addSquare(m.board)
		m.finished = checkLost(m.board)
	}

	log.Info("ima tach u", ":drool:", m.board)
}
