package main

import (
	"errors"
	"image"
	"time"
)

var errOnlySingleFrameOutputSupported = errors.New("this encoder only supports single-frame output")

type frame struct {
	Image image.Image
	Delay time.Duration
}

type effect interface {
	Process(image.Image) []frame
}

type encoder interface {
	Encode(frames []frame, out *output) (err error)
}
