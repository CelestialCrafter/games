package common

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg struct {
	Err        error
	Action     tea.Cmd
	ActionText string
	Fatal      bool
}

func ErrorWithBack(err error) tea.Msg {
	return ErrorMsg{
		Err: err,
		Action: func() tea.Msg {
			return BackMsg{}
		},
		ActionText: "Back",
	}
}
