package common

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg struct {
	Err        error
	Action     tea.Cmd
	ActionText string
	Fatal      bool
}
