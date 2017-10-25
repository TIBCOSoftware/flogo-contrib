package couchbase

import (
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
	"time"
)

/*
To setup the testing environment you must follow the three following steps:

1) Run an instance of Couchbase cluster 5.0.0-beta2.
Docker command: docker run -d -v ~/couchbase/node1:/opt/couchbase/var -p 8091-8094:8091-8094 -p 11210:11210 couchbase/server:5.0.0-beta2

2) You must add in the etc hosts file the host "couchbase" referencing the Couchbase server ip

3) You must create a bucket in Couchbase called "test"

The tests are using the default Couchbase username, password and port
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

func randomString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func insert(id string, data string, t *testing.T) (act activity.Activity, tc *test.TestActivityContext) {
	act = NewActivity(getActivityMetadata())
	tc = test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("key", id)
	if data == "" {
		tc.SetInput("data", `{"name":"foo"}`)
	} else {
		tc.SetInput("data", data)
	}

	tc.SetInput("method", `Insert`)
	tc.SetInput("expiry", 0)
	tc.SetInput("server", "couchbase://couchbase")
	tc.SetInput("username", "Administrator")
	tc.SetInput("password", "password")
	tc.SetInput("bucket", "test")
	tc.SetInput("bucketPassword", "")

	_, insertError := act.Eval(tc)
	if insertError != nil {
		t.Error("Document not inserted")
		t.Fail()
	}
	return
}

/*
Test a create activity
*/
func TestCreate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
	log.Println("TestCreate successful")
}

/*
Insert test
*/
func TestInsert(t *testing.T) {
	id := randomString(5)
	insert(id, "", t)
	log.Println("TestInsert successful")
}

/*
Remove test
*/
func TestRemove(t *testing.T) {
	id := randomString(5)
	act, tc := insert(id, "", t)

	tc.SetInput("method", `Remove`)

	_, RemoveError := act.Eval(tc)
	if RemoveError != nil {
		t.Error("Document not removed")
		t.Fail()
	}
	log.Println("TestRemove successful")
}

/*
Upsert test
*/
func TestUpsert(t *testing.T) {
	id := randomString(5)
	act, tc := insert(id, "", t)

	tc.SetInput("method", `Upsert`)

	_, RemoveError := act.Eval(tc)
	if RemoveError != nil {
		t.Error("Document not upserted")
		t.Fail()
	}
	log.Println("TestUpsert successful")
}

/*
Get test
*/
func TestGet(t *testing.T) {
	id := randomString(5)
	data := `{"name":"foo"}`
	act, tc := insert(id, data, t)

	tc.SetInput("method", `Get`)

	_, GetError := act.Eval(tc)
	if GetError != nil {
		t.Error("Document not retrieved")
		t.Fail()
	}
	s := tc.GetOutput("output").(string)

	if s != data {
		t.Error("The retrieved document is not equals")
		t.Fail()
	}

	log.Println("TestGet successful")
}
