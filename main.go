package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"

	"github.com/takatoh/respspec/response"
	"github.com/takatoh/seismicwave"
	"github.com/takatoh/synthwv/envelope"
	"github.com/takatoh/synthwv/fitting"
	"github.com/takatoh/synthwv/inspector"
	"github.com/takatoh/synthwv/iterator"
	"github.com/takatoh/synthwv/phase"
	"github.com/takatoh/synthwv/synthesizer"
	"github.com/takatoh/synthwv/tuner"
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

	// Show version and exit
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

	// dt : time delta
	// n : number of synthesized wave
	dt := *optDt
	n := int(*optLength / dt)

	// synthPeriod : period points for synthesize
	synthPeriod := synthesizer.DefaultPeriod()
	// m : number of component waves
	m := len(synthPeriod)
	// synthOmega : circular frequency points for synthesize
	synthOmega := make([]float64, m)
	for i := range m {
		synthOmega[i] = 2.0 * math.Pi / synthPeriod[i]
	}
	// Phase angles for synthesize
	synthPhase := phase.RandomPhaseAngles(m)
	// Set envelope function
	envl := envelope.GetEnveolope(*optEnvelope)
	if envl == nil {
		fmt.Printf("Error: Not found envelope named '%s'\n", *optEnvelope)
		os.Exit(1)
	}
	// Synthesizer
	synthszr := synthesizer.New(dt, n, synthOmega, synthPhase, envl)

	// Period points for fitting judgement, descending order
	fittingPeriod := response.DefaultPeriod()
	slices.Reverse(fittingPeriod)
	// Spectra (Sa) for fitting judgement
	_, fittingSa := utils.Interpolate(dsaT, dsaVal, fittingPeriod, true)
	// Fitting tests and inspector
	fittingTestr := fitting.New(fittingPeriod, fittingSa)
	tests := [](func(*seismicwave.Wave) bool){
		fittingTestr.MinSpecRatio,   // minimum spectra retio
		fittingTestr.VariationCoeff, // coefficient of variation
		fittingTestr.ErrAverage,     // error average
		fittingTestr.SIRatio,        // SI ratio
	}
	inspectr := inspector.New(tests)

	// Sa for synthesize and tuning
	_, synthSa := utils.Interpolate(dsaT, dsaVal, synthPeriod, true)
	// Tuner
	tuner := tuner.New(synthPeriod, synthSa)
	// Initial values of amplitude for sysnthesize
	ampInitial := tuner.InitAmplitude(synthPeriod)

	// Synthesize a wave
	itertr := iterator.New(synthszr, inspectr, 3)
	timehist := itertr.Iterate(ampInitial)

	// Output a result wave time history
	if *optCsv {
		printResultAsCsv(n, dt, timehist)
	} else {
		printResult(n, dt, timehist)
	}
}

// Print result
func printResult(n int, dt float64, timehist []float64) {
	fmt.Println("   TIME         ACC.")
	t := 0.0
	for i := range n {
		fmt.Printf("%7.2f   %10.3f\n", t, timehist[i])
		t += dt
	}
}

// Print result as CSV format
func printResultAsCsv(n int, dt float64, timehist []float64) {
	fmt.Println("Time,Acc.")
	t := 0.0
	for i := range n {
		fmt.Printf("%.2f,%.3f\n", t, timehist[i])
		t += dt
	}
}
