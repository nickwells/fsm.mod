package fsm_test

import (
	"fmt"
	"os"

	"github.com/nickwells/fsm.mod/fsm"
)

// these consts are already defined in example 1
//
// const (
// 	FixInProgress  = "FixInProgress"
// 	ReadyToFix     = "ReadyToFix"
// 	ReadyToReview  = "ReadyToReview"
// 	ReadyToTest    = "ReadyToTest"
// 	Rejected       = "Rejected"
// 	Released       = "Released"
// 	TestFailed     = "TestFailed"
// 	TestInProgress = "TestInProgress"
// 	TestPassed     = "TestPassed"
// 	UnderReview    = "UnderReview"
// )

// Example_ex2 is an example of how the StateTrans struct can be used
func Example_ex2() {

	// now construct the set of allowed state transitions
	st, err := fsm.NewStateTrans("my \"best\" ST graph", []fsm.STPair{
		{fsm.InitState, ReadyToReview},
		{ReadyToReview, UnderReview},
		{UnderReview, ReadyToReview},
		{UnderReview, Rejected},
		{UnderReview, ReadyToFix},
		{ReadyToFix, UnderReview},
		{ReadyToFix, FixInProgress},
		{FixInProgress, ReadyToFix},
		{FixInProgress, ReadyToTest},
		{ReadyToTest, TestInProgress},
		{TestInProgress, TestPassed},
		{TestInProgress, TestFailed},
		{TestFailed, FixInProgress},
		{TestPassed, Released},
	}...)
	if err != nil {
		fmt.Println("There was a problem initialising the transitions:", err)
		return
	}

	st.PrintDot(os.Stdout)
	// Output:
	// // A state transition graph for
	// //       my "best" ST graph
	// digraph st {
	//     node [shape = doublecircle
	//           style=filled fillcolor=lightblue];
	//         "init";
	//     node [shape = doublecircle
	//           style=filled fillcolor=grey85];
	//         "Rejected", "Released";
	//     node [shape = circle style=solid];
	//     { rank = same;
	//         "Rejected" "Released" }
	//     "FixInProgress" -> "ReadyToFix"
	//     "FixInProgress" -> "ReadyToTest"
	//     "ReadyToFix" -> "FixInProgress"
	//     "ReadyToFix" -> "UnderReview"
	//     "ReadyToReview" -> "UnderReview"
	//     "ReadyToTest" -> "TestInProgress"
	//     "TestFailed" -> "FixInProgress"
	//     "TestInProgress" -> "TestFailed"
	//     "TestInProgress" -> "TestPassed"
	//     "TestPassed" -> "Released"
	//     "UnderReview" -> "ReadyToFix"
	//     "UnderReview" -> "ReadyToReview"
	//     "UnderReview" -> "Rejected"
	//     "init" -> "ReadyToReview"
	//     fontsize=22
	//     label = "
	// my \"best\" ST graph
	// "
	// }
}
