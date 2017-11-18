package rest

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
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
  "settings": {
    "port": "8091"
  },
  "handlers": [
    {
      "actionId": "my_test_flow",
      "settings": {
        "method": "POST",
        "path": "/device/:id/reset"
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

/*
//TODO fix this test
func TestInitOk(t *testing.T) {
	// New  factory
	f := &RestFactory{}
	tgr := f.New("tibco-rest")

	runner := &TestRunner{}

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(config, runner)
}
*/
/*
//TODO fix this test
func TestHandlerOk(t *testing.T) {

	// New  factory
	f := &RestFactory{}
	tgr := f.New("tibco-rest")

	runner := &TestRunner{}

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	uri := "http://127.0.0.1:8091/device/12345/reset"

	req, err := http.NewRequest("POST", uri, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Debug("response Status:", resp.Status)

	if resp.StatusCode >= 300 {
		t.Fail()
	}
}
*/
