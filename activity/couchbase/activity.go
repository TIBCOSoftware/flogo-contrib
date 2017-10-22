package couchbase

import (
	"fmt"
	"gopkg.in/couchbase/gocb.v1"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-tibco-couchbase")

const (
	methodInsert = "Insert"
	methodUpsert = "Upsert"

	ivKey            = "key"
	ivData           = "data"
	ivMethod         = "method"
	ivExpiry         = "expiry"
	ivServer         = "server"
	ivUsername       = "username"
	ivPassword       = "password"
	ivBucket         = "bucket"
	ivBucketPassword = "bucketPassword"

	ovOutput = "output"
	ovStatus = "status"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

// Integration with Couchbase
// inputs: {data, method, expiry, server, username, password, bucket, bucketPassword}
// outputs: {output, status}
type CouchbaseActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CouchbaseActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *CouchbaseActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Couchbase integration
func (a *CouchbaseActivity) Eval(context activity.Context) (done bool, err error) {

	key, _ := context.GetInput(ivKey).(string)
	data, _ := context.GetInput(ivData).(string)
	method, _ := context.GetInput(ivMethod).(string)
	expiry, _ := context.GetInput(ivExpiry).(int)
	server, _ := context.GetInput(ivServer).(string)
	username, _ := context.GetInput(ivUsername).(string)
	password, _ := context.GetInput(ivPassword).(string)
	bucketName, _ := context.GetInput(ivBucket).(string)
	bucketPassword, _ := context.GetInput(ivBucketPassword).(string)

	cluster, connectError := gocb.Connect(server)
	if connectError != nil {
		activityLog.Errorf("Connection error: %v", connectError)
		return false, connectError
	}

	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: username,
		Password: password,
	})

	bucket, openBucketError := cluster.OpenBucket(bucketName, bucketPassword)
	if openBucketError != nil {
		activityLog.Errorf("Error while opening the bucked with the specified credentials: ", openBucketError)
		return false, openBucketError
	}

	switch method {
	case methodInsert:
		cas, error := bucket.Insert(key, data, uint32(expiry))
		if error != nil {
			activityLog.Errorf("Insert error: %v", error)
			return false, error
		} else {
			context.SetOutput(ovOutput, cas)
			return true, nil
		}
	case methodUpsert:
		cas, error := bucket.Upsert(key, data, uint32(expiry))
		if error != nil {
			activityLog.Errorf("Upsert error: %v", error)
			return false, error
		} else {
			context.SetOutput(ovOutput, cas)
			return true, nil
		}
	default:
		activityLog.Errorf("Method %v not recognized.", method)
		return false, fmt.Errorf("method %v not recognized", method)
	}

	return true, nil
}
