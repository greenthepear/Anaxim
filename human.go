// Processing cellular automata of human cells
package main

import (
	"fmt"
	"math/rand"
)

var upperPopCap int = 100000
var baseGrowthRate float32 = 0.02
var baseMigrationRate float32 = 0.05

type humanCell struct {
	x             int
	y             int
	adjacentCells []*humanCell
	population    int
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
	generation  int
	areaWorld   *mapGrid
}

func (w *HumanGrid) CellAt(x, y int) *humanCell {
	return &w.area[y*w.width+x]
}

func (w *HumanGrid) ChangesCellAt(x, y int) *humanCell {
	return &w.areaChanges[y*w.width+x]
}

func (w *HumanGrid) MapCellAt(x, y int) *mapCell {
	return &w.areaWorld.area[y*w.width+x]
}

func (w *HumanGrid) CorrChangesCellOf(cell *humanCell) *humanCell {
	return w.ChangesCellAt(cell.x, cell.y)
}

func (w *HumanGrid) CorrWorldCellOf(cell *humanCell) *mapCell {
	return w.MapCellAt(cell.x, cell.y)
}

func getNeighborsCoordinates(width, height, x, y int) []coordinate {
	coords := make([]coordinate, 0, 4)
	if y < height-1 {
		coords = append(coords, coordinate{x, y + 1})
	}
	if y > 0 {
		coords = append(coords, coordinate{x, y - 1})
	}
	if x < width-1 {
		coords = append(coords, coordinate{x + 1, y})
	}
	if x > 0 {
		coords = append(coords, coordinate{x - 1, y})
	}
	return coords
}

func (c *humanCell) GenNeighbors(world *HumanGrid) {
	adjCoords := getNeighborsCoordinates(world.width, world.height, c.x, c.y)
	c.adjacentCells = make([]*humanCell, 0, len(adjCoords))
	for _, adjC := range adjCoords {
		c.adjacentCells = append(c.adjacentCells, world.CellAt(adjC.x, adjC.y))
	}
}

func (c *humanCell) initCell(x, y int, world *HumanGrid) {
	c.x, c.y = x, y
	c.GenNeighbors(world)
}

// init inits humanGrid with a random population
func (w *HumanGrid) init(maxLiveCells int) {
	width := w.width
	height := w.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.area[y*width+x].initCell(x, y, w)
			w.areaChanges[y*width+x].initCell(x, y, w)
		}
	}

	//Populate randomly
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		mc := w.MapCellAt(x, y)
		if mc.isLand {
			w.CellAt(x, y).population = rand.Intn(upperPopCap / 2)
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

func (w *HumanGrid) getNeighborsForMigration(cell *humanCell, printDebugInfo bool) []*humanCell {
	validNeighbors := make([]*humanCell, 0, 4)

	for _, n := range cell.adjacentCells {
		//TODO: Make this nicer or just remove the debug
		if n.population > upperPopCap {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: OVERPOP\n", n.x, n.y)
			}
			continue
		}

		if !w.CorrWorldCellOf(n).isLand {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: NOTLAND\n", n.x, n.y)
			}
			continue
		}

		validNeighbors = append(validNeighbors, n)
	}
	return validNeighbors
}

func (w *HumanGrid) calcPopChange(cell *humanCell) int {
	return int((w.CorrWorldCellOf(cell).habitability - rand.Float32()) * baseGrowthRate * float32(cell.population))
}

func (w *HumanGrid) updatePopGrowthOf(cell *humanCell) {
	pop := cell.population
	if pop > 2 && pop < upperPopCap {
		popChange := w.calcPopChange(cell)
		if popChange >= pop {
			popChange = pop
		}
		w.CorrChangesCellOf(cell).population += popChange
	}
}

func (w *HumanGrid) updateMigrationOf(cell *humanCell) {
	pop := cell.population
	corrChangesCell := w.CorrChangesCellOf(cell)
	if pop < 0 {
		fmt.Printf("Something has went terribly wrong...\n")
	}
	if pop < 20 {
		corrChangesCell.population += cell.population
		return
	}

	validNeighbors := w.getNeighborsForMigration(cell, false)
	if len(validNeighbors) == 0 {
		return
	}

	cc := int(float32(pop) * baseMigrationRate)
	if cc <= 0 {
		return
	}

	cIndex := rand.Intn(len(validNeighbors))
	chosenCell := validNeighbors[cIndex]
	chosenChangesCell := w.ChangesCellAt(chosenCell.x, chosenCell.y)

	peopleMoving := rand.Intn(cc)
	chosenChangesCell.population += peopleMoving
	corrChangesCell.population -= peopleMoving
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

	//Sinusoidal pop growth for fun
	//baseGrowthRate = float32((math.Sin(float64(w.generation) / 250)) * 0.05)

	width := w.width
	height := w.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := w.CellAt(x, y)
			w.updatePopGrowthOf(c)
			w.updateMigrationOf(c)
		}
	}
	w.applyChangesArea()
	w.generation++
}
