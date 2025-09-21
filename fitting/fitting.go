package fitting

import (
	"math"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
)

type Fitting struct {
	Period []float64
	DSa    []float64
	DSv    []float64
}

func New(t, dsa []float64) *Fitting {
	p := new(Fitting)
	p.Period = t
	p.DSa = dsa
	var dsv []float64
	for i := range dsa {
		w := 2 * math.Pi / t[i]
		dsv = append(dsv, w*dsa[i])
	}
	p.DSv = dsv
	return p
}

func (f *Fitting) MinSpecRatio(acc []float64) bool {
	wave := seismicwave.Make("", 0.01, acc)
	resp := response.Spectrum(wave, f.Period, 0.05)
	minRatio := 1.0
	for i := range resp {
		ratio := f.DSa[i] / resp[i].Sa
		if ratio < minRatio {
			minRatio = ratio
		}
	}
	return minRatio >= 0.85
}
