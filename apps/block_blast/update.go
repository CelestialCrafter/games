package blockblast

import (
	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/exp/rand"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	position := m.position

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

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

				m.pieces = append(m.pieces[:m.piece], m.pieces[m.piece+1:]...)
				m.piece = 0
				if len(m.pieces) < 1 {
					m.pieces = randomPieces()
				}
			}
		}
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}

	if position != m.position {
		row := position[0] < len(m.board) && position[0] >= 0
		col := position[1] < len(m.board[0]) && position[1] >= 0
		if row && col {
			m.position = position
		}
	}

	return m, nil
}
