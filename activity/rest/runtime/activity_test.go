package rest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
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

	petID = res["id"].(json.Number).String()
	fmt.Println("petID:", petID)
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

	pathParams := map[string]string{
		"id": petID,
	}
	tc.SetInput("pathParams", pathParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestSimpleGetQP(t *testing.T) {

	act := activity.Get("tibco-rest")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "http://petstore.swagger.io/v2/pet/findByStatus")

	queryParams := map[string]string{
		"status": "ava",
	}
	tc.SetInput("queryParams", queryParams)

	//eval
	act.Eval(tc)

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestBuildURI(t *testing.T) {

	uri:= "http://localhost:7070/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI2(t *testing.T) {

	uri:= "https://127.0.0.1:7070/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI3(t *testing.T) {

	uri:= "http://localhost/flow/:id"

	params := map[string]string{
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}

func TestBuildURI4(t *testing.T) {

	uri:= "https://127.0.0.1/:cmd/:id/test"

	params := map[string]string{
		"cmd": "flow",
		"id": "1234",
	}

	newURI := BuildURI(uri, params)

	fmt.Println(newURI)
}