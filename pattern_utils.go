package whistler

func filter(waves []SineWave, pred func(SineWave) bool) []SineWave {
	filteredWaves := make([]SineWave, 0)
	for _, wave := range waves {
		if pred(wave) {
			filteredWaves = append(filteredWaves, wave)
		}
	}
	return filteredWaves
}

func filterLowAmplitudes(waves []SineWave, min float64) []SineWave {
	return filter(waves, func(wave SineWave) bool {
		return wave.Amplitude > min
	})
}

func filterToFrequencyRange(waves []SineWave, low, high float64) []SineWave {
	return filter(waves, func(wave SineWave) bool {
		return wave.Frequency > low && wave.Frequency < high
	})
}
func filterToMinFrequency(waves []SineWave, min float64) []SineWave {
	return filter(waves, func(wave SineWave) bool {
		return wave.Frequency > min
	})
}
func filterToMaxFrequency(waves []SineWave, max float64) []SineWave {
	return filter(waves, func(wave SineWave) bool {
		return wave.Frequency < max
	})
}

func has(waves []SineWave, pred func(SineWave) bool) bool {
	for _, wave := range waves {
		if pred(wave) {
			return true
		}
	}
	return false
}

func hasFrequencyInRange(min, max float64, waves []SineWave) bool {
	return has(waves, func(wave SineWave) bool {
		return min < wave.Frequency && max > wave.Frequency
	})
}
func hasFrequencyAbove(freq float64, waves []SineWave) bool {
	return has(waves, func(wave SineWave) bool {
		return wave.Frequency > freq
	})
}
func hasFrequencyBelow(freq float64, waves []SineWave) bool {
	return has(waves, func(wave SineWave) bool {
		return wave.Frequency < freq
	})
}
func maximize(waves []SineWave, fn func(acc, wave SineWave) SineWave) SineWave {
	if len(waves) == 0 {
		return SineWave{}
	}

	base := waves[0]
	for _, wave := range waves[1:] {
		base = fn(base, wave)
	}
	return base
}
func highestFrequency(waves []SineWave) float64 {
	return maximize(waves, func(acc, wave SineWave) SineWave {
		if wave.Frequency > acc.Frequency {
			return wave
		}
		return acc
	}).Frequency
}
func waveWithStrongestAmplitude(waves []SineWave) SineWave {
	return maximize(waves, func(acc, wave SineWave) SineWave {
		if wave.Amplitude > acc.Amplitude {
			return wave
		}
		return acc
	})
}

func highestFrequencyBelow(freq float64, waves []SineWave) float64 {
	return highestFrequency(filterToFrequencyRange(waves, 0, freq))
}
func strongestFrequencyInRange(min float64, max float64, waves []SineWave) float64 {
	return waveWithStrongestAmplitude(filterToFrequencyRange(waves, min, max)).Frequency
}
func strongestFrequencyAbove(min float64, waves []SineWave) float64 {
	return waveWithStrongestAmplitude(filterToMinFrequency(waves, min)).Frequency
}
func strongestFrequencyBelow(max float64, waves []SineWave) float64 {
	return waveWithStrongestAmplitude(filterToMaxFrequency(waves, max)).Frequency
}
