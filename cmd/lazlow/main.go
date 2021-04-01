package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"

	// TODO - webp input support since Discord now uses that for a few static images…?

	"log"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/icedream/lazlow"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

var (
	cli = kingpin.New("lazlow", "Lazy, low-effort meme generator tool")

	argEffect     = cli.Arg("effect", "Effect to use").Required().Enum(lazlow.GetRegisteredEffectIDs()...)
	argInputFile  = cli.Arg("input-file", "Input file name").Required().ExistingFile()
	argOutputFile = cli.Arg("output-file", "Output file name").Required().String()

	flagEncoder = cli.Flag("output-type", "Output file type").Short('t').Default("auto").Enum(append([]string{"auto"}, lazlow.GetRegisteredEncoderIDs()...)...)
)

func detectOutputType(effect lazlow.LazlowEffect, ext string) string {
	encoderID := lazlow.DetectOutputType(effect, ext)
	if len(encoderID) < 1 {
		log.Fatalf("Can't automatically detect output file type for path: %s", *argOutputFile)
	}
	return encoderID
}

func getEncoder(effect lazlow.LazlowEffect, ext string) (string, lazlow.LazlowEncoder) {
	encoderID := *flagEncoder
	if *flagEncoder == "auto" {
		encoderID = detectOutputType(effect, ext)
	}

	encoder, ok := lazlow.GetEncoder(encoderID)
	if !ok {
		return "", nil
	}
	log.Printf("Using encoder: %s", encoderID)
	return encoderID, encoder
}

func registerDynamicOption(prefix string, name string, option lazlow.LazlowOption) {
	fl := cli.Flag(prefix+"-"+name, option.Description()).
		Default(fmt.Sprint(option.DefaultValue())).
		Action(func(element *kingpin.ParseElement, context *kingpin.ParseContext) error {
			return option.SetValue(element.Value)
		})
	switch option.(type) {
	case *lazlow.LazlowBoolOption:
		fl.Bool()
	case *lazlow.LazlowDurationOption:
		fl.Duration()
	case *lazlow.LazlowIntegerOption:
		fl.Int64()
	}
}

func registerDynamicOptions(prefix string, options map[string]lazlow.LazlowOption) {
	for name, option := range options {
		registerDynamicOption(prefix, name, option)
	}
}

func main() {
	// register flags for the effect options
	effectOptions := map[string]map[string]lazlow.LazlowOption{}
	for effectName, effect := range lazlow.GetRegisteredEffects() {
		effectOptions[effectName] = effect.Options()
		registerDynamicOptions(effectName, effectOptions[effectName])
	}
	// register flags for the encoder options
	encoderOptions := map[string]map[string]lazlow.LazlowOption{}
	for encoderName, encoder := range lazlow.GetRegisteredEncoders() {
		encoderOptions[encoderName] = encoder.Options()
		registerDynamicOptions(encoderName, encoderOptions[encoderName])
	}
	kingpin.MustParse(cli.Parse(os.Args[1:]))

	// get effect
	effect, ok := lazlow.GetEffect(*argEffect)
	if !ok {
		log.Fatalf("Could not find effect named %s", *argEffect)
	}
	effectOpts := effectOptions[*argEffect]

	inputFile, err := os.Open(*argInputFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer inputFile.Close()

	inputImage, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	out := lazlow.NewLazlowOutput(*argOutputFile)

	frames := effect.Process(inputImage, effectOpts)

	encoderID, encoder := getEncoder(effect, filepath.Ext(*argOutputFile))
	if encoder == nil {
		log.Fatalf("Could not find encoder for %s", *argOutputFile)
	}
	encoderOpts := encoderOptions[encoderID]
	encoder.Encode(frames, out, encoderOpts)

	writtenFiles, err := out.WrittenFiles()
	if err != nil {
		log.Fatalf("File writing succeeded but fetching file statistics immediately afterwards failed: %s", err)
		return
	}

	log.Println("Written files:")
	for _, stat := range writtenFiles {
		log.Printf("- %s (%s)", stat.Name(), humanize.Bytes(uint64(stat.Size())))
	}

	log.Println("Done.")
}
