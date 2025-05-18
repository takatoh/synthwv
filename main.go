package main

import (
	"fmt"

	"github.com/takatoh/synthwv/envelope"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/iterator"
	"github.com/takatoh/synthwv/phase"
	"github.com/takatoh/synthwv/synthesizer"
)

func main() {
	n := 4096
	m := n / 2
	dt := 0.01
	ndt := float64(n) * dt

	omega := make([]float64, m)
	for i := range m {
		omega[i] = float64(i) / ndt
	}
	phi := phase.RandomPhaseAngles(m)

	amplitude := make([]float64, m)
	for i := range m {
		amplitude[i] = 1.0
	}

	synthszr := synthesizer.New(dt, omega, phi, envelope.Identity)
	tests := [](func([]float64) bool){test1, test2}
	inspectr := inspector.New(tests)
	itertr := iterator.New(synthszr, inspectr)
	timehist := itertr.Iterate(n, amplitude)

	t := 0.0
	for i := range n {
		fmt.Printf("%7.2f %8.3f\n", t, timehist[i])
		t += dt
	}
}

func test1(values []float64) bool {
	return true
}

func test2(values []float64) bool {
	return false
}
