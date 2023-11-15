// Main ebiten stuff, main() and Update() is here
package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
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
	biggestPop   int
}

func (g *Game) Prerun(generations int) {
	for g.humanGrid.generation < generations {
		g.humanGrid.Update()
		gen := g.humanGrid.generation
		if gen%(generations/50) == 0 {
			fmt.Printf("Prerunning simulation... %d/%d (%d%%)\n",
				gen, generations, int(100*(float32(gen)/float32(generations))))
		}
	}
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
		if g.frameCounter%5 == 0 {
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
	mapPath := flag.String("mappath", "./defmap.png", "Path to the map PNG file.")
	prerunGenerations := flag.Int("prerun", 0, "Generations to simulate before launching, min 50")
	flag.Parse()

	preloadedMap, err := NewMapGrid(*mapPath) //Needed to set screen size
	if err != nil {
		log.Fatalf("Error creating MapGrid: %v", err)
	}

	g := &Game{
		mapGrid:      preloadedMap,
		humanGrid:    NewHumanGrid(*preloadedMap, screenWidth, screenHeight, (screenWidth*screenHeight)/5000),
		images:       makeImagesMap(),
		speed:        Paused,
		frameCounter: 0,
		biggestPop:   upperPopCap,
	}

	if *prerunGenerations > 49 {
		g.Prerun(*prerunGenerations)
	}

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("Anaxi")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
