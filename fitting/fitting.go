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
	p.Period = t                    // Target period
	p.DSa = dsa                     // Target Sa
	psv, psd := pSvSd(t, dsa)       // Calc pSv, pSd
	p.DSv = psv                     // Target pDv
	dspec := spec(t, dsa, psv, psd) // Target spectrun: Sa, pSv, pSd
	p.DSI = response.CalcSI(dspec)  // Target SI
	return p
}

// Minumum spectrum ratio: 0.02-5 sec
func (f *Fitting) MinSpecRatio(acc *seismicwave.Wave) bool {
	period, sa := fittingRange(0.02, 5.0, f.Period, f.DSa)
	resp := response.Spectrum(acc, period, 0.05)
	minRatio := 1.0
	for i := range resp {
		ratio := resp[i].Sa / sa[i]
		if ratio < minRatio {
			minRatio = ratio
		}
	}
	return minRatio >= 0.85
}

// SI ratio: 1-5 sec
func (f *Fitting) SIRatio(acc *seismicwave.Wave) bool {
	period, _ := fittingRange(1.0, 5.0, f.Period, f.DSa)
	resp := response.Spectrum(acc, period, 0.05)
	si := response.CalcSI(resp)
	return f.DSI/si >= 1.0
}

// Coefficient of variation: 0.02-5 sec
func (f *Fitting) VariationCoeff(acc *seismicwave.Wave) bool {
	period, sa := fittingRange(0.02, 5.0, f.Period, f.DSa)
	resp := response.Spectrum(acc, period, 0.05)
	eTotal := 0.0
	for i := range resp {
		e := resp[i].Sa / sa[i]
		eTotal += math.Pow((e - 1.0), 2)
	}
	variationCoeff := math.Sqrt(eTotal / float64(len(resp)))
	return variationCoeff <= 0.05
}

// Error average: 0.02-5 sec
func (f *Fitting) ErrAverage(acc *seismicwave.Wave) bool {
	period, sa := fittingRange(0.02, 5.0, f.Period, f.DSa)
	resp := response.Spectrum(acc, period, 0.05)
	eTotal := 0.0
	for i := range resp {
		e := resp[i].Sa / sa[i]
		eTotal += e
	}
	errAve := eTotal / float64(len(resp))
	return math.Abs(1.0-errAve) <= 0.02
}

// pseudo Sv, Sd
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

func fittingRange(tmin, tmax float64, t, sa []float64) ([]float64, []float64) {
	tr := make([]float64, 0)
	sar := make([]float64, 0)
	for i := range t {
		if t[i] < tmin {
			continue
		} else if t[i] > tmax {
			continue
		} else {
			tr = append(tr, t[i])
			sar = append(sar, sa[i])
		}
	}
	return tr, sar
}
