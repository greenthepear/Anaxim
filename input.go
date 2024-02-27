// Handling of most user input,
// excluding only ones that don't need abstraction out of Layout in ui.go
package main

import (
	"image"

	"github.com/AllenDang/giu"
)

func (a *Anaxim) setSpeedToUnlimited() {
	a.speed = Unlimited
	a.speedWidgets.pause = createBaseSpeedButton("Pause", func() { clickPause(a) }, a.mapWidth)
	a.speedWidgets.max = createBaseSpeedButton("Disable max", func() { clickMax(a) }, a.mapWidth)
}

func (a *Anaxim) setSpeedToCustom() {
	a.speed = Custom
	a.speedWidgets.pause = createBaseSpeedButton("Pause", func() { clickPause(a) }, a.mapWidth)
	a.speedWidgets.max = createBaseSpeedButton("Enable max", func() { clickMax(a) }, a.mapWidth)
}

func (a *Anaxim) setSpeedToPaused() {
	a.speed = Paused
	a.speedWidgets.pause = createBaseSpeedButton("Resume", func() { clickPause(a) }, a.mapWidth)
	a.speedWidgets.max = createBaseSpeedButton("Enable max", func() { clickMax(a) }, a.mapWidth)
}

func clickPause(a *Anaxim) {
	switch a.speed {
	case Paused:
		a.setSpeedToUnlimited()
	case Custom:
		a.setSpeedToPaused()
	case Unlimited:
		a.setSpeedToPaused()
	}
	giu.Update()
}

func clickMax(a *Anaxim) {
	switch a.speed {
	case Paused:
		a.setSpeedToUnlimited()
	case Custom:
		a.setSpeedToUnlimited()
	case Unlimited:
		a.setSpeedToCustom()
	}
	giu.Update()
}

func snapPointToGrid(pt image.Point, gridSize int) image.Point {
	//Div floors
	return pt.Div(mapResize).Mul(mapResize).Add(image.Pt(gridSize, gridSize))
}

func (a *Anaxim) mapInputEvents() giu.Widget {
	return giu.Event().OnHover(func() {
		drawCursorPos := giu.GetCursorPos()
		mousePos := giu.GetMousePos()

		a.howeringOverCellCanvasPoint = snapPointToGrid(mousePos, mapResize)

		overImagePos := image.Point{
			X: mousePos.X - (drawCursorPos.X - a.mapWidth*mapResize) + 8, //magic padding number
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
