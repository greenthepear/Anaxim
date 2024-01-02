// Drawing the simulation
package main

import (
	"image"
	"image/color"
	"math"
)

func (s *Sim) popMapModePixelDrawer(cell humanCell) func() {
	return color.RGBA{}
}

func GenGridImage(s *Sim, mapModeFunc func(c humanCell)) image.Image {
	width, height := s.humanGrid.width, s.humanGrid.height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	maxPop := s.humanGrid.biggestPopCell.population
	for i, cell := range s.humanGrid.area {
		var landValue byte = 0
		if s.mapGrid.area[i].isLand {
			landValue = 255
		}
		pop := cell.population

		popRange := float64(pop) / float64(maxPop)
		if cell.population != 0 { //To see places where there is ANY population
			popRange += 0.1
		}
		red := byte(0xff * (math.Min(1.0, popRange)))

		col := color.RGBA{
			R: red,
			G: landValue - red,
			B: 100 - landValue,
			A: 0xff,
		}

		//Fancy coordinates from index of 1d slice, probably slow though
		img.Set(i%width, i/width, col)
	}
	return img
}

// Called by canvas.Refresh
func (a *Anaxi) updateGridImage(w, h int) image.Image {
	return GenGridImage(a.simulation, a.mapModeDrawers[a.mapMode])
}
