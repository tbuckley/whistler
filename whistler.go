package main

import (
	"fmt"
	"math"
	"time"

	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	w, err := NewWhistler()
	chk(err)
	defer w.Close()

	chk(w.Start())
	time.Sleep(5 * time.Second)
}

type Whistler struct {
	*portaudio.Stream

	SampleRate float64
	Buffer     []float64
	Index      int

	Matcher *Matcher
}

func NewWhistler() (*Whistler, error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Device = nil

	whistler := new(Whistler)
	whistler.SampleRate = p.SampleRate
	whistler.Buffer = make([]float64, 4000)
	whistler.Matcher = NewMatcher()

	whistler.Stream, err = portaudio.OpenStream(p, whistler.ProcessAudio)
	return whistler, err
}

func (w *Whistler) ProcessAudio(in []float32) {
	for i := range in {
		w.Buffer[w.Index] = float64(in[i])
		w.Index = (w.Index + 1) % len(w.Buffer)
		if w.Index == 0 {
			waves := w.CalculateFFT()
			isMatch := w.Matcher.Match(waves)
			if isMatch {
				fmt.Println("MATCH!!!")
			} else {
				fmt.Println("No match...")
			}
		}
	}
}

func (w *Whistler) CalculateFFT() []SineWave {
	ys := fft.FFTReal(w.Buffer)
	ys = ys[:len(ys)/2]

	waves := make([]SineWave, len(ys))
	for index, y := range ys {
		waves[index] = interpretIndividualFFT(y, index, len(w.Buffer), w.SampleRate)
	}

	WaveSet(waves).Sort()
	filteredWaves := make([]SineWave, 0)
	for _, wave := range waves {
		if wave.Amplitude > 1.0e-3 {
			// fmt.Println(wave.String())
			filteredWaves = append(filteredWaves, wave)
		}
	}

	return filteredWaves
}

func interpretIndividualFFT(x complex128, index, n int, samplingRate float64) SineWave {
	r := real(x)
	i := imag(x)

	freq := float64(index) * samplingRate / float64(n)
	amp := math.Sqrt((r*r)+(i*i)) * 2 / float64(n)
	phase := math.Atan(i/r) + math.Pi/2

	return SineWave{
		Amplitude:   amp,
		Frequency:   freq,
		PhaseOffset: phase,
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
