package synthesizer

import "math"

type Synthesizer struct {
	N        int
	Dt       float64
	omega    []float64
	phi      []float64
	envelope func(float64) float64
}

func New(dt float64, omega, phi []float64, envelope func(float64) float64) *Synthesizer {
	p := new(Synthesizer)
	p.Dt = dt
	p.omega = omega
	p.phi = phi
	p.envelope = envelope
	return p
}

func (s *Synthesizer) Synthesize(a []float64) []float64 {
	t := 0.0
	y := make([]float64, s.N)
	m := len(s.omega)
	for j := range s.N {
		for i := range m {
			y[j] += a[i] * math.Sin(s.omega[i]*t+s.phi[i])
		}
		y[j] = s.envelope(t) * y[j]
		t += s.Dt
	}
	return y
}
