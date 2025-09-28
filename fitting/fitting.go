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
	p.DSI = calcSI(dspec)           // Target SI
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
	si := calcSI(resp)
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

func calcSI(resp []*response.Response) float64 {
	var si float64 = 0.0
	for i := 1; i < len(resp); i++ {
		r0 := resp[i-1]
		r1 := resp[i]
		if 1.0 < r1.Period && r1.Period <= 5.0 {
			si = si + (r0.Sv+r1.Sv)*(r1.Period-r0.Period)/2.0
		}
	}

	return si / 2.4
}

// Default fitting periods: 300 points
func DefaultPeriod() []float64 {
	return []float64{
		5.000000,
		4.908515,
		4.818704,
		4.730537,
		4.643982,
		4.559011,
		4.475595,
		4.393705,
		4.313314,
		4.234393,
		4.156917,
		4.080858,
		4.006191,
		3.932889,
		3.860930,
		3.790286,
		3.720936,
		3.652854,
		3.586018,
		3.520404,
		3.455992,
		3.392757,
		3.330680,
		3.269739,
		3.209913,
		3.151181,
		3.093524,
		3.036922,
		2.981355,
		2.926806,
		2.873254,
		2.820682,
		2.769072,
		2.718407,
		2.668668,
		2.619839,
		2.571904,
		2.524846,
		2.478649,
		2.433298,
		2.388776,
		2.345068,
		2.302161,
		2.260038,
		2.218686,
		2.178091,
		2.138239,
		2.099115,
		2.060708,
		2.023003,
		1.985988,
		1.949651,
		1.913978,
		1.878958,
		1.844579,
		1.810829,
		1.777696,
		1.745170,
		1.713238,
		1.681891,
		1.651118,
		1.620907,
		1.591250,
		1.562135,
		1.533552,
		1.505493,
		1.477947,
		1.450905,
		1.424358,
		1.398296,
		1.372712,
		1.347595,
		1.322938,
		1.298733,
		1.274970,
		1.251642,
		1.228740,
		1.206258,
		1.184187,
		1.162520,
		1.141250,
		1.120368,
		1.099869,
		1.079745,
		1.059989,
		1.040594,
		1.021554,
		1.002863,
		0.984514,
		0.966500,
		0.948816,
		0.931456,
		0.914413,
		0.897682,
		0.881257,
		0.865133,
		0.849303,
		0.833764,
		0.818508,
		0.803532,
		0.788830,
		0.774397,
		0.760228,
		0.746318,
		0.732662,
		0.719257,
		0.706097,
		0.693177,
		0.680494,
		0.668043,
		0.655820,
		0.643821,
		0.632041,
		0.620476,
		0.609123,
		0.597978,
		0.587037,
		0.576296,
		0.565752,
		0.555400,
		0.545238,
		0.535262,
		0.525468,
		0.515854,
		0.506415,
		0.497149,
		0.488053,
		0.479123,
		0.470356,
		0.461750,
		0.453302,
		0.445008,
		0.436865,
		0.428872,
		0.421025,
		0.413322,
		0.405759,
		0.398335,
		0.391047,
		0.383892,
		0.376868,
		0.369972,
		0.363203,
		0.356557,
		0.350033,
		0.343629,
		0.337341,
		0.331169,
		0.325110,
		0.319161,
		0.313321,
		0.307589,
		0.301961,
		0.296436,
		0.291012,
		0.285687,
		0.280460,
		0.275328,
		0.270291,
		0.265345,
		0.260490,
		0.255724,
		0.251045,
		0.246452,
		0.241942,
		0.237516,
		0.233170,
		0.228903,
		0.224715,
		0.220604,
		0.216567,
		0.212605,
		0.208715,
		0.204896,
		0.201147,
		0.197466,
		0.193853,
		0.190307,
		0.186824,
		0.183406,
		0.180050,
		0.176756,
		0.173522,
		0.170347,
		0.167230,
		0.164170,
		0.161167,
		0.158218,
		0.155323,
		0.152481,
		0.149691,
		0.146952,
		0.144263,
		0.141624,
		0.139032,
		0.136489,
		0.133991,
		0.131540,
		0.129133,
		0.126770,
		0.124451,
		0.122173,
		0.119938,
		0.117744,
		0.115589,
		0.113474,
		0.111398,
		0.109360,
		0.107359,
		0.105395,
		0.103466,
		0.101573,
		0.099715,
		0.097890,
		0.096099,
		0.094341,
		0.092614,
		0.090920,
		0.089256,
		0.087623,
		0.086020,
		0.084446,
		0.082901,
		0.081384,
		0.079895,
		0.078433,
		0.076998,
		0.075589,
		0.074206,
		0.072849,
		0.071516,
		0.070207,
		0.068923,
		0.067661,
		0.066423,
		0.065208,
		0.064015,
		0.062844,
		0.061694,
		0.060565,
		0.059457,
		0.058369,
		0.057301,
		0.056253,
		0.055223,
		0.054213,
		0.053221,
		0.052247,
		0.051291,
		0.050353,
		0.049431,
		0.048527,
		0.047639,
		0.046767,
		0.045912,
		0.045072,
		0.044247,
		0.043437,
		0.042643,
		0.041862,
		0.041096,
		0.040345,
		0.039606,
		0.038882,
		0.038170,
		0.037472,
		0.036786,
		0.036113,
		0.035452,
		0.034804,
		0.034167,
		0.033542,
		0.032928,
		0.032326,
		0.031734,
		0.031153,
		0.030583,
		0.030024,
		0.029475,
		0.028935,
		0.028406,
		0.027886,
		0.027376,
		0.026875,
		0.026383,
		0.025900,
		0.025427,
		0.024961,
		0.024505,
		0.024056,
		0.023616,
		0.023184,
		0.022760,
		0.022343,
		0.021935,
		0.021533,
		0.021139,
		0.020752,
		0.020373,
		0.020000,
	}
}
