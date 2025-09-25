package tuner

import "math"

type Tuner struct {
	T   []float64
	Sa  []float64
	PSv []float64
}

func New(t, sa []float64) *Tuner {
	p := new(Tuner)
	p.T = t
	p.Sa = sa
	psv := make([]float64, len(sa))
	for i := range sa {
		w := 2.0 * math.Pi / t[i]
		psv[i] = w * sa[i]
	}
	p.PSv = psv
	return p
}

// Initial values of amplitude for synthesize
func (tnr *Tuner) InitAmplitude() []float64 {
	amp := make([]float64, len(tnr.PSv))
	for i := range tnr.PSv {
		amp[i] = 2.0 * tnr.PSv[i]
	}
	return amp
}
