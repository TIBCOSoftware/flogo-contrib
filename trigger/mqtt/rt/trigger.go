package mqtt

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/engine/ext/trigger"
	"github.com/TIBCOSoftware/flogo-lib/engine/starter"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/op/go-logging"
)

// log is the default package logger
var log = logging.MustGetLogger("mqtt-trigger")

// MqttTrigger is simple MQTT trigger
type MqttTrigger struct {
	metadata       *trigger.Metadata
	processStarter starter.ProcessStarter
	client         mqtt.Client
	config         map[string]string
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
func (t *MqttTrigger) Init(processStarter starter.ProcessStarter, config map[string]string) {

	t.processStarter = processStarter
	t.config = config
}

// Start implements ext.Trigger.Start
func (t *MqttTrigger) Start() {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(t.config["broker"])
	opts.SetClientID(t.config["id"])
	opts.SetUsername(t.config["user"])
	opts.SetPassword(t.config["password"])
	b, err := strconv.ParseBool(t.config["cleansess"])
	if err != nil {
		log.Error("Error converting \"cleansess\" to a boolean ", err.Error())
		return
	}
	opts.SetCleanSession(b)
	if t.config["store"] != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(t.config["store"]))
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())

		// Match suffix of topic
		if strings.HasSuffix(topic, "start") {
			t.StartProcess(payload)
		} else if strings.HasSuffix(topic, "restart") {
			t.RestartProcess(payload)
		} else if strings.HasSuffix(topic, "resume") {
			t.ResumeProcess(payload)
		}
	})

	client := mqtt.NewClient(opts)
	t.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	i, err := strconv.Atoi(t.config["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}

	if token := t.client.Subscribe(t.config["topic"], byte(i), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		panic(token.Error())
	}
}

// Stop implements ext.Trigger.Stop
func (t *MqttTrigger) Stop() {
	//unsubscribe from topic
	if token := t.client.Unsubscribe(t.config["topic"]); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		log.Info("Unsubcribing from topic: ", t.config["topic"])
	}

	t.client.Disconnect(250)
}

// StartProcess starts a new Process Instance
func (t *MqttTrigger) StartProcess(payload string) {

	req := &starter.StartRequest{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Error Starting process ", err.Error())
		return
	}

	log.Info("Process URI ", req.ProcessURI)
	log.Info("processStarter.StartProcess ", t.processStarter)
	id := t.processStarter.StartProcess(req)
	log.Info("Start process id: ", id)
	t.publishMessage(req.ReplyTo, id)
}

// RestartProcess restarts a Process Instance
func (t *MqttTrigger) RestartProcess(payload string) {

	req := &starter.RestartRequest{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Error Restarting process ", err.Error())
		return
	}

	id := t.processStarter.RestartProcess(req)
	log.Info("Restart process id: ", id)
}

// ResumeProcess resumes a Process Instance
func (t *MqttTrigger) ResumeProcess(payload string) {

	defer func() {
		if r := recover(); r != nil {
			log.Error("Unable to resume process, make sure definition registered")
		}
	}()

	req := &starter.ResumeRequest{}
	err := json.NewDecoder(strings.NewReader(payload)).Decode(req)
	if err != nil {
		log.Error("Error Resuming process ", err.Error())
		return
	}

	t.processStarter.ResumeProcess(req)
}

func (t *MqttTrigger) publishMessage(topic string, message string) {

	log.Info("ReplyTo topic: ", topic)
	log.Info("Publishing message: ", message)

	qos, err := strconv.Atoi(t.config["qos"])
	if err != nil {
		log.Error("Error converting \"qos\" to an integer ", err.Error())
		return
	}
	token := t.client.Publish(topic, byte(qos), false, message)
	token.Wait()
}
