package saves

import (
	"fmt"

	common "github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/db"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listStyle = lipgloss.NewStyle().Margin(1, 2)

type SetupMsg struct {
	saves []Save
}

type Save struct {
	Id      string `db:"game_id"`
	OwnerId string `db:"owner_id"`
	GameId  uint   `db:"game"`
	Save    uint   `db:"save"`
	Data    string `db:"data"`
}

func (s Save) Title() string {
	return fmt.Sprintf("%v", common.Games[s.GameId].Name)
}

func (s Save) Description() string {
	return fmt.Sprintf("File %v", s.Save)
}

func (s Save) FilterValue() string { return s.Title() }

type KeyMap struct {
	common.ArrowsKeyMap
	Delete key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Delete, k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Delete, k.Help, k.Quit},
	}
}

type Model struct {
	keys   KeyMap
	help   help.Model
	list   list.Model
	userId string
}

func NewModel(userId string) Model {
	return Model{
		keys: KeyMap{
			ArrowsKeyMap: common.NewArrowsKeyMap(),
			Delete:       key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete save")),
			Help:         key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
			Quit:         key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "back")),
		},
		help:   help.New(),
		list:   list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		userId: userId,
	}
}

func (m Model) setup() tea.Msg {
	saves := []Save{}

	err := db.DB.Select(&saves, "SELECT game_id, game FROM games WHERE owner_id=$1", m.userId)
	// life if messages were commands 🤤
	if err != nil {
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
	switch msg := msg.(type) {
	case SetupMsg:
		m.list.Title = "Saves"

		items := make([]list.Item, len(msg.saves))
		for i, v := range msg.saves {
			items[i] = v
		}

		m.list.SetItems(items)

	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		m.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return listStyle.Render(m.list.View())
}
