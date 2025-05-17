package phase

import (
	"math"
	"math/rand"
)

func RandomPhaseAngles(n int) []float64 {
	phi := make([]float64, n)
	for i := range n {
		phi[i] = rand.Float64() * 2.0 * math.Pi
	}
	return phi
}
