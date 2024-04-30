package styles

import "github.com/charmbracelet/lipgloss"

// @TODO have standard color definition & use 256 color ansi

var (
	Button = lipgloss.NewStyle().
		Background(lipgloss.Color("8")).
		Padding(0, 2).
		Margin(1)

	ButtonSelected = Button.Copy().Background(lipgloss.Color("4"))

	DialogBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("2")).
			Padding(1, 2)
)
