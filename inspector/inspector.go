package inspector

import "github.com/takatoh/seismicwave"

type Inspector struct {
	tests []func(*seismicwave.Wave, float64) bool
	h     float64
}

func New(tests []func(*seismicwave.Wave, float64) bool, h float64) *Inspector {
	p := new(Inspector)
	p.tests = tests
	return p
}

func (ins *Inspector) Inspect(y *seismicwave.Wave) bool {
	for _, test := range ins.tests {
		if !(test(y, ins.h)) {
			return false
		}
	}
	return true
}
