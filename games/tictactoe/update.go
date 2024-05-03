package tictactoe

import (
	"github.com/CelestialCrafter/games/common"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(
			msg,
			m.keys.One,
			m.keys.Two,
			m.keys.Three,
			m.keys.Four,
			m.keys.Five,
			m.keys.Six,
			m.keys.Seven,
			m.keys.Eight,
			m.keys.Nine,
		):
			m.process(msg)
		case key.Matches(msg, m.keys.Save):

		}
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	}
	return m, nil
}
