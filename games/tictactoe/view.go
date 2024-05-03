package tictactoe

func (m Model) View() string {
	return m.help.View(m.keys)
}
