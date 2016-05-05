package timer

import (
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"time"
)

const testConfig3 string = `{
  "name": "timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "false"
      }
    }
  ]
}`

const testConfig string = `{
  "name": "timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "false",
        "startDate" : "2016-05-03T19:25:00Z-04:00"
      }
    }
  ]
}`

const testConfig2 string = `{
  "name": "timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
      	"notImmediate": "false",
        "repeating": "true",
        "seconds": "5"
      }
    }
  ]
}`

const testConfig4 string = `{
  "name": "timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "repeating": "false",
        "startDate" : "05/01/2016, 12:25:01"
      }
    },
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "hours": "24"
      }
    },
    {
      "flowURI": "local://testFlow3",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "minutes": "60"
      }
    },
    {
      "flowURI": "local://testFlow3",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "seconds": "30"
      }
    }
  ]
}`

type TestStarter struct {
}

// StartFlowInstance implements flowinst.Starter.StartFlowInstance
func (ts *TestStarter) StartFlowInstance(flowURI string, startAttrs []*data.Attribute, replyHandler flowinst.ReplyHandler, execOptions *flowinst.ExecOptions) (instanceID string, startError error) {
	log.Debugf("Started Flow with data: %v", startAttrs)
	return "dummyid", nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("timer")

	if act == nil {
		t.Error("Timer Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("timer")

	starter := &TestStarter{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(starter, config)
}

func TestTimer(t *testing.T) {

	log.Debugf("TestTimer")
	tgr := trigger.Get("timer")

	tgr.Start()
	time.Sleep(time.Second * 2000)
	defer tgr.Stop()


	log.Debug("Test timer done")
}
