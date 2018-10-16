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

	// Unit test of Estimator Classifier model
	fmt.Println("Unit test of Estimator Classifier model")
	tc.SetInput("model", "Archive_estDNNClf.zip")
	tc.SetInput("inputName", "inputs")
	var estInputsA = make(map[string]interface{})
	estInputsA["one"] = 0.140586
	estInputsA["two"] = 0.140586
	estInputsA["three"] = 0.140586
	estInputsA["label"] = 0

	var featuresA []map[string]interface{}
	featuresA = append(featuresA, make(map[string]interface{}))
	featuresA[0]["name"] = "inputs"
	featuresA[0]["data"] = estInputsA
	// fmt.Println(featuresA)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", featuresA)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	// Unit test of Estimator DNN Regressor model
	fmt.Println("Unit test of Estimator Regressor model")
	tc.SetInput("model", "Archive_estDNNrgr.zip")
	tc.SetInput("inputName", "inputs")
	var estInputsB = make(map[string]interface{})
	estInputsB["one"] = 0.140586
	estInputsB["two"] = 0.140586
	estInputsB["three"] = 0.140586
	estInputsB["label"] = 0.

	var featuresB []map[string]interface{}
	featuresB = append(featuresB, make(map[string]interface{}))
	featuresB[0]["name"] = "inputs"
	featuresB[0]["data"] = estInputsB
	// fmt.Println(featuresB)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", featuresB)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	// Unit test of Estimator Linear Regressor model
	fmt.Println("Unit test of Linear Regressor Estimator model")
	tc.SetInput("model", "Archive_LinReg.zip")
	tc.SetInput("inputName", "inputs")
	var estInputsC = make(map[string]interface{})
	estInputsC["one"] = 0.140586
	estInputsC["two"] = 0.140586
	estInputsC["three"] = 0.140586
	estInputsC["label"] = 0.

	var featuresC []map[string]interface{}
	featuresC = append(featuresC, make(map[string]interface{}))
	featuresC[0]["name"] = "inputs"
	featuresC[0]["data"] = estInputsC
	// fmt.Println(featuresC)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", featuresC)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	// Unit test of Pairwaise Multiplication model
	fmt.Println("Unit test of Pairwaise Multiplication model")
	tc.SetInput("model", "Archive_pairwise_multi.zip")
	var features2 []map[string]interface{}
	features2 = append(features2, make(map[string]interface{}))
	features2[0]["name"] = "X1"
	features2[0]["data"] = [][]float32{{0.23, 4.5, -3.1}, {7.1, 3.14159, -0.00123}}
	features2 = append(features2, make(map[string]interface{}))
	features2[1]["name"] = "X2"
	features2[1]["data"] = [][]float32{{4.34782608, 0.2222222222, -0.3225806451612903},
		{0.14084507042253522, 0.31831015504887655, -813.0081300813008}}
	fmt.Println(features2)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", features2)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	// Unit test ofSimple CNN model
	fmt.Println("Unit test of simple CNN model")
	tc.SetInput("model", "Archive_simpleCNN.zip")
	var features3 []map[string]interface{}
	features3 = append(features3, make(map[string]interface{}))
	features3[0]["name"] = "X"
	features3[0]["data"] = [][][][]float32{
		{{{0.0000000856947568}}, {{0.00000331318370}}, {{0.0000858655563}}, {{0.00149167657}}, {{0.0173705094}}, {{0.135591557}}, {{0.709471493}}, {{2.48839579}}, {{5.85040827}}, {{9.22008867}}},
		{{{9.22008867}}, {{5.85040827}}, {{2.48839579}}, {{00.709471493}}, {{0.135591557}}, {{0.00149167657}}, {{0.0000858655563}}, {{0.00000331318370}}, {{0.0000000856947568}}, {{0.}}},
		{{{0.0173705094}}, {{0.135591557}}, {{0.709471493}}, {{2.48839579}}, {{5.85040827}}, {{9.22008867}}, {{5.85040827}}, {{2.48839579}}, {{0.709471493}}, {{0.135591557}}},
	}
	fmt.Println(features3)

	tc.SetInput("inputName", "inputs")
	tc.SetInput("framework", "Tensorflow")
	tc.SetInput("sigDefName", "serving_default")
	tc.SetInput("tag", "serve")
	tc.SetInput("features", features3)

	done, err = act.Eval(tc)
	if done == false {
		assert.Fail(t, fmt.Sprintf("Error raised: %s", err))
	} else {
		assert.True(t, done, fmt.Sprintf("Evaluation came back: %t", done))
	}

	//check result attr
	fmt.Println(tc.GetOutput("result"))
}
