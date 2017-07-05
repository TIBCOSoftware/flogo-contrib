package kafkapub

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// log is the default package logger
var flogoLogger = logger.GetLogger("activity-tibco-kafkapub")

// MyActivity is a stub for your Activity implementation
type KafkaPubActivity struct {
	metadata     *activity.Metadata
	kafkaConfig  *sarama.Config
	brokers      []string
	topic        string
	syncProducer sarama.SyncProducer
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &KafkaPubActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *KafkaPubActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *KafkaPubActivity) Eval(context activity.Context) (done bool, err error) {
	err = initParms(a, context)
	if err != nil {
		flogoLogger.Errorf("Kafkapub parameters initialization got error: [%s]", err.Error())
		return false, err
	}
	defer func() {
		if err := a.syncProducer.Close(); err != nil {
			flogoLogger.Errorf("Kafkapub producer close got error: [%s]", err.Error())
		}
	}()
	if context.GetInput("Message") != nil {
		msg := &sarama.ProducerMessage{
			Topic: a.topic,
			Value: sarama.StringEncoder(context.GetInput("Message").(string)),
		}
		partition, offset, err := a.syncProducer.SendMessage(msg)
		if err != nil {
			return false, fmt.Errorf("kafkapub failed to send message for reason [%s]", err.Error())
		}
		context.SetOutput("partition", partition)
		context.SetOutput("offset", offset)
		flogoLogger.Debugf("Kafkapub message [%s] sent successfully on partition [%d] and offset [%d]",
			context.GetInput("Message").(string), partition, offset)
		return true, nil
	}
	return false, fmt.Errorf("kafkapub called without a message to publish")
}

func initParms(a *KafkaPubActivity, context activity.Context) error {
	if context.GetInput("BrokerUrls") != nil && context.GetInput("BrokerUrls").(string) != "" {
		a.kafkaConfig = sarama.NewConfig()
		a.kafkaConfig.Producer.Return.Errors = true
		a.kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
		a.kafkaConfig.Producer.Retry.Max = 5
		a.kafkaConfig.Producer.Return.Successes = true
		brokerUrls := strings.Split(context.GetInput("BrokerUrls").(string), ",")
		brokers := make([]string, len(brokerUrls))
		for brokerNo, broker := range brokerUrls {
			error := validateBrokerUrl(broker)
			if error != nil {
				return fmt.Errorf("BrokerUrl [%s] format invalid for reason: [%s]", broker, error.Error())
			}
			brokers[brokerNo] = broker
		}
		a.brokers = brokers
		flogoLogger.Debugf("Kafkapub brokers [%v]", brokers)
	} else {
		return fmt.Errorf("Kafkapub activity is not configured with at least one BrokerUrl")
	}
	if context.GetInput("Topic") != nil && context.GetInput("Topic").(string) != "" {
		a.topic = context.GetInput("Topic").(string)
		flogoLogger.Debugf("Kafkapub topic [%s]", a.topic)
	} else {
		return fmt.Errorf("Topic input parameter not provided and is required")
	}

	//clientKeystore
	/*
		Its worth mentioning here that when the keystore for kafka is created it must support RSA keys via
		the -keyalg RSA option.  If not then there will be ZERO overlap in supported cipher suites with java.
		see:   https://issues.apache.org/jira/browse/KAFKA-3647
		for more info
	*/
	if context.GetInput("truststore") != nil {
		trustStore := context.GetInput("truststore")
		if trustStore != nil && len(trustStore.(string)) > 0 {
			trustPool, err := getCerts(trustStore.(string))
			if err != nil {
				return err
			}
			config := tls.Config{
				RootCAs:            trustPool,
				InsecureSkipVerify: true}
			a.kafkaConfig.Net.TLS.Enable = true
			a.kafkaConfig.Net.TLS.Config = &config
		}
		flogoLogger.Debugf("Kafkapub initialized truststore from [%s]", trustStore)
	}
	// SASL
	if context.GetInput("user") != nil && context.GetInput("user").(string) != "" {
		var password string
		user := context.GetInput("user").(string)
		if context.GetInput("password") != nil {
			password = context.GetInput("password").(string)
		}
		a.kafkaConfig.Net.SASL.Enable = true
		a.kafkaConfig.Net.SASL.User = user
		a.kafkaConfig.Net.SASL.Password = password
		flogoLogger.Debugf("Kafkapub SASL parms initialized; user [%s]  password[########]", user)
	}

	syncProducer, err := sarama.NewSyncProducer(a.brokers, a.kafkaConfig)
	if err != nil {
		return fmt.Errorf("Kafkapub failed to create a SyncProducer.  Check any TLS or SASL parameters carefully.  Reason given: [%s]", err)
	}
	a.syncProducer = syncProducer

	flogoLogger.Debug("Kafkapub synchronous producer created")
	return nil
}

//Ensure that this string meets the host:port definition of a kafka hostspec
//Kafka calls it a url but its really just host:port, which for numeric ip addresses is not a valid URI
//technically speaking.
func validateBrokerUrl(broker string) error {
	hostport := strings.Split(broker, ":")
	if len(hostport) != 2 {
		return fmt.Errorf("BrokerUrl must be composed of sections like \"host:port\"")
	}
	i, err := strconv.Atoi(hostport[1])
	if err != nil || i < 0 || i > 32767 {
		return fmt.Errorf("Port specification [%s] is not numeric and between 0 and 32767", hostport[1])
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
			flogoLogger.Warnf("Failed to read trusted certificate [%s] ... continuing", trustCertFile.Name())
		} else if trustCertBytes != nil {
			certPool.AppendCertsFromPEM(trustCertBytes)
		}
	}
	if len(certPool.Subjects()) < 1 {
		return certPool, fmt.Errorf("Failed to read trusted certificates from [%s]  After processing all files in the directory no valid trusted certs were found", trustStore)
	}
	return certPool, nil
}
