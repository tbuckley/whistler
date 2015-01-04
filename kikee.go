package whistler

import (
	"fmt"
)

var (
	MATCH SuccessState = 0
)

type kikeeFactory struct {
	start *StartState
}

func (f *kikeeFactory) New() *Matcher {
	return &Matcher{
		start: f.start,
		match: &MATCH,
	}
}

var (
	Kikee = &kikeeFactory{
		start: new(StartState),
	}
)

// In range 950-1150
// Then, optional short pause
// Then, hits above 1300 and decreases over X+ points to below 950

// START:
//   Freq in range 950-1150 -> INIT
// INIT:
//   Freq in range 950-1150 -> INIT
//   Nothing -> QUIET(0)
//   >1300 -> PEAK(Freq, 0)
// QUIET(N)
//   >1300 -> PEAK(Freq, 0)
//   Nothing && N < 5 -> QUIET(N+1)
//   Nothing && N == 5 -> FAIL
// PEAK(Freq, N)
//   <Freq && N < 4 -> PEAK(Freq, N+1)
//   <Freq && N == 4 -> SUCCESS
//   >Freq -> FAIL
// SUCCESS
//   Anything -> SUCCESS
//   Nothing -> MATCH

type StartState int
type InitState int
type QuietState struct {
	N int
}
type PeakState struct {
	N    int
	Freq float64
}
type SuccessState int

func (s *StartState) Name() string {
	return "START"
}
func (s *InitState) Name() string {
	return "INIT"
}
func (s *QuietState) Name() string {
	return fmt.Sprintf("QUIET(%d)", s.N)
}
func (s *PeakState) Name() string {
	return fmt.Sprintf("PEAK(%d, %0.2f)", s.N, s.Freq)
}
func (s *SuccessState) Name() string {
	return "SUCCESS"
}

func (s *StartState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyInRange(950, 1150, waves):
		return new(InitState)
	default:
		return nil
	}
}

func (s *InitState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyInRange(950, 1150, waves):
		// fmt.Printf("%v: %v", len(waves), hasFrequencyInRange(950, 1150, waves))
		return new(InitState)
	case len(waves) == 0:
		return &QuietState{N: 0}
	case hasFrequencyAbove(1300, waves):
		return &PeakState{N: 0, Freq: highestFrequency(waves)}
	default:
		return nil
	}
}

func (s *QuietState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyAbove(1300, waves):
		return &PeakState{N: 0, Freq: highestFrequency(waves)}
	case len(waves) == 0 && s.N < 5:
		return &QuietState{N: s.N + 1}
	default:
		return nil
	}
}

func (s *PeakState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyBelow(s.Freq, waves) && s.N == 3:
		return new(SuccessState)
	case hasFrequencyBelow(s.Freq, waves):
		return &PeakState{N: s.N + 1, Freq: highestFrequency(waves)}
	default:
		return nil
	}
}

func (s *SuccessState) Handle(waves []SineWave) State {
	switch {
	case len(waves) > 0:
		return new(SuccessState)
	default:
		return &MATCH
	}
}
