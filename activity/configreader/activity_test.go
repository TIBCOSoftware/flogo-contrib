package configreader

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

func TestStringConfig(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("configFile", "config.json")
	tc.SetInput("configName", "string_config")
	tc.SetInput("configType", "string")
	tc.SetInput("readEachTime", "true")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("configValue")
	fmt.Printf("String result : %v\n", val)

}

func TestIntConfig(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("configFile", "config.json")
	tc.SetInput("configName", "int_config")
	tc.SetInput("configType", "int")
	tc.SetInput("readEachTime", "true")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("configValue")
	fmt.Printf("Int result : %v\n", val)

}

func TestFloatConfig(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("configFile", "config.json")
	tc.SetInput("configName", "float_config")
	tc.SetInput("configType", "float")
	tc.SetInput("readEachTime", "true")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("configValue")
	fmt.Printf("Float result : %v\n", val)

}

func TestBoolConfig(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("configFile", "config.json")
	tc.SetInput("configName", "bool_config")
	tc.SetInput("configType", "bool")
	tc.SetInput("readEachTime", "true")

	act.Eval(tc)

	//check result attr
	val := tc.GetOutput("configValue")
	fmt.Printf("Bool result : %v\n", val)

}
