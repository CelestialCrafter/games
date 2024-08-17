package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/CelestialCrafter/games/common"
	mainmodel "github.com/CelestialCrafter/games/mainModel"
	"github.com/CelestialCrafter/games/multiplayer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"golang.org/x/crypto/ssh"
)

func startProgram() {
	lipgloss.DefaultRenderer().SetColorProfile(termenv.Profile(*common.ColorProfile))

	keyBytes, err := os.ReadFile(".ssh/id_ed25519.pub")
	if err != nil {
		log.Fatal("Could not read ed25519 key file", "error", err)
	}

	key, _, _, _, err := ssh.ParseAuthorizedKey(keyBytes)
	if err != nil {
		log.Fatal("Could not parse ed25519 key", "error", err)
	}
	keyType := key.Type()
	keyData := hex.EncodeToString(key.Marshal())
	id := fmt.Sprintf("%v-%v", keyType, keyData)

	m := mainmodel.NewModel(id)
	program := tea.NewProgram(m, programOpts...)
	multiplayer.Players.Store(id, &multiplayer.Player{
		Program: program,
		ID:      id,
	})

	// logging messes with the TUI, so we have to write logs to a file when not running as a ssh server
	f, err := tea.LogToFileWith("program.log", "debug", log.Default())
	if err != nil {
		log.Fatal("could not open log file", "error", err)
	}
	defer f.Close()

	_, err = program.Run()
	if err != nil {
		log.Fatal("Could not start program", "error", err)
	}
}
