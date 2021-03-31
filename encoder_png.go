package main

import "image/png"

type pngEncoder struct {
}

func (encoder *pngEncoder) Encode(frames []frame, out *output) (err error) {
	if len(frames) != 1 {
		err = errOnlySingleFrameOutputSupported
		return
	}

	f, err := out.CreateFile()
	if err != nil {
		return
	}
	defer f.Close()

	err = png.Encode(f, frames[0].Image)
	if err != nil {
		return
	}

	return
}
