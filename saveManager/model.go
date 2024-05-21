package saveManager

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/CelestialCrafter/games/apps/saves"
	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/db"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type TryLoad struct {
	ID uint
}

type SaveMsg struct {
	Data []byte
	ID   uint
}

type LoadMsg struct {
	Data []byte
}

type Model struct {
	userId string
}

func NewModel(userId string) Model {
	return Model{
		userId,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SaveMsg:
		saveFile := 0
		h := sha1.New()
		io.WriteString(h, fmt.Sprintf("%v-%v-%v", m.userId, msg.ID, saveFile))

		_, err := db.DB.Exec(`
				INSERT INTO games(game_id,owner_id,game,data,save,last_save_time) VALUES($1,$2,$3,$4,$5,$6)
					ON CONFLICT(game_id) DO UPDATE SET data=$4;`,
			hex.EncodeToString(h.Sum(nil)),
			m.userId,
			msg.ID,
			string(msg.Data),
			saveFile,
			time.Now(),
		)

		if err != nil {
			log.Error("couldnt save to database", "error", err)
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err:    fmt.Errorf("couldnt save: %v", err),
					Action: nil,
				}
			}
		}
	case TryLoad:
		// @TODO multiple save files
		save := []saves.Save{}
		saveFile := 0

		err := db.DB.Select(&save, "SELECT data FROM games WHERE owner_id=$1 AND game=$2 AND save=$3;", m.userId, msg.ID, saveFile)
		if err != nil {
			log.Error("couldnt load save from database", "error", err)
			return m, func() tea.Msg {
				return common.ErrorMsg{
					Err:    fmt.Errorf("couldnt load save: %v", err),
					Action: nil,
				}
			}
		}

		if len(save) > 0 {
			return m, func() tea.Msg {
				return LoadMsg{
					Data: []byte(save[0].Data),
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	return ""
}
