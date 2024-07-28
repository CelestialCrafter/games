package tictactoe

import "testing"

const player = 1

func testWin(t *testing.T, board [][]uint8) {

	model := Model{
		turn:  player,
		board: board,
	}

	winner := model.checkGameState()
	if winner != player {
		t.Errorf("checkGameState(%v) = %d; want %d", model.board, winner, player)
	}
}

func TestWins(t *testing.T) {
	// horizontal
	testWin(t, [][]uint8{
		{player, 0, 0},
		{player, 0, 0},
		{player, 0, 0},
	})

	// vertical
	testWin(t, [][]uint8{
		{player, player, player},
		{0, 0, 0},
		{0, 0, 0},
	})

	// diagonal
	testWin(t, [][]uint8{
		{player, 0, 0},
		{0, player, 0},
		{0, 0, player},
	})
}
