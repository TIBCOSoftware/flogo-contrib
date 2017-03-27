package coap

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
	"io/ioutil"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
	if err != nil{
		panic("No Json Metadata found for activity.json path")
	}
	return string(jsonMetadataBytes)
}

func TestRegistered(t *testing.T) {
	act := activity.Get("tibco-coap")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

const reqPostStr string = `{
  "name": "my pet"
}
`

var petID string

func TestSimplePost(t *testing.T) {

	act := activity.Get("tibco-coap")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "POST")
	tc.SetInput("uri", "coap://blah:5683/device")
	tc.SetInput("type", "CONFIRMABLE")
	tc.SetInput("content", reqPostStr)

	//eval
	_, err := act.Eval(tc)

	if err != nil {
		t.Error(err)
		return
	}

	val := tc.GetOutput("result")

	fmt.Printf("result: %v\n", val)

	res := val.(map[string]interface{})

	petID = res["id"].(json.Number).String()
	fmt.Println("petID:", petID)
}

func TestSimpleGet(t *testing.T) {

	act := activity.Get("tibco-coap")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "coap://blah:5683/getpet")

	//eval
	_, err := act.Eval(tc)

	if err != nil {
		t.Error(err)
		return
	}

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}

func TestSimpleGetQP(t *testing.T) {

	act := activity.Get("tibco-coap")
	tc := test.NewTestActivityContext(act.Metadata())

	//setup attrs
	tc.SetInput("method", "GET")
	tc.SetInput("uri", "coap://blah:5683/getpet")

	queryParams := map[string]string{
		"petId": "12345",
	}
	tc.SetInput("queryParams", queryParams)

	//eval
	_, err := act.Eval(tc)

	if err != nil {
		t.Error(err)
		return
	}

	val := tc.GetOutput("result")
	fmt.Printf("result: %v\n", val)
}
