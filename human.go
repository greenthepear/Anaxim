// Processing cellular automata of human cells
package main

import (
	"fmt"
	"math/rand"
)

var upperPopCap int = 100000
var baseGrowthRate float32 = 0.05

type humanCell struct {
	population int
	//IDEA: make a slice of every neighbors coordinate and precompute that
}

type coordinate struct {
	x int
	y int
}

// Human grid handles human activity
type HumanGrid struct {
	area        []humanCell
	areaChanges []humanCell
	width       int
	height      int
	generation  uint64
	areaWorld   *mapGrid
}

// init inits humanGrid with a random population
func (w *HumanGrid) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		if w.areaWorld.area[y*w.width+x].isLand {
			w.area[y*w.width+x].population = rand.Intn(upperPopCap / 2)
		}
	}
}

// NewHumanGrid creates a new humanGrid
func NewHumanGrid(m mapGrid, width, height int, maxInitLiveCells int) *HumanGrid {
	w := &HumanGrid{
		area:        make([]humanCell, width*height),
		areaChanges: make([]humanCell, width*height),
		width:       width,
		height:      height,
		generation:  0,
		areaWorld:   &m,
	}
	w.init(maxInitLiveCells)
	return w
}

func getNeighborsCoordinates(world []humanCell, width, height, x, y int) []coordinate {
	coords := make([]coordinate, 0, 4)
	if y < height {
		coords = append(coords, coordinate{x, y + 1})
	}
	if y > 0 {
		coords = append(coords, coordinate{x, y - 1})
	}
	if x < width {
		coords = append(coords, coordinate{x + 1, y})
	}
	if x > 0 {
		coords = append(coords, coordinate{x - 1, y})
	}
	return coords
}

//lint:ignore U1000 Might switch
func getNeighborsCoordinatesMoore(world []humanCell, width, height, x, y int) []coordinate {
	coords := make([]coordinate, 0, 8)
	for j := -1; j <= 1; j++ {
		for i := -1; i <= 1; i++ {
			if i == 0 && j == 0 {
				continue
			}
			x2 := x + i
			y2 := y + j
			if x2 < 0 || y2 < 0 || width <= x2 || height <= y2 {
				continue
			}
			//cells = append(cells, world[y2*width+x2])

			//fmt.Printf("Coordinate for [%d,%d] found at [%d,%d]\t%d\n", x, y, x2, y2, y2*width+x2)
			coords = append(coords, coordinate{x2, y2})
		}
	}
	return coords
}

func (w *HumanGrid) getNeighborsForMigration(x, y int, printDebugInfo bool) []coordinate {
	gridNeighbors := getNeighborsCoordinates(w.area, w.width, w.height, x, y)

	validNeighbors := make([]coordinate, 0, 9)
	//mainCellPopulation := w.area[y*w.width+x].population
	for _, n := range gridNeighbors {
		nCoord := n.y*w.width + n.x

		//TODO: Make this nicer or just remove the debug
		if w.area[nCoord].population > upperPopCap {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: OVERPOP\n", n.x, n.y)
			}
			continue
		}

		if !w.areaWorld.area[nCoord].isLand {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: NOTLAND\n", n.x, n.y)
			}
			continue
		}

		validNeighbors = append(validNeighbors, n)
	}
	return validNeighbors
}

func (w *HumanGrid) updatePopGrowthAt(x, y int) {
	pop := w.area[y*w.width+x].population
	if pop > 2 && pop < upperPopCap {
		w.areaChanges[y*w.width+x].population += int(rand.Float32() * baseGrowthRate * float32(pop))
	}
}

func (w *HumanGrid) updateMigrationAt(x, y int) {
	width := w.width
	shortMainCoord := y*width + x
	pop := w.area[shortMainCoord].population
	if pop < 0 {
		fmt.Printf("Something has went terribly wrong...\n")
	}
	if pop < 20 {
		w.areaChanges[shortMainCoord].population += w.area[shortMainCoord].population
		return
	}

	validNeighbors := w.getNeighborsForMigration(x, y, false)
	if len(validNeighbors) == 0 {
		return
	}
	chosenDirection := rand.Intn(len(validNeighbors))
	chosenDirectionCoord := validNeighbors[chosenDirection].y*width + validNeighbors[chosenDirection].x

	cc := int(float32(pop) * 0.05)

	if cc <= 0 {
		return
	}

	peopleMoving := rand.Intn(cc)
	w.areaChanges[chosenDirectionCoord].population += peopleMoving
	w.areaChanges[shortMainCoord].population -= peopleMoving
}

func (w *HumanGrid) applyChangesArea() {
	width := w.width
	height := w.height
	worldpop := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			coord := y*width + x
			w.area[coord].population += w.areaChanges[coord].population
			w.areaChanges[coord].population = 0
			worldpop += w.area[coord].population
		}
	}
	if w.generation%64 == 0 {
		fmt.Printf("Gen %d | World Population: %d\n", w.generation, worldpop)
	}
}

// Update game state by one tick.
func (w *HumanGrid) Update(m mapGrid) {
	width := w.width
	height := w.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.updatePopGrowthAt(x, y)
			w.updateMigrationAt(x, y)
		}
	}
	w.applyChangesArea()
	w.generation++
}
