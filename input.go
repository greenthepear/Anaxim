// Handling of most user input,
// excluding only ones that don't need abstraction out of Layout in ui.go
package main

import (
	"image"

	"github.com/AllenDang/giu"
)

func (a *Anaxim) unpauseInto(s Speed) {
	previousSpeed := a.speed
	a.speed = s
	if previousSpeed == Paused {
		go a.runSim()
	}
}

func (a *Anaxim) setSpeedToUnlimited() {
	a.unpauseInto(Unlimited)
	a.speedWidgets.pause = a.PauseButton("Pause")
	a.speedWidgets.max = a.MaxButton("Disable max")
}

func (a *Anaxim) setSpeedToCustom() {
	a.unpauseInto(Custom)
	a.speedWidgets.pause = a.PauseButton("Pause")
	a.speedWidgets.max = a.MaxButton("Enable max")
}

func (a *Anaxim) setSpeedToPaused() {
	a.speed = Paused
	a.speedWidgets.pause = a.PauseButton("Resume")
	a.speedWidgets.max = a.MaxButton("Enable max")
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
	return image.Pt(
		pt.X/mapResize*mapResize+gridSize,
		pt.Y/mapResize*mapResize+gridSize,
	)
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
		a.inspectingCanvasPoint = a.howeringOverCellCanvasPoint
		a.inspectingCell = a.simulation.humanGrid.CellAt(
			a.inspectingCellAt.X, a.inspectingCellAt.Y)
		a.updateMapTexture()
	})
}
