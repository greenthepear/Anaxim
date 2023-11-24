// Drawing the simulation
package main

import (
	"image"
	"image/color"
	"math"
)

func GenGridImage(s *Sim) image.Image {
	width, height := s.humanGrid.width, s.humanGrid.height
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	maxPop := s.biggestPop
	for i, v := range s.humanGrid.area {
		var landValue byte = 0
		if s.mapGrid.area[i].isLand {
			landValue = 255
		}
		pop := v.population

		popRange := float64(pop) / float64(maxPop)
		if v.population != 0 { //To see places where there is ANY population
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

		if pop > maxPop {
			maxPop = pop
		}
	}
	s.biggestPop = maxPop
	return img
}

// Called by canvas.Refresh
func (a *Anaxi) updateGridImage(w, h int) image.Image {
	return GenGridImage(a.simulation)
}
