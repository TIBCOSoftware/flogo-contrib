package kafkasub

/*
This is the Kafka server setup to support these tests:

sasl.enabled.mechanisms=PLAIN
sasl.mechanism.inter.broker.protocol=PLAIN
advertised.listeners=PLAINTEXT://bilbo:9092,SSL://bilbo:9093,SASL_PLAINTEXT://bilbo:9094,SASL_SSL://bilbo:9095

ssl.keystore.location=/local/opt/kafka/kafka_2.11-0.10.2.0/keys/kafka.jks
ssl.keystore.password=sauron
ssl.key.password=sauron
ssl.truststore.location=/local/opt/kafka/kafka_2.11-0.10.2.0/keys/kafka.jks
ssl.truststore.password=sauron
ssl.client.auth=none
ssl.enabled.protocols=TLSv1.2,TLSv1.1,TLSv1


*/
import (
	"context"
	"encoding/json"
	"os/signal"
	"syscall"
	"testing"

	"io/ioutil"

	"time"

	"log"
	"os"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJsonMetadata()
var listentime time.Duration = 10

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "name": "tibco-kafkasub",
  "settings": {
    "BrokerUrl": "bilbo:9092"
  },
  "handlers": [
    {
      "actionId": "kafka_message",
      "settings": {
        "Topic": "syslog"
      }
    }
  ],
	"outputs": [
    {
      "name": "message",
      "type": "string"
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Printf("Ran Action: %v", uri)
	return 0, nil, nil
}

func consoleHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received console interrupt.  Exiting.")
		time.Sleep(time.Second * 3)
		os.Exit(1)
	}()
}
func TestInit(t *testing.T) {
	consoleHandler()
	f := &KafkasubFactory{}
	config := trigger.Config{}
	error := json.Unmarshal([]byte(testConfig), &config)
	if error != nil {
		log.Printf("Failed to unmarshal the config args:%s", error)
		t.Fail()
	}
	tgr := f.New(&config)
	log.Printf("TestInit: Successfully created the trigger object")
	runner := &TestRunner{}
	tgr.Init(runner)
	log.Printf("TestInit: Successfully initialized the trigger object")
}

func TestEndpoint(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	log.Printf("TestEndpoint: Successfully started the trigger on a non-authenticated plain text port.  It will consume and print messages for %d seconds", listentime)

	time.Sleep(time.Second * listentime)
	tgr.Stop()
	log.Printf("TestEndpoint: Successfully stopped the trigger.")
	time.Sleep(time.Second * 2)
}

func TestMultiBrokers(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	config.Settings["BrokerUrl"] = "bilbo:9092,bilbo:9092"
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	tgr.Start()
	log.Printf("TestEndpoint: Successfully started the trigger on a non-authenticated plain text port.  It will consume and print messages for %d seconds", listentime)

	time.Sleep(time.Second * listentime)
	tgr.Stop()
	log.Printf("TestEndpoint: Successfully stopped the trigger.")
	time.Sleep(time.Second * 2)
}

func TestFailingEndpoint(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	config.Handlers[0].Settings["partitions"] = "21,31"
	tgr := f.New(&config)
	runner := &TestRunner{}
	tgr.Init(runner)
	log.Printf("TestFailingEndpoint: Should detect that none of the specified partitions exist and shutdown.")
	tgr.Start()
	defer tgr.Stop()
	time.Sleep(time.Second * 2)
}

func TestTLS(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	config.Handlers[0].Settings["truststore"] = "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust"
	config.Settings["BrokerUrl"] = "bilbo:9093"
	tgr := f.New(&config)
	if tgr == nil {
		log.Printf("Failed to create trigger")
		return
	}
	runner := &TestRunner{}
	tgr.Init(runner)

	err := tgr.Start()
	if err != nil {
		log.Printf("Trigger Star failed: %s", err)
		return
	}
	log.Printf("TestTLS: Successfully started the trigger on a TLS port.  It will consume and print messages for %d seconds", listentime)

	time.Sleep(time.Second * listentime)
	defer tgr.Stop()
	log.Printf("TestTLS: Successfully stopped the trigger.")
	time.Sleep(time.Second * 2)
}

func TestSASL(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	config.Handlers[0].Settings["user"] = "wcn00"
	config.Handlers[0].Settings["password"] = "sauron"
	config.Settings["BrokerUrl"] = "bilbo:9094"
	tgr := f.New(&config)
	if tgr == nil {
		log.Printf("Failed to create trigger")
		return
	}
	runner := &TestRunner{}
	tgr.Init(runner)

	err := tgr.Start()
	if err != nil {
		log.Printf("Trigger Start failed: %s", err)
		return
	}
	log.Printf("TestSASL: Successfully started the trigger on a plaintext port using SASL.  It will consume and print messages for %d seconds", listentime)

	time.Sleep(time.Second * listentime)
	defer tgr.Stop()
	log.Printf("TestSASL: Successfully stopped the trigger.")
	time.Sleep(time.Second * 2)
}

func TestSASL_SSL(t *testing.T) {
	f := &KafkasubFactory{}
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	config.Handlers[0].Settings["truststore"] = "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust"
	config.Handlers[0].Settings["user"] = "wcn00"
	config.Handlers[0].Settings["password"] = "sauron"
	config.Settings["BrokerUrl"] = "bilbo:9095"

	tgr := f.New(&config)
	if tgr == nil {
		log.Printf("TestSASL_SSL Failed to create trigger")
		return
	}
	runner := &TestRunner{}
	tgr.Init(runner)

	err := tgr.Start()
	if err != nil {
		log.Printf("Trigger start failed: %s", err)
		return
	}
	log.Printf("TestSASL_SSL: Successfully started the trigger.  It will consume and print messages for %d seconds", listentime)

	time.Sleep(time.Second * listentime)
	defer tgr.Stop()
	log.Printf("TestSASL_SSL: Successfully stopped the trigger.")
	time.Sleep(time.Second * 2)
}
