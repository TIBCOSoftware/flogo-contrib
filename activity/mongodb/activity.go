package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-mongodb")

const (
	methodGet     = "GET"
	methodDelete  = "DELETE"
	methodInsert  = "INSERT"
	methodReplace = "REPLACE"
	methodUpdate  = "UPDATE"

	ivConnectionURI = "uri"
	ivDbName        = "dbName"
	ivCollection    = "collection"
	ivMethod        = "method"

	ivKeyName  = "keyName"
	ivKeyValue = "keyValue"
	ivData     = "data"

	ovOutput = "output"
	ovCount  = "count"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

/*
Integration with MongoDb
inputs: {uri, dbName, collection, method, [keyName, keyValue, value]}
outputs: {output, count}
*/
type MongoDbActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MongoDbActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MongoDbActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - MongoDb integration
func (a *MongoDbActivity) Eval(ctx activity.Context) (done bool, err error) {

	//mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	connectionURI, _ := ctx.GetInput(ivConnectionURI).(string)
	dbName, _ := ctx.GetInput(ivDbName).(string)
	collectionName, _ := ctx.GetInput(ivCollection).(string)
	method, _ := ctx.GetInput(ivMethod).(string)
	keyName, _ := ctx.GetInput(ivKeyName).(string)
	keyValue, _ := ctx.GetInput(ivKeyValue).(string)
	value := ctx.GetInput(ivData)

	//todo implement shared sessions
	// client, err := mongo.NewClient(connectionURI)
	/*
		The above function was giving below error;
		"data not inserted topology is closed"
	*/
	bCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(bCtx, options.Client().ApplyURI(connectionURI))

	defer cancel()

	defer client.Disconnect(context.Background())
	if err != nil {
		activityLog.Errorf("Connection error: %v", err)
		return false, err
	}

	db := client.Database(dbName)
	coll := db.Collection(collectionName)

	switch strings.ToUpper(method) {
	case methodGet:
		result := coll.FindOne(bCtx, bson.M{keyName: keyValue})
		val := make(map[string]interface{})
		err := result.Decode(val)
		if err != nil {
			activityLog.Debug("Error during getting data ..", err)
			return false, err
		}

		activityLog.Debugf("Get Results $#v", result)

		ctx.SetOutput(ovOutput, val)
	case methodDelete:
		result, err := coll.DeleteOne(bCtx, bson.M{keyName: keyValue}, nil)
		if err != nil {
			activityLog.Debug("Error during deleting data ..", err)
			return false, err
		}

		activityLog.Debugf("Delete Results $#v", result)

		ctx.SetOutput(ovCount, result.DeletedCount)
	case methodInsert:
		if value == nil && keyValue == "" {
			// should we throw an error or warn?
			activityLog.Warnf("Nothing to insert")
			return true, nil
		}

		var result *mongo.InsertOneResult

		if value != nil && keyValue == "" {
			result, err = coll.InsertOne(bCtx, value)
			if err != nil {
				activityLog.Debug("Error during adding data ..", err)
				return false, err
			}

		} else {
			result, err = coll.InsertOne(bCtx, bson.M{keyName: keyValue})
			if err != nil {
				activityLog.Debug("Error during adding data ..", err)
				return false, err
			}
		}

		activityLog.Debugf("Insert Results $#v", result)

		ctx.SetOutput(ovOutput, result.InsertedID)
	case methodReplace:
		result, err := coll.ReplaceOne(bCtx, bson.M{keyName: keyValue}, value)
		if err != nil {
			activityLog.Debug("Error during replacing data ..", err)
			return false, err
		}

		activityLog.Debugf("Replace Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)

	case methodUpdate:
		result, err := coll.UpdateOne(bCtx, bson.M{keyName: keyValue}, bson.M{"$set": value})
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Update Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)
	default:
		activityLog.Errorf("unsupported method '%s'", method)
		return false, fmt.Errorf("unsupported method '%s'", method)
	}

	return true, nil
}
