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
		StartState: f.start,
		MatchState: &MATCH,
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
//   <1100 && N < 2 -> SUCCESS
//   <Freq -> PEAK(Freq, N+1)
//   >Freq -> FAIL
// SUCCESS
//   Anything -> SUCCESS
//   Nothing -> MATCH

type StartState int
type SustainState struct {
	Freq float64
}
type RisingState struct {
	N          int
	Freq       float64
	TargetFreq float64
}
type QuietState struct {
	Freq float64
	N    int
}
type FallingState struct {
	N          int
	Freq       float64
	TargetFreq float64
}
type SuccessState int

func (s *StartState) Name() string {
	return "START"
}
func (s *SustainState) Name() string {
	return fmt.Sprintf("SUSTAIN(%0.2f)", s.Freq)
}
func (s *RisingState) Name() string {
	return fmt.Sprintf("RISING(%d, %0.2f, %0.2f)", s.N, s.Freq, s.TargetFreq)
}
func (s *QuietState) Name() string {
	return fmt.Sprintf("QUIET(%d, %0.2f)", s.N, s.Freq)
}
func (s *FallingState) Name() string {
	return fmt.Sprintf("FALLING(%d, %0.2f, %0.2f)", s.N, s.Freq, s.TargetFreq)
}
func (s *SuccessState) Name() string {
	return "SUCCESS"
}

func (s *StartState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyInRange(900, 1300, waves):
		return &SustainState{strongestFrequencyInRange(900, 1300, waves)}
	default:
		return nil
	}
}

func (s *SustainState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyAbove(s.Freq+25, waves):
		return &RisingState{N: 0, Freq: strongestFrequencyAbove(s.Freq+25, waves), TargetFreq: s.Freq + 200}
	case hasFrequencyInRange(s.Freq-25, s.Freq+25, waves):
		return &SustainState{s.Freq}
	case len(waves) == 0:
		return &QuietState{N: 0, Freq: s.Freq}
	default:
		return nil
	}
}

func (s *RisingState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyAbove(s.Freq, waves):
		return &RisingState{N: s.N + 1, Freq: strongestFrequencyAbove(s.Freq, waves), TargetFreq: s.TargetFreq}
	case hasFrequencyBelow(s.Freq, waves) && s.Freq > s.TargetFreq:
		return &FallingState{N: 0, Freq: strongestFrequencyBelow(s.Freq, waves), TargetFreq: s.Freq - 500}
	default:
		return nil
	}
}

func (s *QuietState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyAbove(s.Freq+25, waves):
		return &RisingState{N: 0, Freq: strongestFrequencyAbove(s.Freq+25, waves), TargetFreq: s.Freq + 200}
	case len(waves) == 0 && s.N < 1:
		return &QuietState{N: s.N + 1, Freq: s.Freq}
	default:
		return nil
	}
}

func (s *FallingState) Handle(waves []SineWave) State {
	switch {
	case hasFrequencyBelow(s.TargetFreq, waves) && s.N >= 1:
		return new(SuccessState)
	case hasFrequencyInRange(s.Freq-400, s.Freq-25, waves):
		return &FallingState{N: s.N + 1, Freq: strongestFrequencyBelow(s.Freq, waves), TargetFreq: s.TargetFreq}
	case hasFrequencyInRange(s.Freq-25, s.Freq, waves):
		return &FallingState{N: s.N, Freq: strongestFrequencyBelow(s.Freq, waves), TargetFreq: s.TargetFreq}
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
