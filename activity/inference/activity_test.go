package inference

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework/tf"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var _ tf.TensorflowModel

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

	//We need to get a small model here so we can actually test this without having to customize this everytime
	//   i.e. a real unit test
	//setup attrs
	tc.SetInput("model", "/Users/avanders/working/working_python/box_drop_demo/Archive.zip")
	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")

	var features = make(map[string]interface{})
	features["0_0"] = 0.140586
	features["1_0"] = 0.140586
	features["2_0"] = 0.140586
	features["amag_0"] = 0.140586
	features["0_1"] = 0.140586
	features["1_1"] = 0.140586
	features["2_1"] = 0.140586
	features["amag_1"] = 0.140586
	features["0_2"] = 0.140586
	features["1_2"] = 0.140586
	features["2_2"] = 0.140586
	features["amag_2"] = 0.140586
	features["0_3"] = 0.140586
	features["1_3"] = 0.140586
	features["2_3"] = 0.140586
	features["amag_3"] = 0.140586
	features["0_4"] = 0.140586
	features["1_4"] = 0.140586
	features["2_4"] = 0.140586
	features["amag_4"] = 0.140586
	features["0_5"] = 0.140586
	features["1_5"] = 0.140586
	features["2_5"] = 0.140586
	features["amag_5"] = 0.140586
	features["0_6"] = 0.140586
	features["1_6"] = 0.140586
	features["2_6"] = 0.140586
	features["amag_6"] = 0.140586
	features["0_7"] = 0.140586
	features["1_7"] = 0.140586
	features["2_7"] = 0.140586
	features["amag_7"] = 0.140586
	features["0_8"] = 0.140586
	features["1_8"] = 0.140586
	features["2_8"] = 0.140586
	features["amag_8"] = 0.140586
	features["0_9"] = 0.140586
	features["1_9"] = 0.140586
	features["2_9"] = 0.140586
	features["amag_9"] = 0.140586
	features["0_10"] = 0.140586
	features["1_10"] = 0.140586
	features["2_10"] = 0.140586
	features["amag_10"] = 0.140586
	features["word_label"] = 0

	tc.SetInput("features", features)

	done, err := act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
