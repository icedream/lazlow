package lazlow

import (
	"strings"
	"time"

	"github.com/icedream/lazlow/effects"
	"github.com/sizeofint/webpanimation"
)

const (
	lazlowWebPEncoderOptionLossless = "lossless"
)

type webpEncoder struct {
}

func (encoder *webpEncoder) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{
		"lossless": effects.NewLazlowBoolOption("Lossless", "Enable lossless output", false),
	}
}

func (encoder *webpEncoder) SupportsFileExtension(ext string) bool {
	return strings.EqualFold(ext, ".webp") || strings.EqualFold(ext, ".webm")
}

func (encoder *webpEncoder) SupportsFrames(frameCount int) bool {
	return frameCount > 0
}

func (encoder *webpEncoder) Encode(frames []effects.LazlowFrame, out *LazlowOutput, options map[string]effects.LazlowOption) (err error) {
	// Create animated WebP out of generated frames
	webpanim := webpanimation.NewWebpAnimation(
		frames[0].Image.Bounds().Dx(),
		frames[0].Image.Bounds().Dy(),
		0)
	webpanim.WebPAnimEncoderOptions.SetKmin(9)
	webpanim.WebPAnimEncoderOptions.SetKmax(17)
	defer webpanim.ReleaseMemory()

	webpConfig := webpanimation.NewWebpConfig()
	if options[lazlowWebPEncoderOptionLossless].(*effects.LazlowBoolOption).TypedValue() {
		webpConfig.SetLossless(1)
	}

	var currentDuration time.Duration

	for _, frame := range frames {
		currentDuration += frame.Delay
		webpanim.AddFrame(
			frame.Image,
			int(currentDuration.Milliseconds()),
			webpConfig)
	}

	// Create output file
	outputFile, err := out.CreateFile()
	if err != nil {
		return
	}
	defer outputFile.Close()

	err = webpanim.Encode(outputFile)
	return
}

func init() {
	RegisterEncoder("webp", new(webpEncoder))
}
