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

	"github.com/AllenDang/giu"
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

type Sim struct {
	mapGrid   *mapGrid
	humanGrid *HumanGrid
}

type Anaxi struct {
	simulation *Sim
	mapImage   image.Image
	mapTexture *giu.Texture

	speed          Speed
	speedCustomTPS int32
	lastTick       time.Time
	lastRefresh    time.Time

	speedWidgets *SpeedWidgets
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
		speed:          Unlimited,
		lastTick:       time.Now(),
		lastRefresh:    time.Now(),
		speedCustomTPS: 1,
	}
	return a
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
		//a.updateMapImage()
		a.updateMapTexture()
		giu.Update()
		a.lastRefresh = time.Now()
	}
	a.lastTick = time.Now()
}

func (a *Anaxi) loop() {
	a.createLayout()
}

func (a *Anaxi) updateMapTexture() {
	giu.EnqueueNewTextureFromRgba(GenGridImage(a.simulation), func(tex *giu.Texture) {
		a.mapTexture = tex
	})
}

func (a *Anaxi) runSim() {
	go func() {
		for { //Bad for performance, should use tick channels instead for custom speed
			switch a.speed {
			case Paused:
				time.Sleep(time.Microsecond)
			case Unlimited:
				a.Update()
			default:
				if a.TimeSinceLastTick() > time.Second/time.Duration(a.speedCustomTPS) {
					a.Update()
				}
				time.Sleep(time.Microsecond)
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
		mapGrid:   preloadedMap,
		humanGrid: NewHumanGrid(*preloadedMap, mapWidth, mapHeight, (mapWidth*mapHeight)/5000),
	}

	if *prerunGenerations > 49 {
		s.Prerun(*prerunGenerations)
	}

	anaxi := NewAnaxi(s)

	wnd := giu.NewMasterWindow("Anaxi", mapWidth*4+200, mapHeight*4+50, giu.MasterWindowFlagsFloating)
	giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)

	anaxi.updateMapTexture()
	anaxi.initUI()

	anaxi.runSim()

	wnd.Run(anaxi.loop)
}
