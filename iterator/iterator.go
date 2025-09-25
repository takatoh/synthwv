package iterator

import (
	"github.com/takatoh/seismicwave"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/synthesizer"
	"github.com/takatoh/synthwv/tuner"
)

type Iterator struct {
	synthesizer *synthesizer.Synthesizer
	inspector   *inspector.Inspector
	tuner       *tuner.Tuner
	iter_limit  int
}

func New(synthesizer *synthesizer.Synthesizer, inspector *inspector.Inspector, tuner *tuner.Tuner, iter_limit int) *Iterator {
	p := new(Iterator)
	p.synthesizer = synthesizer
	p.inspector = inspector
	p.tuner = tuner
	p.iter_limit = iter_limit
	return p
}

func (itr *Iterator) Iterate(amp []float64) []float64 {
	var y []float64
	count := 0
	for {
		count += 1
		y = itr.synthesizer.Synthesize(amp)
		wave := seismicwave.Make("", itr.synthesizer.Dt, y)
		if itr.inspector.Inspect(wave) {
			return y
		} else if count == itr.iter_limit {
			return y
		} else {
			amp = itr.tuner.Tune(amp, wave)
		}
	}
}
