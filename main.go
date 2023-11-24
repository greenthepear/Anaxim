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

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

var ( //set by NewMapGrid
	mapWidth  = 0
	mapHeight = 0
)

type Speed int

const (
	Paused Speed = iota
	Slow
	Faster
	Fastest
)

type Sim struct {
	mapGrid    *mapGrid
	humanGrid  *HumanGrid
	biggestPop int
}

type Anaxi struct {
	widget.BaseWidget

	simulation   *Sim
	mapImage     image.Image
	mapCanvas    *canvas.Raster
	speed        Speed
	speedTimes   map[Speed]uint64
	frameCounter uint64
	speedButtons []*widget.Button
}

func (s *Sim) Prerun(generations int) {
	for s.humanGrid.generation < generations {
		s.humanGrid.Update()
		gen := s.humanGrid.generation
		if gen%(generations/50) == 0 {
			fmt.Printf("Prerunning simulation... %d/%d (%d%%)\n",
				gen, generations, int(100*(float32(gen)/float32(generations))))
		}
	}
}

func (s *Sim) Update() error {
	s.humanGrid.Update()
	return nil
}

func NewAnaxi(s *Sim) *Anaxi {
	ar := Anaxi{
		simulation: s,
		mapImage:   GenGridImage(s),
		speed:      Paused,
		speedTimes: map[Speed]uint64{
			Slow:   60,
			Faster: 10,
		},
		frameCounter: 0,
	}
	ar.speedButtons = ar.genSpeedControls()
	ar.mapCanvas = canvas.NewRaster(ar.draw)
	ar.mapCanvas.ScaleMode = canvas.ImageScalePixels

	return &ar
}

func (a *Anaxi) Update() {
	err := a.simulation.Update()
	if err != nil {
		log.Fatalf("Simulation error: %v", err)
	}
	a.updateGridImage()
	canvas.Refresh(a.mapCanvas)
}

func (a *Anaxi) animate() { //TODO: There's gotta be a way to use a ticker instead of the frame counter things
	go func() {
		tick := time.NewTicker(time.Second / 60)

		for range tick.C {
			a.frameCounter++
			fc := a.frameCounter
			switch a.speed {
			case Paused:
				continue
			case Slow:
				if fc%a.speedTimes[Slow] != 0 {
					continue
				}
			case Faster:
				if fc%a.speedTimes[Faster] != 0 {
					continue
				}
			case Fastest:
				//Pass
			}
			a.Update()
		}
	}()
}

func main() {
	mapPath := flag.String("mappath", "./defmap.png", "Path to the map PNG file.")
	prerunGenerations := flag.Int("prerun", 0, "Generations to simulate before launching, min 50")
	flag.Parse()

	preloadedMap, err := NewMapGrid(*mapPath) //Needed to set screen size
	if err != nil {
		log.Fatalf("Error creating MapGrid: %v", err)
	}

	s := &Sim{
		mapGrid:    preloadedMap,
		humanGrid:  NewHumanGrid(*preloadedMap, mapWidth, mapHeight, (mapWidth*mapHeight)/5000),
		biggestPop: upperPopCap,
	}

	if *prerunGenerations > 49 {
		s.Prerun(*prerunGenerations)
	}

	ar := NewAnaxi(s)

	a := app.New()
	w := a.NewWindow("Anaxi")

	ar.animate()

	w.SetContent(ar.buildUI())

	w.Resize(fyne.NewSize(float32(mapWidth*2), float32(mapHeight*2)))

	w.ShowAndRun()
}
