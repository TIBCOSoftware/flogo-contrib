package databasequery

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"testing"
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

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	fmt.Println("===============================")
	fmt.Println("Unit Test ===> MySQL Connection")
	fmt.Println("===============================")
	fmt.Println("")
	tc.SetInput("driverName", "mysql")
	tc.SetInput("datasourceName", "username:password@tcp(hostserver:port)/dbName")
	//tc.SetInput("query", "insert into user_details (mobile) values (8) ")
	//tc.SetInput("query", "delete from user_details where mobile=8 ")
	//tc.SetInput("query", "select * from user_details where mobile=8 ")
	act.Eval(tc)

	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}
