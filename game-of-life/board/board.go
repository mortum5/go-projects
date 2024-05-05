package board

import (
	"fmt"
	"strings"
)

type PrintFunc func(int) rune

type Board struct {
	Width  int
	Height int
	Data   []byte
}

func NewBoard(width, height int) Board {
	d := width * height / 8
	if (width*height)%8 != 0 {
		d++
	}
	return Board{
		Width:  width,
		Height: height,
		Data:   make([]byte, d),
	}
}

func (b Board) Copy() Board {
	return b
}

func (b Board) GetCell(x, y int) (int, bool) {
	if x < 0 {
		x = b.Height - 1
	}

	if y < 0 {
		y = b.Width - 1
	}

	if x > b.Height-1 {
		x = 0
	}

	if y > b.Width-1 {
		y = 0
	}
	pos := x*b.Width + y
	byteOffset := pos / 8
	bitOffset := pos % 8

	return int((b.Data[byteOffset] >> bitOffset) & 1), true
}

func (b Board) SetCell(x, y int, state int) {
	pos := x*b.Width + y
	byteOffset := pos / 8
	bitOffset := pos % 8
	b.Data[byteOffset] = (b.Data[byteOffset] & ^(1 << bitOffset)) | (byte(state) << bitOffset)
}

func (b Board) Print() {
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			state, _ := b.GetCell(i, j)
			if state > 0 {
				fmt.Print(1)
			} else {
				fmt.Print(0)
			}
		}
		fmt.Println("")
	}
	fmt.Println("")
}

func (b Board) String(f PrintFunc) string {
	var sb strings.Builder
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			state, _ := b.GetCell(i, j)
			sb.WriteRune(f(state))
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
