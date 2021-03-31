package main

import (
	"time"

	"github.com/sizeofint/webpanimation"
)

var (
	flagWebPLossless = cli.Flag("webp-lossless", "Enable lossless output for WebP").Default("false").Bool()
)

type webpEncoder struct {
}

func (encoder *webpEncoder) Encode(frames []frame, out *output) (err error) {
	// Create animated WebP out of generated frames
	webpanim := webpanimation.NewWebpAnimation(
		frames[0].Image.Bounds().Dx(),
		frames[0].Image.Bounds().Dy(),
		0)
	webpanim.WebPAnimEncoderOptions.SetKmin(9)
	webpanim.WebPAnimEncoderOptions.SetKmax(17)
	defer webpanim.ReleaseMemory()

	webpConfig := webpanimation.NewWebpConfig()
	if *flagWebPLossless {
		webpConfig.SetLossless(1)
	}

	for i, frame := range frames {
		webpanim.AddFrame(
			frame.Image,
			int((*flagDelay * time.Duration(i)).Milliseconds()),
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
