package pat

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"time"

	"github.com/icedream/lazlow/effects"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/nfnt/resize"
)

var (
	colorTransparent = color.RGBA{0, 0, 0, 0xff}

	heightPercentages = []float64{
		// TODO - use sin() for this?
		1.0,
		0.85,
		0.75,
		0.85,
		0.99,
	}
)

func getFrameFiles() (files []pkging.File, err error) {
	for _, name := range []string{
		pkger.Include("/effects/pat/assets/frame0.gif"),
		pkger.Include("/effects/pat/assets/frame1.gif"),
		pkger.Include("/effects/pat/assets/frame2.gif"),
		pkger.Include("/effects/pat/assets/frame3.gif"),
		pkger.Include("/effects/pat/assets/frame4.gif"),
	} {
		var file pkging.File
		file, err = pkger.Open(name)
		if err != nil {
			return
		}
		files = append(files, file)
	}
	return
}

type LazlowPatEffect struct {
}

func (effect *LazlowPatEffect) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{}
}

func (effect *LazlowPatEffect) IsAnimated() bool {
	return true
}

func (effect *LazlowPatEffect) Process(inputImage image.Image, options map[string]effects.LazlowOption) (images []effects.LazlowFrame, err error) {
	images = make([]effects.LazlowFrame, 0)

	frames, err := getFrameFiles()
	if err != nil {
		log.Fatalf("Failed to getFrameFiles: %s", err.Error())
		return
	}

	for frameIndex, f := range frames {
		defer func(f pkging.File) { f.Close() }(f)

		log.Println("PAT > Frame", frameIndex)

		// get the original handpat frame
		var patFrame image.Image
		patFrame, err = gif.Decode(f)
		if err != nil {
			log.Fatalf("Failed to decode gif: %s", err.Error())
			return
		}

		// squish the original frame
		newInputImageHeight := float64(inputImage.Bounds().Dy()) *
			heightPercentages[frameIndex]
		// imageHeightDifference := float64(inputImage.Bounds().Dy()) - newInputImageHeight
		newInputImageWidth := (1 / heightPercentages[frameIndex]) *
			float64(inputImage.Bounds().Dx())

		// TODO - squish
		resizedImage := resize.Resize(
			uint(newInputImageWidth),
			uint(newInputImageHeight),
			inputImage,
			resize.NearestNeighbor)
		offsetX := (inputImage.Bounds().Dx() - resizedImage.Bounds().Dx()) / -2
		offsetY := inputImage.Bounds().Dy() - resizedImage.Bounds().Dy()
		log.Printf("  offsetY = %d", offsetY)

		// overlay on top of original image
		finalFrame := image.NewRGBA(inputImage.Bounds())
		draw.Draw(finalFrame,
			resizedImage.Bounds(),
			resizedImage,
			image.Pt(offsetX, -offsetY),
			draw.Src) // input image
		draw.Draw(finalFrame,
			patFrame.Bounds(),
			patFrame,
			image.Pt(0, 0),
			draw.Over) // handpat overlay

		images = append(images, effects.LazlowFrame{
			Image: finalFrame,
			Delay: time.Duration(frameIndex) * 10 * time.Millisecond,
		})

	}

	log.Println("PAT > Images count", len(images))

	return
}

func init() {
	effects.RegisterEffect("pat", new(LazlowPatEffect))
}
