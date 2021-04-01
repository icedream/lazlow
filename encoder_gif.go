package lazlow

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"strings"

	"github.com/esimov/colorquant"
)

var dither = colorquant.Dither{
	Filter: [][]float32{
		{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
		{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
		{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
	},
}

func quantify(src image.Image, numColors int) (dst *image.Paletted, quant image.Image) {
	dst = image.NewPaletted(
		image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()),
		palette.Plan9)
	// if noDither {
	// 	quant = colorquant.NoDither.Quantize(src, dst, numColors, false, true)
	// } else {
	quant = dither.Quantize(src, dst, len(dst.Palette), true, true)
	// }

	return
}

type gifEncoder struct {
}

func (encoder *gifEncoder) SupportsFileExtension(ext string) bool {
	return strings.EqualFold(ext, ".gif")
}

func (encoder *gifEncoder) SupportsFrames(frameCount int) bool {
	return frameCount > 0
}

func (encoder *gifEncoder) Options() map[string]LazlowOption {
	return map[string]LazlowOption{}
}

func (encoder *gifEncoder) Encode(frames []LazlowFrame, out *LazlowOutput, options map[string]LazlowOption) (err error) {
	quantizedImage, _ := quantify(frames[0].Image, 256)

	var images []*image.Paletted
	var delays []int
	var disposal []byte

	// Create shaken frames
	for _, frame := range frames {
		img := image.NewPaletted(
			image.Rect(
				0,
				0,
				frame.Image.Bounds().Dx(),
				frame.Image.Bounds().Dy()),
			quantizedImage.Palette)
		draw.Draw(img, frame.Image.Bounds(), frame.Image, image.Point{0, 0}, draw.Src)

		// images = append(images, img)
		images = append(images, img)
		delays = append(delays, delays[0])
		disposal = append(disposal, disposal[0])
	}

	// Create output file
	outputFile, err := out.CreateFile()
	if err != nil {
		return
	}
	defer outputFile.Close()

	// Create animated GIF out of generated frames
	err = gif.EncodeAll(outputFile, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposal,
		Config: image.Config{
			ColorModel: images[0].ColorModel(),
			Width:      images[0].Rect.Dx(),
			Height:     images[0].Rect.Dy(),
		},
	})
	return
}

func init() {
	RegisterEncoder("gif", new(gifEncoder))
}
