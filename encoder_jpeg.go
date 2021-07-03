package lazlow

import (
	"image/jpeg"
	"strings"

	"github.com/icedream/lazlow/effects"
)

const (
	lazlowJPEGEncoderOptionQuality = "quality"
)

type jpegEncoder struct {
}

func (encoder *jpegEncoder) SupportsFrames(frameCount int) bool {
	return frameCount == 1
}

func (encoder *jpegEncoder) SupportsFileExtension(ext string) bool {
	return strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg")
}

func (encoder *jpegEncoder) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{
		lazlowJPEGEncoderOptionQuality: effects.NewLazlowEncoderIntegerOption("JPEG quality level", "The quality level with which JPEG output will be written, where 100 = lossless and lower will be increasingly lossy.", 90, 0, 100, 1),
	}
}

func (encoder *jpegEncoder) Encode(frames []effects.LazlowFrame, out *LazlowOutput, options map[string]effects.LazlowOption) (err error) {
	// TODO - just a safety check that can be removed later once the plugin framework does the check itself
	if !encoder.SupportsFrames(len(frames)) {
		err = ErrOnlySingleFrameOutputSupported
		return
	}

	f, err := out.CreateFile()
	if err != nil {
		return
	}
	defer f.Close()

	err = jpeg.Encode(f, frames[0].Image, &jpeg.Options{
		Quality: int(options[lazlowJPEGEncoderOptionQuality].(*effects.LazlowIntegerOption).TypedValue()),
	})
	if err != nil {
		return
	}

	return
}

func init() {
	RegisterEncoder("jpeg", new(jpegEncoder))
}
