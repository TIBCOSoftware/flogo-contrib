package coap

import (
	"context"
	//"encoding/json"
	//"net/http"
	//"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	//"net/http"
	//"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
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

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action: %v", act.Config().Id)
	return nil, nil
}
/*
// TODO Fix this test

func TestHandlerOk(t *testing.T) {
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)

	// New  factory
	f := &CoapFactory{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

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
*/
