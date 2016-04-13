package mqtt

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/processinst"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const testConfig string = `{
  "name": "mqtt",
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
      "processURI": "local://testProcess",
      "settings": {
        "topic": "test_start"
      }
    }
  ]
}`

type TestStarter struct {
}

// StartProcessInstance implements processinst.Starter.StartProcessInstance
func (ts *TestStarter) StartProcessInstance(processURI string, startData map[string]string, replyHandler processinst.ReplyHandler, execOptions *processinst.ExecOptions) string {
	fmt.Printf("Started Process with data: %v", startData)
	return "dummyid"
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("mqtt")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	tgr := trigger.Get("mqtt")

	starter := &TestStarter{}

	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	fmt.Println(config)
	tgr.Init(starter, config)
}

func TestEndpoint(t *testing.T) {

	tgr := trigger.Get("mqtt")

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
	fmt.Println("---- doing publish ----")
	token := client.Publish("flogo/test_start", 0, false, "Test message payload")
	token.Wait()

	client.Disconnect(250)
	fmt.Println("Sample Publisher Disconnected")

	//req, err := http.NewRequest("POST", uri, nil)
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//log.Debug("response Status:", resp.Status)
	//
	//if resp.StatusCode >= 300 {
	//	t.Fail()
	//}
}
