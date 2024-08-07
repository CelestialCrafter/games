package chess

import (
	"fmt"
	"slices"
	"strings"

	"github.com/CelestialCrafter/games/common"
	"github.com/CelestialCrafter/games/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/notnil/chess"
)

const boardSize = 8

var cellStyle = lipgloss.NewStyle().Padding(1, 3)

func cellCheckerboard(style lipgloss.Style, i int) (changed bool) {
	row := i / boardSize
	col := i % boardSize

	if (row+col)%2 == 0 {
		// black
		style.Background(styles.Colors.Secondary)
	} else {
		// white
		style.Background(styles.Colors.Accent)
	}

	return true
}

func (m Model) piecePossibleMoves(style lipgloss.Style, currentSquare chess.Square) (changed bool) {
	position := m.game.Position()
	board := position.Board()
	if !m.ready {
		return
	}

	for _, move := range position.ValidMoves() {
		colorsMatch := board.Piece(move.S1()).Color() == m.color
		destIsCurrent := move.S2() == currentSquare
		srcIsSelectedSquare := move.S1() == m.selectedSquare
		srcIsSelectedPiece := m.selectedPiece != nil && move.S1() == m.selectedPiece.square

		if colorsMatch && destIsCurrent && (srcIsSelectedSquare || srcIsSelectedPiece) {
			style.Background(styles.CellColors[0])
			return true
		}
	}

	return
}

func pieceSelected(style lipgloss.Style, current chess.Square, selected *pieceSquare) (changed bool) {
	if selected != nil && selected.square == current {
		style.Background(styles.CellColors[1])
		changed = true
	}

	return
}

func cellSelected(style lipgloss.Style, current chess.Square, selected chess.Square) (changed bool) {
	if selected == current {
		style.Background(styles.CellColors[2])
		changed = true
	}

	return
}

func cellMargin(style lipgloss.Style, i int) (changed bool) {
	row := i / boardSize
	col := i % boardSize

	if row == 0 {
		style.MarginTop(3)
		changed = true
	}
	if col == boardSize-1 {
		style.MarginRight(6)
		changed = true
	}

	return
}

func pieceColor(style lipgloss.Style, piece chess.Piece) (bool, chess.Piece) {
	blackPiece := piece
	if blackPiece.Color() == chess.White {
		style.Foreground(lipgloss.Color("0"))
		blackPiece += 6
	} else {
		style.Foreground(lipgloss.Color("7"))
	}

	return true, blackPiece
}

// sorry for the mess (not sorry)
func (m Model) View() string {
	// status
	statusSlice := []string{fmt.Sprint("you are ", strings.ToLower(m.color.Name()))}

	outcome := m.game.Outcome()
	method := m.game.Method()
	if outcome != chess.NoOutcome {
		statusSlice = append(statusSlice, fmt.Sprint(outcome), fmt.Sprint(method))
	}

	var status string
	if m.ready {
		status = styles.StatusStyle.Render(strings.Join(statusSlice, " • "))
	} else {
		status = ""
	}

	// board rendering
	sliceChessBoard := common.CreateBoard[*pieceSquare](boardSize, boardSize)
	chessBoard := m.game.Position().Board()

	for r := boardSize - 1; r >= 0; r-- {
		for f := range boardSize {
			square := chess.NewSquare(chess.File(f), chess.Rank(r))
			piece := chessBoard.Piece(square)

			if piece == chess.NoPiece {
				sliceChessBoard[square.File()][square.Rank()] = nil
			}

			sliceChessBoard[square.File()][square.Rank()] = &pieceSquare{
				piece,
				square,
			}
		}
	}

	if m.color == chess.White {
		for i := range sliceChessBoard {
			slices.Reverse(sliceChessBoard[i])
		}
	} else {
		slices.Reverse(sliceChessBoard)
	}

	i := -1
	board := common.RenderBoard(sliceChessBoard, func(piece *pieceSquare) string {
		i++
		var changed bool
		newCellStyle := cellStyle.Copy()

		cellMargin(newCellStyle, i)
		if piece == nil {
			cellCheckerboard(newCellStyle, i)
			return newCellStyle.Render(" ")
		}

		changed = pieceSelected(newCellStyle, piece.square, m.selectedPiece)

		if !changed {
			changed = cellSelected(newCellStyle, piece.square, m.selectedSquare)
		}

		if !changed {
			changed = m.piecePossibleMoves(
				newCellStyle,
				piece.square,
			)
		}

		if !changed {
			cellCheckerboard(newCellStyle, i)
		}

		_, newPiece := pieceColor(newCellStyle, piece.piece)

		return newCellStyle.Render(newPiece.String())
	})

	// ranks & files
	ranks := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	files := []string{"8", "7", "6", "5", "4", "3", "2", "1"}

	if m.color == chess.Black {
		slices.Reverse(ranks)
		slices.Reverse(files)
	}

	for i, rank := range ranks {
		ranks[i] = cellStyle.Render(rank)
	}

	for i, file := range files {
		files[i] = cellStyle.Render(file)
	}

	// merging everything together
	file := lipgloss.JoinVertical(lipgloss.Left, files...)
	rank := lipgloss.JoinHorizontal(lipgloss.Left, ranks...)

	board = lipgloss.JoinHorizontal(lipgloss.Bottom, file, board)
	board = lipgloss.JoinVertical(lipgloss.Center, board, rank)

	board = lipgloss.NewStyle().
		BorderForeground(styles.Colors.Accent).
		Border(lipgloss.RoundedBorder()).
		Render(board)
	board = lipgloss.JoinVertical(
		lipgloss.Top,
		board,
		status,
	)
	board = lipgloss.Place(
		m.width,
		lipgloss.Height(board),
		lipgloss.Center,
		lipgloss.Top,
		board,
	)

	// rendering
	help := m.help.View(m.keys)
	var multiplayer string
	if m.multiplayer.Lobby != nil {
		multiplayer = m.multiplayer.View()
	} else {
		multiplayer = ""
	}

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)
	availableHeight -= lipgloss.Height(multiplayer)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(board),
		multiplayer,
		help,
	)
}
