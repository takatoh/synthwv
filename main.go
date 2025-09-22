package main

import (
	"flag"
	"fmt"
	"math"
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

const progVersion = "v0.2.0"

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
	optEnvelope := flag.String("envelope", "id", "Specify envelope function.")
	optCsv := flag.Bool("csv", false, "Output as CSV format.")
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

	m := len(dsaT)
	omega := make([]float64, m)
	for i := range m {
		omega[i] = 2.0 * math.Pi / dsaT[i]
	}

	dt := *optDt
	n := int(*optLength / dt)
	//m := n / 2
	//ndt := float64(n) * dt

	phi := phase.RandomPhaseAngles(m)

	//ampInitial := make([]float64, m)
	//for i := range m {
	//	ampInitial[i] = 1.0
	//}
	ampInitial := initAmplitude(dsaVal, omega)

	// Set envelope function
	env := envelope.GetEnveolope(*optEnvelope)
	if env == nil {
		fmt.Printf("Error: Not found envelope named '%s'\n", *optEnvelope)
		os.Exit(1)
	}

	// Synthesize a wave
	synthszr := synthesizer.New(dt, n, omega, phi, env)
	fittingTestr := fitting.New(dsaT, dsaVal)
	tests := [](func(*seismicwave.Wave) bool){
		fittingTestr.MinSpecRatio,
		fittingTestr.VariationCoeff,
		fittingTestr.ErrAverage,
		fittingTestr.SIRatio,
	}
	inspectr := inspector.New(tests)
	itertr := iterator.New(synthszr, inspectr, 3)
	timehist := itertr.Iterate(ampInitial)

	// Output a result wave time history
	if *optCsv {
		printResultAsCsv(n, dt, timehist)
	} else {
		printResult(n, dt, timehist)
	}
}

func initAmplitude(sa, w []float64) []float64 {
	m := len(sa)
	amp := make([]float64, m)
	for i := range m {
		amp[i] = 2.0 * w[i] * sa[i]
	}
	return amp
}

func printResult(n int, dt float64, timehist []float64) {
	t := 0.0
	for i := range n {
		fmt.Printf("%7.2f %8.3f\n", t, timehist[i])
		t += dt
	}
}

func printResultAsCsv(n int, dt float64, timehist []float64) {
	fmt.Println("Time,Acc.")
	t := 0.0
	for i := range n {
		fmt.Printf("%.2f,%.3f\n", t, timehist[i])
		t += dt
	}
}
