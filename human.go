// Processing cellular automata of human cells
package main

import (
	"fmt"
	"math/rand"
)

var upperPopCap int = 100000
var baseGrowthRate float32 = 0.05

type humanCell struct {
	population         int
	x                  int
	y                  int
	adjacentCellCoords []coordinate
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

func (c *humanCell) GenNeighbors(world *HumanGrid) {
	c.adjacentCellCoords = getNeighborsCoordinates(world.area, world.width, world.height, c.x, c.y)
}

// init inits humanGrid with a random population
func (w *HumanGrid) init(maxLiveCells int) {
	width := w.width
	height := w.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.area[y*width+x].x, w.area[y*width+x].y = x, y
			w.area[y*width+x].GenNeighbors(w)
		}
	}

	//Populate randomly
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		if w.areaWorld.area[y*width+x].isLand {
			w.area[y*width+x].population = rand.Intn(upperPopCap / 2)
		}
	}
}

// NewHumanGrid creates a new HumanGrid
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

func (w *HumanGrid) CellAt(x int, y int) *humanCell {
	return &w.area[y*w.width+x]
}

func (w *HumanGrid) ChangesCellAt(x int, y int) *humanCell {
	return &w.areaChanges[y*w.width+x]
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

			coords = append(coords, coordinate{x2, y2})
		}
	}
	return coords
}

func (w *HumanGrid) getNeighborsForMigration(x, y int, printDebugInfo bool) []coordinate {
	validNeighbors := make([]coordinate, 0, 4)

	for _, n := range w.CellAt(x, y).adjacentCellCoords {
		//TODO: Make this nicer or just remove the debug
		if w.CellAt(n.x, n.y).population > upperPopCap {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: OVERPOP\n", n.x, n.y)
			}
			continue
		}

		if !w.areaWorld.CellAt(n.x, n.y).isLand {
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
	pop := w.CellAt(x, y).population
	if pop > 2 && pop < upperPopCap {
		w.ChangesCellAt(x, y).population +=
			int((w.areaWorld.CellAt(x, y).habitability - rand.Float32()) * baseGrowthRate * float32(pop))
	}
}

func (w *HumanGrid) updateMigrationAt(x, y int) {
	pop := w.CellAt(x, y).population
	if pop < 0 {
		fmt.Printf("Something has went terribly wrong...\n")
	}
	if pop < 20 {
		w.CellAt(x, y).population += w.CellAt(x, y).population
		return
	}

	validNeighbors := w.getNeighborsForMigration(x, y, false)
	if len(validNeighbors) == 0 {
		return
	}
	chosenDirection := rand.Intn(len(validNeighbors))
	chosenDirectionX, chosenDirectionY := validNeighbors[chosenDirection].x, validNeighbors[chosenDirection].y

	cc := int(float32(pop) * 0.05)

	if cc <= 0 {
		return
	}

	peopleMoving := rand.Intn(cc)
	w.ChangesCellAt(chosenDirectionX, chosenDirectionY).population += peopleMoving
	w.ChangesCellAt(x, y).population -= peopleMoving
}

func (w *HumanGrid) applyChangesArea() {
	width := w.width
	height := w.height
	worldpop := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.CellAt(x, y).population += w.ChangesCellAt(x, y).population
			w.ChangesCellAt(x, y).population = 0
			worldpop += w.CellAt(x, y).population
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
