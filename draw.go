// Drawing stuff on screen, and things like setting pixels of HumanGrid, Draw() is here
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var fontPressStart font.Face

func getPixels(file *os.File) ([]color.Color, int, int) {
	img, err := png.Decode(file)
	if err != nil {
		log.Fatalf("getPixels() | %v", err)
	}
	file.Close()
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	r := make([]color.Color, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r[y*width+x] = img.At(x, y)
		}
	}
	return r, width, height
}

func loadImage(path string) *os.File {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Image not found")
	}

	return file
}

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

func (g *Game) DrawGrid(pix []byte) {
	for i, v := range g.humanGrid.area {
		var landValue byte = 0
		if g.mapGrid.area[i].isLand {
			landValue = 255
		}

		popRange := float64(v.population) / float64(upperPopCap)
		red := byte(254.0 * (math.Min(1.0, popRange)))
		pix[4*i] = red
		pix[4*i+1] = landValue - red
		pix[4*i+2] = 0
		pix[4*i+3] = 0
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.DrawGrid(g.pixels)
	screen.WritePixels(g.pixels)

	text.Draw(screen, g.humanGrid.genCellInfoAtCursor(), fontPressStart, 1, 10, color.White)
	performanceInfoMsg := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\n", ebiten.ActualTPS(), ebiten.ActualFPS())
	text.Draw(screen, performanceInfoMsg, fontPressStart, screenWidth-80, screenHeight-8, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
