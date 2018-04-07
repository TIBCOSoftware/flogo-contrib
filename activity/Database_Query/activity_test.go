package Database_Query

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

//debugging

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attr
	// tc.SetInput("driverName", "mysql")
	// tc.SetInput("datasourceName", "flogo:password@tcp(localhost:3306)/testdb")
	// tc.SetInput("query", "select * from person")

	//setup attr for postgres
	tc.SetInput("driverName", "postgres")
	tc.SetInput("datasourceName", "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable")
	tc.SetInput("query", "select * from company")
	act.Eval(tc)

	//check result attr
	result := tc.GetOutput("result")
	fmt.Println("result: ", result)

	if result == nil {
		t.Fail()
	}

}
