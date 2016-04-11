package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
)

const testConfig string = `{
  "name": "rest",
  "settings": {
    "port": "8091"
  },
  "endpoints": [
    {
      "processURI": "local://testProcess",
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

// StartProcessInstance implements processinst.Starter.StartProcessInstance
func (ts *TestStarter) StartProcessInstance(processURI string, startData map[string]string, replyHandler processinst.ReplyHandler, execOptions *processinst.ExecOptions) string {
	fmt.Printf("Started Process with data: %v", startData)
	return "dummyid"
}

func TestRegistered(t *testing.T) {
	tgr := trigger.Get("rest")

	if tgr == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("rest")

	starter := &TestStarter{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)

	tgr.Init(starter, config)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("rest")

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
