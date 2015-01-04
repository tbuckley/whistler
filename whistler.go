package whistler

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

func Initialize() {
	portaudio.Initialize()
}

func Terminate() {
	portaudio.Terminate()
}

type handler struct {
	matcher   *Matcher
	matchChan chan bool
}

type Whistler struct {
	stream *portaudio.Stream

	sampleRate float64
	buffer     []float64
	index      int

	handlers []handler
}

func New() (*Whistler, error) {
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	p := portaudio.LowLatencyParameters(h.DefaultInputDevice, h.DefaultOutputDevice)
	p.Input.Channels = 1
	p.Output.Device = nil

	whistler := new(Whistler)
	whistler.sampleRate = p.SampleRate
	whistler.buffer = make([]float64, 4000)

	whistler.stream, err = portaudio.OpenStream(p, whistler.processAudio)
	return whistler, err
}

func (w *Whistler) Add(whistle MatchFactory) chan bool {
	matchChan := make(chan bool)
	matcher := whistle.New()
	w.handlers = append(w.handlers, handler{
		matcher:   matcher,
		matchChan: matchChan,
	})
	return matchChan
}

func (w *Whistler) Listen() error {
	return w.stream.Start()
}

func (w *Whistler) Close() error {
	return w.stream.Close()
}

func (w *Whistler) processAudio(in []float32) {
	for i := range in {
		w.buffer[w.index] = float64(in[i])
		w.index = (w.index + 1) % len(w.buffer)
		if w.index == 0 {
			w.checkMatches()
		}
	}
}

func (w *Whistler) checkMatches() {
	waves := w.calculateFFT()
	for _, handler := range w.handlers {
		if handler.matcher.Match(waves) {
			handler.matchChan <- true
		}
	}
}

func (w *Whistler) calculateFFT() []SineWave {
	ys := fft.FFTReal(w.buffer)
	ys = ys[:len(ys)/2]

	waves := make([]SineWave, len(ys))
	for index, y := range ys {
		waves[index] = interpretIndividualFFT(y, index, len(w.buffer), w.sampleRate)
	}

	WaveSet(waves).Sort()
	filteredWaves := make([]SineWave, 0)
	for _, wave := range waves {
		if wave.Amplitude > 1.0e-3 {
			filteredWaves = append(filteredWaves, wave)
		}
	}

	return filteredWaves
}
