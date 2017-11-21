package mqtt

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	//MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"io/ioutil"
	//"time"
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
  "name": "tibco-mqtt",
  "settings": {
    "topic": "flogo/#",
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "handlers": [
    {
      "actionId": "device_info",
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

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	log.Debugf("Ran Action: %v", act.Config().Id)
	return nil, nil
}

func TestInit(t *testing.T) {

	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)
}

/*
// TODO Fix this test
func TestEndpoint(t *testing.T) {

	// New  factory
	md := trigger.NewMetadata(jsonMetadata)
	f := NewFactory(md)

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID("flogo_test")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	log.Debug("---- doing first publish ----")

	token := client.Publish("test_start", 0, false, "Test message payload!")
	token.Wait()

	duration2 := time.Duration(2)*time.Second
	time.Sleep(duration2)

	log.Debug("---- doing second publish ----")

	token = client.Publish("test_start", 0, false, "Test message payload!")
	token.Wait()

	duration5 := time.Duration(5)*time.Second
	time.Sleep(duration5)

	client.Disconnect(250)
	log.Debug("Sample Publisher Disconnected")
}
*/
