package rest

import (
	"testing"

	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/activity"
	"github.com/TIBCOSoftware/flogo-lib/test"
)

const reqPostStr string = `{
  "name": "my pet"
}
`

var petID string

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-rest")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestSimplePost(t *testing.T) {

	act := activity.Get("tibco-rest")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "POST")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet")
	tc.SetInput("content", reqPostStr)

	//eval
	act.Eval(tc)
	val := tc.GetOutput("result")

	fmt.Printf("result: %v\n", val)

	res := val.(map[string]interface{})

	petIDInt := int64(res["id"].(float64))
	petID = strconv.FormatInt(petIDInt, 10)
}

func TestSimpleGet(t *testing.T) {

	act := activity.Get("tibco-rest")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/"+petID)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestParamGet(t *testing.T) {

	act := activity.Get("tibco-rest")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/:id")

	params := map[string]string{
		"id": petID,
	}
	tc.SetInput("params", params)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}
