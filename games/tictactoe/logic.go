package tictactoe

import (
	"errors"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func getBoardPosition(position uint) (uint, uint) {
	x := (position - 1) % 3
	y := 2 - (position-1)/3

	return x, y
}

// ./wincondition-explanation.png
func (m Model) checkGameState() int {
	// win condition 1
	for i := 0; i < len(m.board); i++ {
		initial := int(m.board[i][0])
		failed := false

		for j := 0; j < len(m.board[0]); j++ {
			current := int(m.board[i][j])
			if current == 0 {
				failed = true
				break
			}

			if initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}
	}

	// win condition 2
	// AAAAAAAAAAA JAVASCRIPT :mikudead::mikudead::mikudead:
	rv, ok := (func() (int, bool) {
		initial := int(m.board[0][0])
		failed := false

		for i := 0; i < len(m.board); i++ {
			current := int(m.board[i][i])
			if current == 0 || initial != current {
				failed = true
				break
			}
		}

		initial = int(m.board[0][len(m.board)-1])
		failed = false

		for i := len(m.board) - 1; i >= 0; i-- {
			current := int(m.board[len(m.board)-1-i][i])
			if current == 0 || initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial, true
		}

		return 0, false
	})()

	if ok {
		return rv
	}

	// win condition 3
	for i := 0; i < len(m.board[0]); i++ {
		initial := int(m.board[0][i])
		failed := false

		for j := 0; j < len(m.board); j++ {
			current := int(m.board[j][i])
			if current == 0 || initial != current {
				failed = true
				break
			}
		}

		if !failed {
			return initial
		}
	}

	return 0
}

func (m *Model) place(position uint) error {
	x, y := getBoardPosition(position)
	cell := &m.board[x][y]

	if *cell != 0 {
		m.err = errors.New("spot is taken")
		return m.err
	}

	*cell = m.turn

	m.winner = m.checkGameState()

	return nil
}

func (m *Model) process(msg tea.Msg) {
	var err error

	switch {
	// top to bottom
	case key.Matches(msg.(tea.KeyMsg), m.keys.One):
		err = m.place(1)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Two):
		err = m.place(2)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Three):
		err = m.place(3)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Four):
		err = m.place(4)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Five):
		err = m.place(5)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Six):
		err = m.place(6)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Seven):
		err = m.place(7)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Eight):
		err = m.place(8)
	case key.Matches(msg.(tea.KeyMsg), m.keys.Nine):
		err = m.place(9)
	}

	if err == nil {
		if m.turn == 2 {
			m.turn = 1
		} else {
			m.turn = 2
		}
	}
}
