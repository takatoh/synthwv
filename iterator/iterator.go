package iterator

import (
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/synthesizer"
)

type Iterator struct {
	synthesizer *synthesizer.Synthesizer
	inspector   *inspector.Inspector
}

func New(synthesizer *synthesizer.Synthesizer, inspector *inspector.Inspector) *Iterator {
	p := new(Iterator)
	p.synthesizer = synthesizer
	p.inspector = inspector
	return p
}

func (itr *Iterator) Iterate(n int, amp []float64) []float64 {
	var y []float64
	iter_limit := 3
	count := 0
	for {
		count += 1
		y := itr.synthesizer.Synthesize(n, amp)
		if itr.inspector.Inspect(y) {
			break
		} else if count == iter_limit {
			break
		}
	}
	return y
}
