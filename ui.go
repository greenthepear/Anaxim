// Drawing the Fyne UI
package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SpeedWidgets struct {
	widget.BaseWidget

	pauseButton     *widget.Button
	speedSlider     *widget.Slider
	unlimitedButton *widget.Button
}

func (sw *SpeedWidgets) whilePausedSetTheme() {
	sw.pauseButton.SetIcon(theme.MediaPlayIcon())
	sw.pauseButton.SetText("Resume")

	sw.unlimitedButton.SetIcon(theme.MediaFastForwardIcon())
	sw.unlimitedButton.SetText("Paused\nClick to enable max speed")
}

func (sw *SpeedWidgets) whileUnlimitedSetTheme() {
	sw.unlimitedButton.SetIcon(theme.MediaPlayIcon())
	sw.unlimitedButton.SetText("On max speed\nClick to enable slider speed")

	sw.pauseButton.SetIcon(theme.MediaPauseIcon())
	sw.pauseButton.SetText("Pause")
}

func (sw *SpeedWidgets) whileCustomSetTheme() {
	sw.unlimitedButton.SetIcon(theme.MediaFastForwardIcon())
	sw.unlimitedButton.SetText("On slider speed\nClick to enable max speed")

	sw.pauseButton.SetIcon(theme.MediaPauseIcon())
	sw.pauseButton.SetText("Pause")
}

func whenTappedPause(a *Anaxi) func() {
	return func() {
		switch {
		case a.speed == Paused:
			a.speed = Custom
			a.speedWidgets.whileCustomSetTheme()
			a.speedWidgets.speedSlider.SetValue(50.0)
		case a.speed == Custom:
			a.speed = Paused
			a.speedWidgets.whilePausedSetTheme()
		case a.speed == Unlimited:
			a.speed = Paused
			a.speedWidgets.whilePausedSetTheme()
		}
	}
}

func whenTappedUnlimited(a *Anaxi) func() {
	return func() {
		switch {
		case a.speed == Paused:
			a.speed = Unlimited
			a.speedWidgets.whileUnlimitedSetTheme()
		case a.speed == Custom:
			a.speed = Unlimited
			a.speedWidgets.whileUnlimitedSetTheme()
		case a.speed == Unlimited:
			a.speed = Custom
			a.speedWidgets.whileCustomSetTheme()
		}
	}
}

func whenSpeedSliderDragEnd(a *Anaxi) func(float64) {
	return func(float64) {
		val := a.speedWidgets.speedSlider.Value
		if val == 0 {
			a.speed = Paused
			a.speedWidgets.whilePausedSetTheme()
			return
		}
		a.speed = Custom
		a.speedWidgets.whileCustomSetTheme()
		a.speedCustomTPS = time.Duration(time.Second / time.Duration(val))
	}
}

func NewSpeedWidgets(a *Anaxi) *SpeedWidgets {
	sw := &SpeedWidgets{
		pauseButton: widget.NewButtonWithIcon(
			"Pause", theme.MediaPauseIcon(), whenTappedPause(a)),
		unlimitedButton: widget.NewButtonWithIcon(
			"Paused\nClick to enable max speed", theme.MediaFastForwardIcon(), whenTappedUnlimited(a)),
		speedSlider: widget.NewSlider(0.0, 100.0),
	}

	sw.speedSlider.Step = 5.0
	sw.speedSlider.OnChanged = whenSpeedSliderDragEnd(a)

	return sw
}

func (a *Anaxi) initUI() {
	a.speedWidgets = NewSpeedWidgets(a)
	a.mapCanvas = canvas.NewRaster(a.updateGridImage)
	a.mapCanvas.ScaleMode = canvas.ImageScalePixels
}

func (a *Anaxi) buildSpeedControls() fyne.CanvasObject {
	return container.NewGridWithColumns(
		3,
		a.speedWidgets.pauseButton,
		a.speedWidgets.speedSlider,
		a.speedWidgets.unlimitedButton,
	)
}

func (a *Anaxi) buildUI() fyne.CanvasObject {
	return container.NewBorder(a.buildSpeedControls(), nil, nil, nil, a.mapCanvas)
}
