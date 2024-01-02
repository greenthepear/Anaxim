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
	Custom
	Unlimited
)

type MapMode int

const (
	PopMode MapMode = iota
	DevMode
)

type Sim struct {
	mapGrid   *mapGrid
	humanGrid *HumanGrid
}

type Anaxi struct {
	widget.BaseWidget

	simulation *Sim
	mapImage   image.Image
	mapCanvas  *canvas.Raster

	speed          Speed
	speedCustomTPS time.Duration
	lastTick       time.Time
	lastRefresh    time.Time

	speedWidgets    *SpeedWidgets
	leftInfoWidgets *LeftInfoWidgets
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
	return s.humanGrid.Update()
}

func NewAnaxi(s *Sim) *Anaxi {
	a := &Anaxi{
		simulation:     s,
		mapImage:       GenGridImage(s),
		speed:          Paused,
		lastTick:       time.Now(),
		lastRefresh:    time.Now(),
		speedCustomTPS: 0,
	}
	a.initUI()
	a.ExtendBaseWidget(a)
	return a
}

func (a *Anaxi) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(a.buildUI())
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

	//Refresh image and info only now and then
	if a.TimeSinceLastRefresh() > time.Second/24 {
		canvas.Refresh(a.mapCanvas)
		a.updateGlobalStatsWidgets()
		a.lastRefresh = time.Now()
	}
	a.lastTick = time.Now()
}

func (a *Anaxi) runSim() {
	go func() {
		for { //Bad for performance, should use tick channels instead for custom speed
			switch a.speed {
			case Paused:
				return
			case Unlimited:
				a.Update()
			default:
				if a.TimeSinceLastTick() > a.speedCustomTPS {
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

	a := app.New()

	s := &Sim{
		mapGrid:   preloadedMap,
		humanGrid: NewHumanGrid(*preloadedMap, mapWidth, mapHeight, (mapWidth*mapHeight)/5000),
	}

	if *prerunGenerations > 49 {
		s.Prerun(*prerunGenerations)
	}

	anaxi := NewAnaxi(s)

	w := a.NewWindow("Anaxi")

	w.SetContent(anaxi)

	w.Resize(fyne.NewSize(float32(mapWidth*2)+100, float32(mapHeight*2)+50))

	w.ShowAndRun()
}
