package main

import (
	"flag"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"

	"github.com/CelestialCrafter/games/common"
)

const host = "0.0.0.0"

var programOpts = []tea.ProgramOption{tea.WithAltScreen()}

func main() {
	flag.Parse()

	if *common.EnableSsh {
		startSSH()
		return
	}

	startProgram()
}
