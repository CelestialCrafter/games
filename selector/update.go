package selector

import (
	common "github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Play):
			return m, func() tea.Msg {
				return PlayMsg{
					ID: common.Games[uint(m.selectedGame)].ID,
				}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Up):
			m.selectedGame -= m.rowLength
		case key.Matches(msg, m.keys.Down):
			m.selectedGame += m.rowLength
		case key.Matches(msg, m.keys.Right):
			m.selectedGame++
		case key.Matches(msg, m.keys.Left):
			m.selectedGame--
		}

		m.selectedGame = min(max(m.selectedGame, 0), len(common.Games)-1)
	case tea.WindowSizeMsg:
		// -1 is to account for margin
		m.rowLength = msg.Width/common.ICON_WIDTH - 1
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}
