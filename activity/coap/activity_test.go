package coap

import (
	//"encoding/json"
	//"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	//"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"io/ioutil"
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

const reqPostStr string = `{
  "name": "my pet"
}
`

//var petID string
/*
//TODO fix this test
func TestSimplePost(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

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
*/

/*
//TODO fix this test
func TestSimpleGet(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

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
*/

/*
//TODO fix this test
func TestSimpleGetQP(t *testing.T) {

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

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
*/
