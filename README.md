# LazLow Image & Animation Generator

This tool is made to let you easily modify or animate images with preprogrammed effects.

Currently supported effects:

- Animated violent image shaking
- Discord ping count overlay

Currently supported output formats:

- GIF (static and animated), you want to use this format for when you upload animations as emojis to Discord
- WebP (static and animated), this is better for static images on Discord, WebP animations somehow don't work there yet
- PNG (static and animated), however Animated PNG tends to be unsupported by most software
- JPEG (static only)

## Usage

    usage: lazlow [<flags>] <effect> <input-file> <output-file>

    Lazy, low-effort meme generator tool


    Flags:
        --help                   Show context-sensitive help.
    -t, --output-type=auto       Output file type
        --discord-ping-number=1  What the red counter at the bottom right should display
        --shake-delay=20ms       Delay between frames
        --shake-frames=12        How many frames to generate
        --shake-percentage=20    How much to shake the picture in percent
        --jpeg-quality=90        The quality level with which JPEG output will be written, where 100 = lossless and lower will be increasingly lossy.
        --webp-lossless          Enable lossless output

    Args:
    <effect>       Effect to use
    <input-file>   Input file name
    <output-file>  Output file name

## Building

All you need to compile is [Go](https://golang.org), you don't even have to clone the source code since Go will do that for you. Just run this command:

    go get github.com/icedream/lazlow/cmd/lazlow

This is enough to build the binary and it will be put into your `${GOPATH}/bin` by Go.
