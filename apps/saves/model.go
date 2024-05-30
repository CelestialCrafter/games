package saves

import (
	"fmt"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/db"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)
var deleteBinding = key.NewBinding(key.WithKeys("delete", "d"), key.WithHelp("delete/d", "delete save"))

type SetupMsg struct {
	saves []Save
}

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
	listModel := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	listModel.DisableQuitKeybindings()
	modifiedKeyMap := list.DefaultKeyMap()
	modifiedKeyMap.Quit = common.NewBackBinding()
	listModel.KeyMap = modifiedKeyMap
	listModel.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			deleteBinding,
		}
	}

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
		return func() tea.Msg {
			return common.ErrorMsg{
				Err: err,
				Action: func() tea.Msg {
					return common.BackMsg{}
				},
				ActionText: "Back",
			}
		}
	}

	return SetupMsg{
		saves,
	}
}

func (m Model) Init() tea.Cmd {
	return m.setup
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SetupMsg:
		m.list.Title = "Saves"

		items := make([]list.Item, len(msg.saves))
		for i, v := range msg.saves {
			items[i] = v
		}

		m.list.SetItems(items)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, deleteBinding):
			item := m.list.Items()[m.list.Index()]
			save := item.(Save)

			_, err := db.DB.Exec("DELETE FROM saves WHERE save_id=$1wawa;", save.Id)
			if err != nil {
				log.Error("couldn't delete save from database", "error", err)
				cmd = func() tea.Msg {
					return common.ErrorMsg{
						Err: err,
					}
				}
			}
		case key.Matches(msg, m.list.KeyMap.Quit):
			return m, func() tea.Msg {
				return common.BackMsg{}
			}
		}
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var listCmd tea.Cmd
	m.list, listCmd = m.list.Update(msg)
	log.Warn("cmd", "cmd", cmd)
	return m, tea.Batch(cmd, listCmd)
}

func (m Model) View() string {
	return m.list.View()
}
