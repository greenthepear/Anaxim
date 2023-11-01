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
}

type coordinate struct {
	x int
	y int
}

// HumanGrid represents the game state
type HumanGrid struct {
	area        []humanCell
	areaChanges []humanCell
	width       int
	height      int
	generation  uint64
}

// init inits humanGrid with a random population
func (w *HumanGrid) init(maxLiveCells int) {
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x].population = rand.Intn(upperPopCap / 2)
	}
}

// NewHumanGrid creates a new humanGrid
func NewHumanGrid(width, height int, maxInitLiveCells int) *HumanGrid {
	w := &HumanGrid{
		area:        make([]humanCell, width*height),
		areaChanges: make([]humanCell, width*height),
		width:       width,
		height:      height,
		generation:  0,
	}
	w.init(maxInitLiveCells)
	return w
}

//lint:ignore U1000 Will try this one later maybe
func getNeighborsCoordinates(world []humanCell, width, height, x, y int) []coordinate {
	coords := make([]coordinate, 0, 4)
	if y < height {
		coords = append(coords, coordinate{y + 1, x})
	}
	if y > 0 {
		coords = append(coords, coordinate{y - 1, x})
	}
	if x < width {
		coords = append(coords, coordinate{y, x + 1})
	}
	if x > 0 {
		coords = append(coords, coordinate{y, x - 1})
	}
	return coords
}

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
func (w *HumanGrid) updatePopGrowthAt(x int, y int) {
	pop := w.area[y*w.width+x].population
	if pop > 2 && pop < upperPopCap {
		w.areaChanges[y*w.width+x].population += int(rand.Float32() * baseGrowthRate * float32(pop))
	}
}

func (w *HumanGrid) updateMigrationAt(x int, y int) {
	width := w.width
	height := w.height
	shortMainCoord := y*width + x
	pop := w.area[shortMainCoord].population
	if pop < 0 {
		fmt.Printf("Something has went terribly wrong...\n")
	}
	if pop < 20 {
		w.areaChanges[shortMainCoord] = w.area[shortMainCoord]
		return
	}

	for _, c := range getNeighborsCoordinatesMoore(w.area, width, height, x, y) {
		mainCellPopulation :=
			w.area[shortMainCoord].population - w.areaChanges[shortMainCoord].population

		//fmt.Printf("x: %d\ty: %d\n", c.x, c.y)
		shortNeighborCoord := c.y*width + c.x
		currentNeighborCell := w.area[shortNeighborCoord]
		if currentNeighborCell.population < upperPopCap {
			cc := int(float32(mainCellPopulation) * 0.05)
			if cc <= 0 { //bad
				continue
			}
			peopleMoving := rand.Intn(cc)
			w.areaChanges[shortNeighborCoord].population += peopleMoving
			w.areaChanges[shortMainCoord].population -= peopleMoving
		}
	}
}

func (w *HumanGrid) applyChangesArea() {
	width := w.width
	height := w.height
	worldpop := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.area[y*width+x].population += w.areaChanges[y*width+x].population
			w.areaChanges[y*width+x].population = 0
			worldpop += w.area[y*width+x].population
		}
	}
	if w.generation%64 == 0 {
		fmt.Printf("Gen %d | World Population: %d\n", w.generation, worldpop)
	}
}

// Update game state by one tick.
func (w *HumanGrid) Update() {
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
