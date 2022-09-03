package fsm

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const InitState = "init"

// StateTrans records the valid state changes.
type StateTrans struct {
	name   string
	states map[string]*state
}

// StateDesc records a state name and an associated description
type StateDesc struct {
	Name string
	Desc string
}

// STPair represents a state transition. A slice of these should be passed to
// the NewStateTrans method in order to set the valid transitions.
type STPair struct {
	From, To string
}

// NewStateTrans creates a new set of State transitions. The allowed
// transitions must be set at creation time by passing STPair's. If setting
// the transitions from the passed slice returns an error then this function
// will return a nil StateTrans and the error, otherwise the newly created
// StateTrans is returned and a nil error.
//
// Note that the order of creation of states is important - states must
// be created before a transition can be created from them. Every set of
// state transitions has a starting state given by the fsm.InitState
// const. New states will be created automatically if the 'to' state doesn't
// exist but the 'from' state must always exist. The state named from the
// InitState constant is always present in each StateTrans.
//
// The name has no semantic meaning and is only used for documentation
// purposes.
func NewStateTrans(name string, transitions ...STPair) (*StateTrans, error) {
	st := &StateTrans{
		name:   name,
		states: make(map[string]*state),
	}

	is := newState(InitState)
	is.desc = "the initial state"
	st.states[InitState] = is

	err := st.set(transitions...)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// HasState return true if the StateTrans object contains a state with the
// given name.
func (st StateTrans) HasState(name string) bool {
	_, ok := st.states[name]
	return ok
}

// add adds a new transition from one state in the FSM to another.
//
// The 'from' state must already exist in the FSM so the order of adding
// changes is important. If the 'from' state doesn't exist an error will be
// returned. This is to ensure that every state can be reached.
func (st *StateTrans) add(from, to string) error {
	fromState, ok := st.states[from]
	if !ok {
		return fmt.Errorf(
			"%s: state: '%s' does not exist. Add('%s', '%s') failed",
			st.name, from, from, to)
	}

	toState, ok := st.states[to]
	if !ok {
		toState = newState(to)
		st.states[to] = toState
	}
	fromState.nextState[toState.name] = toState

	return nil
}

// set sets the transitions from the slice of state transition pairs. The same
// rules apply as for the add func about the existence of the 'From' state
// before trying to switch to the 'To' state
func (st *StateTrans) set(transitions ...STPair) error {
	for _, stp := range transitions {
		err := st.add(stp.From, stp.To)
		if err != nil {
			return err
		}
	}
	return nil
}

// Name returns the name of the collection of StateTrans
func (st StateTrans) Name() string {
	return st.name
}

// StateCount returns a count of the number of states
func (st StateTrans) StateCount() int {
	return len(st.states)
}

// SetStateDesc sets the state description. It will return an error if the
// named state does not exist.
func (st *StateTrans) SetStateDesc(name, desc string) error {
	s, ok := st.states[name]
	if !ok {
		return fmt.Errorf("%s: state: %q does not exist", st.name, name)
	}

	s.desc = desc
	return nil
}

// SetDescriptions sets the state descriptions from the values given in the
// slice of state descriptions. It will return an error if any
// state does not exist and will set the error state on the StateTrans. It will
// also return an error if the StateTrans already has its error state set.
func (st *StateTrans) SetDescriptions(descriptions ...StateDesc) error {
	for _, sd := range descriptions {
		err := st.SetStateDesc(sd.Name, sd.Desc)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrintDot prints the state transitions as a directed graph in the
// graphviz DOT language. The output of this func can be interpreted by the
// dot command (on Linux). To generate a png file from this you could write
// the output to a file called, for instance, stateTrans.gv and then use the
// following command to generate the png image and write it to a file called
// graph.png
//
//	dot -Tpng -ograph.png stateTrans.gv
//
// This might be useful for generating documentation for your package.
func (st StateTrans) PrintDot(w io.Writer) {
	safeNames := make(map[string]string)
	for stateName := range st.states {
		safeNames[stateName] = strings.ReplaceAll(stateName, "\"", "\\\"")
	}

	namesInOrder := make([]string, 0, len(st.states))
	for name := range st.states {
		namesInOrder = append(namesInOrder, name)
	}
	sort.Strings(namesInOrder)

	fmt.Fprintln(w, "// A state transition graph for")
	fmt.Fprintln(w, "//      ", st.name)
	fmt.Fprintln(w, "digraph st {")

	fmt.Fprintln(w, "    node [shape = doublecircle")
	fmt.Fprintln(w, "          style=filled fillcolor=lightblue];")
	fmt.Fprintf(w, "        \"%s\"", InitState)
	fmt.Fprintln(w, ";")
	fmt.Fprintln(w, "    node [shape = doublecircle")
	fmt.Fprintln(w, "          style=filled fillcolor=grey85];")
	sep := ""
	fmt.Fprint(w, "        ")
	for _, name := range namesInOrder {
		s := st.states[name]
		if s.isTerminal() {
			fmt.Fprintf(w, "%s\"%s\"", sep, safeNames[name])
			sep = ", "
		}
	}
	fmt.Fprintln(w, ";")
	fmt.Fprintln(w, "    node [shape = circle style=solid];")

	fmt.Fprintln(w, "    { rank = same;")
	fmt.Fprint(w, "        ")
	for _, name := range namesInOrder {
		s := st.states[name]
		if s.isTerminal() {
			fmt.Fprintf(w, "\"%s\" ", safeNames[name])
		}
	}
	fmt.Fprintln(w, "}")

	for _, name := range namesInOrder {
		s := st.states[name]
		nextNamesInOrder := make([]string, 0, len(s.nextState))
		for _, ns := range s.nextState {
			nextNamesInOrder = append(nextNamesInOrder, ns.name)
		}
		sort.Strings(nextNamesInOrder)

		for _, nextName := range nextNamesInOrder {
			fmt.Fprintf(w, "    \"%s\" -> \"%s\"\n",
				safeNames[name],
				safeNames[nextName])
		}
	}

	fmt.Fprintln(w, "    fontsize=22")

	fmt.Fprintf(w, "    label = \"\n%s\n\"\n",
		strings.ReplaceAll(st.name, "\"", "\\\""))

	fmt.Fprintln(w, "}")
}
