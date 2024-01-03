// Drawing the Fyne UI
package main

import (
	"github.com/AllenDang/giu"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type SpeedWidgets struct {
	pause  *giu.ButtonWidget
	slider *giu.SliderIntWidget
	max    *giu.ButtonWidget
}

func createButtonWithSize(label string, onclick func(), width, height float32) *giu.ButtonWidget {
	butt := giu.Button(label).OnClick(onclick)
	butt.Size(width, height)
	return butt
}

func createBaseSpeedButton(label string, onlick func()) *giu.ButtonWidget {
	return createButtonWithSize(label, onlick, float32(mapWidth*mapResize/8), 20)
}

func (a *Anaxi) setSpeedToUnlimited() {
	a.speed = Unlimited
	a.speedWidgets.pause = createBaseSpeedButton("Pause", func() { clickPause(a) })
	a.speedWidgets.max = createBaseSpeedButton("Disable max", func() { clickMax(a) })
}

func (a *Anaxi) setSpeedToCustom() {
	a.speed = Custom
	a.speedWidgets.pause = createBaseSpeedButton("Pause", func() { clickPause(a) })
	a.speedWidgets.max = createBaseSpeedButton("Enable max", func() { clickMax(a) })
}

func (a *Anaxi) setSpeedToPaused() {
	a.speed = Paused
	a.speedWidgets.pause = createBaseSpeedButton("Resume", func() { clickPause(a) })
	a.speedWidgets.max = createBaseSpeedButton("Enable max", func() { clickMax(a) })
}

func clickPause(a *Anaxi) {
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

func clickMax(a *Anaxi) {
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

func (a *Anaxi) initUI() {
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

func NewSpeedWidgets(a *Anaxi) *SpeedWidgets {
	slider := giu.SliderInt(&a.speedCustomTPS, 1, 100).OnChange(func() { a.setSpeedToCustom() })
	slider.Size(float32(mapWidth*mapResize) * 0.74) //makes it roughly fit but should find a better way for this
	return &SpeedWidgets{
		createBaseSpeedButton("Pause", func() { clickPause(a) }),
		slider,
		createBaseSpeedButton("Enable max", func() { clickMax(a) }),
	}
}

func (a *Anaxi) genSpeedText() string {
	str := "Speed controls. Currently "
	suffixes := map[Speed]string{
		Paused:    "paused",
		Custom:    "on custom (slider) speed.",
		Unlimited: "on unlimited (max) speed.",
	}
	return str + suffixes[a.speed]
}

func (a *Anaxi) genGlobalStatsString() string {
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

func (a *Anaxi) createLayout() {
	img := giu.Image(a.mapTexture)
	img.Size(float32(a.mapImage.Bounds().Dx()*mapResize), float32(a.mapImage.Bounds().Dy()*mapResize))

	statLabel := giu.Label(a.genGlobalStatsString())
	mapAndSpeedCol := giu.Column(
		giu.Row(
			giu.Label(a.genSpeedText()),
		),
		giu.Row(
			a.speedWidgets.pause,
			a.speedWidgets.slider,
			a.speedWidgets.max,
		),
		giu.Row(
			img,
		),
	)

	giu.SingleWindow().Layout(
		giu.Row(
			giu.Column(
				statLabel,
			),
			mapAndSpeedCol,
		),
	)
}
