package fsm_test

import (
	"fmt"

	"github.com/nickwells/fsm.mod/fsm"
)

const (
	FixInProgress  = "FixInProgress"
	ReadyToFix     = "ReadyToFix"
	ReadyToReview  = "ReadyToReview"
	ReadyToTest    = "ReadyToTest"
	Rejected       = "Rejected"
	Released       = "Released"
	TestFailed     = "TestFailed"
	TestInProgress = "TestInProgress"
	TestPassed     = "TestPassed"
	UnderReview    = "UnderReview"
)

var transitionToChecks = map[string]func(b *BugReport, f *fsm.FSM) error{
	TestPassed: checkTestEvidence,
}

// BugReport is a toy structure representing a bug. The key point of
// interest is that we embed the FSM in the structure. This allows us to
// call FSM methods on the BugReport objects directly
type BugReport struct {
	*fsm.FSM
	name         string
	description  string
	testEvidence string
}

// checkTestEvidence returns an error if the bug report doesn't have test
// evidence
func checkTestEvidence(b *BugReport, _ *fsm.FSM) error {
	if b.testEvidence == "" {
		return fmt.Errorf("missing test evidence")
	}
	return nil
}

// TransitionAllowed satisfies the Underlying interface. It can be used to
// perform extra checks on state transitions in addition to the checks that
// the FSM will perform. Here we make sure that there is test evidence
// attached to the bug before we allow the state to be changed to TestPassed
func (b *BugReport) TransitionAllowed(f *fsm.FSM, newState string) error {
	if cf, ok := transitionToChecks[newState]; ok {
		err := cf(b, f)
		if err != nil {
			return err
		}
	}
	if newState == TestPassed {
		if b.testEvidence == "" {
			return fmt.Errorf("missing test evidence")
		}
	}
	return nil
}

// OnTransition satisfies the Underlying interface. It can be used to
// perform custom actions when certain states are reached or when states are
// left. Here we notify the release team when the tests for a bugfix have
// all passed
func (b *BugReport) OnTransition(f *fsm.FSM) {
	if f.CurrentState() == TestPassed {
		notifyReleaseTeam(b)
	}
}

// SetFSM satisfies the Underlying interface. It can be used to set the
// underlying's embedded FSM as here
func (b *BugReport) SetFSM(f *fsm.FSM) {
	b.FSM = f
}

// notifyReleaseTeam will tell the release team that the bugfix is ready
// to be released
func notifyReleaseTeam(b *BugReport) {
	fmt.Println("The fix to bug:", b.name,
		"(", b.description, ") is ready to release")
}

// Example_ex1 is an example of how the fsm can be used to manage the
// lifecycle of some business object; in this case, a bug.
func Example_ex1() {

	// now construct the set of allowed state transitions
	st, err := fsm.NewStateTrans("BugReport", []fsm.STPair{
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
	_ = st.SetStateDesc(ReadyToReview,
		"The bug is ready to review")
	_ = st.SetStateDesc(UnderReview,
		"The bug is being reviewed")
	_ = st.SetStateDesc(Rejected,
		"The bug has been rejected")
	_ = st.SetStateDesc(ReadyToFix,
		"The bug has been accepted and work can start")
	_ = st.SetStateDesc(FixInProgress,
		"The bug is being worked on")
	_ = st.SetStateDesc(ReadyToTest,
		"The bug has been fixed and is ready to test")
	_ = st.SetStateDesc(TestInProgress,
		"The bug is being tested")
	_ = st.SetDescriptions(
		fsm.StateDesc{TestPassed, "The bug fix has passed the tests"},
		fsm.StateDesc{TestFailed, "The bug fix has failed the tests"},
		fsm.StateDesc{Released, "The bug fix has been released"})

	b := BugReport{name: "ID-1234"}

	// now we can create a new FSM which allows state transitions as given
	// by 'st' and with an underlying of 'b'. The SetFSM function on 'b'
	// will be called
	fsm.New(st, &b)

	sampleStateChanges := []string{
		ReadyToReview,
		"Any",
		Released,
		UnderReview,
		Rejected,
	}

	fmt.Println("Current state:", b.CurrentState())
	for _, newState := range sampleStateChanges {
		fmt.Println("Changing to:  ", newState)
		err := b.ChangeState(newState)
		if err != nil {
			fmt.Println("ERROR:         Cannot change from",
				b.CurrentState(), "to", newState, ": ", err)
		}
		fmt.Println("Current state:", b.CurrentState())
		if b.IsInTerminalState() {
			fmt.Println("done")
			break
		}
	}
	// Output:
	// Current state: init
	// Changing to:   ReadyToReview
	// Current state: ReadyToReview
	// Changing to:   Any
	// ERROR:         Cannot change from ReadyToReview to Any :  FSM: "BugReport": "Any" is not a known state
	// Current state: ReadyToReview
	// Changing to:   Released
	// ERROR:         Cannot change from ReadyToReview to Released :  FSM: "BugReport": There is no valid transition from "ReadyToReview" to "Released"
	// Current state: ReadyToReview
	// Changing to:   UnderReview
	// Current state: UnderReview
	// Changing to:   Rejected
	// Current state: Rejected
	// done
}
