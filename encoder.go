package lazlow

import (
	"errors"

	"github.com/icedream/lazlow/effects"
)

var ErrOnlySingleFrameOutputSupported = errors.New("this encoder only supports single-frame output")

type LazlowEncoder interface {
	Options() map[string]effects.LazlowOption
	SupportsFileExtension(ext string) bool
	SupportsFrames(frameCount int) bool
	Encode(frames []effects.LazlowFrame, out *LazlowOutput, options map[string]effects.LazlowOption) (err error)
}

var registeredEncoders = map[string]LazlowEncoder{}

func RegisterEncoder(id string, encoder LazlowEncoder) {
	// TODO - implement safety check
	registeredEncoders[id] = encoder
}

func GetEncoder(id string) (effect LazlowEncoder, ok bool) {
	effect, ok = registeredEncoders[id]
	return
}

func GetRegisteredEncoderIDs() (retval []string) {
	// TODO - locking
	retval = make([]string, len(registeredEncoders))
	i := 0
	for key := range registeredEncoders {
		retval[i] = key
		i++
	}
	return
}

func GetRegisteredEncoders() (retval map[string]LazlowEncoder) {
	// TODO - locking
	return registeredEncoders
}

func DetectOutputType(effect effects.LazlowEffect, ext string) (encoderID string) {
	// TODO - locking
	var encoder LazlowEncoder
	for currentEncoderID, currentEncoder := range registeredEncoders {
		if !currentEncoder.SupportsFileExtension(ext) {
			continue
		}

		if effect.IsAnimated() && currentEncoder.SupportsFrames(2) {
			encoder = currentEncoder
			encoderID = currentEncoderID
			continue
		}

		if !effect.IsAnimated() && currentEncoder.SupportsFrames(1) {
			if encoder != nil && !encoder.SupportsFrames(2) {
				// skip other encoders if we already got one that is optimized for static images
				continue
			}
			encoder = currentEncoder
			encoderID = currentEncoderID
			continue
		}
	}
	return encoderID
}
