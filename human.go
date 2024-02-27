// Processing the cell simulation of human cells
package main

import (
	"fmt"
	"math"
	"math/rand"
)

var upperPopCap int64 = 10000        //Cap of the population growth or population suitable for migration
var baseGrowthRate float64 = 0.01    //Max growth rate per tick
var baseMigrationRate float64 = 0.05 //Max portion of people migrating from one cell per tick

var devGrowthChance float64 = 0.001 //Chance dev of human cell with increase, scaled by population range
var devGrowthScale float64 = 0.8

type humanCell struct {
	x             int
	y             int
	adjacentCells []*humanCell
	mapCell       *mapCell
	population    int64
	development   float64
}

type coordinate struct {
	x int
	y int
}

// Human grid handles human activity
type HumanGrid struct {
	area           []humanCell
	areaChanges    []humanCell
	width          int
	height         int
	generation     int
	areaWorld      *mapGrid
	globalPop      int64
	biggestPopCell humanCell
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

// Get corresponding `areaChanges` cell of `area` cell
func (w *HumanGrid) ChangesCellOf(cell *humanCell) *humanCell {
	return w.ChangesCellAt(cell.x, cell.y)
}

// Get corresponding `areaWorld.area` cell of `area` cell
func (w *HumanGrid) MapCellOf(cell *humanCell) *mapCell {
	return w.MapCellAt(cell.x, cell.y)
}

// Get coordinates of neighboring cells up down left right of [x,y] within bounds
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

// Generate pointers to neighbors of cell c in world
func (c *humanCell) GenNeighbors(world *HumanGrid) {
	adjCoords := getNeighborsCoordinates(world.width, world.height, c.x, c.y)
	c.adjacentCells = make([]*humanCell, 0, len(adjCoords))
	for _, adjC := range adjCoords {
		c.adjacentCells = append(c.adjacentCells, world.CellAt(adjC.x, adjC.y))
	}
}

// Intialize cells by giving them their coordinates and generating neighbor pointers
func (c *humanCell) initCell(x, y int, world *HumanGrid, startingDev float64) {
	c.x, c.y = x, y
	c.development = startingDev
	c.GenNeighbors(world)
	c.mapCell = world.MapCellOf(c)
}

// init inits humanGrid by initing cells and placing with a population in random spots
//
// maxLiveCells sets the number of "rolls" for the random spots
func (w *HumanGrid) init(maxLiveCells int) {
	width := w.width
	height := w.height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			w.area[y*width+x].initCell(x, y, w, 1.0)
			w.areaChanges[y*width+x].initCell(x, y, w, 0.0)
		}
	}

	//Populate randomly
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		mc := w.MapCellAt(x, y)
		if mc.isLand {
			w.CellAt(x, y).population = rand.Int63n(upperPopCap / 2)
		}
	}
}

// NewHumanGrid creates a new HumanGrid
func NewHumanGrid(m mapGrid, width, height int, maxInitLiveCells int) *HumanGrid {
	w := &HumanGrid{
		area:           make([]humanCell, width*height),
		areaChanges:    make([]humanCell, width*height),
		width:          width,
		height:         height,
		generation:     0,
		areaWorld:      &m,
		globalPop:      0,
		biggestPopCell: humanCell{0, 0, nil, nil, 0, 0.0}, //Temporary
	}
	w.init(maxInitLiveCells)
	return w
}

func (w *HumanGrid) calcCapacityOfCell(cell *humanCell) float64 {
	hab := w.MapCellOf(cell).habitability
	if hab == 0 {
		return 0
	}

	//Logistic function with K being capacity
	return hab * float64(upperPopCap) * cell.development
}

// Gets pointers to neighbor cells of `cell` valid for migration:
// ones that are land and have small enough population (under K).
func (w *HumanGrid) getNeighborsForMigration(cell *humanCell, printDebugInfo bool) []*humanCell {
	validNeighbors := make([]*humanCell, 0, 4)

	for _, n := range cell.adjacentCells {

		if !w.MapCellOf(n).isLand {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: NOTLAND\n", n.x, n.y)
			}
			continue
		}

		if n.population > upperPopCap {
			if printDebugInfo {
				fmt.Printf("Neighbor at [%d,%d] skipped: OVERPOP\n", n.x, n.y)
			}
			continue
		}

		validNeighbors = append(validNeighbors, n)
	}
	return validNeighbors
}

// https://en.wikipedia.org/wiki/Logistic_function#In_ecology:_modeling_population_growth
func logisticFunction(r, P, K float64) float64 {
	return r * P * (1 - P/K)
}

// Calculates the random population growth of a cell, taking into account the habitability and baseGrowthRate
func (w *HumanGrid) calcPopChange(cell *humanCell) int64 {
	K := w.calcCapacityOfCell(cell)
	if K == 0 {
		return 0
	}

	change := logisticFunction(float64(baseGrowthRate), float64(cell.population), K) * rand.Float64()
	return int64(change)
}

func (w *HumanGrid) calcDevChance(cell *humanCell) float64 {
	dev := cell.development
	var limiter float64 = 1
	if dev < 1.5 { //Delay starting development
		limiter = 0.01
	}
	K := w.MapCellOf(cell).habitability * dev * float64(upperPopCap)
	return math.Min(1.0, float64(cell.population)/K) * devGrowthChance * limiter
}

// Applies population growth (calcPopChange()) of a cell into areaChanges
func (w *HumanGrid) updatePopAndDevGrowthOf(cell *humanCell) {
	pop := cell.population

	if pop > 2 {
		popChange := w.calcPopChange(cell)
		if popChange >= pop {
			popChange = pop
		}
		w.ChangesCellOf(cell).population += popChange

		//Roll for growth rate
		if rand.Float64() < w.calcDevChance(cell) {
			w.ChangesCellOf(cell).development += rand.Float64() * devGrowthScale
		}
	} else {
		w.ChangesCellOf(cell).population -= pop
	}
}

// Moves random population of cell in (one) random direction, applies that to areaChanges
func (w *HumanGrid) updateMigrationOf(cell *humanCell) {
	pop := cell.population
	if pop <= 20 {
		return
	}

	corrChangesCell := w.ChangesCellOf(cell)
	//To avoid having to apply the changes grid before doing migrations we check if the pop has
	//changed negatively, to limit migrations.
	if cccp := corrChangesCell.population; cccp < 0 {
		pop -= cccp
		if pop <= 0 {
			return
		}
	}

	validNeighbors := w.getNeighborsForMigration(cell, false)
	if len(validNeighbors) == 0 {
		return
	}

	cc := int64(float64(pop) * baseMigrationRate)
	if cc <= 0 {
		return
	}

	cIndex := rand.Intn(len(validNeighbors))
	chosenCell := validNeighbors[cIndex]
	chosenChangesCell := w.ChangesCellAt(chosenCell.x, chosenCell.y)

	peopleMoving := rand.Int63n(cc)
	chosenChangesCell.population += peopleMoving
	corrChangesCell.population -= peopleMoving
}

// Applies areaChanges into area by matrix addition, resets areaChanges and updates global stats
func (w *HumanGrid) applyChangesArea() error {
	var worldpop int64 = 0
	for i := range w.area {
		c := &w.area[i]
		corrC := &w.areaChanges[i]

		c.population += corrC.population
		c.development += corrC.development

		corrC.population = 0
		corrC.development = 0.0

		worldpop += c.population

		if c.population > w.biggestPopCell.population {
			w.biggestPopCell = *c
		}

		if c.population < 0 {
			return fmt.Errorf("negative population @ (%d,%d): %d.\nFull cell info:%+v", c.x, c.y, c.population, c)
		}
	}
	w.globalPop = worldpop
	return nil
}

// Update grid state by one tick. First pop growth, then migrations
func (w *HumanGrid) Update() error {

	//Sinusoidal pop growth for fun
	//baseGrowthRate = float32((math.Sin(float64(w.generation) / 250)) * 0.05)

	for i := range w.area {
		cell := &w.area[i]
		w.updatePopAndDevGrowthOf(cell)
		w.updateMigrationOf(cell)
	}

	err := w.applyChangesArea()
	w.generation++
	return err
}

// Within cells interlinked
