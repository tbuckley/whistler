package whistler

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type SineWave struct {
	Amplitude   float64
	Frequency   float64
	PhaseOffset float64
}

func (w SineWave) String() string {
	return fmt.Sprintf("Amp=%v, Freq=%v, Phase=%v", w.Amplitude, w.Frequency, w.PhaseOffset)
}

func (w SineWave) Value(t float64) float64 {
	return w.Amplitude * math.Sin(t*2*math.Pi*w.Frequency+w.PhaseOffset)
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

// WaveSet is a series of SineWaves.
type WaveSet []SineWave

// String converts a WaveSet into a readable representation.
func (s WaveSet) String() string {
	strs := make([]string, len(s))
	for i, w := range s {
		strs[i] = w.String()
	}
	return strings.Join(strs, "\n")
}

// Value calculates the value of the combined SineWaves at the given time.
func (s WaveSet) Value(t float64) float64 {
	v := 0.0
	for _, w := range s {
		v += w.Value(t)
	}
	return v
}

// Len provides the number of waves in the set. Used for sort.Interface.
func (s WaveSet) Len() int {
	return len(s)
}

// Swap exchanges the positions of two waves in the set. Used for sort.Interface.
func (s WaveSet) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less compares the values of two waves in the set. Used for sort.Interface.
func (s WaveSet) Less(i, j int) bool {
	return s[i].Amplitude > s[j].Amplitude
}

// Sort orders the waves such that the highest amplitudes appear first.
func (s WaveSet) Sort() {
	sort.Sort(s)
}
