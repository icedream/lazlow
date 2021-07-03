package lazlow

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"strings"
	"time"

	"github.com/esimov/colorquant"
	"github.com/icedream/lazlow/effects"
)

var dither = colorquant.Dither{
	Filter: [][]float32{
		{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
		{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
		{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
	},
}

func quantify(src image.Image, palette color.Palette) (dst *image.Paletted) {
	dst = image.NewPaletted(
		image.Rect(0, 0, src.Bounds().Dx(), src.Bounds().Dy()),
		palette)
	// if noDither {
	// 	colorquant.NoDither.Quantize(src, dst, numColors, false, false)
	// } else {
	dither.Quantize(src, dst, 0 /*unused*/, true, false)
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

func (encoder *gifEncoder) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{}
}

func (encoder *gifEncoder) Encode(frames []effects.LazlowFrame, out *LazlowOutput, options map[string]effects.LazlowOption) (err error) {
	var images []*image.Paletted
	var delays []int
	var disposal []byte

	// Create color model/palette for all frames
	var totalWidth int
	var totalHeight int
	for _, frame := range frames {
		totalWidth += frame.Image.Bounds().Dx()
		if frame.Image.Bounds().Dy() > totalHeight {
			totalHeight = frame.Image.Bounds().Dy()
		}
	}
	colorPaletteSrc := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	colorPaletteX := 0
	for _, frame := range frames {
		draw.Draw(colorPaletteSrc,
			frame.Image.Bounds(),
			frame.Image,
			image.Point{colorPaletteX, 0},
			draw.Src)
	}
	// NOTE - we reserve one color out of the 256 for transparency
	// TODO - allow multiple transparency levels
	colorPaletteImg := colorquant.Quant{}.Quantize(colorPaletteSrc, 255)
	palette := color.Palette{
		color.Transparent,
	}
	palette = append(palette, colorPaletteImg.(*image.Paletted).Palette...)

	// Create shaken frames
	for _, frame := range frames {
		quantizedImage := quantify(frame.Image, palette)
		img := image.NewPaletted(
			image.Rect(
				0,
				0,
				frame.Image.Bounds().Dx(),
				frame.Image.Bounds().Dy()),
			quantizedImage.Palette)
		draw.Draw(img, quantizedImage.Bounds(), quantizedImage, image.Point{0, 0}, draw.Src)

		// images = append(images, img)
		images = append(images, img)
		delays = append(delays, int(frame.Delay/(10*time.Millisecond)))
		disposal = append(disposal, gif.DisposalBackground)
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
