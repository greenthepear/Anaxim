// Drawing stuff on screen, and things like setting pixels of HumanGrid, Draw() is here
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontPressStart font.Face

func getPixels(path string) ([]color.Color, int, int) {
	_, img, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	r := make([]color.Color, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r[y*width+x] = img.At(x, y)
		}
	}
	return r, width, height
}

func initFonts() { //global, used in init
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

func makeImagesMap() map[string]*ebiten.Image {
	images := make(map[string]*ebiten.Image)
	imageNames := []string{"speed0", "speed1", "speed2", "speed3"}
	for _, n := range imageNames {
		img, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("./Graphics/%v.png", n))
		if err != nil {
			log.Fatal(err)
		}
		images[n] = img
	}
	return images
}

func (g *Game) DrawGrid(pix []byte) {
	for i, v := range g.humanGrid.area {
		var landValue byte = 0
		if g.mapGrid.area[i].isLand {
			landValue = 255
		}

		popRange := float64(v.population) / float64(upperPopCap)
		if popRange != 0 { //To see places where there is ANY population
			popRange += 0.1
		}
		red := byte(254.0 * (math.Min(1.0, popRange)))
		pix[4*i] = red
		pix[4*i+1] = landValue - red
		pix[4*i+2] = 100 - landValue
		pix[4*i+3] = 0
	}
}

func drawTextWithDropShadow(destination *ebiten.Image, contents string, face font.Face, x, y int, clr color.Color) {
	text.Draw(destination, contents, face, x+1, y+1, color.Black)
	text.Draw(destination, contents, face, x, y, clr)
}

func (g Game) drawSpeedControls(screen *ebiten.Image) {
	speedControlImgOp := &ebiten.DrawImageOptions{}
	speedControlImageKeys := []string{"speed0", "speed1", "speed2", "speed3"}

	var speedControlx float64 = 0
	for i, k := range speedControlImageKeys {
		speedControlImgOp.ColorScale.Reset()
		speedControlImgOp.GeoM.Reset()
		speedControlImgOp.GeoM.Translate(speedControlx, float64(screenHeight-16))
		if int(g.speed) != i {
			speedControlImgOp.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(g.images[k], speedControlImgOp)
		speedControlx += 16
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.DrawGrid(g.pixels)
	screen.WritePixels(g.pixels)

	drawTextWithDropShadow(screen, g.humanGrid.genCellInfoAtCursor(), fontPressStart, 1, 10, color.White)
	performanceInfoMsg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	drawTextWithDropShadow(screen, performanceInfoMsg, fontPressStart, screenWidth-80, screenHeight-8, color.White)

	g.drawSpeedControls(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
