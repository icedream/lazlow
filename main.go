package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

var (
	cli = kingpin.New("lazlow", "Lazy, low-effort meme generator tool")

	argEffect     = cli.Arg("effect", "Effect to use").Required().Enum("discord-ping", "shake")
	argInputFile  = cli.Arg("input-file", "Input file name").Required().ExistingFile()
	argOutputFile = cli.Arg("output-file", "Output file name").Required().String()

	flagEncoder = cli.Flag("output-type", "Output file type").Short('t').Default("auto").Enum("auto", "apng", "webp", "gif")
	flagDelay   = cli.Flag("delay", "Frame delay").Default("16ms").Duration()
)

func getEffect() effect {
	switch *argEffect {
	case "shake":
		return new(shakeEffect)
	case "discord-ping":
		return new(discordPingEffect)
	default:
		panic("logic error when selecting effect")
	}
}

func detectOutputType(frameCount int) string {
	switch strings.ToLower(filepath.Ext(*argOutputFile)) {
	case ".png":
		if frameCount > 1 {
			return "apng"
		}
		return "png"
	case ".jpeg", ".jpg":
		return "jpeg"
	case ".webm", ".webp":
		return "webp"
	case ".gif":
		return "gif"
	default:
		log.Fatalf("Can't automatically detect output file type for path: %s", *argOutputFile)
		return ""
	}
}

func getEncoder(frameCount int) encoder {
	encoder := *flagEncoder
	if *flagEncoder == "auto" {
		encoder = detectOutputType(frameCount)
	}

	switch encoder {
	case "apng":
		return new(apngEncoder)
	case "png":
		return new(pngEncoder)
	case "jpeg":
		return new(jpegEncoder)
	case "gif":
		return new(gifEncoder)
	case "webp":
		return new(webpEncoder)
	default:
		return nil
	}
}

func main() {
	kingpin.MustParse(cli.Parse(os.Args[1:]))

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

	out := newOutput(*argOutputFile)

	effect := getEffect()
	if effect == nil {
		panic("effect is nil")
	}

	frames := effect.Process(inputImage)

	encoder := getEncoder(len(frames))
	if encoder == nil {
		panic("encoder is nil")
	}
	encoder.Encode(frames, out)

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
