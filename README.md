<!-- Code generated by mkbadge; DO NOT EDIT. START -->
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-green?logo=go)](https://pkg.go.dev/mod/github.com/nickwells/fsm.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/nickwells/fsm.mod)](https://goreportcard.com/report/github.com/nickwells/fsm.mod)
![GitHub License](https://img.shields.io/github/license/nickwells/fsm.mod)
<!-- Code generated by mkbadge; DO NOT EDIT. END -->

# fsm.mod
A finite state machine.

You create first a set of allowed state transitions and then you create an
FSM that takes an object through the allowed states. This manages a lifecycle
and lets you collect the rules into a single place where it is easy to
review. It saves you having logic about allowed states dispersed throughout
the codebase.
