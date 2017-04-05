package coap

import (
	"context"
	"encoding/json"
	//"net/http"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/types"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"net/http"
)

const testConfig string = `{
  "id": "tibco-coap",
  "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/coap",
  "settings": {
    "port": "5683"
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

func TestInit(t *testing.T) {
	// New  factory
	f := &CoapFactory{}
	tgr := f.New("tibco-coap")

	runner := &TestRunner{}

	config := types.TriggerConfig{}
	err := json.Unmarshal([]byte(testConfig), &config)
	if err != nil{
		t.Error(err)
	}
	tgr.Init(config, runner)
}

func TestHandlerOk(t *testing.T) {

	// New  factory
	f := &CoapFactory{}
	tgr := f.New("tibco-coap")

	runner := &TestRunner{}

	config := types.TriggerConfig{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(config, runner)

	tgr.Start()
	defer tgr.Stop()

	uri := "http://127.0.0.1:5683/device/12345/reset"

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
