package tuner

import (
	"math"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
)

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

// Values of amplitude for next
func (tnr *Tuner) Tune(currAmp []float64, currWave *seismicwave.Wave) []float64 {
	resp := response.Spectrum(currWave, tnr.T, 0.05)
	currSa := make([]float64, len(resp))
	for i := range resp {
		currSa[i] = resp[i].Sa
	}
	amp := make([]float64, len(currAmp))
	for i := range amp {
		amp[i] = currAmp[i] * tnr.Sa[i] / currSa[i]
	}
	return amp
}
