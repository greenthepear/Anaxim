package main

import (
	"fmt"
	"image"

	"github.com/AllenDang/giu"
)

func (a *Anaxi) mapInputEvents() giu.Widget {
	return giu.Event().OnHover(func() {
		//empty for now
	}).OnClick(giu.MouseButtonLeft, func() {
		imgPos := giu.GetCursorPos()
		curPos := giu.GetMousePos()

		actualPos := image.Point{
			X: curPos.X - (imgPos.X - mapWidth*mapResize) + 8, //magic padding number
			Y: curPos.Y - imgPos.Y,
		}
		fmt.Printf("\n%v -> %v\n%v\n", imgPos, curPos, actualPos)
	})
}
