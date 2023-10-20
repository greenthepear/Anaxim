package main

import (
	"fmt"
	"math/rand"
)

type humanCell struct {
	population int
}

type coordinate struct {
	x int
	y int
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
	coords := make([]coordinate, 0, 9)
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
			coords = append(coords, coordinate{x2, y2})
		}
	}
	return coords
}
func (w *HumanGrid) updatePopGrowth() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			pop := w.area[y*w.width+x].population
			if pop != 0 && pop < 1000 {
				w.area[y*w.width+x].population += int(rand.Float32() * 0.5 * float32(pop))
			}
		}
	}
}

func (w *HumanGrid) updateMigration() {
	width := w.width
	height := w.height
	next := make([]humanCell, width*height)
	worldpop := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			shortMainCoord := y*width + x
			pop := w.area[shortMainCoord].population
			worldpop += pop
			if pop < 20 {
				next[shortMainCoord] = w.area[shortMainCoord]
				continue
			}
			for _, c := range getNeighborsCoordinatesMoore(w.area, width, height, x, y) {
				mainCellPopulation := w.area[shortMainCoord].population

				//fmt.Printf("x: %d\ty: %d\n", c.x, c.y)
				shortNeighborCoord := c.y*width + c.x
				currentNeighborCell := w.area[shortNeighborCoord]
				if currentNeighborCell.population < 1000 {
					peopleMoving := rand.Intn(int(float32(mainCellPopulation) * 0.05))
					//fmt.Printf("Moving %d people...\n", peopleMoving)
					next[shortNeighborCoord].population += peopleMoving
					next[shortMainCoord].population -= peopleMoving
					w.area[shortMainCoord].population -= peopleMoving
				}
			}
		}
	}
	fmt.Printf("World Population: %d\n", worldpop)
	w.area = next
}
