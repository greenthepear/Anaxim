// Handling input
package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (w HumanGrid) genCellInfoAtCursor() string {
	cursorX, cursorY := ebiten.CursorPosition()
	if cursorX >= 0 && cursorX < w.width && cursorY >= 0 && cursorY < w.height {
		pop := w.CellAt(cursorX, cursorY).population
		return fmt.Sprintf("[%d,%d]:\n%d", cursorX, cursorY, pop)
	}
	return ""
}

func (w HumanGrid) clickDebug() {
	cursorX, cursorY := ebiten.CursorPosition()
	if cursorX >= 0 && cursorX < screenWidth && cursorY >= 0 && cursorY < screenHeight {
		c := w.CellAt(cursorX, cursorY)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			fmt.Printf("* Click at [%d,%d]:\nHuman: %v\nMap: %v\n*\n",
				cursorX, cursorY, c, w.CorrMapCellOf(c))
		}
	}
}

func (g *Game) handleSpeedControls() {
	cursorX, cursorY := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if cursorY > screenHeight-16 && cursorX < 70 {
			s := int(cursorX / 16)
			if s > 3 {
				s = 3
			}
			g.speed = Speed(s)
		}
	}

}
