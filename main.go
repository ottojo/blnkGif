package main

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/ottojo/blnk2"
	"github.com/ottojo/blnk2/client"
	"github.com/ottojo/blnk2/vector"
	"image"
	"image/gif"
	"math"
	"os"
	"time"
)

func main() {

	f, _ := os.Open(os.Args[2])
	g, _ := gif.DecodeAll(f)

	blnksystem := blnk2.CreateFromFile(os.Args[1])
	go blnksystem.Discovery()
	defer blnksystem.Disconnect()
	time.Sleep(7 * time.Second)

	for true {
		for i, im := range g.Image {
			RenderBitmap(im, &blnksystem.Stage, vector.Vec3{0.5, -5, 5},
				vector.Vec3{0, 1, -.5},
				0.5*math.Pi)
			go blnksystem.Commit()
			time.Sleep(time.Duration(g.Delay[i]*10) * time.Millisecond)
		}
	}
}

// Right handed coordinate system, Z=UP
// TODO: Find out why image is mirrored
func RenderBitmap(bmp image.Image, list *client.LedList, pStart vector.Vec3, pDir vector.Vec3, horizontalAngle float64) {

	bmp.Bounds()

	imageWidth := bmp.Bounds().Max.X - bmp.Bounds().Min.X
	imageHeight := bmp.Bounds().Max.Y - bmp.Bounds().Min.Y

	w := 2 * math.Tan(horizontalAngle/2)
	h := w * (float64(imageHeight) / float64(imageWidth))

	pixel := list.First
	for pixel != nil {
		pixelPos := pixel.Data.Position

		pLoc := pixelPos.Minus(pStart)

		phi := pLoc.Phi() - pDir.Phi()

		if phi < 0 {
			phi += 2 * math.Pi
		}

		plTheta := pLoc.Theta()
		pdTheta := pDir.Theta()
		theta := plTheta - pdTheta
		if theta < 0 {
			theta += math.Pi
		}

		x := (w / 2) - math.Tan(phi)
		x = float64(imageWidth) * (x / w)

		y := (h / 2) - math.Tan(theta)
		y = float64(imageHeight) * (y / h)

		if x < 0 || y < 0 || int(x) >= imageWidth || int(y) >= imageHeight {
			pixel = pixel.Next
			continue
		}
		r, g, b, _ := bmp.At(int(x), int(y)).RGBA()
		red := float64(r) / float64(0xffff)
		green := float64(g) / float64(0xffff)
		blue := float64(b) / float64(0xffff)

		pixel.Data.Color = colorful.Color{R: red, G: green, B: blue}
		pixel = pixel.Next
	}
}
