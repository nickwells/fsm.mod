/*
Package fsm manages the behaviour of a Finite State Machine.

# Basics

We create an FSM in two stages - firstly we create a set of allowed state
transitions:

	st, err := fsm.NewStateTrans("name",
	                 fsm.STPair{fsm.InitState, "start"},
	                 fsm.STPair{"start", "finish"})

Then we create the FSM itself giving it the set of allowed state transitions
and an associated Underlying (any type which satisfies the Underlying
interface):

	f := fsm.New(st, myThing)

This FSM will then enforce the rules given in the Transitions
object. Additional rules can be applied through the TransitionAllowed
function on the Underlying object.

The intention is that there can be multiple instances of an FSM all sharing
the same set of allowed state transitions

Each FSM has an associated underlying (which can be nil). This represents an
object whose state is being managed by the FSM. This underlying, if present,
must satisfy the Underlying interface.

The default starting state of every FSM is given by the fsm.InitState const
*/
package fsm
