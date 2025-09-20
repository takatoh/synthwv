package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/takatoh/synthwv/envelope"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/iterator"
	"github.com/takatoh/synthwv/phase"
	"github.com/takatoh/synthwv/synthesizer"
)

func main() {
	progName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s [options]

Options:
`, progName)
		flag.PrintDefaults()
	}
	optLength := flag.Float64("length", 60.0, "Time-history length(sec).")
	optDt := flag.Float64("dt", 0.01, "DT")
	optLevel := flag.Int("level", 2, "Specify level 1 or 2.")
	flag.Parse()

	dt := *optDt
	n := int(*optLength / dt)
	m := n / 2
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

	var env func(float64) float64
	if *optLevel == 1 {
		env = envelope.Level1
	} else {
		env = envelope.Level2
	}
	synthszr := synthesizer.New(dt, omega, phi, env)
	tests := [](func([]float64) bool){test1, test2}
	inspectr := inspector.New(tests)
	itertr := iterator.New(synthszr, inspectr, 3)
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
