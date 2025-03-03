package fsm

// state represents a state in a Finite State Machine. A terminal state is one
// with an empty nextState map
type state struct {
	name      string
	desc      string
	nextState map[string]*state
}

// newState returns a newly constructed state
func newState(name string) *state {
	return &state{
		name:      name,
		nextState: make(map[string]*state),
	}
}

// isTerminal returns whether the state is terminal. A state is terminal when
// it has no next states
func (s state) isTerminal() bool {
	return len(s.nextState) == 0
}

// (s state)String returns a string describing the state
func (s state) String() string {
	str := s.name
	if s.desc != "" {
		str += " [" + s.desc + "]"
	}

	return str
}
