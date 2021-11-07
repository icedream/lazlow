package cas

import (
	"image"
	"image/draw"
	"log"
	"time"

	"github.com/esimov/caire"
	"github.com/icedream/lazlow/effects"
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

func (effect *LazlowContentAwareScalingEffect) Options() map[string]effects.LazlowOption {
	return map[string]effects.LazlowOption{
		lazlowContentAwareScalingEffectOptionDelay:  effects.NewLazlowEncoderDurationOption("Frame delay", "Delay between frames", 40*time.Millisecond, 0, 0xff*10*time.Millisecond, 10*time.Millisecond),
		lazlowContentAwareScalingEffectOptionFrames: effects.NewLazlowEncoderIntegerOption("Frame count", "How many frames to generate", 15*2, 0, 0xff, 1),
		lazlowContentAwareScalingEffectOptionStep:   effects.NewLazlowEncoderIntegerOption("Pixel delta step", "By how many pixels to scale with each frame", 4, 0, 0xffff, 1),
	}
}

func (effect *LazlowContentAwareScalingEffect) Process(inputImage image.Image, options map[string]effects.LazlowOption) (images []effects.LazlowFrame, err error) {
	delay := options[lazlowContentAwareScalingEffectOptionDelay].(*effects.LazlowDurationOption).TypedValue()
	frameCount := int(options[lazlowContentAwareScalingEffectOptionFrames].(*effects.LazlowIntegerOption).TypedValue())
	pixelDelta := int(options[lazlowContentAwareScalingEffectOptionStep].(*effects.LazlowIntegerOption).TypedValue())

	images = make([]effects.LazlowFrame, frameCount)

	// Add initial frame (always unprocessed for preview purposes)
	images[0] = effects.LazlowFrame{
		Image: inputImage,
		Delay: delay,
	}

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
			// can't shrink further
			images = images[0 : i-1]
			break
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

		images[i] = effects.LazlowFrame{
			Image: img,
			Delay: delay,
		}

		log.Println("Making frame #", i)
	}

	log.Println("Expecting to make this many frames:", frameCount)
	log.Println("Actually made this many frames:", len(images))

	return
}

func init() {
	effects.RegisterEffect("cas", new(LazlowContentAwareScalingEffect))
}
