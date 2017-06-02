package kafkapub

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

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
	tc.SetInput("BrokerUrls", "bilbo:9092")
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
	tc.SetInput("BrokerUrls", "bilbo:9093")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("truststore", "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust")
	tc.SetInput("Message", "######### TLS ###########  Mary had a little lamb, its fleece was white as snow.")
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
	tc.SetInput("BrokerUrls", "bilbo:9094")
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
	tc.SetInput("BrokerUrls", "bilbo:9095")
	tc.SetInput("truststore", "/opt/kafka/kafka_2.11-0.10.2.0/keys/trust")
	tc.SetInput("Topic", "syslog")
	tc.SetInput("user", "wcn00")
	tc.SetInput("password", "sauron")
	tc.SetInput("Message", "######### SASL_TLS ###########  Mary had a little lamb, its fleece was white as snow.")
	act.Eval(tc)
	log.Printf("TestEval successfull.  partition [%d]  offset [%d]", tc.GetOutput("partition"), tc.GetOutput("offset"))
}
