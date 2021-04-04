package lazlow

import (
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/esimov/caire"
	"github.com/nfnt/resize"
)

const (
	lazlowContentAwareScalingEffectOptionDelay  = "delay"
	lazlowContentAwareScalingEffectOptionFrames = "frames"
	lazlowContentAwareScalingEffectOptionStep   = "step"
)

type LazlowContentAwareScalingEffect struct {
}

func (effect *LazlowContentAwareScalingEffect) IsAnimated() bool {
	return true
}

func (effect *LazlowContentAwareScalingEffect) Options() map[string]LazlowOption {
	return map[string]LazlowOption{
		lazlowContentAwareScalingEffectOptionDelay:  NewLazlowEncoderDurationOption("Frame delay", "Delay between frames", 40*time.Millisecond, 0, 0xff*10*time.Millisecond, 10*time.Millisecond),
		lazlowContentAwareScalingEffectOptionFrames: NewLazlowEncoderIntegerOption("Frame count", "How many frames to generate", 15*2, 0, 0xff, 1),
		lazlowContentAwareScalingEffectOptionStep:   NewLazlowEncoderIntegerOption("Pixel delta step", "By how many pixels to scale with each frame", 4, 0, 0xffff, 1),
	}
}

func (effect *LazlowContentAwareScalingEffect) Process(inputImage image.Image, options map[string]LazlowOption) (images []LazlowFrame) {
	delay := options[lazlowContentAwareScalingEffectOptionDelay].(*LazlowDurationOption).TypedValue()
	frameCount := int(options[lazlowContentAwareScalingEffectOptionFrames].(*LazlowIntegerOption).TypedValue())
	pixelDelta := int(options[lazlowContentAwareScalingEffectOptionStep].(*LazlowIntegerOption).TypedValue())

	images = make([]LazlowFrame, 1)

	// Add initial frame (always unprocessed for preview purposes)
	images[0] = LazlowFrame{inputImage, delay}

	p := &caire.Processor{
		// Initialize struct variables
		NewWidth:  inputImage.Bounds().Dx(),
		NewHeight: inputImage.Bounds().Dy(),
	}

	// Create shaken frames
	for i := 1; i < frameCount; i++ {
		log.Printf("Generating CAS frame %d", i)

		p.NewHeight -= pixelDelta
		if p.NewHeight < 1 {
			p.NewHeight = 1
		}

		p.NewWidth -= pixelDelta
		if p.NewWidth < 1 {
			p.NewWidth = 1
		}

		if p.NewHeight == 1 && p.NewWidth == 1 {
			break // can't shrink further
		}

		// Content-aware resizing
		src := image.NewNRGBA(
			image.Rect(
				0,
				0,
				inputImage.Bounds().Dx(),
				inputImage.Bounds().Dy()))
		draw.Draw(src, inputImage.Bounds(), inputImage, image.Point{}, draw.Src)
		casOutputImg, err := p.Resize(src)
		if err != nil {
			panic(err) // TODO
		}

		// Resize back to normal size
		img := resize.Resize(
			uint(inputImage.Bounds().Dx()),
			uint(inputImage.Bounds().Dy()),
			casOutputImg, resize.Lanczos3)

		images = append(images, LazlowFrame{img, delay})
	}

	return
}

func init() {
	RegisterEffect("cas", new(LazlowContentAwareScalingEffect))
}
