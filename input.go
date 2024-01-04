package main

import (
	"image"

	"github.com/AllenDang/giu"
)

func snapPointToGrid(pt image.Point, gridSize int) image.Point {
	//Div floors
	return pt.Div(mapResize).Mul(mapResize).Add(image.Pt(mapResize, mapResize))
}

func (a *Anaxi) mapInputEvents() giu.Widget {
	return giu.Event().OnHover(func() {
		drawCursorPos := giu.GetCursorPos()
		mousePos := giu.GetMousePos()

		a.howeringOverCellCanvasPoint = snapPointToGrid(mousePos, mapResize)

		overImagePos := image.Point{
			X: mousePos.X - (drawCursorPos.X - mapWidth*mapResize) + 8, //magic padding number
			Y: mousePos.Y - drawCursorPos.Y,
		}
		pixelPos := image.Point{
			X: overImagePos.X / mapResize,
			Y: overImagePos.Y / mapResize,
		}

		a.howeringOverCellAt = pixelPos
	}).OnClick(giu.MouseButtonLeft, func() {
		a.inspectingCellAt = a.howeringOverCellAt
		realPos := giu.GetMousePos()

		if mapResize != 1 {
			realPos = snapPointToGrid(realPos, mapResize)
		}

		a.inspectingCanvasPoint = realPos

		a.inspectingCell = a.simulation.humanGrid.CellAt(
			a.inspectingCellAt.X, a.inspectingCellAt.Y)
	})
}
