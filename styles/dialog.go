package styles

import "github.com/charmbracelet/lipgloss"

var (
	Button = lipgloss.NewStyle().
		Background(Colors.Muted).
		Padding(0, 2).
		Margin(1)

	ButtonSelected = Button.Copy().Background(Colors.Primary)

	DialogBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Colors.Accent).
			Padding(1, 2)

	StatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(Colors.Muted)).
			Italic(true)
)
