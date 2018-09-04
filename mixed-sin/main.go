package main

import (
	"fmt"
	"math"
	"os"

	"github.com/youpy/go-wav"
)

const (
	channelCount = 2
	samplingRate = 44110
	bitPerSample = 16
)

const tuningToneFrq = 440

var scale = [7]string{"A", "B", "C", "D", "E", "F", "G"}
var intonation = map[string]float64{
	"A": math.Pow(2.0, 0.0/12.0),
	"B": math.Pow(2.0, 2.0/12.0),
	"C": math.Pow(2.0, 3.0/12.0),
	"D": math.Pow(2.0, 5.0/12.0),
	"E": math.Pow(2.0, 7.0/12.0),
	"F": math.Pow(2.0, 8.0/12.0),
	"G": math.Pow(2.0, 10.0/12.0),
}

func main() {
	file, err := os.Create("output.wav")
	must(err)

	const seconds = 3
	const speed = 1
	const scaleOffset = 2
	const debug = false

	maxValue := math.Pow(2.0, bitPerSample-1) - 1
	minValue := -1 * maxValue

	t := 0
	samples := make([]wav.Sample, seconds*samplingRate)
	for i := range samples {
		if i%(samplingRate/speed) == 0 {
			t++
		}

		oct := (t - 1 + scaleOffset) / len(scale)
		note := scale[(t-1+scaleOffset)%len(scale)]
		baseToneFrq := math.Pow(2.0, float64(oct)) * tuningToneFrq
		toneFrq := intonation[note] * baseToneFrq

		if debug {
			if i%(samplingRate/speed) == 0 {
				fmt.Printf("%d: oct=%d(%.2f Hz) note=%s freq=%.2f Hz\n", t, oct, baseToneFrq, note, toneFrq)
			}
		}

		x := float64(i) * toneFrq / samplingRate
		y := (maxValue*0.5)*math.Sin(x*2.0*math.Pi) + // 基音
			(maxValue*0.3)*math.Sin((60*2.0*math.Pi/360.0)+2.0*x*2.0*math.Pi) // 倍音

		samples[i].Values[0] = int(math.Round(math.Min(math.Max(y, minValue), maxValue)))
		samples[i].Values[1] = samples[i].Values[0]
	}

	w := wav.NewWriter(file, uint32(len(samples)), channelCount, samplingRate, bitPerSample)
	err = w.WriteSamples(samples)
	must(err)

	err = file.Close()
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
