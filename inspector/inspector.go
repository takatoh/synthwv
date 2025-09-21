package inspector

import "github.com/takatoh/seismicwave"

type Inspector struct {
	tests []func(*seismicwave.Wave) bool
}

func New(tests []func(*seismicwave.Wave) bool) *Inspector {
	p := new(Inspector)
	p.tests = tests
	return p
}

func (ins *Inspector) Inspect(y *seismicwave.Wave) bool {
	for _, test := range ins.tests {
		if !(test(y)) {
			return false
		}
	}
	return true
}
