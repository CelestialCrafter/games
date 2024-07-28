package saves

import (
	"errors"
	"fmt"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/db"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var deleteBinding = key.NewBinding(key.WithKeys("delete", "x"), key.WithHelp("delete/x", "delete save"))

type SavesMsg []Save

type Save struct {
	Id      string `db:"save_id"`
	OwnerId string `db:"owner_id"`
	GameId  uint   `db:"game_id"`
	File    uint   `db:"file"`
	Data    string `db:"data"`
}

func (s Save) Title() string {
	return fmt.Sprintf("%v", common.Games[s.GameId].Name)
}

func (s Save) Description() string {
	return fmt.Sprintf("File %v", s.File)
}

func (s Save) FilterValue() string { return s.Title() }

type Model struct {
	list   list.Model
	userId string
}

func NewModel(userId string) Model {
	// @TODO list styles
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	modifiedKeyMap := list.DefaultKeyMap()
	modifiedKeyMap.Quit = common.NewBackBinding()
	listModel.KeyMap = modifiedKeyMap
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			deleteBinding,
		}
	}
	listModel.Title = ""
	listModel.Styles.Title = lipgloss.NewStyle()
	listModel.SetItems([]list.Item{})

	return Model{
		list:   listModel,
		userId: userId,
	}
}

func (m Model) setup() tea.Msg {
	saves := []Save{}

	err := db.DB.Select(&saves, "SELECT save_id, game_id FROM saves WHERE owner_id=$1", m.userId)
	// life if messages were commands 🤤
	if err != nil {
		log.Error("couldn't load saves from database", "error", err)
		return common.ErrorWithBack(fmt.Errorf("could not delete save: %v", err))
	}

	return SavesMsg(saves)
}

func (m Model) Init() tea.Cmd {
	return m.setup
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SavesMsg:
		if len(msg) < 1 {
			return m, func() tea.Msg {
				return common.ErrorWithBack(errors.New("no save files exist"))
			}
		}

		items := make([]list.Item, len(msg))
		for i, v := range msg {
			items[i] = v
		}

		m.list.SetItems(items)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		}
		switch {
		case key.Matches(msg, deleteBinding):
			save := m.list.SelectedItem().(Save)
			index := m.list.Index()

			cmds = append(
				cmds,
				func() tea.Msg {
					// @FIX if this errors the item does not get inserted back into the list
					_, err := db.DB.Exec("DELETE FROM saves WHERE save_id=$1;", save.Id)
					m.list.RemoveItem(index)

					if err != nil {
						log.Error("couldn't delete save from database", "error", err)

						return common.ErrorWithBack(fmt.Errorf("couldn't delete save: %v", err))
					}

					if len(m.list.Items()) < 1 {
						return common.ErrorWithBack(errors.New("no save files exist"))
					}

					return nil
				})
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	cmds = append(cmds, listCmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}
