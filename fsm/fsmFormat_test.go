package fsm_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/fsm.mod/fsm"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFormat(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		formatStr string
		expOut    string
	}{
		{
			ID:        testhelper.MkID("simple 'v' format"),
			formatStr: "%v",
			expOut:    "State: finish",
		},
		{
			ID:        testhelper.MkID("expanded 'v' format"),
			formatStr: "%#v",
			expOut:    "FSMType: lifecycle: State: finish: PriorState: start",
		},
		{
			ID:        testhelper.MkID("simple 's' format"),
			formatStr: "%s",
			expOut:    "finish",
		},
		{
			ID:        testhelper.MkID("expanded 's' format"),
			formatStr: "%#s",
			expOut:    "lifecycle: finish (was: start [the first state])",
		},
		{
			ID:        testhelper.MkID("unsupported format"),
			formatStr: "%d",
			expOut:    "%!d(FSM=finish)",
		},
		{
			ID:        testhelper.MkID("unsupported format - expanded"),
			formatStr: "%#d",
			expOut: "%!d(FSM=lifecycle: finish" +
				" (was: start [the first state]))",
		},
	}

	st, _ := fsm.NewStateTrans("lifecycle",
		fsm.STPair{fsm.InitState, "start"},
		fsm.STPair{"start", "finish"})
	_ = st.SetStateDesc("start", "the first state")
	f := fsm.New(st, nil)
	_ = f.ChangeState("start")
	_ = f.ChangeState("finish")

	for _, tc := range testCases {
		out := fmt.Sprintf(tc.formatStr, f)
		testhelper.DiffString(t,
			tc.IDStr()+": format string: "+tc.formatStr, "output",
			out, tc.expOut)
	}
}
