package output

import (
	"fmt"
	"os"

	escapes "github.com/snugfox/ansi-escapes"
)

type Color int

const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func FprintColor(f *os.File, text string, color Color) {
	var colorSequence string
	switch color {
	case Black:
		colorSequence = escapes.TextColorBlack
	case Red:
		colorSequence = escapes.TextColorRed
	case Green:
		colorSequence = escapes.TextColorGreen
	case Yellow:
		colorSequence = escapes.TextColorYellow
	case Blue:
		colorSequence = escapes.TextColorBlue
	case Magenta:
		colorSequence = escapes.TextColorMagenta
	case Cyan:
		colorSequence = escapes.TextColorCyan
	case White:
		colorSequence = escapes.TextColorWhite
	}

	_, err := fmt.Fprintf(f, "%s%s%s", colorSequence, text, escapes.TextColorWhite)
	if err != nil {
		// Since we're only printing to stdout/stderr we don't expect errors.
		// Just panic if one happens.
		panic(err)
	}
}

func PrintColor(text string, color Color) {
	FprintColor(os.Stdout, text, color)
}

func PrintColorLn(text string, color Color) {
	FprintColor(os.Stdout, text+"\n", color)
}

func Ln() {
	fmt.Println()
}

func ClearScreen() {
	fmt.Print(escapes.ClearScreen)
}
