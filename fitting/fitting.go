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

func (f *Fitting) MinSpecRatio(acc *seismicwave.Wave) bool {
	resp := response.Spectrum(acc, f.Period, 0.05)
	minRatio := 1.0
	for i := range resp {
		ratio := resp[i].Sa / f.DSa[i]
		if ratio < minRatio {
			minRatio = ratio
		}
	}
	return minRatio >= 0.85
}

func (f *Fitting) VariationCoeff(acc *seismicwave.Wave) bool {
	resp := response.Spectrum(acc, f.Period, 0.05)
	eTotal := 0.0
	for i := range resp {
		e := resp[i].Sa / f.DSa[i]
		eTotal += math.Pow((e - 1.0), 2)
	}
	variationCoeff := math.Sqrt(eTotal / float64(len(resp)))
	return variationCoeff <= 0.05
}
