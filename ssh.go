package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/multiplayer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

func createTeaHandler(sess ssh.Session) *tea.Program {
	key := sess.PublicKey()
	if key == nil {
		log.Error("Key was nil (enable PublicKeyAuth middleware?)")
		_, err := sess.Write([]byte("You need to ssh in with a public key!"))
		if err != nil {
			log.Error("Could not write error message to session", "error", err)
		}

		err = sess.Close()
		if err != nil {
			log.Error("Could not close session due to nil key", "error", err)
		}
	}

	keyType := key.Type()
	keyData := hex.EncodeToString(key.Marshal())
	id := fmt.Sprintf("%v-%v", keyType, keyData)

	m := NewModel(id)
	renderer := bubbletea.MakeRenderer(sess)
	lipgloss.SetDefaultRenderer(renderer)
	program := tea.NewProgram(m, append(bubbletea.MakeOptions(sess), programOpts...)...)

	multiplayer.Players.Store(id, &multiplayer.Player{
		Program: program,
		ID:      id,
	})

	// wait for disconnect and handle it
	go func() {
		program.Wait()
		multiplayer.Players.Compute(id, func(player *multiplayer.Player, loaded bool) (_ *multiplayer.Player, delete bool) {
			multiplayer.Cleanup(player.Lobby, id)
			return nil, true
		})
	}()

	return program
}

func startSSH() {
	addr := net.JoinHostPort(host, port)

	s, err := wish.NewServer(
		wish.WithAddress(addr),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(
				createTeaHandler,
				termenv.Profile(*common.ColorProfile),
			),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Fatal("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatal("Could not stop server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Fatal("Could not stop server", "error", err)
	}
}
