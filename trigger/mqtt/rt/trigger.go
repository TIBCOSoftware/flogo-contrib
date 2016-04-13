package mqtt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/core/flow"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-mqtt")

// todo: switch to use endpoint registration

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata    *trigger.Metadata
	flowStarter flowinst.Starter
	client      mqtt.Client
	settings    map[string]string
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MqttTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *MqttTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *MqttTrigger) Init(flowStarter flowinst.Starter, config *trigger.Config) {

	t.flowStarter = flowStarter
	t.settings = config.Settings
}

// Start implements ext.Trigger.Start
func (t *MqttTrigger) Start() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.settings["broker"])
	opts.SetClientID(t.settings["id"])
	opts.SetUsername(t.settings["user"])
	opts.SetPassword(t.settings["password"])
	b, err := strconv.ParseBool(t.settings["cleansess"])
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return
	}
	opts.SetCleanSession(b)
	if t.settings["store"] != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.settings["store"]))
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())

		// Match suffix of topic
		if strings.HasSuffix(topic, "start") {
			t.StartFlow(payload)
		}
	})

	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	i, err := strconv.Atoi(t.settings["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}

	if token := t.client.Subscribe(t.settings["topic"], byte(i), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() {
	//unsubscribe from topic
	if token := t.client.Unsubscribe(t.settings["topic"]); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		log.Debug("Unsubcribing from topic: ", t.settings["topic"])
	}

	t.client.Disconnect(250)
}

// StartFlow starts a new Flow Instance
func (t *MqttTrigger) StartFlow(payload string) {

	req := &StartRequest{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Error Starting flow ", err.Error())
		return
	}

	log.Debug("Flow URI ", req.FlowURI)
	log.Debug("flowStarter.StartFlow ", t.flowStarter)
	id := t.flowStarter.StartFlowInstance(req.FlowURI, req.Data, nil, nil)
	log.Debug("Start flow id: ", id)
	t.publishMessage(req.ReplyTo, id)
}

func (t *MqttTrigger) publishMessage(topic string, message string) {

	log.Debug("ReplyTo topic: ", topic)
	log.Debug("Publishing message: ", message)

	qos, err := strconv.Atoi(t.settings["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	token.Wait()
}

// StartRequest describes a request for starting a FlowInstance
type StartRequest struct {
	FlowURI     string            `json:"flowUri"`
	Data        map[string]string `json:"data"`
	Interceptor *flow.Interceptor `json:"interceptor"`
	Patch       *flow.Patch       `json:"patch"`
	ReplyTo     string            `json:"replyTo"`
}
