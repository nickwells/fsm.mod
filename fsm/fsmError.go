package fsm

import "fmt"

// FSMErrUnknownState is an error type that represents an unknown state
type FSMErrUnknownState struct {
	FSMName string
	State   string
}

// mkErrUnknownState constructs and returns an unknown state error
func (f FSM) mkErrUnknownState(s string) FSMErrUnknownState {
	return FSMErrUnknownState{
		FSMName: f.Name(),
		State:   s,
	}
}

// Error returns a string form of the error
func (fe FSMErrUnknownState) Error() string {
	return fmt.Sprintf("FSM: %q: %q is not a known state",
		fe.FSMName, fe.State)
}

// FSMErrNoTransition is an error type that represents a non-existent
// transition between states. This is the case where the fsm has no single
// step between the two states.
type FSMErrNoTransition struct {
	FSMName   string
	FromState string
	ToState   string
}

// mkErrNoTransition constructs and returns an unknown state error
func (f FSM) mkErrNoTransition(s string) FSMErrNoTransition {
	return FSMErrNoTransition{
		FSMName:   f.Name(),
		FromState: f.current.name,
		ToState:   s,
	}
}

// Error returns a string form of the error
func (fe FSMErrNoTransition) Error() string {
	return fmt.Sprintf("FSM: %q: There is no valid transition from %q to %q",
		fe.FSMName, fe.FromState, fe.ToState)
}

// FSMErrForbiddenChange is an error type that represents a forbidden transition
// between states. This is the case where the fsm would allow the transition
// but the underlying check prevents it (TransitionAllowed returns a non-nil
// error)
type FSMErrForbiddenChange struct {
	FSMName   string
	FromState string
	ToState   string
	UndError  error
}

// mkErrForbiddenChange constructs and returns an unknown state error
func (f FSM) mkErrForbiddenChange(s string, ue error) FSMErrForbiddenChange {
	return FSMErrForbiddenChange{
		FSMName:   f.Name(),
		FromState: f.current.name,
		ToState:   s,
		UndError:  ue,
	}
}

// Error returns a string form of the error
func (fe FSMErrForbiddenChange) Error() string {
	return fmt.Sprintf("FSM: %q: The change from %q to %q is forbidden: %s",
		fe.FSMName, fe.FromState, fe.ToState, fe.UndError)
}
