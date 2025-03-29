package fsm

import "fmt"

// Error is the type of an error from this package
type Error interface {
	error

	// FSMError is a no-op function but it serves to
	// distinguish errors from this package from other errors
	FSMError()
}

// UnknownState is an error type that represents an unknown state
type UnknownState struct {
	FSMName string
	State   string
}

// mkErrUnknownState constructs and returns an UnknownState error
func (f FSM) mkErrUnknownState(s string) UnknownState {
	return UnknownState{
		FSMName: f.Name(),
		State:   s,
	}
}

// Error returns a string form of the error
func (fe UnknownState) Error() string {
	return fmt.Sprintf("FSM: %q: %q is not a known state",
		fe.FSMName, fe.State)
}

// FSMError is a no-op function used to distinguish between an fsm package
// error and other error types.
func (UnknownState) FSMError() {}

// NoTransition is an error type that represents a non-existent
// transition between states. This is the case where the fsm has no single
// step between the two states.
type NoTransition struct {
	FSMName   string
	FromState string
	ToState   string
}

// mkErrNoTransition constructs and returns a NoTransition error
func (f FSM) mkErrNoTransition(s string) NoTransition {
	return NoTransition{
		FSMName:   f.Name(),
		FromState: f.current.name,
		ToState:   s,
	}
}

// Error returns a string form of the error
func (fe NoTransition) Error() string {
	return fmt.Sprintf("FSM: %q: There is no valid transition from %q to %q",
		fe.FSMName, fe.FromState, fe.ToState)
}

// FSMError is a no-op function used to distinguish between an fsm package
// error and other error types.
func (NoTransition) FSMError() {}

// ForbiddenChange is an error type that represents a forbidden transition
// between states. This is the case where the fsm would allow the transition
// but the underlying check prevents it (TransitionAllowed returns a non-nil
// error)
type ForbiddenChange struct {
	FSMName   string
	FromState string
	ToState   string
	UndError  error
}

// mkErrForbiddenChange constructs and returns a ForbiddenChange error
func (f FSM) mkErrForbiddenChange(s string, ue error) ForbiddenChange {
	return ForbiddenChange{
		FSMName:   f.Name(),
		FromState: f.current.name,
		ToState:   s,
		UndError:  ue,
	}
}

// Unwrap returns the underlying error
func (fe ForbiddenChange) Unwrap() error {
	return fe.UndError
}

// Error returns a string form of the error
func (fe ForbiddenChange) Error() string {
	return fmt.Sprintf("FSM: %q: The change from %q to %q is forbidden: %s",
		fe.FSMName, fe.FromState, fe.ToState, fe.UndError)
}

// FSMError is a no-op function used to distinguish between an fsm package
// error and other error types.
func (ForbiddenChange) FSMError() {}
