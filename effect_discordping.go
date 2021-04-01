package lazlow

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var (
	colorTransparent       = color.RGBA{0, 0, 0, 0xff}
	colorDiscordCountBadge = color.RGBA{240, 71, 71, 0xff}
)

const (
	discordServerIconSize            float64 = 48
	discordRoundedBadgeMinWidth      float64 = 10
	discordRoundedBadgeDigitWidth    float64 = 6
	discordRoundedBadgeHeight        float64 = 16
	discordRoundedBadgeBorderRadiusX float64 = 8
	discordRoundedBadgeBorderRadiusY float64 = discordRoundedBadgeBorderRadiusX
)

const (
	lazlowDiscordPingEffectOptionNumber = "number"
)

type LazlowDiscordPingEffect struct {
}

func (effect *LazlowDiscordPingEffect) Options() map[string]LazlowOption {
	return map[string]LazlowOption{
		lazlowDiscordPingEffectOptionNumber: NewLazlowEncoderIntegerOption("Ping number", "What the red counter at the bottom right should display", 1, math.MinInt64, math.MaxInt64, 1),
	}
}

func (effect *LazlowDiscordPingEffect) IsAnimated() bool {
	return false
}

func (effect *LazlowDiscordPingEffect) Process(inputImage image.Image, options map[string]LazlowOption) (images []LazlowFrame) {
	images = make([]LazlowFrame, 1)

	inputNumber := options[lazlowDiscordPingEffectOptionNumber]

	/*
		Rounded corner for badge and mask: rx = 12, ry = 12

		Widths for discord digits:

		1 digit = red badge width 16px = mask width 24px
		2 digits = red badge width 22px = mask width 30px
		3 digits = red badge width 28px = mask width 36px
	*/

	inputWidth := float64(inputImage.Bounds().Dx())
	inputHeight := float64(inputImage.Bounds().Dy())

	size := inputWidth
	if inputHeight > size {
		size = inputHeight
	}

	scale := size / discordServerIconSize

	imagePosX := (size - inputWidth) / 2
	imagePosY := (size - inputHeight) / 2

	discordRoundedBadgeText := fmt.Sprintf("%d", inputNumber.(*LazlowIntegerOption).TypedValue())
	discordRoundedBadgeWidth := scale * (discordRoundedBadgeMinWidth + (float64(len(discordRoundedBadgeText)) * discordRoundedBadgeDigitWidth))
	discordRoundedBadgeX := imagePosX + inputWidth - discordRoundedBadgeWidth
	discordRoundedBadgeY := imagePosY + inputHeight - (scale * discordRoundedBadgeHeight)
	discordRoundedBadgeFontSizePixels := scale * 12

	discordRoundedBadgeMaskWidth := discordRoundedBadgeWidth + (scale * 8)
	discordRoundedBadgeMaskHeight := scale * (discordRoundedBadgeHeight + 8)
	discordRoundedBadgeMaskX := discordRoundedBadgeX - (scale * 4)
	discordRoundedBadgeMaskY := discordRoundedBadgeY - (scale * 4)
	discordRoundedBadgeMaskBorderRadiusX := scale * (discordRoundedBadgeBorderRadiusX + 4)
	discordRoundedBadgeMaskBorderRadiusY := scale * (discordRoundedBadgeBorderRadiusY + 4)

	// Calculate text width
	f, err := os.Open("whitneybold.otf")
	if err != nil {
		return
	}
	defer f.Close()
	otf, err := opentype.ParseReaderAt(f)
	if err != nil {
		return
	}
	face, err := opentype.NewFace(otf, &opentype.FaceOptions{
		Size:    discordRoundedBadgeFontSizePixels,
		DPI:     70,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	defer face.Close()
	var badgeTextWidth float64 = 0
	for _, x := range discordRoundedBadgeText {
		awidth, _ := face.GlyphAdvance(rune(x))
		badgeTextWidth += float64(awidth) / 64
	}
	discordRoundedBadgeTextX := discordRoundedBadgeX + (discordRoundedBadgeWidth / 2) - (badgeTextWidth / 2)
	discordRoundedBadgeTextY := discordRoundedBadgeY + discordRoundedBadgeFontSizePixels

	// Generate mask
	mask := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	draw.Draw(mask, inputImage.Bounds(), image.White, image.Pt(0, 0), draw.Src)
	mgc := draw2dimg.NewGraphicContext(mask)
	mgc.SetFillColor(color.Black)
	mgc.SetStrokeColor(color.Black)
	draw2dkit.RoundedRectangle(mgc,
		discordRoundedBadgeMaskX,
		discordRoundedBadgeMaskY,
		discordRoundedBadgeMaskX+discordRoundedBadgeMaskWidth,
		discordRoundedBadgeMaskY+discordRoundedBadgeMaskHeight,
		discordRoundedBadgeMaskBorderRadiusX*2,
		discordRoundedBadgeMaskBorderRadiusY*2)
	mgc.FillStroke()
	mgc.Close()
	for x := 0; x < mask.Rect.Dx(); x++ {
		for y := 0; y < mask.Rect.Dy(); y++ {
			c := mask.At(x, y)
			r32, _, _, _ := c.RGBA()
			r := uint8(r32)
			mask.SetRGBA(x, y, color.RGBA{r, r, r, r})
		}
	}

	// Generate final image
	dest := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	gc := draw2dimg.NewGraphicContext(dest)
	gc.SetFillColor(colorDiscordCountBadge)
	gc.SetStrokeColor(colorDiscordCountBadge)
	draw.DrawMask(
		// dest
		dest,
		image.Rect(
			int(imagePosX), int(imagePosY),
			int(imagePosX+inputWidth), int(imagePosY+inputHeight),
		),
		// src
		inputImage,
		image.Point{},
		// mask
		mask,
		image.Point{},
		// op
		draw.Src)
	draw2dkit.RoundedRectangle(gc,
		discordRoundedBadgeX, discordRoundedBadgeY,
		discordRoundedBadgeX+discordRoundedBadgeWidth, discordRoundedBadgeY+(scale*discordRoundedBadgeHeight),
		scale*discordRoundedBadgeBorderRadiusX*2,
		scale*discordRoundedBadgeBorderRadiusY*2)
	gc.FillStroke()

	d := font.Drawer{
		Dst:  dest,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(int(discordRoundedBadgeTextX), int(discordRoundedBadgeTextY)),
	}
	d.DrawString(discordRoundedBadgeText)

	images[0] = LazlowFrame{
		Image: dest,
	}

	return
}

func init() {
	RegisterEffect("discord-ping", new(LazlowDiscordPingEffect))
}
