package inspector

type Inspector struct {
	tests []func([]float64) bool
}

func New(tests []func([]float64) bool) *Inspector {
	p := new(Inspector)
	p.tests = tests
	return p
}

func (ins *Inspector) Inspect(y []float64) bool {
	for _, test := range ins.tests {
		if !(test(y)) {
			return false
		}
	}
	return true
}
