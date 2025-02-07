package snake

import (
	"math"
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/key"
)

type Point struct {
	X, Y int
}

func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Add(other Point) Point {
	return Point{
		p.X + other.X,
		p.Y + other.Y,
	}
}

func (m *Model) nextDirection() Point {
	if len(m.inputBuffer) == 0 {
		return m.direction
	}

	lastInput := m.inputBuffer[len(m.inputBuffer)-1]

	var newDirection Point
	switch {
	case key.Matches(lastInput, m.keys.Up):
		newDirection = Point{X: 0, Y: -1}
	case key.Matches(lastInput, m.keys.Down):
		newDirection = Point{X: 0, Y: 1}
	case key.Matches(lastInput, m.keys.Left):
		newDirection = Point{X: -1, Y: 0}
	case key.Matches(lastInput, m.keys.Right):
		newDirection = Point{X: 1, Y: 0}
	}

	// if new dir == opposite of current dir, ignore the input
	if newDirection.Equals(Point{X: m.direction.X * -1, Y: m.direction.Y * -1}) {
		return m.direction
	}

	return newDirection
}

func getEmpty(board boardType) []*uint8 {
	empty := make([]*uint8, 0)

	for x := 0; x < len(board); x++ {
		for y := 0; y < len(board[0]); y++ {
			c := &board[x][y]
			if *c == 0 {
				empty = append(empty, c)
			}
		}
	}

	return empty
}

func addApple(board boardType) boardType {
	empty := getEmpty(board)

	if len(empty) < 1 {
		return board
	}

	*empty[rand.Intn(len(empty))] = apple

	return board
}

func (m *Model) detectCollision(newHead Point) (collided bool) {
	// wall
	if newHead.X < 0 || newHead.X >= boardWidth || newHead.Y < 0 || newHead.Y >= boardHeight {
		return true
	}

	// self
	if len(m.inputBuffer) < 1 {
		return
	}
	for _, segment := range m.snake {
		if segment.Equals(newHead) {
			return true
		}
	}

	return
}

func (m *Model) process(lastTick time.Time) {
	m.direction = m.nextDirection()

	deltaTime := time.Since(lastTick)
	speedMulti := 1.0 + math.Log(float64(len(m.snake)))*0.2
	m.progress += deltaTime.Seconds() * baseSnakeSpeed * speedMulti

	for m.progress >= 1 {
		m.progress--

		newHead := m.snake[0].Add(m.direction)

		// collisions
		m.Finished = m.detectCollision(newHead)
		if m.Finished {
			return
		}

		// apples
		var ateApple bool
		if m.Board[newHead.X][newHead.Y] == apple {
			m.Board[newHead.X][newHead.Y] = empty
			ateApple = true
			m.Board = addApple(m.Board)
			m.progress++
			m.score++
		}

		m.snake = append([]Point{newHead}, m.snake...)
		if !ateApple {
			m.snake = m.snake[:len(m.snake)-1]
		}
	}
}
