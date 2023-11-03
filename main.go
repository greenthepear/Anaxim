// Main ebiten stuff, Update() is here
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	initFonts()
}

var (
	screenWidth  = 320
	screenHeight = 180
)

type Game struct {
	mapGrid   *mapGrid
	humanGrid *HumanGrid
	pixels    []byte
}

func (g *Game) Update() error {
	g.humanGrid.Update(*g.mapGrid)
	g.humanGrid.clickDebug()
	return nil
}

func main() {
	preloadedMap := NewMapGrid("./defmap.png") //Needed to set screen size

	g := &Game{
		mapGrid:   preloadedMap,
		humanGrid: NewHumanGrid(*preloadedMap, screenWidth, screenHeight, int((screenWidth*screenHeight)/1000)),
	}

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Anaxi")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
