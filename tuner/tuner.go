package tuner

import (
	"math"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
)

type Tuner struct {
	T  []float64
	Sa []float64
}

func New(t, sa []float64) *Tuner {
	p := new(Tuner)
	p.T = t
	p.Sa = sa
	return p
}

// Initial values of amplitude for synthesize
func (tnr *Tuner) InitAmplitude(length float64) []float64 {
	amp := make([]float64, len(tnr.Sa))
	for i := range tnr.Sa {
		w := 2.0 * math.Pi / tnr.T[i] // omega
		psv := w * tnr.Sa[i]          // Sv(h=0.05)
		psv0 := (2.0 / length) * psv  // Sv(h=0.0)
		amp[i] = 2.0 * psv0           // initial value of amplitude
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
