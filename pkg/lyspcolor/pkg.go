package lyspcolor

import (
	"encoding/hex"
	"fmt"
)

// from https://github.com/teacat/noire/blob/master/noire.go#L276
func hexToRGB(hexColor string) (r float64, g float64, b float64, err error) {

	if len(hexColor) != 7 || hexColor[0] != '#' {
		return 0, 0, 0, fmt.Errorf("hexColor must be in the form '#AABBCC'")
	}

	byteArray, err := hex.DecodeString(hexColor[1:])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("hex.DecodeString failed for string: %s: %w", hexColor[1:], err)
	}

	r = float64(byteArray[0])
	g = float64(byteArray[1])
	b = float64(byteArray[2])

	return r, g, b, nil
}

// from https://github.com/teacat/noire/blob/master/noire.go#L582
func isLight(r float64, g float64, b float64) bool {
	darkness := 1 - (0.299*r+0.587*g+0.114*b)/255
	return darkness < 0.5
}

// HexIsLight returns true if the hex color is a light scheme
// hexColor must be a hex RGB color with prefix, e.g. #B2DFDB
func HexIsLight(hexColor string) (ret bool, err error) {

	r, g, b, err := hexToRGB(hexColor)
	if err != nil {
		return false, fmt.Errorf("hexToRGB failed: %w", err)
	}

	return isLight(r, g, b), nil
}
