package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// HumanGrid represents the game state
type HumanGrid struct {
	area   []humanCell
	width  int
	height int
}

// NewHumanGrid creates a new humanGrid
func NewHumanGrid(width, height int, maxInitLiveCells int) *HumanGrid {
	w := &HumanGrid{
		area:   make([]humanCell, width*height),
		width:  width,
		height: height,
	}
	w.init(maxInitLiveCells)
	return w
}

// init inits humanGrid with a random population
func (w *HumanGrid) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x].population = rand.Intn(200)
	}
}

// Update game state by one tick.
func (w *HumanGrid) Update() {
	w.updatePopGrowth()
	w.updateMigration()
}

// Draw paints current game state.
func (w *HumanGrid) Draw(pix []byte) {
	for i, v := range w.area {
		if v.population != 0 {
			pix[4*i] = byte(255 * (float32(v.population) / 1000))
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
	screenWidth  = 100
	screenHeight = 100
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {

	g := &Game{
		humanGrid: NewHumanGrid(screenWidth, screenHeight, int((screenWidth*screenHeight)/50)),
	}

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Anexi")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
