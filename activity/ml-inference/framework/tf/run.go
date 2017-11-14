package tf

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	models "github.com/TIBCOSoftware/flogo-contrib/activity/ml-inference/model"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

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
		examplePb, err := createInputExampleTensor(inputName, model)
		if err != nil {
			fmt.Println("err")
		}
		inputs[inputMap.Output(0)] = examplePb
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

func getTensorValue(tensor *tf.Tensor) interface{} {
	switch tensor.Value().(type) {
	case [][]string:
		return tensor.Value().([][]string)
	case [][]float32:
		return tensor.Value().([][]float32)
	}

	return nil
}

func createInputExampleTensor(inputName string, model *models.Model) (*tf.Tensor, error) {
	pb, _ := Example(model.Inputs[inputName])
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
