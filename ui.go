// Drawing the Fyne UI
package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *Anaxi) initUI() {
	a.speedButtons = a.genSpeedControls()
	a.mapCanvas = canvas.NewRaster(a.updateGridImage)
	a.mapCanvas.ScaleMode = canvas.ImageScalePixels
}

// "example" -> "[ example ]"
func enbracked(str string) string {
	return fmt.Sprintf("[ %s ]", str)
}

type speedSymbol struct { //Maps are iterated in random order, so I use this instead
	spd    Speed
	symbol string
}

func (a *Anaxi) genSpeedButtonFunction(button *widget.Button, assignedSpd Speed, symbols []speedSymbol) func() {
	return func() {
		a.speed = assignedSpd
		for i, b := range a.speedButtons {
			if i != int(a.speed) {
				b.SetText(symbols[i].symbol)
			} else {
				b.SetText(enbracked(symbols[i].symbol))
			}
		}
	}
}

func (a *Anaxi) genSpeedControls() []*widget.Button {
	speedSymbols := []speedSymbol{
		{Paused, "⏸"},
		{Slow, "⏩"},
		{Faster, "⏩⏩"},
		{Unlimited, "⏩⏩⏩"},
	}
	buttons := make([]*widget.Button, 0, len(speedSymbols))
	var butt *widget.Button
	for _, e := range speedSymbols {
		butt = widget.NewButton(e.symbol, a.genSpeedButtonFunction(butt, e.spd, speedSymbols))
		buttons = append(buttons, butt)
	}
	return buttons
}

func (a *Anaxi) buildSpeedControls() fyne.CanvasObject {
	//Set current speed button to be visually selected
	a.speedButtons[int(a.speed)].SetText(enbracked(a.speedButtons[int(a.speed)].Text)) //What is this, JavaScript code?
	return container.NewGridWithColumns(
		len(a.speedButtons),
		//Bad. No clue why I can't just do `a.speedButtons...`
		a.speedButtons[0],
		a.speedButtons[1],
		a.speedButtons[2],
		a.speedButtons[3],
	)
}

func (a *Anaxi) buildUI() fyne.CanvasObject {
	return container.NewBorder(a.buildSpeedControls(), nil, nil, nil, a.mapCanvas)
}
