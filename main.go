package main

import (
	"fmt"

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

	a := make([]float64, m)
	for i := range m {
		a[i] = 1.0
	}

	synthszr := synthesizer.New(dt, omega, phi)
	y := synthszr.Synthesize(n, a)

	t := 0.0
	for i := range n {
		fmt.Printf("%7.2f %8.3f\n", t, y[i])
		t += dt
	}
}
