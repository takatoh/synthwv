package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
	"github.com/takatoh/synthwv/envelope"
	"github.com/takatoh/synthwv/fitting"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/iterator"
	"github.com/takatoh/synthwv/phase"
	"github.com/takatoh/synthwv/synthesizer"
	"github.com/takatoh/synthwv/utils"
)

const progVersion = "v0.1.0"

func main() {
	progName := filepath.Base(os.Args[0])
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s [options] <DSa.csv>

Options:
`, progName)
		flag.PrintDefaults()
	}
	optLength := flag.Float64("length", 60.0, "Time-history length(sec).")
	optDt := flag.Float64("dt", 0.01, "DT")
	optLevel := flag.Int("level", 2, "Specify level 1 or 2.")
	optVersion := flag.Bool("version", false, "Show version.")
	flag.Parse()

	if *optVersion {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	// Load a target spectrum for design
	dsaFile := flag.Arg(0)
	dsaT, dsaVal, err := utils.LoadDesignSpectrum(dsaFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defaultPeriod := response.DefaultPeriod()
	dsaT, dsaVal = utils.Interpolate(dsaT, dsaVal, defaultPeriod)

	dt := *optDt
	n := int(*optLength / dt)
	m := n / 2
	ndt := float64(n) * dt

	omega := make([]float64, m)
	for i := range m {
		omega[i] = float64(i) / ndt
	}

	phi := phase.RandomPhaseAngles(m)

	initialAmplitude := make([]float64, m)
	for i := range m {
		initialAmplitude[i] = 1.0
	}

	// Set envelope function
	var env func(float64) float64
	if *optLevel == 1 {
		env = envelope.Level1
	} else {
		env = envelope.Level2
	}

	// Synthesize a wave
	synthszr := synthesizer.New(dt, n, omega, phi, env)
	fittingTestr := fitting.New(dsaT, dsaVal)
	tests := [](func(*seismicwave.Wave) bool){
		fittingTestr.MinSpecRatio,
		fittingTestr.VariationCoeff,
		fittingTestr.MeanErr,
		fittingTestr.SIRatio,
	}
	inspectr := inspector.New(tests)
	itertr := iterator.New(synthszr, inspectr, 3)
	timehist := itertr.Iterate(initialAmplitude)

	// Output a result wave time history
	printResult(n, dt, timehist)
}

func printResult(n int, dt float64, timehist []float64) {
	t := 0.0
	for i := range n {
		fmt.Printf("%7.2f %8.3f\n", t, timehist[i])
		t += dt
	}
}
