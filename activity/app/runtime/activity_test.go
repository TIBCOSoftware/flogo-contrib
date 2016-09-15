package app

import (
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-app")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestAdd(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &AppActivity{metadata: md}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "ADD")
	tc.SetInput(ivType, "string")
	tc.SetInput(ivValue, "test")

	act.Eval(tc)

	value, _ := app.GetContext().GetAttrValue("myAttr")

	if value != "test" {
		fmt.Println("Bad Value: " + value.(string))
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
	act := &AppActivity{metadata: md}

	tc := test.NewTestActivityContext(md)

	app.GetContext().AddAttr("myAttr", data.STRING, "test2")

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

	md := activity.NewMetadata(jsonMetadata)
	act := &AppActivity{metadata: md}

	tc := test.NewTestActivityContext(md)

	//setup attrs
	tc.SetInput(ivAttrName, "myAttr")
	tc.SetInput(ivOp, "UPDATE")
	tc.SetInput(ivValue, "test3")

	act.Eval(tc)

	value, _ := app.GetContext().GetAttrValue("myAttr")

	if value != "test3" {
		fmt.Println("Bad Value: " + value.(string))
		t.Fail()
	}
}
