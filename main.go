// Main ebiten stuff, Update() is here
package main

import (
	"image"
	"image/png"
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

type Speed int

const (
	Paused Speed = iota
	Slow
	Faster
	Fastest
)

type Game struct {
	mapGrid      *mapGrid
	humanGrid    *HumanGrid
	pixels       []byte
	images       map[string]*ebiten.Image
	speed        Speed
	frameCounter uint64
}

func (g *Game) Update() error {
	switch g.speed {
	case Paused:
		//Pass
	case Slow:
		if g.frameCounter%60 == 0 {
			g.humanGrid.Update()
		}
	case Faster:
		if g.frameCounter%10 == 0 {
			g.humanGrid.Update()
		}
	case Fastest:
		g.humanGrid.Update()
	}
	g.handleSpeedControls()
	g.humanGrid.clickDebug()
	g.frameCounter++
	return nil
}

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	preloadedMap := NewMapGrid("./defmap.png") //Needed to set screen size

	g := &Game{
		mapGrid:      preloadedMap,
		humanGrid:    NewHumanGrid(*preloadedMap, screenWidth, screenHeight, int((screenWidth*screenHeight)/1000)),
		images:       makeImagesMap(),
		speed:        Faster,
		frameCounter: 0,
	}

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Anaxi")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
