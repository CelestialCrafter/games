package twenty48

import (
	"bytes"
	"encoding/gob"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/save"
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
		case key.Matches(msg, m.keys.Up, m.keys.Down, m.keys.Left, m.keys.Right):
			m.process(msg)
		case key.Matches(msg, m.keys.Save):
			return m, func() tea.Msg {
				gameSave := gameSave{
					Board: m.board,
				}

				bytes := bytes.Buffer{}
				encoder := gob.NewEncoder(&bytes)
				err := encoder.Encode(gameSave)

				if err != nil {
					return common.ErrorMsg{
						Err: err,
					}
				}

				return save.SaveMsg{
					Data: bytes.Bytes(),
					ID:   GetMetadata().ID,
				}
			}
		}
	case save.LoadMsg:
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
	}

	return m, nil
}