package lambda

import (
	"context"
	//"encoding/json"
	//"net/http"
	//"testing"
	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	//"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "id": "tibco-rest",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/lambda",
  "settings": {
  },
  "handlers": [
    {
      "actionId": "my_test_flow",
      "settings": {
      }
    }
  ]
}
`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action: %v", act.Config().Id)
	return nil, nil
}