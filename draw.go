// Drawing stuff on screen, and things like setting pixels of HumanGrid, Draw() is here
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontPressStart font.Face

func initFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	fontPressStart, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    8,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (w *HumanGrid) Draw(pix []byte) {
	for i, v := range w.area {
		if v.population != 0 {
			pix[4*i] = byte(254.0 * (math.Min(1.0, float64(v.population)/float64(upperPopCap))))
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0xff
		} else {
			pix[4*i] = 0
			pix[4*i+1] = 0
			pix[4*i+2] = 0
			pix[4*i+3] = 0
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.humanGrid.Draw(g.pixels)
	screen.WritePixels(g.pixels)

	text.Draw(screen, g.humanGrid.genCellInfoAtCursor(), fontPressStart, 1, 10, color.White)
	performanceInfoMsg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	text.Draw(screen, performanceInfoMsg, fontPressStart, screenWidth-80, screenHeight-8, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
