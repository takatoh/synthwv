package phase

import (
	"math"
	"math/rand"
)

func RandomPhaseAngles(m int) []float64 {
	phi := make([]float64, m)
	for i := range m {
		phi[i] = rand.Float64() * 2.0 * math.Pi
	}
	return phi
}
