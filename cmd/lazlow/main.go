package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"
	"strconv"
	"time"

	// TODO - webp input support since Discord now uses that for a few static images…?

	"log"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/icedream/lazlow"
	"github.com/icedream/lazlow/effects"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

var (
	cli = kingpin.New("lazlow", "Lazy, low-effort meme generator tool")

	argEffect     = cli.Arg("effect", "Effect to use").Required().Enum(effects.GetRegisteredEffectIDs()...)
	argInputFile  = cli.Arg("input-file", "Input file name").Required().ExistingFile()
	argOutputFile = cli.Arg("output-file", "Output file name").Required().String()

	flagEncoder = cli.Flag("output-type", "Output file type").Short('t').Default("auto").Enum(append([]string{"auto"}, lazlow.GetRegisteredEncoderIDs()...)...)
)

var Version = "dev"

func detectOutputType(effect effects.LazlowEffect, ext string) string {
	encoderID := lazlow.DetectOutputType(effect, ext)
	if len(encoderID) < 1 {
		log.Fatalf("Can't automatically detect output file type for path: %s", *argOutputFile)
	}
	return encoderID
}

func getEncoder(effect effects.LazlowEffect, ext string) (string, lazlow.LazlowEncoder) {
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

func registerDynamicOption(prefix string, name string, option effects.LazlowOption) {
	fl := cli.Flag(prefix+"-"+name, option.Description()).
		Default(fmt.Sprint(option.DefaultValue()))
	switch option.(type) {
	case *effects.LazlowBoolOption:
		fl.
			Action(func(element *kingpin.ParseElement, context *kingpin.ParseContext) (err error) {
				v, err := strconv.ParseBool(*element.Value)
				if err != nil {
					return
				}
				return option.SetValue(v)
			}).
			Bool()
	case *effects.LazlowDurationOption:
		fl.
			Action(func(element *kingpin.ParseElement, context *kingpin.ParseContext) (err error) {
				v, err := time.ParseDuration(*element.Value)
				if err != nil {
					return
				}
				return option.SetValue(time.Duration(v))
			}).
			Duration()
	case *effects.LazlowIntegerOption:
		fl.
			Action(func(element *kingpin.ParseElement, context *kingpin.ParseContext) (err error) {
				v, err := strconv.ParseInt(*element.Value, 10, 64)
				if err != nil {
					return
				}
				return option.SetValue(v)
			}).
			Int64()
	default:
		panic(fmt.Sprintf("unknown type for option %s-%s, %+v", prefix, name, option))
	}
}

func registerDynamicOptions(prefix string, options map[string]effects.LazlowOption) {
	for name, option := range options {
		registerDynamicOption(prefix, name, option)
	}
}

func main() {
	fmt.Printf("LazLow %s\n", Version)
	fmt.Println("\tby Carl Kittelberger <icedream@icedream.pw>")
	fmt.Println("\thttps://github.com/icedream/lazlow")
	fmt.Println()

	// register flags for the effect options
	effectOptions := map[string]map[string]effects.LazlowOption{}
	for effectName, effect := range effects.GetRegisteredEffects() {
		effectOptions[effectName] = effect.Options()
		registerDynamicOptions(effectName, effectOptions[effectName])
	}
	// register flags for the encoder options
	encoderOptions := map[string]map[string]effects.LazlowOption{}
	for encoderName, encoder := range lazlow.GetRegisteredEncoders() {
		encoderOptions[encoderName] = encoder.Options()
		registerDynamicOptions(encoderName, encoderOptions[encoderName])
	}
	kingpin.MustParse(cli.Parse(os.Args[1:]))

	// get effect
	effect, ok := effects.GetEffect(*argEffect)
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

	frames, err := effect.Process(inputImage, effectOpts)
	if err != nil {
		log.Fatal("Failed to process image:", err)
		return
	}

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
