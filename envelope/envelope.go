package envelope

import "math"

func Identity(t float64) float64 {
	return 1.0
}

func Level1(t float64) float64 {
	if t < 5.0 {
		return math.Pow(t/5.0, 2.0)
	} else if t < 25.0 {
		return 1.0
	} else {
		return math.Exp(-0.066 * (t - 25.0))
	}
}

func Level2(t float64) float64 {
	if t < 5.0 {
		return math.Pow(t/5.0, 2.0)
	} else if t < 35.0 {
		return 1.0
	} else {
		return math.Exp(-0.027 * (t - 35.0))
	}
}

var Envelopes = map[string]func(float64) float64{
	"id":     Identity,
	"level1": Level1,
	"level2": Level2,
}

func GetEnveolope(key string) func(float64) float64 {
	env := Envelopes[key]
	return env
}
