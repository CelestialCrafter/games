package twenty48

import (
	"reflect"
	"testing"
)

func TestLose(t *testing.T) {
	// no adjacent cells
	board := boardType{
		{2, 4, 6, 8},
		{8, 6, 4, 2},
		{2, 4, 6, 8},
		{8, 6, 4, 2},
	}

	lost := checkLost(board)
	if !lost {
		t.Errorf("checkLost(%v) = %t; want true", board, lost)
	}

	// 1 adjacent cell at [0][0] and [0][1]
	board = boardType{
		{2, 4, 6, 8},
		{2, 6, 4, 2},
		{2, 4, 6, 8},
		{8, 6, 4, 2},
	}

	lost = checkLost(board)
	if lost {
		t.Errorf("checkLost(%v) = %t; want false", board, lost)
	}
}

func testDirection(t *testing.T, display string, board boardType, want boardType, movement func(boardType)) {
	movement(board)
	if !reflect.DeepEqual(board, want) {
		t.Errorf("%s = %v; want %v", display, board, want)
	}
}

func TestDirections(t *testing.T) {
	// right
	testDirection(t, "right()", boardType{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, boardType{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{2, 0, 0, 0},
	}, right)

	// left
	testDirection(t, "left()", boardType{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{2, 0, 0, 0},
	}, boardType{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, left)

	// up
	testDirection(t, "up()", boardType{
		{0, 0, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, boardType{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, up)

	// down
	testDirection(t, "down()", boardType{
		{2, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, boardType{
		{0, 0, 0, 2},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}, down)
}
