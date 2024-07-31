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

func cellCheckerboard(style lipgloss.Style, i int) {
	row := i / boardSize
	col := i % boardSize

	if (row+col)%2 == 0 {
		// black
		style.Background(styles.Colors.Secondary)
	} else {
		// white
		style.Background(styles.Colors.Accent)
	}
}

func (m Model) piecePossibleMoves(style lipgloss.Style, currentSquare chess.Square) {
	position := m.game.Position()
	board := position.Board()
	for _, move := range position.ValidMoves() {
		colorsMatch := board.Piece(move.S1()).Color() == m.color
		destIsCurrent := move.S2() == currentSquare
		srcIsSelectedSquare := move.S1() == m.selectedSquare
		srcIsSelectedPiece := m.selectedPiece != nil && move.S1() == m.selectedPiece.square

		if colorsMatch && destIsCurrent && (srcIsSelectedSquare || srcIsSelectedPiece) {
			style.Background(styles.CellColors[0])
		}
	}
}

func pieceSelected(style lipgloss.Style, current chess.Square, selected *pieceSquare) {
	if selected != nil && selected.square == current {
		style.Background(styles.CellColors[1])
	}
}

func cellSelected(style lipgloss.Style, current chess.Square, selected chess.Square) {
	if selected == current {
		style.Background(styles.CellColors[2])
	}
}

func cellMarign(style lipgloss.Style, i int) {
	row := i / boardSize
	col := i % boardSize

	if row == 0 {
		style.MarginTop(3)
	}
	if col == boardSize-1 {
		style.MarginRight(6)
	}
}

func pieceColor(style lipgloss.Style, piece chess.Piece) chess.Piece {
	blackPiece := piece
	if blackPiece.Color() == chess.White {
		style.Foreground(lipgloss.Color("0"))
		blackPiece += 6
	} else {
		style.Foreground(lipgloss.Color("7"))
	}

	return blackPiece
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
		newCellStyle := cellStyle.Copy()

		cellCheckerboard(newCellStyle, i)
		if piece == nil {
			return newCellStyle.Render(" ")
		}

		m.piecePossibleMoves(
			newCellStyle,
			piece.square,
		)
		pieceSelected(newCellStyle, piece.square, m.selectedPiece)
		cellSelected(newCellStyle, piece.square, m.selectedSquare)
		cellMarign(newCellStyle, i)
		newPiece := pieceColor(newCellStyle, piece.piece)

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
		styles.StatusStyle.Render(strings.Join(statusSlice, " • ")),
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

	availableHeight := m.height
	availableHeight -= lipgloss.Height(help)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(availableHeight).Render(board),
		help,
	)
}
