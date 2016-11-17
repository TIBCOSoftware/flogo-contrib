package device

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-device")

// DeviceTrigger is simple MQTT trigger
type DeviceTrigger struct {
	metadata  *trigger.Metadata
	runner    action.Runner
	client    mqtt.Client
	settings  map[string]string
	config    *trigger.Config
	endpoints map[string]*trigger.EndpointConfig

	pubTopic  string
	subTopic  string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&DeviceTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *DeviceTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *DeviceTrigger) Init(config *trigger.Config, runner action.Runner) {

	t.config = config
	t.settings = config.Settings
	t.runner = runner
	t.pubTopic = "flogo/" + t.settings["device:name"] + "/in"
	t.subTopic = "flogo/" + t.settings["device:name"] + "/out"
}

// Start implements ext.Trigger.Start
func (t *DeviceTrigger) Start() error {

	opts := mqtt.NewClientOptions()

	broker := "tcp://" + t.settings["mqtt_server"] + ":1883"
	opts.AddBroker(broker)
	clientId := strconv.FormatInt(time.Now().Unix(), 10)
	opts.SetClientID("flogo-" + t.settings["device:name"] + "-" + clientId)
	opts.SetUsername(t.settings["mqtt_user"])
	opts.SetPassword(t.settings["mqtt_password"])

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {

		payload := string(msg.Payload())
		log.Debug("Received msg:", payload)

		req := &DeviceRequest{}
		err := json.Unmarshal(msg.Payload(), req)
		if err != nil {
			log.Errorf("Enable to parse request: %s", err.Error())
			return
		}

		if req.Status == "READY" {
			log.Info("Device %s - READY", t.settings["device:name"])
			return
		}

		t.RunAction(req)
	})

	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	t.endpoints = make(map[string]*trigger.EndpointConfig)

	//todo support multiple endpoints per device
	for _, endpoint := range t.config.Endpoints {

		t.endpoints["0"] = endpoint
	}

	if token := t.client.Subscribe(t.subTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Errorf("Error subscribing to topic %s: %s", t.subTopic, token.Error())
		panic(token.Error())
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *DeviceTrigger) Stop() error {

	log.Debug("Unsubcribing from topic: ", t.subTopic)

	if token := t.client.Unsubscribe(t.subTopic); token.Wait() && token.Error() != nil {
		log.Errorf("Error unsubscribing from topic %s: %s", t.subTopic, token.Error())
	}

	t.client.Disconnect(250)

	return nil
}

// RunAction starts a new Process Instance
func (t *DeviceTrigger) RunAction(req *DeviceRequest) {

	epConfig := t.endpoints[req.Endpoint]

	val, _ := strconv.Atoi(req.Value);

	data := map[string]interface{}{
		"value": val,
	}

	//todo handle error
	startAttrs, _ := t.metadata.OutputsToAttrs(data, false)

	action := action.Get(epConfig.ActionType)
	context := trigger.NewContext(context.Background(), startAttrs)

	//todo handle error
	_, _, err := t.runner.Run(context, action, epConfig.ActionURI, nil)
	if err != nil {
		log.Error(err)
	}

	log.Debugf("Ran action: [%s-%s]", epConfig.ActionType, epConfig.ActionURI)

	//todo convert reply data to response
	//if replyData != nil {
	//	data, err := json.Marshal(replyData)
	//	if err != nil {
	//		log.Error(err)
	//	} else {
	//		t.publishMessage(req.ReplyTo, string(data))
	//	}
	//}
}


func (t *DeviceTrigger) publishMessage(ep string, value string) {

	log.Debug("ReplyTo endpoint: ", ep)
	log.Debug("Publishing message: ", ep + value)

	token := t.client.Publish(t.pubTopic, 0, false, ep+value)
	token.Wait()
}

type DeviceRequest struct {
	Status   string `json:"status,omitempty"`
	Endpoint string `json:"ep,omitempty"`
	Value    string `json:"value,omitempty"`
}