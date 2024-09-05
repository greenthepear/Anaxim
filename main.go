// Initialization and general application loop
package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/AllenDang/giu"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

var mapResize = 4 //TODO: calculate it instead

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

type Anaxim struct {
	simulation          *Sim
	mapImage            image.Image
	mapTexture          *giu.Texture
	mapWidth, mapHeight int

	speed          Speed
	speedCustomTPS int32
	lastTick       time.Time
	lastRefresh    time.Time

	speedWidgets *SpeedWidgets

	howeringOverCellAt          image.Point
	howeringOverCellCanvasPoint image.Point
	inspectingCellAt            image.Point
	inspectingCell              *humanCell
	inspectingCanvasPoint       image.Point
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

func NewAnaxim(s *Sim) *Anaxim {
	a := &Anaxim{
		simulation: s,
		mapImage: GenGridImage(s,
			image.NewRGBA(image.Rect(0, 0, s.mapGrid.width, s.mapGrid.height))),
		mapWidth:                    s.mapGrid.width,
		mapHeight:                   s.mapGrid.height,
		speed:                       Unlimited,
		lastTick:                    time.Now(),
		lastRefresh:                 time.Now(),
		speedCustomTPS:              1,
		howeringOverCellAt:          image.Pt(0, 0),
		howeringOverCellCanvasPoint: image.Pt(0, 0),
		inspectingCellAt:            image.Pt(0, 0),
		inspectingCell:              nil,
	}
	return a
}

func (a *Anaxim) TimeSinceLastTick() time.Duration {
	return time.Since(a.lastTick)
}

func (a *Anaxim) TimeSinceLastRefresh() time.Duration {
	return time.Since(a.lastRefresh)
}

func (a *Anaxim) Update() {
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

func (a *Anaxim) loop() {
	a.createLayout()
}

func (a *Anaxim) updateMapTexture() {
	giu.EnqueueNewTextureFromRgba(
		GenGridImage(a.simulation, a.mapImage),
		func(tex *giu.Texture) {
			a.mapTexture = tex
		})
}

func (a *Anaxim) runSim() {
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

	mapPath := flag.String("mappath", "./Maps/oldworld.png", "Path to the map PNG file")
	prerunGenerations := flag.Int("prerun", 0, "Generations to simulate before launching, min 50")
	flag.IntVar(&mapResize, "mapsize", 4, "How much to resize the map (needs to be between 1 and 8)")
	cpuprof := flag.Bool("pprof", false, "Enable cpu profiling")

	flag.Parse()

	if *cpuprof {
		f, err := os.Create("anaxim.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if mapResize < 1 || mapResize > 8 {
		log.Fatalf("Flag mapsize out of range [1,8]: %d", mapResize)
	}

	preloadedMap, err := NewMapGrid(*mapPath)
	if err != nil {
		log.Fatalf("Error creating MapGrid: %v", err)
	}
	mapWidth, mapHeight := preloadedMap.width, preloadedMap.height

	s := &Sim{
		mapGrid:   preloadedMap,
		humanGrid: NewHumanGrid(*preloadedMap, mapWidth, mapHeight, (mapWidth*mapHeight)/5000),
	}

	if *prerunGenerations > 49 {
		s.Prerun(*prerunGenerations)
	}

	anaxim := NewAnaxim(s)

	wnd := giu.NewMasterWindow("Anaxim", mapWidth*mapResize+leftColumnWidth+20, mapHeight*mapResize+100, giu.MasterWindowFlagsNotResizable)
	giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)

	anaxim.updateMapTexture()
	anaxim.initUI()

	anaxim.runSim()

	wnd.Run(anaxim.loop)
}
