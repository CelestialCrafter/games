package common

import tea "github.com/charmbracelet/bubbletea"

// @TODO offload this to save model

type SaveMsg struct {
	Data []byte
}

type LoadMsg struct {
	Data []byte
}

type ErrorMsg struct {
	Err        error
	Action     tea.Cmd
	ActionText string
	Fatal      bool
}
