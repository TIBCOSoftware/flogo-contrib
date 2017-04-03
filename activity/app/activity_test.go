package app

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"io/ioutil"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil{
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestAdd(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "ADD")
	tc.SetInput(ivType, "string")
	tc.SetInput(ivValue, "test")

	act.Eval(tc)

	value, _ := tc.GetAttrValue("value")

	if value != "test" {
		fmt.Println("Bad Value: " + value)
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

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//add attribute
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "GET")

	act.Eval(tc)

	value, _ := tc.GetOutput(ovValue).(string)

	if value != "test2" {
		fmt.Println("Bad Value: " + value)
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "UPDATE")
	tc.SetInput(ivValue, "test3")

	act.Eval(tc)

	value, _ := tc.GetAttrValue("myAttr")

	if value != "test3" {
		fmt.Println("Bad Value: " + value)
		t.Fail()
	}
}
