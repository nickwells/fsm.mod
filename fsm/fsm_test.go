package fsm_test

import (
	"errors"
	"testing"

	"github.com/nickwells/fsm.mod/fsm"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

type underlying struct {
	allowChange             bool
	transitionAllowedCalled bool
	onTransitionCalled      bool
	setFSMCallCount         int
}

const undErrStr = "the underlying forbids the change"

// (u underlying)TransitionAllowed ...
func (u *underlying) TransitionAllowed(_ *fsm.FSM, _ string) error {
	u.transitionAllowedCalled = true
	if u.allowChange {
		return nil
	}

	return errors.New(undErrStr)
}

// (u underlying)OnTransition ...
func (u *underlying) OnTransition(_ *fsm.FSM) {
	u.onTransitionCalled = true
}

// (u underlying)SetFSM ...
func (u *underlying) SetFSM(_ *fsm.FSM) {
	u.setFSMCallCount++
}

// (u underlying)Reset ...
func (u *underlying) Reset() {
	u.allowChange = false
	u.transitionAllowedCalled = false
	u.onTransitionCalled = false
	u.setFSMCallCount = 0
}

func TestStateTransAdd(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		states []fsm.STPair
		testhelper.ExpErr
	}{
		{
			ID:     testhelper.MkID("bad transition"),
			states: []fsm.STPair{{From: "nonesuch", To: "any"}},
			ExpErr: testhelper.MkExpErr(
				"testStateTrans: state: \"nonesuch\" does not exist", "failed"),
		},
		{
			ID:     testhelper.MkID("init transition"),
			states: []fsm.STPair{{From: fsm.InitState, To: "state1"}},
		},
		{
			ID: testhelper.MkID("multi transition"),
			states: []fsm.STPair{
				{From: fsm.InitState, To: "state1"},
				{From: "state1", To: "state2"},
			},
		},
	}

	for _, tc := range testCases {
		_, err := fsm.NewStateTrans("testStateTrans", tc.states...)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestStateTransSet(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		states []fsm.STPair
		testhelper.ExpErr
		statesExpected []string
	}{
		{
			ID: testhelper.MkID("good - one tran"),
			states: []fsm.STPair{
				{fsm.InitState, "A"},
			},
			statesExpected: []string{fsm.InitState, "A"},
		},
		{
			ID: testhelper.MkID("good - multi tran"),
			states: []fsm.STPair{
				{fsm.InitState, "A"},
				{fsm.InitState, "B"},
				{fsm.InitState, "C"},
				{fsm.InitState, "D"},
				{"D", "E"},
			},
			statesExpected: []string{fsm.InitState, "A", "B", "C", "D", "E"},
		},
		{
			ID: testhelper.MkID("bad - multi tran"),
			states: []fsm.STPair{
				{"X", "A"},
				{fsm.InitState, "B"},
				{fsm.InitState, "C"},
				{fsm.InitState, "D"},
				{"D", "E"},
			},
			ExpErr: testhelper.MkExpErr(
				"testStateTrans: state: \"X\" does not exist", "failed"),
			statesExpected: []string{fsm.InitState},
		},
	}

	for _, tc := range testCases {
		st, err := fsm.NewStateTrans("testStateTrans", tc.states...)
		if testhelper.CheckExpErr(t, err, tc) && err == nil {
			testhelper.DiffInt(t, tc.IDStr(), "number of states",
				st.StateCount(), len(tc.statesExpected))

			for _, state := range tc.statesExpected {
				if !st.HasState(state) {
					t.Log(tc.IDStr())
					t.Errorf("\t: state: %q was expected but not found", state)
				}
			}
		}
	}
}

func TestFSM_BadST(t *testing.T) {
	badST, err := fsm.NewStateTrans("badStateTrans",
		fsm.STPair{"nonesuch", "state1"})
	if err == nil {
		t.Error("making a StateTrans with bad transitions should give an error")
	} else if badST != nil {
		t.Error("making a StateTrans with bad transitions should return a nil")
	} else {
		f := fsm.New(badST, nil)
		if f != nil {
			t.Errorf("a bad ST should not create an FSM\n")
		}
	}
}

func TestFsm(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		uState             bool
		newState           string
		expectedUnderlying underlying
	}{
		{
			ID:       testhelper.MkID("bad transition - no such state"),
			uState:   true,
			newState: "noSuchState",
			ExpErr:   testhelper.MkExpErr("is not a known state"),
			expectedUnderlying: underlying{
				allowChange:             true,
				transitionAllowedCalled: false,
				onTransitionCalled:      false,
			},
		},
		{
			ID:       testhelper.MkID("bad transition - no valid transition"),
			uState:   true,
			newState: "final",
			ExpErr: testhelper.MkExpErr(
				"There is no valid transition from"),
			expectedUnderlying: underlying{
				allowChange:             true,
				transitionAllowedCalled: false,
				onTransitionCalled:      false,
			},
		},
		{
			ID:       testhelper.MkID("bad transition - not allowed"),
			uState:   false,
			newState: "state1",
			ExpErr: testhelper.MkExpErr(
				"The change",
				"is forbidden",
				undErrStr),
			expectedUnderlying: underlying{
				allowChange:             false,
				transitionAllowedCalled: true,
				onTransitionCalled:      false,
			},
		},
		{
			ID:       testhelper.MkID("good transition"),
			uState:   true,
			newState: "state1",
			expectedUnderlying: underlying{
				allowChange:             true,
				transitionAllowedCalled: true,
				onTransitionCalled:      true,
			},
		},
	}

	var err error

	st, err := fsm.NewStateTrans("testFSM",
		[]fsm.STPair{
			{fsm.InitState, "state1"},
			{fsm.InitState, "state2"},
			{"state2", "final"},
			{"state1", "final"},
		}...)
	if err != nil {
		t.Fatal("couldn't setup the test:", err)
	}

	if stateCount := st.StateCount(); stateCount != 4 {
		t.Errorf("Unexpected StateCount: expected 4, got %d\n", stateCount)
	}

	u := underlying{}
	f := fsm.New(st, &u)

	if u.setFSMCallCount != 1 {
		t.Errorf("Unexpected SetFSMCallCount: expected 1, got %d\n",
			u.setFSMCallCount)
	}

	nextStates := f.NextStates()
	if len(nextStates) != 2 {
		t.Fatalf("Unexpected number of next states: expected 2, got %d\n",
			len(nextStates))
	}

	if !f.IsInInitialState() {
		t.Fatalf("initially, the fsm should be in the initial state. %+v\n", f)
	}

	if f.CurrentState() != f.PriorState() {
		t.Fatalf(
			"initially, current and prior states should be equal: %+v\n", f)
	}

	for _, tc := range testCases {
		u.Reset()
		u.allowChange = tc.uState
		err = f.ChangeState(tc.newState)
		testhelper.CheckExpErr(t, err, tc)

		if u != tc.expectedUnderlying {
			t.Log(tc.IDStr())
			t.Errorf("\t: Unexpected underlying value. Expected: %v, got: %v",
				tc.expectedUnderlying, u)
		}
	}
}
