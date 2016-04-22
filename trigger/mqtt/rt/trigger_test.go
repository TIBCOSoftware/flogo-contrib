package mqtt

import (
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const testConfig string = `{
  "name": "tibco-mqtt",
  "settings": {
    "topic": "flogo/#",
    "broker": "tcp://192.168.1.12:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "topic": "test_start"
      }
    }
  ]
}`

type TestStarter struct {
}

// StartFlowInstance implements flowinst.Starter.StartFlowInstance
func (ts *TestStarter) StartFlowInstance(flowURI string, startData map[string]interface{}, replyHandler flowinst.ReplyHandler, execOptions *flowinst.ExecOptions) (instanceID string, startError error) {
	log.Debugf("Started Flow with data: %v", startData)
	return "dummyid", nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("tibco-mqtt")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("tibco-mqtt")

	starter := &TestStarter{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(starter, config)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("tibco-mqtt")

	tgr.Start()
	defer tgr.Stop()

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://192.168.1.12:1883")
	opts.SetClientID("flogoEngine")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Debug("---- doing publish ----")
	token := client.Publish("flogo/test_start", 0, false, "Test message payload!")
	token.Wait()

	client.Disconnect(250)
	log.Debug("Sample Publisher Disconnected")
}
