package main

import (
	"image/jpeg"
)

var flagJPEGQuality = cli.Flag("jpeg-quality", "The quality with which JPEG files will be written").Default("90").Int()

type jpegEncoder struct {
}

func (encoder *jpegEncoder) Encode(frames []frame, out *output) (err error) {
	if len(frames) != 1 {
		err = errOnlySingleFrameOutputSupported
		return
	}

	f, err := out.CreateFile()
	if err != nil {
		return
	}
	defer f.Close()

	err = jpeg.Encode(f, frames[0].Image, &jpeg.Options{
		Quality: *flagJPEGQuality,
	})
	if err != nil {
		return
	}

	return
}
