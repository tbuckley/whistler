package whistler

import (
	"log"
)

type State interface {
	Handle([]SineWave) State
	Name() string
}

type Matcher struct {
	states     []State
	StartState State
	MatchState State
}

type MatchFactory interface {
	New() *Matcher
}

type stateTransition struct {
	Start State
	End   State
}

func (m *Matcher) Match(point []SineWave) bool {
	match := false
	states := make([]State, 0)

	transitions := make([]stateTransition, 0)

	// Iterate over existing states
	for _, state := range m.states {
		if out := state.Handle(point); out != nil {
			transitions = append(transitions, stateTransition{state, out})
			if out == m.MatchState {
				match = true
			} else {
				states = append(states, out)
			}
		}
	}

	// Always apply the START state
	if out := m.StartState.Handle(point); out != nil {
		transitions = append(transitions, stateTransition{m.StartState, out})
		states = append(states, out)
	}

	if len(transitions) > 0 {
		log.Printf("[STATE] =======================")
		for _, transition := range transitions {
			log.Printf("[STATE] %s -> %s", transition.Start.Name(), transition.End.Name())
		}
	}

	// Remove other state machines if there's a match
	if match {
		m.states = nil
	} else {
		m.states = states
	}

	return match
}
