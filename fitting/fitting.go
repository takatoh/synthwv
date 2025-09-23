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
	DSI    float64
}

func New(t, dsa []float64) *Fitting {
	p := new(Fitting)
	p.Period = t
	p.DSa = dsa
	psv, psd := pSvSd(t, dsa)
	p.DSv = psv
	dspec := spec(t, dsa, psv, psd)
	p.DSI = response.CalcSI(dspec)
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

func (f *Fitting) SIRatio(acc *seismicwave.Wave) bool {
	resp := response.Spectrum(acc, f.Period, 0.05)
	si := response.CalcSI(resp)
	return f.DSI/si >= 1.0
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

func (f *Fitting) ErrAverage(acc *seismicwave.Wave) bool {
	resp := response.Spectrum(acc, f.Period, 0.05)
	eTotal := 0.0
	for i := range resp {
		e := resp[i].Sa / f.DSa[i]
		eTotal += e
	}
	errAve := eTotal / float64(len(resp))
	return math.Abs(1.0-errAve) <= 0.02
}

func pSvSd(t, sa []float64) ([]float64, []float64) {
	n := len(sa)
	sv := make([]float64, n)
	sd := make([]float64, n)
	for i := range n {
		w := 2 * math.Pi / t[i]
		sv[i] = w * sa[i]
		sd[i] = w * w * sa[i]
	}
	return sv, sd
}

func spec(t, sa, sv, sd []float64) []*response.Response {
	n := len(t)
	s := make([]*response.Response, n)
	for i := range t {
		r := response.NewResponse(t[i], sa[i], sv[i], sd[i])
		s[i] = r
	}
	return s
}
