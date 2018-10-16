package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/TIBCOSoftware/flogo-contrib/activity/inference/utils"
)

// Framework interface used to implement specific ml framework implementations
type Framework interface {
	Load(model *Model, flags ModelFlags) (err error)
	Run(model *Model) (out map[string]interface{}, err error)
	FrameworkTyp() string
}

// Model represents a ML model. The metadata defines inputs, outputs, shapes, etc.
// The Instance object is a pointer to the actual framework obj. This will be used by the framework impl to execute models
type Model struct {
	Metadata *Metadata
	Instance interface{}
	Inputs   map[string]interface{}
}

// ModelFlags Contains flags to add to metadata and to aid in loading
type ModelFlags struct {
	Tag       string
	SigDef    string
	ModelFile string
	ModelPath string
}

// Load is unzipping the file and then passing it to be read
func Load(modelArchive string, framework Framework, flags ModelFlags) (*Model, error) {
	f, _ := os.Open(modelArchive)
	defer f.Close()
	var outDir string
	if fi, err := f.Stat(); err == nil && fi.IsDir() {
		outDir = modelArchive
	} else if err == nil && !fi.IsDir() {
		tmpDir := os.TempDir()
		outDir = filepath.Join(tmpDir, utils.PseudoUuid())
		if err := utils.Unzip(modelArchive, outDir); err != nil {
			return nil, fmt.Errorf("Failed to extract model archive: %v", err)
		}
	} else {
		return nil, fmt.Errorf("%s does not exist", modelArchive)
	}

	modelFilename := filepath.Join(outDir, "saved_model.pb")
	if _, err := os.Stat(modelFilename); err != nil {
		// This if is here for when we can read pbtxt files
		if _, err2 := os.Stat(modelFilename + "txt"); err2 == nil {
			modelFilename = modelFilename + "txt"
			return nil, errors.New("Currently loading saved_model.pbtxt is not supported")
			//comment the return when we can read pbtxt
		} else {
			return nil, errors.New("saved_model.pb does not exist")
		}
	}

	flags.ModelPath = outDir
	flags.ModelFile = modelFilename
	var model Model
	err := framework.Load(&model, flags)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// This is un-used (un-needed?) and now m.Inputs is of the form of map[string]interface() with both maps and arrays so this doesn't work so well
// func (m *Model) AppendInput(inputName string, feature string, in interface{}) {
// 	if m.Inputs == nil {
// 		m.Inputs = make(map[string]map[string]interface{})
// 	}
//
// 	m.Inputs[inputName][feature] = in
// }

func (m *Model) SetInputs(in map[string]interface{}) {
	m.Inputs = in
}

func (m *Model) RemoveInput(featureName string) {

}

func (m *Model) Run(framework Framework) (map[string]interface{}, error) {
	// check if inputs are available. maybe validate inputs against metadata?

	// run the model
	out, err := framework.Run(m)
	if err != nil {
		return nil, err
	}

	return out, nil
}
