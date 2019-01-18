package tf

import (
	"fmt"
	"reflect"
	"strings"

	models "github.com/TIBCOSoftware/flogo-contrib/activity/inference/model"
	"github.com/golang/protobuf/proto"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-inference")

// Run is used to execute a Tensorflow model with the model input data
func (i *TensorflowModel) Run(model *models.Model) (out map[string]interface{}, err error) {
	// Grab native tf SavedModel
	savedModel := model.Instance.(*tf.SavedModel)

	var inputOps = make(map[string]*tf.Operation)
	var outputOps []tf.Output

	// Validate that the operations exsist and create operation
	for k, v := range model.Metadata.Inputs.Params {
		if validateOperation(v.Name, savedModel) == false {
			return nil, fmt.Errorf("Invalid operation %s", v.Name)
		}
		inputOps[k] = savedModel.Graph.Operation(v.Name)
	}

	// Create output operations
	var outputOrder []string
	for k, o := range model.Metadata.Outputs {
		outputOps = append(outputOps, savedModel.Graph.Operation(o.Name).Output(0))
		outputOrder = append(outputOrder, k)
	}

	// create input tensors and add to map
	inputs := make(map[tf.Output]*tf.Tensor)
	for inputName, inputMap := range inputOps {
		v := reflect.ValueOf(model.Inputs[inputName])
		switch v.Kind() {
		case reflect.Map:
			// Need to check names against pb structure, right now just assume it
			examplePb, err := createInputExampleTensor(model.Inputs[inputName])
			if err != nil {
				return nil, err
			}
			inputs[inputMap.Output(0)] = examplePb

		case reflect.Slice, reflect.Array:
			shape := model.Metadata.Inputs.Features[inputName].Shape
			typ := model.Metadata.Inputs.Features[inputName].Type
			data, err := checkDataTypes(model.Inputs[inputName], shape, typ, inputName)
			if err != nil {
				return nil, err
			}

			inputs[inputMap.Output(0)], err = tf.NewTensor(data)
			if err != nil {
				return nil, err
			}

		case reflect.Ptr:
			if val, ok := model.Inputs[inputName].(*tf.Tensor); ok {
				inputs[inputMap.Output(0)] = val
			} else {
				if val2, ok2 := model.Inputs[inputName].(*[]byte); ok2 {
					inputs[inputMap.Output(0)], err = tf.NewTensor(val2)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("Interface not casting to Tensor or byte object. Is your pointer a tensor?")
				}

			}

		default:
			log.Info("Type not a Slice, Array, Map, or Pointer/Tensor, but still trying to make a tf.Tensor.")
			inputs[inputMap.Output(0)], err = tf.NewTensor(model.Inputs[inputName])
			if err != nil {
				return nil, err
			}
		}
	}

	results, err := savedModel.Session.Run(inputs, outputOps, nil)
	if err != nil {
		return nil, err
	}

	// Iterate over the expected outputs, find the actual and map into map
	out = make(map[string]interface{})
	for k := range model.Metadata.Outputs {
		for i := 0; i < len(outputOrder); i++ {
			if outputOrder[i] == k {
				out[k] = getTensorValue(results[i])
			}
		}
	}

	return out, nil

}

func checkDataTypes(data interface{}, shape []int64, typ string, inputName string) (outdata interface{}, err error) {
	t := fmt.Sprintf("%T", data)
	outdata = data
	switch typ {
	case "DT_FLOAT":
		// if strings.Contains(t, "float64") {
		// 	outdata, err = float64TensorTofloat32Tensor(data, nil)  //location of coerce functions to be deteremined
		// 	if err != nil {
		// 		return nil, fmt.Errorf("Data conversion for %s had error: %s", inputName, err)
		// 	}
		// 	fmt.Println("Coerceing FLoat to Double")
		// } else
		if !strings.Contains(t, "float32") {
			return nil, fmt.Errorf("Data for %s not of the right type. should be tensor of %s (TF type) but is array of %s (go type)", inputName, typ, t)
		}
	case "DT_DOUBLE":
		if !strings.Contains(t, "float64") {
			return nil, fmt.Errorf("Data for %s not of the right type. should be tensor of %s (TF type) but is array of %s (go type)", inputName, typ, t)
		}
	case "DT_INT32":
		if !strings.Contains(t, "int32") {
			return nil, fmt.Errorf("Data for %s not of the right type. should be tensor of %s (TF type) but is array of %s (go type)", inputName, typ, t)
		}
	case "DT_INT64":
		if !strings.Contains(t, "int64") {
			return nil, fmt.Errorf("Data for %s not of the right type. should be tensor of %s (TF type) but is array of %s (go type)", inputName, typ, t)
		}
	}
	return outdata, nil
}

func getTensorValue(tensor *tf.Tensor) interface{} {
	switch tensor.Value().(type) {
	case [][]string:
		return tensor.Value().([][]string)
	case []string:
		return tensor.Value().([]string)
	case []float32:
		return tensor.Value().([]float32)
	case [][]float32:
		return tensor.Value().([][]float32)
	case []float64:
		return tensor.Value().([]float64)
	case [][]float64:
		return tensor.Value().([][]float64)
	case []int64:
		return tensor.Value().([]int64)
	case [][]int64:
		return tensor.Value().([][]int64)
	case []int32:
		return tensor.Value().([]int32)
	case [][]int32:
		return tensor.Value().([][]int32)
	case []byte:
		return tensor.Value().([]byte)
	case [][]byte:
		return tensor.Value().([][]byte)
	case []int:
		return tensor.Value().([]int)
	}
	return nil
}

func createInputExampleTensor(featMap interface{}) (*tf.Tensor, error) {
	pb, err := Example(featMap.(map[string]interface{}))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Example: %s", err)
	}

	byteList, err := proto.Marshal(pb)
	if err != nil {
		return nil, fmt.Errorf("marshaling error: %s", err)
	}

	newTensor, err := tf.NewTensor([]string{string(byteList)})
	if err != nil {
		return nil, err
	}

	return newTensor, nil
}

func validateOperation(op string, savedModel *tf.SavedModel) bool {

	tfOp := savedModel.Graph.Operation(op)
	if tfOp == nil {
		return false
	}
	return true
}
