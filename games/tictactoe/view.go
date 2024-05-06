package tictactoe

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	board := ""
	status := ""

	cellStyle := lipgloss.NewStyle().
		Padding(1, 0).
		Width(7).
		Align(lipgloss.Center)

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		Italic(true)

	xStyle := lipgloss.NewStyle().
		Inherit(statusStyle).
		Foreground(lipgloss.Color("4"))
	yStyle := lipgloss.NewStyle().
		Inherit(statusStyle).
		Foreground(lipgloss.Color("2"))
	errorStyle := lipgloss.NewStyle().
		Inherit(statusStyle).
		Foreground(lipgloss.Color("1"))

	var turn string
	if m.turn == 1 {
		turn = xStyle.Render("x")
	} else {
		turn = yStyle.Render("y")
	}

	_ = turn
	status = fmt.Sprintf("%v%v", turn, statusStyle.Render("'s turn"))

	if m.err != nil {
		status = errorStyle.Render(m.err.Error())
	}

	if m.winner != 0 {
		var winner string
		if m.winner == 1 {
			winner = xStyle.Render("x")
		} else {
			winner = yStyle.Render("y")
		}
		status = fmt.Sprintf("%v %v", winner, statusStyle.Render("wins!"))
	}

	var boardRows []string
	for y := range m.board[0] {
		var row []string
		for x := range m.board {
			cell := m.board[x][y]
			var color lipgloss.Color
			if cell == 0 {
				color = lipgloss.Color("0")
			} else if cell == 1 {
				color = lipgloss.Color("4")
			} else {
				color = lipgloss.Color("2")
			}

			var cellString string
			if cell == 0 {
				cellString = "-"
			} else if cell == 1 {
				cellString = "x"
			} else {
				cellString = "o"
			}

			row = append(row, cellStyle.Background(color).Render(cellString))
		}

		boardRows = append(boardRows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
		board = lipgloss.JoinVertical(lipgloss.Left, boardRows...)
	}

	return fmt.Sprintf("%v\n\n%v\n%v", board, status, m.help.View(m.keys))
}
