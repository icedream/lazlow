package lazlow

import (
	"image"
	"image/draw"
	"math/rand"
	"time"
)

const (
	lazlowShakeEffectOptionDelay      = "delay"
	lazlowShakeEffectOptionFrames     = "frames"
	lazlowShakeEffectOptionPercentage = "percentage"
)

type LazlowShakeEffect struct {
}

func (effect *LazlowShakeEffect) IsAnimated() bool {
	return true
}

func (effect *LazlowShakeEffect) Options() map[string]LazlowOption {
	return map[string]LazlowOption{
		lazlowShakeEffectOptionDelay:      NewLazlowEncoderDurationOption("Frame delay", "Delay between frames", 20*time.Millisecond, 0, 0xff*10*time.Millisecond, 10*time.Millisecond),
		lazlowShakeEffectOptionFrames:     NewLazlowEncoderIntegerOption("Frame count", "How many frames to generate", 12, 0, 0xff, 1),
		lazlowShakeEffectOptionPercentage: NewLazlowEncoderIntegerOption("Shake percentage", "How much to shake the picture in percent", 20, 0, 100, 1),
	}
}

func (effect *LazlowShakeEffect) Process(inputImage image.Image, options map[string]LazlowOption) (images []LazlowFrame) {
	delay := options[lazlowShakeEffectOptionDelay].(*LazlowDurationOption).TypedValue()
	frameCount := int(options[lazlowShakeEffectOptionFrames].(*LazlowIntegerOption).TypedValue())
	shakeAmount := float64(options[lazlowShakeEffectOptionPercentage].(*LazlowIntegerOption).TypedValue()) / 100

	images = make([]LazlowFrame, frameCount)

	// Add initial frame (always unshaken for preview purposes)
	images[0] = LazlowFrame{inputImage, delay}

	shakeX := int(float64(inputImage.Bounds().Dx()) * shakeAmount)
	shakeY := int(float64(inputImage.Bounds().Dy()) * shakeAmount)

	// Create shaken frames
	for i := 1; i < frameCount; i++ {
		img := image.NewRGBA(
			image.Rect(
				0,
				0,
				inputImage.Bounds().Dx(),
				inputImage.Bounds().Dy()))
		x := rand.Intn(shakeX) - (shakeX / 2)
		y := rand.Intn(shakeY) - (shakeY / 2)

		draw.Draw(img, inputImage.Bounds(), inputImage, image.Point{x, y}, draw.Src)

		images[i] = LazlowFrame{img, delay}
	}

	return
}

func init() {
	RegisterEffect("shake", new(LazlowShakeEffect))
}
