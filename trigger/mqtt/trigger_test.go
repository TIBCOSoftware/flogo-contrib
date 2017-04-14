package mqtt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/types"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string{
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil{
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

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

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}


func TestInit(t *testing.T) {
	// New  factory
	f := &MQTTFactory{}
	tgr := f.New("tibco-mqtt")

	runner := &TestRunner{}

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr.Init(config, runner)
}

func TestEndpoint(t *testing.T) {
	// New  factory
	f := &MQTTFactory{}
	tgr := f.New("tibco-mqtt")

	runner := &TestRunner{}

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr.Init(config, runner)

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
