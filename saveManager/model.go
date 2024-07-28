package saveManager

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

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

func saveError(err error) func() tea.Msg {
	return func() tea.Msg {

		return common.ErrorMsg{
			Err:    fmt.Errorf("couldnt save: %v", err),
			Action: nil,
		}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SaveMsg:
		file := 0
		h := sha1.New()
		_, err := io.WriteString(h, fmt.Sprintf("%v-%v-%v", m.userId, msg.ID, file))
		if err != nil {
			log.Error("couldn't write save id to hash")
			return m, saveError(err)
		}

		_, err = db.DB.Exec(`
				INSERT INTO saves(save_id,owner_id,game_id,data,file) VALUES($1,$2,$3,$4,$5)
					ON CONFLICT(save_id) DO UPDATE SET data=$4;`,
			hex.EncodeToString(h.Sum(nil)),
			m.userId,
			msg.ID,
			string(msg.Data),
			file,
		)

		if err != nil {
			log.Error("couldnt save to database", "error", err)
			return m, saveError(err)
		}
	case TryLoad:
		// @TODO multiple save files
		save := []saves.Save{}
		saveFile := 0

		err := db.DB.Select(&save, "SELECT data FROM saves WHERE owner_id=$1 AND game_id=$2 AND file=$3;", m.userId, msg.ID, saveFile)
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
