package main

import (
	common "github.com/CelestialCrafter/games/common"
	tea "github.com/charmbracelet/bubbletea"
)

type EmptyModel struct {
}

func (m EmptyModel) Init() tea.Cmd {
	return func() tea.Msg {
		return common.BackMsg{}
	}
}

func (m EmptyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		return common.BackMsg{}
	}
}

func (m EmptyModel) View() string {
	return ""
}
