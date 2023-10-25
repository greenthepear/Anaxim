package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func (w HumanGrid) cellAt(x int, y int) humanCell {
	return w.area[y*w.width+x]
}

func (w HumanGrid) genCellInfoAtCursor() string {
	cursorX, cursorY := ebiten.CursorPosition()
	if cursorX >= 0 && cursorX < w.width && cursorY >= 0 && cursorY < w.height {
		pop := w.cellAt(cursorX, cursorY).population
		return fmt.Sprintf("[%d,%d]: %d", cursorX, cursorY, pop)
	}
	return ""
}
