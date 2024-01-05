// Map handles the underlying map grid, but not the image itself or drawing it
package main

import (
	"image"
	"image/png"
	"os"
)

type mapCell struct {
	isLand       bool
	habitability float64
}

type mapGrid struct {
	area   []mapCell
	width  int
	height int
}

// Return cell at [x,y] of mapGrid (!)
func (m *mapGrid) CellAt(x, y int) *mapCell {
	return &m.area[y*m.width+x]
}

func calcHabitability(colorLevel uint32) float64 {
	return float64(0xffff-colorLevel) / 0xffff
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	img, err := png.Decode(file)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Creates the mapGrid from a .png image under `path`
func NewMapGrid(path string) (*mapGrid, error) {
	img, err := loadImage(path)
	if err != nil {
		return nil, err
	}
	imgWidth, imgHeight := img.Bounds().Max.X, img.Bounds().Max.Y

	mGrid := &mapGrid{
		area:   make([]mapCell, imgWidth*imgHeight),
		width:  imgWidth,
		height: imgHeight,
	}

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			red, green, blue, _ := img.At(x, y).RGBA()
			mc := mGrid.CellAt(x, y)
			if blue == 0xffff && red == 0 && green == 0 { //RGB of [255, 0, 0] means water
				mc.isLand = false
			} else {
				mc.isLand = true

				//The inverse of how white the cell is determines habilability, both habitable and
				//inhabitable cells have the same green value so we must use either blue or red
				//to determine the "whiteness". This will probably be changed in the future as we
				//have 4 values to work with which can be used to determine more things from a
				//single bitmap, like maybe rare resources. It is like this for now because it
				//makes the base map image clear and easy to edit with basic editing tools.
				mc.habitability = calcHabitability(red)
			}
		}
	}

	return mGrid, nil
}
