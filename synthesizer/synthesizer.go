package synthesizer

import "math"

type Synthesizer struct {
	dt    float64
	omega []float64
	phi   []float64
}

func New(dt float64, omega, phi []float64) *Synthesizer {
	p := new(Synthesizer)
	p.dt = dt
	p.omega = omega
	p.phi = phi
	return p
}

func (s *Synthesizer) Synthesize(n int, a []float64) []float64 {
	t := 0.0
	y := make([]float64, n)
	for i := range n {
		y[i] = a[i] * math.Sin(s.omega[i]*t+s.phi[i])
		t += s.dt
	}
	return y
}
