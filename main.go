// Initialization and general application loop
package main

import (
	"encoding/csv"
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
	mapImage            *image.RGBA
	mapTexture          *giu.Texture
	mapWidth, mapHeight int

	speed          Speed
	speedCustomTPS float32
	previousTick   time.Time
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
	mapImg :=
		image.NewRGBA(image.Rect(0, 0, s.mapGrid.width, s.mapGrid.height))

	UpdateGridImage(s, mapImg, image.Pt(-1, -1))

	a := &Anaxim{
		simulation:                  s,
		mapImage:                    mapImg,
		mapWidth:                    s.mapGrid.width,
		mapHeight:                   s.mapGrid.height,
		speed:                       Unlimited,
		lastTick:                    time.Now(),
		lastRefresh:                 time.Now(),
		speedCustomTPS:              1,
		howeringOverCellAt:          image.Pt(0, 0),
		howeringOverCellCanvasPoint: image.Pt(0, 0),
		inspectingCellAt:            image.Pt(-1, -1),
		inspectingCell:              nil,
	}
	return a
}

func (a *Anaxim) TimeSinceLastTick() time.Duration {
	return time.Since(a.lastTick)
}

func (a *Anaxim) TimeBetweenLastTicks() time.Duration {
	return a.lastTick.Sub(a.previousTick)
}

func (a *Anaxim) TimeSinceLastRefresh() time.Duration {
	return time.Since(a.lastRefresh)
}

func (a *Anaxim) WriteRecord() {
	gen := a.simulation.humanGrid.generation
	UpdateGridImage(a.simulation, a.mapImage, a.inspectingCellAt)
	pngFile, fileerr := os.Create(
		fmt.Sprintf("%v/gen_%v.png", recordFolder, gen))
	if fileerr != nil {
		log.Fatalf("While creating png file: %v", fileerr)
	}
	defer pngFile.Close()

	if fileerr = png.Encode(pngFile, a.mapImage); fileerr != nil {
		log.Fatalf("While encoding png file: %v", fileerr)
	}

	data := []string{
		pngFile.Name(),
		fmt.Sprint(gen),
		fmt.Sprint(a.simulation.humanGrid.globalPop),
		fmt.Sprint(a.simulation.humanGrid.biggestPopCell.population),
		fmt.Sprintf("(%d,%d)",
			a.simulation.humanGrid.biggestPopCell.x,
			a.simulation.humanGrid.biggestPopCell.y,
		),
	}
	if csverr := csvWriter.Write(data); csverr != nil {
		log.Fatal(csverr)
	}
	csvWriter.Flush()
	if csverr := csvWriter.Error(); csverr != nil {
		log.Fatal(csverr)
	}
}

func (a *Anaxim) Update() {
	err := a.simulation.Update()
	if err != nil {
		log.Fatalf("Simulation error: %v", err)
	}

	// Add csv record if needed
	if recordInterval != -1 &&
		(a.simulation.humanGrid.generation == 1 ||
			a.simulation.humanGrid.generation%recordInterval == 0) {

		a.WriteRecord()
	}

	//Refresh image and info only now and then
	if a.TimeSinceLastRefresh() > time.Second/24 {
		a.updateMapTexture()
	}

	a.previousTick = a.lastTick
	a.lastTick = time.Now()
}

func (a *Anaxim) loop() {
	a.createLayout()
}

func (a *Anaxim) updateMapTexture() {
	UpdateGridImage(a.simulation, a.mapImage, a.inspectingCellAt)
	giu.EnqueueNewTextureFromRgba(
		a.mapImage,
		func(tex *giu.Texture) {
			a.mapTexture = tex
		})
	giu.Update()
	a.lastRefresh = time.Now()
}

func (a *Anaxim) runSim() {
	for {
		switch a.speed {
		case Paused:
			return
		case Unlimited:
			a.Update()
		case Custom:
			a.Update()
			time.Sleep(time.Duration(
				(1 - a.speedCustomTPS) * float32(time.Second)))
		}
	}
}

var csvWriter *csv.Writer
var recordFolder string
var recordInterval int

func main() {

	mapPath := flag.String("mappath", "./Maps/oldworld.png", "Path to the map PNG file")
	prerunGenerations := flag.Int("prerun", 0, "Generations to simulate before launching, min 50")
	flag.IntVar(&mapResize, "mapsize", 4, "How much to resize the map (needs to be between 1 and 8)")
	cpuprof := flag.Bool("pprof", false, "Enable cpu profiling")

	flag.IntVar(&recordInterval, "ri", -1,
		"Record interval, how often to record the simulation data in a csv file and images")

	flag.Parse()

	if recordInterval > -1 {
		name := time.Now().Format(time.RFC3339)
		recordFolder = "record_" + name
		err := os.Mkdir(recordFolder, 0755)
		if err != nil {
			log.Fatalf("When creating record folder: %v", err)
		}
		csvFile, err := os.Create(
			recordFolder + "/" + name + ".csv")
		if err != nil {
			log.Fatalf("When creating csv file: %v", err)
		}
		csvWriter = csv.NewWriter(csvFile)
		csvWriter.Write([]string{"Map", "Generation", "World pop", "Record pop", "Record at"})
		defer csvFile.Close()
	}

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
		mapGrid: preloadedMap,
		humanGrid: NewHumanGrid(*preloadedMap, mapWidth, mapHeight,
			(mapWidth*mapHeight)/5000),
	}

	if *prerunGenerations > 49 {
		s.Prerun(*prerunGenerations)
	}

	anaxim := NewAnaxim(s)

	wnd := giu.NewMasterWindow("Anaxim",
		mapWidth*mapResize+leftColumnWidth+20,
		mapHeight*mapResize+100,
		giu.MasterWindowFlagsFloating)

	anaxim.updateMapTexture()
	anaxim.initUI()

	go anaxim.runSim()
	giu.Context.GetRenderer().
		SetTextureMagFilter(giu.TextureFilterNearest)

	wnd.Run(anaxim.loop)
}
