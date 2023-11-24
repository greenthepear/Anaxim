package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// "example" -> "[ example ]"
func enbracked(str string) string {
	return fmt.Sprintf("[ %s ]", str)
}

func (a *Anaxi) genSpeedButtonFunction(button *widget.Button, speed Speed, symbols map[Speed]string) func() {
	return func() {
		a.speed = speed
		for i, b := range a.speedButtons {
			if i != int(a.speed) {
				b.SetText(symbols[Speed(i)])
			} else {
				b.SetText(enbracked(symbols[Speed(i)]))
			}
		}
	}
}

func (a *Anaxi) genSpeedControls() []*widget.Button {
	speedSymbols := map[Speed]string{
		Paused:  "⏸",
		Slow:    "⏩",
		Faster:  "⏩⏩",
		Fastest: "⏩⏩⏩",
	}
	buttons := make([]*widget.Button, 0, len(speedSymbols))
	var butt *widget.Button
	for speed, symbol := range speedSymbols {
		butt = widget.NewButton(symbol, a.genSpeedButtonFunction(butt, speed, speedSymbols))
		buttons = append(buttons, butt)
	}
	return buttons
}

func (a *Anaxi) buildSpeedControls() fyne.CanvasObject {
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
