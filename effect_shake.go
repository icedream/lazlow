package main

import (
	"image"
	"image/draw"
	"math/rand"
)

var (
	flagShakeFrameCount = cli.Flag("shake-frames", "Shake effect frame count").Default("12").Int()
	flagShakePercentage = cli.Flag("shake", "How much to shake the picture in percent").Default("20").Uint8()
)

type shakeEffect struct {
}

func (effect *shakeEffect) Process(inputImage image.Image) (images []frame) {
	images = make([]frame, *flagShakeFrameCount)

	// Add initial frame (always unshaken for preview purposes)
	images[0] = frame{inputImage, *flagDelay}

	shakeX := int(float64(inputImage.Bounds().Dx()) * float64(*flagShakePercentage) / 100)
	shakeY := int(float64(inputImage.Bounds().Dy()) * float64(*flagShakePercentage) / 100)

	// Create shaken frames
	for i := 1; i < *flagShakeFrameCount; i++ {
		img := image.NewRGBA(
			image.Rect(
				0,
				0,
				inputImage.Bounds().Dx(),
				inputImage.Bounds().Dy()))
		x := rand.Intn(shakeX) - (shakeX / 2)
		y := rand.Intn(shakeY) - (shakeY / 2)

		draw.Draw(img, inputImage.Bounds(), inputImage, image.Point{x, y}, draw.Src)

		images[i] = frame{img, *flagDelay}
	}

	return
}
