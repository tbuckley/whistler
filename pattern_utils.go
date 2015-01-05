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
func highestFrequency(waves []SineWave) float64 {
	max := 0.0
	for _, wave := range waves {
		if wave.Frequency > max {
			max = wave.Frequency
		}
	}
	return max
}
func highestFrequencyBelow(freq float64, waves []SineWave) float64 {
	return highestFrequency(filterToFrequencyRange(waves, 0, freq))
}
