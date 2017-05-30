package kafkasub

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"fmt"
	"net/url"

	"time"

	"crypto/tls"
	"crypto/x509"

	"io/ioutil"

	"github.com/Shopify/sarama"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var flogoLogger = logger.GetLogger("trigger-tibco-kafkasub")

type _kafkaParms struct {
	brokers    []string
	topic      string
	offset     int64
	group      string
	partitions []int32
}

// MqttTrigger is simple MQTT trigger
type KafkaSubTrigger struct {
	metadata           *trigger.Metadata
	runner             action.Runner
	config             *trigger.Config
	kafkaParms         _kafkaParms
	shutdownChan       *chan struct{}
	signals            *chan os.Signal
	kafkaConfig        *sarama.Config
	kafkaConsumer      *sarama.Consumer
	partitionConsumers *map[int32]sarama.PartitionConsumer
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &KafkasubFactory{metadata: md}
}

// KafkasubFactory Kafka subscriber Trigger factory
type KafkasubFactory struct {
	metadata *trigger.Metadata
}

//New Creates a new trigger instance for a given id
func (t *KafkasubFactory) New(config *trigger.Config) trigger.Trigger {
	kafkaTrigger := &KafkaSubTrigger{metadata: t.metadata, config: config}
	kafkaTrigger.kafkaConfig = sarama.NewConfig()
	kafkaTrigger.kafkaConfig.Consumer.Return.Errors = true
	error := initKafkaParms(kafkaTrigger)
	if error != nil {
		flogoLogger.Error("Failed to initialize parms: ", error)
		return nil
	}
	//flogoLogger.Debug("KafkaSubTrigger New complete")
	return kafkaTrigger
}

// Metadata implements trigger.Trigger.Metadata
func (t *KafkaSubTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *KafkaSubTrigger) Init(runner action.Runner) {
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *KafkaSubTrigger) Start() error {
	shutdownChan := make(chan struct{})
	t.shutdownChan = &shutdownChan
	signals := make(chan os.Signal, 1)
	t.signals = &signals
	signal.Notify(*t.signals, os.Interrupt)
	go run(t)
	flogoLogger.Debug("KafkaSubTrigger Started")
	return nil
}

// Stop implements ext.Trigger.Stop
func (t *KafkaSubTrigger) Stop() error {
	//unsubscribe from topic
	flogoLogger.Debugf("Unsubcribing from topic [%s]", t.kafkaParms.topic)
	if t.partitionConsumers == nil {
		flogoLogger.Debug("Closed called for a subscriber with no running consumers")
		flogoLogger.Debug("Stopped")
		return nil
	}

	for id, partitionConsumer := range *t.partitionConsumers {
		partitionConsumer.Close()
		flogoLogger.Debug("Closed partition consumer:", id)
	}
	if t.kafkaConsumer != nil {
		(*t.kafkaConsumer).Close()
		flogoLogger.Debug("Closed kafka consumer")
	}
	flogoLogger.Debug("Stopped")
	return nil
}

func run(t *KafkaSubTrigger) {
	kafkaConsumer, error := sarama.NewConsumer(t.kafkaParms.brokers, t.kafkaConfig)
	if error != nil {
		flogoLogger.Errorf("Failed to create Kafka consumer for reason [%s]", error)
		return
	}
	t.kafkaConsumer = &kafkaConsumer
	thing := make(map[int32]sarama.PartitionConsumer)
	t.partitionConsumers = &thing
	flogoLogger.Debugf("Subscribing to topic [%s]", t.kafkaParms.topic)
	validPartitions, error := kafkaConsumer.Partitions(t.kafkaParms.topic)
	if error != nil {
		flogoLogger.Errorf("Failed to get valid partitions from Kafka Consumer for reason [%s].  Aborting subscriber", error)
		return
	}

	flogoLogger.Debugf("Valid partitions for topic [%s] detected as: [%v]", t.kafkaParms.topic, validPartitions)
	if t.kafkaParms.partitions == nil { //subscribe to all valid partitions
		t.kafkaParms.partitions = validPartitions
		for _, part := range validPartitions {
			startConsumer(t, part)
		}
	} else { //subscribe to a subset of valid partitions
		configPartitions := t.kafkaParms.partitions
		for _, confPart := range configPartitions {
			for _, valPart := range validPartitions {
				if confPart == valPart {
					startConsumer(t, confPart)
					break
				}
			}
			flogoLogger.Errorf("Configured partition [%d] on topic [%s] does not exist", confPart, t.kafkaParms.topic)
		}
	}
	if len(*t.partitionConsumers) < 1 {
		t.Stop()
		flogoLogger.Errorf("Kafka consumer is not configured for any valid partitions and will shut down")
	} else {
		flogoLogger.Debugf("Kafka consumers for topic [%s] started", t.kafkaParms.topic)
	}
	return
}

func startConsumer(t *KafkaSubTrigger, part int32) error {
	flogoLogger.Debugf("Creating PartitionConsumer for valid partition: [%s:%d]", t.kafkaParms.topic, part)
	consumer := *t.kafkaConsumer
	partitionConsumer, error := consumer.ConsumePartition(t.kafkaParms.topic, part, t.kafkaParms.offset)
	if error != nil {
		flogoLogger.Errorf("Creating PartitionConsumer for valid partition: [%s:%d] failed for reason: %s", t.kafkaParms.topic, part, error)
		return error
	}
	(*t.partitionConsumers)[part] = partitionConsumer
	go consumePartition(t, partitionConsumer, part)
	return nil
}

func consumePartition(t *KafkaSubTrigger, consumer sarama.PartitionConsumer, part int32) {
	for {
		select {
		case err := <-consumer.Errors():
			if err == nil {
				//been shutdown
				return
			}
			flogoLogger.Warnf("PartitionConsumer [%d] got error: [%s]", part, err)
			time.Sleep(time.Millisecond * 300)
		case msg := <-consumer.Messages():
			onMessage(t, msg)
		case <-*t.signals:
			flogoLogger.Infof("Partition consumer got SIGINT; exiting")
			*t.shutdownChan <- struct{}{}
			return
		case <-*t.shutdownChan:
			flogoLogger.Infof("Partition consumer got SHUTDOWN signal; exiting")
			return
		}
	}
}

func getTopics(t *KafkaSubTrigger) []string {
	return strings.Split(t.config.Settings["Topics"].(string), ",")
}

func initKafkaParms(t *KafkaSubTrigger) error {
	brokersString := t.config.Settings["BrokerUrl"]
	if brokersString == nil || brokersString.(string) == "" {
		return fmt.Errorf("BrokerUrl not provided")
	}
	brokers := strings.Split(brokersString.(string), ",")
	if len(brokers) < 1 {
		return fmt.Errorf("BrokerUrl [%s] is invalid, require at least one broker", brokersString.(string))
	}
	t.kafkaParms.brokers = make([]string, len(brokers))
	for brokerNo, broker := range brokers {
		_, error := url.Parse(broker)
		if error != nil {
			return fmt.Errorf("BrokerUrl [%s] format invalid for reason: [%s]", t.kafkaParms.brokers[0], error.Error())
		}
		t.kafkaParms.brokers[brokerNo] = broker
	}
	if t.config.Handlers == nil || len(t.config.Handlers) < 1 {
		return fmt.Errorf("Kafka trigger requires at least one handler containing a valid topic name")
	}
	for _, handler := range t.config.Handlers {
		if handler.Settings["Topic"] == nil {
			return fmt.Errorf("Topic string was not provided")
		}
		t.kafkaParms.topic = handler.Settings["Topic"].(string)

		//offset
		if handler.Settings["offset"] != nil {
			i, error := strconv.Atoi(handler.Settings["offset"].(string))
			if error != nil {
				t.kafkaParms.offset = sarama.OffsetNewest
			} else {
				t.kafkaParms.offset = int64(i)
			}
		} else {
			t.kafkaParms.offset = sarama.OffsetNewest
		}
		//partitions
		if handler.Settings["partitions"] != nil {
			partitions := handler.Settings["partitions"].(string)
			if partitions == "" {
				t.kafkaParms.partitions = nil
			} else {
				i := 0
				parts := strings.Split(partitions, ",")
				t.kafkaParms.partitions = make([]int32, len(parts))
				for _, p := range parts {
					n, error := strconv.Atoi(p)
					if error == nil {
						t.kafkaParms.partitions[i] = int32(n)
						i++
					} else {
						flogoLogger.Warn("Discarding non-numeric partition [%s]", p)
					}
				}
			}
		} else {
			t.kafkaParms.partitions = nil
		}

		//group
		if handler.Settings["group"] != nil {
			group := handler.Settings["group"].(string)
			if len(group) > 0 {
				t.kafkaParms.group = group
			}
		} else {
			t.kafkaParms.group = ""
		}
		//user
		if handler.Settings["user"] != nil {
			user := handler.Settings["user"].(string)
			if len(user) > 0 {
				t.kafkaConfig.Net.SASL.Enable = true
				t.kafkaConfig.Net.SASL.User = user
			}
		}
		//password
		if handler.Settings["password"] != nil {
			password := handler.Settings["password"].(string)
			if len(password) > 0 {
				t.kafkaConfig.Net.SASL.Password = password
			}
		}
		//clientKeystore
		/*
			Its worth mentioning here that when the keystore for kafka is created it must support RSA keys via
			the -keyalg RSA option.  If not then there will be ZERO overlap in supported cipher suites with java.
			see:   https://issues.apache.org/jira/browse/KAFKA-3647
			for more info
		*/
		if handler.Settings["truststore"] != nil {
			trustStore := handler.Settings["truststore"]
			if trustStore != nil && len(trustStore.(string)) > 0 {
				trustPool, err := getCerts(trustStore.(string))
				if err != nil {
					return err
				}
				config := tls.Config{
					RootCAs:            trustPool,
					InsecureSkipVerify: true}
				t.kafkaConfig.Net.TLS.Enable = true
				t.kafkaConfig.Net.TLS.Config = &config
			}
		}
		// SASL
		if handler.Settings["user"] != nil {
			var password string
			user := handler.Settings["user"]
			if user != nil {
				if handler.Settings["password"] == nil {
					//TODO  Can't read password from shell reliably on different systems
					//      need a secure way of getting it...
					return fmt.Errorf("Password not provided")
				}
				password = handler.Settings["password"].(string)
				t.kafkaConfig.Net.SASL.Enable = true
				t.kafkaConfig.Net.SASL.User = user.(string)
				t.kafkaConfig.Net.SASL.Password = password
			}
		}
		return nil
	}
	return nil
}

func getCerts(trustStore string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()
	fileInfo, err := os.Stat(trustStore)
	if err != nil {
		return certPool, fmt.Errorf("Truststore [%s] does not exist", trustStore)
	}
	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		break
	case mode.IsRegular():
		return certPool, fmt.Errorf("Truststore [%s] is not a directory.  Must be a directory containing trusted certificates in PEM format",
			trustStore)
	}
	trustedCertFiles, err := ioutil.ReadDir(trustStore)
	if err != nil || len(trustedCertFiles) == 0 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  Must be a directory containing trusted certificates in PEM format", trustStore)
	}
	for _, trustCertFile := range trustedCertFiles {
		fqfName := fmt.Sprintf("%s%c%s", trustStore, os.PathSeparator, trustCertFile.Name())
		trustCertBytes, err := ioutil.ReadFile(fqfName)
		if err != nil {
			flogoLogger.Warnf("Failed to read trusted certificate [%s] ... continueing", trustCertFile.Name())
		}
		certPool.AppendCertsFromPEM(trustCertBytes)
	}
	if len(certPool.Subjects()) < 1 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  After processing all files in the directory no valid trusted certs were found", trustStore)
	}
	return certPool, nil
}

func onMessage(t *KafkaSubTrigger, msg *sarama.ConsumerMessage) {
	if msg == nil {
		return
	}
	flogoLogger.Debugf("Kafka subscriber triggering job from topic [%s] on partition [%d] with key [%s] at offset [%d]",
		msg.Topic, msg.Partition, msg.Key, msg.Offset)
	for _, handler := range t.config.Handlers {
		actionID := action.Get(handler.ActionId)
		flogoLogger.Debugf("Found action: '%+x' for ActionID: %s", actionID, handler.ActionId)
		if t.metadata == nil {
			flogoLogger.Infof("Kafka subscriber running in test mode received message on topic [%s] on partition [%d] with key [%s] at offset [%d]",
				msg.Topic, msg.Partition, msg.Key, msg.Offset)
			flogoLogger.Infof("Content: [%s]", string(msg.Value))

			return
		}
		data := make(map[string]interface{})
		data["message"] = string(msg.Value)
		//if(t.metadata.Metadata.OutPuts
		startAttrs, errorAttrs := t.metadata.OutputsToAttrs(data, true)
		if errorAttrs != nil || startAttrs == nil {
			flogoLogger.Errorf("Failed to create output attributes for kafka message for ActionID [%s] for reason [%s] message lost", handler.ActionId, errorAttrs)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		_, _, err := t.runner.Run(ctx, actionID, handler.ActionId, nil)

		if err != nil {
			flogoLogger.Errorf("Run action for ActionID [%s] failed for reason [%s] message lost", err, handler.ActionId)
		}
	}

}
