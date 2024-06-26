package lyspmath

import "math"

// from https://gosamples.dev/round-float/
func RoundFloat32(val float32, precision uint) float32 {
	ratio := math.Pow(10, float64(precision))
	return float32(math.Round(float64(val)*ratio) / ratio)
}
