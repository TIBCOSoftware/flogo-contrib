package counter

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-counter")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestIncrement(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &CounterActivity{metadata: md, counters: make(map[string]int)}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivCounterName, "messages")
	tc.SetInput(ivIncrement, true)

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 1 {
		t.Fail()
	}
}

func TestGet(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)

	counters := map[string]int{
		"messages": 5,
	}

	act := &CounterActivity{metadata: md, counters: counters}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivCounterName, "messages")

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 5 {
		t.Fail()
	}
}

func TestReset(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	counters := map[string]int{
		"messages": 3,
	}

	act := &CounterActivity{metadata: md, counters: counters}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivCounterName, "messages")
	tc.SetInput(ivReset, true)

	act.Eval(tc)

	value := tc.GetOutput(ovValue).(int)

	if value != 0 {
		t.Fail()
	}
}
