package kafkapub

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

/*
The setup for the kafka server used to run these tests:
section from server.properties:
	listeners=PLAINTEXT://cheetah.na.tibco.com:9092,SSL://cheetah:9093,SASL_PLAINTEXT://cheetah:9094,SASL_SSL://cheetah.com:9095

	sasl.enabled.mechanisms=PLAIN
	sasl.mechanism.inter.broker.protocol=PLAIN

	ssl.keystore.location=/opt/kafka_2.12-0.10.2.1/keys/kafka.jks
	ssl.keystore.password=sauron
	ssl.key.password=sauron
	ssl.truststore.location=/opt/kafka_2.12-0.10.2.1/keys/kafka.jks
	ssl.truststore.password=sauron
	ssl.client.auth=none
	ssl.enabled.protocols=TLSv1.2,TLSv1.1,TLSv1

	advertised.listeners=PLAINTEXT://cheetah.na.tibco.com:9092,SSL://cheetah:9093,SASL_PLAINTEXT://cheetah:9094,SASL_SSL://cheetah:9095

The SASL file:
	KafkaServer {
	org.apache.kafka.common.security.plain.PlainLoginModule required
	username="admin"
	password="admin"
	user_wcn00="sauron"
	user_alice="sissy";
	};
To get kafka to pick up the jaas file add a vm parm like:
	-Djava.security.auth.login.config=/local/opt/kafka/kafka_2.11-0.10.2.0/config/kafka_server_jaas.conf
*/
var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
	log.Println("TestCreate successfull")
}

func TestPlain(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("BrokerUrls", "cheetah:9092")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("Message", "######### PLAIN ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}

func TestSSL(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("BrokerUrls", "cheetah:9093")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("truststore", "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust")
	tc.SetInput("Message", "######### TLS ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}

/*
Run this test with FLOGO_LOG_LEVEL=DEBUG and observe the debug messages as the activity either creates new synpublishers, or reuses cached ones.
*/
func TestCache(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("BrokerUrls", "cheetah:9093")
	tc.SetInput("truststore", "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("Message", "######### TLS ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	tc.SetInput("Topic", "publishpet")
	act.Eval(tc)
	tc.SetInput("Topic", "publishpet")
	act.Eval(tc)
	tc = test.NewTestActivityContext(getActivityMetadata())
	tc.SetInput("Message", "######### TLS ###########  Mary had a little lamb, its fleece was white as snow.")
	tc.SetInput("BrokerUrls", "cheetah:9092")
	tc.SetInput("Topic", "syslog")
	act.Eval(tc)
	tc.SetInput("Topic", "publishpet")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}

func TestSASL_PLAIN(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("BrokerUrls", "cheetah:9094")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("user", "wcn00")
	tc.SetInput("password", "sauron")
	tc.SetInput("Message", "######### SASL_PLAIN ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}

func TestSASL_TLS(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("BrokerUrls", "cheetah:9095")
	tc.SetInput("truststore", "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("user", "wcn00")
	tc.SetInput("password", "sauron")
	tc.SetInput("Message", "######### SASL_TLS ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}
