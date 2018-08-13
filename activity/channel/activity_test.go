package channel

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/engine/channels"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	doneCh := make(chan bool)
	defer close(doneCh)

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	channels.Add("test")
	ch := channels.Get("test")

	//setup attrs
	tc.SetSetting(sChannel, "test")
	tc.SetInput(ivValue, 2)

	var done bool
	var err error

	go func() {
		done, err = act.Eval(tc)
		doneCh <- true
	}()

	<-doneCh // blocks until the input write routine is finished

	expected := 2
	found := <-ch // blocks until the output has contents

	if found != expected {
		t.Errorf("Expected %s, found %s", expected, found)
	}

	if !done {
		t.Error("activity should be done")
		return
	}

	if err != nil {
		t.Error("activity has an error: ", err)
		return
	}
}

func TestProcess(t *testing.T) {
	// GIVEN
	input := make(chan string)
	defer close(input)

	done := make(chan bool)
	defer close(done)

	go func() {
		input <- "hello world"
		done <- true
	}()

	// WHEN
	output := Process(input)
	<-done // blocks until the input write routine is finished

	// THEN
	expected := "(hello world)"
	found := <-output // blocks until the output has contents

	if found != expected {
		t.Errorf("Expected %s, found %s", expected, found)
	}
}