package fsm

import (
	"fmt"
	"sort"
)

// Underlying is an interface representing a set of functions to be called
// around various FSM transitions. The Underlying can be used to represent an
// object whose state is managed by the FSM
type Underlying interface {
	// TransitionAllowed is called before the change of state to the newState.
	// If it returns an error the transition is not performed
	TransitionAllowed(f *FSM, newState string) error
	// OnTransition is called after the state of the FSM has been successsfully
	// changed from the old state to the new state
	OnTransition(f *FSM)
	// SetFSM is called when the Underlying object is set on the FSM. It
	// provides a way for the Underlying object to know the FSM(s) with which
	// it is associated.
	//
	// Note that an Underlying could be associated with multiple FSMs at the
	// same time in which case a map from FSM name to FSM pointer might be
	// appropriate
	//
	// if the Underlying has only a single associated FSM then the FSM pointer
	// could be embedded in the Underlying. This would allow the FSM methods to
	// be called on the Underlying directly
	SetFSM(f *FSM)
}

// FSM represents a Finite State Machine
type FSM struct {
	st      *StateTrans
	prior   *state
	current *state
	und     Underlying
}

// New creates a new Finite State Machine. It returns nil if the StateTrans
// is nil.
//
// The prior and current states are set to InitState.
//
// The SetFSM method on the Underlying is called with the new FSM so that the
// Underlying can store the associated FSM if required.
func New(st *StateTrans, u Underlying) *FSM {
	if st == nil {
		return nil
	}

	f := &FSM{
		st:      st,
		prior:   st.states[InitState],
		current: st.states[InitState],
		und:     u,
	}
	if u != nil {
		u.SetFSM(f)
	}
	return f
}

// Name returns the name of the Finite State Machine
func (f *FSM) Name() string {
	return f.st.name
}

// CurrentState returns the name of the current state of the FSM
func (f *FSM) CurrentState() string {
	return f.current.name
}

// PriorState returns the name of the prior state of the FSM
func (f *FSM) PriorState() string {
	return f.prior.name
}

// IsInTerminalState returns true if the FSM is in a terminal state
func (f *FSM) IsInTerminalState() bool {
	return f.current.isTerminal()
}

// IsInInitialState returns true if the FSM is in the initial state
func (f *FSM) IsInInitialState() bool {
	return f.current.name == InitState
}

// NextStates returns a sorted slice containing the names of the valid next
// states of the FSM
func (f *FSM) NextStates() []string {
	states := make([]string, 0, len(f.current.nextState))
	for _, s := range f.current.nextState {
		states = append(states, s.name)
	}
	sort.StringSlice(states).Sort()
	return states
}

// ChangeState changes the state from the current state to the new state
// provided the new state is a valid transition from the current state of the
// FSM and the transition is allowed by the Underlying TransitionAllowed
// function. Following the change of state the Underlying OnTransition function
// is called
func (f *FSM) ChangeState(newState string) error {
	oldState := f.current
	state, ok := f.current.nextState[newState]

	if !ok {
		if !f.st.HasState(newState) {
			return f.mkErrUnknownState(newState)
		}
		return f.mkErrNoTransition(newState)
	}

	if f.und != nil {
		if err := f.und.TransitionAllowed(f, newState); err != nil {
			return f.mkErrForbiddenChange(newState, err)
		}
	}

	f.prior = oldState
	f.current = state

	if f.und != nil {
		f.und.OnTransition(f)
	}
	return nil
}

// Format is used by the fmt package in the standard library to format the
// FSM. It supports two formats:
//
//	%s which prints the current state
//	%v which prints the current state with a label "State: "
//
// Either of these can be given the '#' flag which causes them to also print
// the FSM name and the previous state. Additionally the %s format will print
// any state descriptions.
func (f FSM) Format(fstate fmt.State, c rune) {
	str := ""

	switch c {
	case 'v':
		str = f.vFormat(fstate)
	case 's':
		str = f.sFormat(fstate)
	default:
		str += "%!" + string(c) +
			"(FSM=" +
			f.sFormat(fstate) +
			")"
	}

	_, _ = fstate.Write([]byte(str))
}

// sFormat returns a string representing the FSM formatted according to the
// 's' verb
func (f FSM) sFormat(fstate fmt.State) string {
	s := ""
	if fstate.Flag('#') {
		s += f.Name() + ": " +
			f.current.String() +
			" (was: " + f.prior.String() + ")"
	} else {
		s += f.current.name
	}
	return s
}

// vFormat returns a string representing the FSM formatted according to the
// 'v' verb
func (f FSM) vFormat(fstate fmt.State) string {
	s := ""
	if fstate.Flag('#') {
		s += "FSMType: " + f.Name() + ": "
	}
	s += "State: " + f.current.name
	if fstate.Flag('#') {
		s += ": PriorState: " + f.prior.name
	}
	return s
}
