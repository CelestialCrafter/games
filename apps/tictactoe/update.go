package tictactoe

import (
	"bytes"
	"encoding/gob"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/saveManager"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type gameSave struct {
	Board [][]uint8
}

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
		case key.Matches(msg, m.keys.Reset):
			newModel := NewModel()
			newModel.width = m.width
			newModel.height = m.height
			m = newModel
		case key.Matches(msg, m.keys.Save):
			return m, func() tea.Msg {
				bytes := bytes.Buffer{}
				encoder := gob.NewEncoder(&bytes)
				err := encoder.Encode(gameSave{
					Board: m.board,
				})

				if err != nil {
					return common.ErrorMsg{
						Err: err,
					}
				}

				return saveManager.SaveMsg{
					Data: bytes.Bytes(),
					ID:   common.TicTacToe.ID,
				}
			}
		}
	case saveManager.LoadMsg:
		saveData := gameSave{}

		bytes := bytes.Buffer{}
		bytes.Write(msg.Data)
		decoder := gob.NewDecoder(&bytes)
		err := decoder.Decode(&saveData)

		if err != nil {
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err: err,
				}
			}
		}

		m.board = saveData.Board
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}
