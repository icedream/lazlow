package main

import (
	"time"

	"github.com/kettek/apng"
)

type apngEncoder struct {
}

func (encoder *apngEncoder) Encode(frames []frame, out *output) (err error) {
	a := apng.APNG{
		Frames: make([]apng.Frame, len(frames)),
	}
	denominator := 100
	for i, frame := range frames {
		a.Frames[i].Image = frame.Image
		a.Frames[i].DelayNumerator = uint16(frame.Delay/(time.Second/time.Duration(denominator))) + 1
		a.Frames[i].DelayDenominator = uint16(denominator) // 1/100th of a second
		a.Frames[i].DisposeOp = 1                          // APNG_DISPOSE_OP_BACKGROUND
	}

	f, err := out.CreateFile()
	if err != nil {
		return
	}
	defer f.Close()

	err = apng.Encode(f, a)
	if err != nil {
		return
	}

	return
}
