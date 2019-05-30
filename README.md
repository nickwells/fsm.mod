[![GoDoc](https://godoc.org/github.com/nickwells/fsm.mod?status.png)](https://godoc.org/github.com/nickwells/fsm.mod)

# fsm.mod
A finite state machine.

You create first a set of allowed state transitions and then you create an
FSM that takes an object through the allowed states. This manages a lifecycle
and lets you collect the rules into a single place where it is easy to
review. It saves you having logic about allowed states dispersed throughout
the codebase.
