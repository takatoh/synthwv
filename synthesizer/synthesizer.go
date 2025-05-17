package synthesizer

import "math"

type Synthesizer struct {
	dt       float64
	omega    []float64
	phi      []float64
	envelope func(float64) float64
}

func New(dt float64, omega, phi []float64, envelope func(float64) float64) *Synthesizer {
	p := new(Synthesizer)
	p.dt = dt
	p.omega = omega
	p.phi = phi
	p.envelope = envelope
	return p
}

func (s *Synthesizer) Synthesize(n int, a []float64) []float64 {
	t := 0.0
	y := make([]float64, n)
	m := len(s.omega)
	for j := range n {
		for i := range m {
			y[j] += a[i] * math.Sin(s.omega[i]*t+s.phi[i])
		}
		y[j] = s.envelope(y[j])
		t += s.dt
	}
	return y
}
