package iterator

import (
	"github.com/takatoh/seismicwave"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/synthesizer"
)

type Iterator struct {
	synthesizer *synthesizer.Synthesizer
	inspector   *inspector.Inspector
	iter_limit  int
	dt          float64
}

func New(synthesizer *synthesizer.Synthesizer, inspector *inspector.Inspector, iter_limit int, dt float64) *Iterator {
	p := new(Iterator)
	p.synthesizer = synthesizer
	p.inspector = inspector
	p.iter_limit = iter_limit
	p.dt = dt
	return p
}

func (itr *Iterator) Iterate(amp []float64) []float64 {
	var y []float64
	count := 0
	for {
		count += 1
		y = itr.synthesizer.Synthesize(amp)
		wave := seismicwave.Make("", itr.dt, y)
		if itr.inspector.Inspect(wave) {
			break
		} else if count == itr.iter_limit {
			break
		}
	}
	return y
}
