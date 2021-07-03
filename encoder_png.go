package lazlow

import (
	"image/png"
	"strings"

	"github.com/icedream/lazlow/effects"
)

type pngEncoder struct {
}

func (encoder *pngEncoder) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{}
}

func (encoder *pngEncoder) SupportsFileExtension(ext string) bool {
	return strings.EqualFold(ext, ".png")
}

func (encoder *pngEncoder) SupportsFrames(frameCount int) bool {
	return frameCount == 1
}

func (encoder *pngEncoder) Encode(frames []effects.LazlowFrame, out *LazlowOutput, options map[string]effects.LazlowOption) (err error) {
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

	err = png.Encode(f, frames[0].Image)
	if err != nil {
		return
	}

	return
}

func init() {
	RegisterEncoder("png", new(pngEncoder))
}
