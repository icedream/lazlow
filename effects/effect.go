package effects

import (
	"image"
	"time"
)

type LazlowFrame struct {
	Image image.Image
	Delay time.Duration
}

type LazlowEffect interface {
	Options() map[string]LazlowOption
	IsAnimated() bool
	Process(input image.Image, options map[string]LazlowOption) (output []LazlowFrame)
}

var registeredEffects = map[string]LazlowEffect{}

func RegisterEffect(id string, effect LazlowEffect) {
	// TODO - implement safety check
	registeredEffects[id] = effect
}

func GetEffect(id string) (effect LazlowEffect, ok bool) {
	effect, ok = registeredEffects[id]
	return
}

func GetRegisteredEffectIDs() (retval []string) {
	// TODO - locking
	retval = make([]string, len(registeredEffects))
	i := 0
	for key := range registeredEffects {
		retval[i] = key
		i++
	}
	return
}

func GetRegisteredEffects() (retval map[string]LazlowEffect) {
	// TODO - locking
	return registeredEffects
}
