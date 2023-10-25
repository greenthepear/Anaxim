package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Draw paints current game state.
func (w *HumanGrid) Draw(pix []byte) {
	for i, v := range w.area {
		if v.population != 0 {
			pix[4*i] = byte(254.0 * (math.Min(1.0, float64(v.population)/float64(upperPopCap))))
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

const (
	screenWidth  = 320
	screenHeight = 180
)

type Game struct {
	humanGrid *HumanGrid
	pixels    []byte
}

func (g *Game) Update() error {
	g.humanGrid.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.humanGrid.Draw(g.pixels)
	screen.WritePixels(g.pixels)
	ebitenutil.DebugPrint(screen, g.humanGrid.genCellInfoAtCursor())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	g := &Game{
		humanGrid: NewHumanGrid(screenWidth, screenHeight, int((screenWidth*screenHeight)/1000)),
	}

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Anaxi")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
