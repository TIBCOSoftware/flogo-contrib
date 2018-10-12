package inference

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/framework/tf"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
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

	var done bool
	var err error

	// Unit test of Estimator model
	fmt.Println("Unit test of Estimator model")
	tc.SetInput("model", "/Users/avanderg@tibco.com/working/working_python/box_drop_demo/Archive.zip")
	tc.SetInput("inputName", "inputs")
	var estInputs = make(map[string]interface{})
	estInputs["0_0"] = 0.140586
	estInputs["1_0"] = 0.140586
	estInputs["2_0"] = 0.140586
	estInputs["amag_0"] = 0.140586
	estInputs["0_1"] = 0.140586
	estInputs["1_1"] = 0.140586
	estInputs["2_1"] = 0.140586
	estInputs["amag_1"] = 0.140586
	estInputs["0_2"] = 0.140586
	estInputs["1_2"] = 0.140586
	estInputs["2_2"] = 0.140586
	estInputs["amag_2"] = 0.140586
	estInputs["0_3"] = 0.140586
	estInputs["1_3"] = 0.140586
	estInputs["2_3"] = 0.140586
	estInputs["amag_3"] = 0.140586
	estInputs["0_4"] = 0.140586
	estInputs["1_4"] = 0.140586
	estInputs["2_4"] = 0.140586
	estInputs["amag_4"] = 0.140586
	estInputs["0_5"] = 0.140586
	estInputs["1_5"] = 0.140586
	estInputs["2_5"] = 0.140586
	estInputs["amag_5"] = 0.140586
	estInputs["0_6"] = 0.140586
	estInputs["1_6"] = 0.140586
	estInputs["2_6"] = 0.140586
	estInputs["amag_6"] = 0.140586
	estInputs["0_7"] = 0.140586
	estInputs["1_7"] = 0.140586
	estInputs["2_7"] = 0.140586
	estInputs["amag_7"] = 0.140586
	estInputs["0_8"] = 0.140586
	estInputs["1_8"] = 0.140586
	estInputs["2_8"] = 0.140586
	estInputs["amag_8"] = 0.140586
	estInputs["0_9"] = 0.140586
	estInputs["1_9"] = 0.140586
	estInputs["2_9"] = 0.140586
	estInputs["amag_9"] = 0.140586
	estInputs["0_10"] = 0.140586
	estInputs["1_10"] = 0.140586
	estInputs["2_10"] = 0.140586
	estInputs["amag_10"] = 0.140586
	estInputs["word_label"] = 0

	var features []map[string]interface{}
	features = append(features, make(map[string]interface{}))
	features[0]["name"] = "inputs"
	features[0]["data"] = estInputs
	// fmt.Println(features)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", features)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	// ///???????NEED TO MAKE SURE DIFFERENT MODELS ARE LOADED////////
	// // Unit test of Pass inputs to Outputs model
	// fmt.Println("Unit test of Pass inputs to Outputs model")
	// tc.SetInput("model", "/Users/avanderg@tibco.com/working/working_python/simplest_model_just_passes_inputs/model/simple_pass/")
	// var features2 []map[string]interface{}
	// features2 = append(features2, make(map[string]interface{}))
	// features2[0]["name"] = "X"
	// features2[0]["data"] = []float64{0.23, 4.5, 234.234}
	// fmt.Println(features2)

	// tc.SetInput("inputName", "inputs")
	// tc.SetInput("framework", "Tensorflow")
	// tc.SetInput("sigDefName", "serving_default")
	// tc.SetInput("tag", "serve")
	// tc.SetInput("features", features2)

	// done, err = act.Eval(tc)
	// if done == false {
	// 	assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	// } else {
	// 	assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	// }

	// // Unit test of Pass inputs to Outputs model
	// fmt.Println("Unit test of Pass inputs to Outputs model")
	// tc.SetInput("model", "/Users/avanderg@tibco.com/sample_tf_models/simpleCNN/")
	// var features3 []map[string]interface{}
	// features3 = append(features3, make(map[string]interface{}))
	// features3[0]["name"] = "X"
	// features3[0]["data"] = [][][][]float32{{{{0.0000000856947568}}, {{0.00000331318370}}, {{0.0000858655563}}, {{0.00149167657}}, {{0.0173705094}}, {{0.135591557}}, {{0.709471493}}, {{2.48839579}}, {{5.85040827}}, {{9.22008867}}}}
	// fmt.Println(features3)

	// tc.SetInput("inputName", "inputs")
	// tc.SetInput("framework", "Tensorflow")
	// tc.SetInput("sigDefName", "serving_default")
	// tc.SetInput("tag", "serve")
	// tc.SetInput("features", features3)

	// done, err = act.Eval(tc)
	// if done == false {
	// 	assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	// } else {
	// 	assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	// }

	// // the below model,multilineinference_test_model_forflogo_contrib, doesn't do anything
	// //     other than only have two inputs to make changing number of inputs easy
	// tc.SetInput("model", "/Users/avanderg@tibco.com/working/working_python/simple_cnn/model/SimpleCNN")
	// tc.SetInput("model", "/Users/avanderg@tibco.com/working/working_python/multilineinference_test_model_forflogo_contrib/Archive.zip")

	// var features = make(map[string]interface{})
	// features["one"] = []float64{0.23}
	// features["two"] = []float64{2.1}
	// features["label"] = []int{0}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
