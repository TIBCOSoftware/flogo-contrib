package rest

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/test"
	"fmt"
)

const reqPostStr string = `{
  "name": "my pet"
}
`

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-rest")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestSimpleGet(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := activity.Get("tibco-rest")
	tc := test.NewTestActivityContext()

	//setup attrs
	//tc.SetOrAddAttrValue("method","GET")
	//tc.SetOrAddAttrValue("uri","http://petstore.swagger.io/v2/pet/1234")

	//eval
	act.Eval(tc)
	val,_ := tc.GetAttrValue("result")

	fmt.Println("result:",val)

	//check result attr
}
