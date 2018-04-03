package mongodb

import (
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	tc.SetInput("server", "mongodb://mongodb")
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

// TestDelete
func TestDelete(t *testing.T) {

}

// TestInsert
func TestInsert(t *testing.T) {

}

// TestReplace
func TestReplace(t *testing.T) {

}

// TestReplace
func TestUpdate(t *testing.T) {

}
