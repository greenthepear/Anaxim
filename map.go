// Map handles the underlying map grid, but not the image itself or drawing it
package main

type mapCell struct {
	isLand       bool
	habitability float32
}

type mapGrid struct {
	area   []mapCell
	width  int
	height int
}

func (m *mapGrid) CellAt(x, y int) *mapCell {
	return &m.area[y*m.width+x]
}

func calcHabitability(colorLevel uint32) float32 {
	return float32(0xffff-colorLevel) / 0xffff
}

func NewMapGrid(path string) *mapGrid {
	colorSlice, imgWidth, imgHeight := getPixels(path)

	//Set screen size
	screenWidth, screenHeight = imgWidth, imgHeight

	mGrid := &mapGrid{
		area:   make([]mapCell, imgWidth*imgHeight),
		width:  imgWidth,
		height: imgHeight,
	}

	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			red, green, blue, _ := colorSlice[y*imgWidth+x].RGBA()
			if blue == 0xffff && red == 0 && green == 0 {
				mGrid.area[y*imgWidth+x].isLand = false
			} else {
				mGrid.area[y*imgWidth+x].isLand = true
				mGrid.area[y*imgWidth+x].habitability = calcHabitability(red)
			}
		}
	}

	return mGrid
}
