package coap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

const testConfig string = `{
  "name": "tibco-rest",
  "settings": {
    "port": "8091"
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "method": "POST",
        "path": "/device/:id/reset"
      }
    }
  ]
}
`

type TestStarter struct {
}

// StartFlowInstance implements flowinst.Starter.StartFlowInstance
func (ts *TestStarter) StartFlowInstance(flowURI string, startAttrs []*data.Attribute, replyHandler flowinst.ReplyHandler, execOptions *flowinst.ExecOptions) (instanceID string, startError error) {
	fmt.Printf("Started Flow with data: %v", startAttrs)
	return "dummyid", nil
}

func TestRegistered(t *testing.T) {
	tgr := trigger.Get("tibco-rest")

	if tgr == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("tibco-rest")

	starter := &TestStarter{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)

	tgr.Init(starter, config)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("tibco-rest")

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
