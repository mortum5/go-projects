package game

import (
	"math/rand"

	"github.com/mortum5/go-projects/game-of-life/board"
)

const (
	DEAD = iota
	ALIVE
)

type State byte

type NextMoveFunc func(int, int) (int, int)

func GenerateMoves() (moves []NextMoveFunc) {
	dir := []int{-1, 0, 1}
	for i := 0; i < len(dir); i++ {
		for j := 0; j < len(dir); j++ {
			i := i
			j := j
			if dir[i] == dir[j] && dir[i] == 0 {
				continue
			}
			moves = append(moves, func(x, y int) (int, int) {
				return x + dir[i], y + dir[j]
			})
		}
	}
	return
}

type GameOfLife struct {
	oldB  board.Board
	newB  board.Board
	rules []NextMoveFunc
}

func New(bF board.Board, bS board.Board, rules []NextMoveFunc) *GameOfLife {
	return &GameOfLife{
		oldB:  bF,
		newB:  bS,
		rules: rules,
	}
}

func (g *GameOfLife) Generate() {
	count := g.oldB.Height * g.oldB.Width * 10 / 100
	for i := 0; i < count; i++ {
		n := rand.Int()
		m := rand.Int()
		x := n % g.oldB.Height
		y := m % g.oldB.Width
		g.oldB.SetCell(x, y, 1)
	}
}

func (g *GameOfLife) GetNeighborsCount(x, y int) (count int) {
	for _, f := range g.rules {
		if st, ok := g.oldB.GetCell(f(x, y)); ok && st == ALIVE {
			count++
		}
	}
	return
}

func (g *GameOfLife) GetFBoard() board.Board {
	return g.oldB
}

func (g *GameOfLife) GetSBoard() board.Board {
	return g.newB
}

func (g *GameOfLife) NextState(x, y int) (state State) {
	currentState, _ := g.oldB.GetCell(x, y)
	switch count := g.GetNeighborsCount(x, y); {
	case count == 3 || (currentState == ALIVE && count == 2):
		state = ALIVE
	default:
		state = DEAD
	}
	return
}

func (g *GameOfLife) Update() {
	for i := 0; i < g.oldB.Height; i++ {
		for j := 0; j < g.oldB.Width; j++ {
			newState := g.NextState(i, j)
			g.newB.SetCell(i, j, int(newState))
		}
	}
	g.oldB.Data, g.newB.Data = g.newB.Data, g.oldB.Data
}

func (g *GameOfLife) String() string {
	return g.oldB.String(func(state int) rune {
		switch state {
		case ALIVE:
			return '*'
		default:
			return '.'
		}

	})
}
