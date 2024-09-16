// Drawing the giu UI
package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/AllenDang/giu"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const leftColumnWidth = 200

type SpeedWidgets struct {
	pause  *giu.ButtonWidget
	slider *giu.SliderFloatWidget
	max    *giu.ButtonWidget
}

func (a *Anaxim) initUI() {
	a.speedWidgets = NewSpeedWidgets(a)

	switch a.speed {
	case Paused:
		a.setSpeedToPaused()
	case Custom:
		a.setSpeedToCustom()
	case Unlimited:
		a.setSpeedToUnlimited()
	}
}

func (a *Anaxim) NewBaseSpeedButton(label string, onClick func(*Anaxim)) *giu.ButtonWidget {
	return giu.Button(label).OnClick(func() { onClick(a) }).Size(100, 20)
}

func (a *Anaxim) PauseButton(label string) *giu.ButtonWidget {
	return a.NewBaseSpeedButton(label, clickPause)
}

func (a *Anaxim) MaxButton(label string) *giu.ButtonWidget {
	return a.NewBaseSpeedButton(label, clickMax)
}

func NewSpeedWidgets(a *Anaxim) *SpeedWidgets {
	return &SpeedWidgets{
		a.PauseButton("Pause"),
		giu.SliderFloat(&a.speedCustomTPS, 0, 1).
			OnChange(func() { a.setSpeedToCustom() }).
			Size(300),
		a.MaxButton("Enable max"),
	}
}

func (a *Anaxim) genSpeedText() string {
	return "Speed controls. Currently " + map[Speed]string{
		Paused:    "paused",
		Custom:    "on custom (slider) speed.",
		Unlimited: "on unlimited (max) speed.",
	}[a.speed]
}

func (a *Anaxim) genGlobalStatsString() string {
	biggestCell := a.simulation.humanGrid.biggestPopCell
	printer := message.NewPrinter(language.English)
	return printer.Sprintf(
		`Simulation generation:
%d
			
World population:
%d
Biggest population:
%d @ (%d,%d)`,
		a.simulation.humanGrid.generation,
		a.simulation.humanGrid.globalPop,
		biggestCell.population, biggestCell.x, biggestCell.y,
	)
}

func (a *Anaxim) genLocalStatsString() string {
	ic := a.inspectingCell
	if ic == nil {
		return "Click on the map to inspect."
	}
	printer := message.NewPrinter(language.English)
	dev := fmt.Sprint(ic.development)
	if dev == "1" {
		dev += " (base)"
	}

	prefix := fmt.Sprintf(
		`Inspecting cell @ %v

Map info:`, a.inspectingCellAt.String())

	mc := ic.mapCell
	if !mc.isLand {
		return prefix + "\n\tWater cell"
	}
	mapInfo := fmt.Sprintf("Land cell\n\tHabitability:\n\t%v", mc.habitability)

	return printer.Sprintf(
		`%v
	%v

Human activity:
	Population:
	%d
	Development:
	%v`,
		prefix,
		mapInfo,
		ic.population, dev,
	)
}

func (a *Anaxim) createLayout() {
	img := giu.Image(a.mapTexture)
	img.Size(
		float32(a.mapImage.Bounds().Dx()*mapResize),
		float32(a.mapImage.Bounds().Dy()*mapResize))

	giu.SingleWindow().Layout(
		giu.Row(
			giu.Column(
				giu.Row(
					giu.Label(a.genGlobalStatsString()).Wrapped(true),
				),
				giu.Row(
					giu.Dummy(leftColumnWidth, 30),
				),
				giu.Row(
					giu.Label(a.genLocalStatsString()).Wrapped(true),
				),
			),
			giu.Column(
				giu.Row(
					giu.Label(a.genSpeedText()),
				),
				giu.Row(
					a.speedWidgets.pause,
					a.speedWidgets.max,
					a.speedWidgets.slider,
				),
				giu.Row(
					img,
					a.mapInputEvents(),
					giu.Custom(func() {
						can := giu.GetCanvas()

						//Inspecting map cursor
						if a.inspectingCell != nil {
							can.AddCircle(
								a.inspectingCanvasPoint.Add(image.Pt(-4, -1)),
								float32(mapResize),
								color.RGBA{0xff, 0xff, 0xff, 200}, 4, 2)
						}

						//Howering map cursor
						can.AddCircle(a.howeringOverCellCanvasPoint.Add(image.Pt(-4, -1)),
							float32(mapResize*2),
							color.RGBA{0xff, 0xff, 0xff, 100}, 4, 2)
					}),
				),
				giu.Row(
					giu.Dummy(0, 0),
				),
				giu.Row(
					giu.Label(a.howeringOverCellAt.String()),
				),
			),
		),
	)
}
