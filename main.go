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
	Unlimited
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
	speedTimes   map[Speed]time.Duration
	lastTick     time.Time
	lastRefresh  time.Time
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
	a := Anaxi{
		simulation: s,
		mapImage:   GenGridImage(s),
		speed:      Paused,
		speedTimes: map[Speed]time.Duration{
			Slow:   time.Second,
			Faster: time.Second / 10,
		},
		lastTick:    time.Now(),
		lastRefresh: time.Now(),
	}
	a.initUI()

	return &a
}

func (a *Anaxi) TimeSinceLastTick() time.Duration {
	return time.Since(a.lastTick)
}

func (a *Anaxi) TimeSinceLastRefresh() time.Duration {
	return time.Since(a.lastRefresh)
}

func (a *Anaxi) Update() {
	err := a.simulation.Update()
	if err != nil {
		log.Fatalf("Simulation error: %v", err)
	}

	//Refresh image only now and then
	if a.TimeSinceLastRefresh() > time.Second/24 {
		canvas.Refresh(a.mapCanvas)
		a.lastRefresh = time.Now()
	}
	a.lastTick = time.Now()
}

func (a *Anaxi) runSim() {
	go func() {
		for {
			switch a.speed {
			case Paused:
				//pass
			case Unlimited:
				a.Update()
			default:
				if a.TimeSinceLastTick() > a.speedTimes[a.speed] {
					a.Update()
				}
			}
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

	anaxi := NewAnaxi(s)

	a := app.New()
	w := a.NewWindow("Anaxi")

	anaxi.runSim()

	w.SetContent(anaxi.buildUI())

	w.Resize(fyne.NewSize(float32(mapWidth*2), float32(mapHeight*2)))

	w.ShowAndRun()
}
