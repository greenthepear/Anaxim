// Map handles the underlying map grid, but not the image itself or drawing it
package main

type mapCell struct {
	isLand bool
}

type mapGrid struct {
	area []mapCell
}

func NewMapGrid(path string) *mapGrid {
	file := loadFile(path)

	colorSlice, width, height := getPixels(file)

	//Set screen size
	screenWidth, screenHeight = width, height

	w := &mapGrid{
		area: make([]mapCell, width*height),
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			_, green, blue, _ := colorSlice[y*width+x].RGBA()
			if green > blue {
				w.area[y*width+x].isLand = true
			}
		}
	}

	return w
}
