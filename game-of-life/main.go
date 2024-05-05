package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mortum5/go-projects/game-of-life/board"
	"github.com/mortum5/go-projects/game-of-life/game-of-life"
)

const (
	N    = 600
	M    = 600
	SIZE = 5
)

// Game implements ebiten.Game intrface.
type Game struct {
	width      int
	height     int
	generation int
	pause      bool
	gol        *game.GameOfLife
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.pause = !g.pause
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		x, y := ebiten.CursorPosition()

		y = y / SIZE
		x = x / SIZE

		state, _ := g.gol.GetFBoard().GetCell(x, y)

		g.gol.GetFBoard().SetCell(x, y, (state+1)%2)
	}

	if !g.pause {
		g.gol.Update()
		g.generation++

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	str := fmt.Sprintf("Generation: %d, FPS: %f, TPS: %f", g.generation, ebiten.ActualFPS(), ebiten.ActualTPS())
	ebitenutil.DebugPrint(screen, str)
	b := g.gol.GetFBoard()
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if state, ok := b.GetCell(i, j); ok && state == game.ALIVE {
				vector.DrawFilledRect(screen, float32(i*SIZE), float32(j*SIZE), SIZE, SIZE, color.White, true)
			}
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return g.width, g.height
}

func main() {
	boardF := board.NewBoard(M, N)
	boardS := board.NewBoard(M, N)
	moves := game.GenerateMoves()
	gol := game.New(boardF, boardS, moves)
	gol.Generate()

	game := &Game{gol: gol}

	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(1200, 800)
	// ebiten.SetTPS(2)
	ebiten.SetWindowTitle("Game of life")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
