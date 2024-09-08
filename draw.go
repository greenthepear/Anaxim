// Drawing the simulation
package main

import (
	"image"
	"image/color"
	"math"
)

func UpdateGridImage(s *Sim, img *image.RGBA) {
	width := s.humanGrid.width

	maxPop := s.humanGrid.biggestPopCell.population
	for i, cell := range s.humanGrid.area {
		var landValue byte = 0
		if s.mapGrid.area[i].isLand {
			landValue = 0xff
		}

		pop := cell.population
		var red byte = 0
		if pop != 0 {
			popRange := 0.1 + float64(pop)/float64(maxPop)
			red = byte(0xff * (math.Min(1.0, popRange)))
		}

		col := color.RGBA{
			R: red,
			G: landValue - red,
			B: 100,
			A: 0xff,
		}

		//Fancy coordinates from index of 1d slice, probably slow though
		img.Set(i%width, i/width, col)
	}
}
