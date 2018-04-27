package main

import (
	"context"
	"encoding/json"
	"flag"

	fl "github.com/TIBCOSoftware/flogo-contrib/trigger/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Handle implements the Flogo Function handler
func Handle(ctx context.Context, evt json.RawMessage) (interface{}, error) {
	err := setupArgs(evt, &ctx)
	if err != nil {
		return nil, err
	}

	// Get the actual lambda event object. Right now just using the evt typ to return the correct response.
	// This does not impact the event in the action (flow)
	evtTyp, _ := getEvtType(evt)

	// Invoke the flogo lambda trigger and handle the event
	result, err := fl.Invoke()
	if err != nil {
		return nil, err
	}

	return coherceResponseObj(result, evtTyp), nil
}

func getEvtType(raw json.RawMessage) (string, interface{}) {
	var evt map[string]interface{}
	json.Unmarshal(raw, &evt)

	if _, ok := evt["RequestContext"]; ok {
		apiGw := events.ApiGatewayProxyRequest{}
		json.Unmarshal(raw, &apiGw)
		return "apigw", apiGw
	}

	return "", nil
}

func coherceResponseObj(result map[string]interface{}, evtTyp string) interface{} {
	var returnObj interface{}

	responseData := result["data"]
	statusCode := result["status"].(int)

	// Marshal the response
	responseRaw, _ := json.Marshal(responseData)

	// Check if API GW request. If so, build the correct response
	switch evtTyp {
	case "apigw":
		returnObj = events.ApiGatewayProxyResponse{
			StatusCode: func() int {
				if statusCode == 0 {
					return 200
				} else {
					return statusCode
				}
			}(),
			Body:            string(responseRaw),
			IsBase64Encoded: false,
		}

	default:
		returnObj = responseData
	}

	return returnObj
}

func setupArgs(evt json.RawMessage, ctx *context.Context) error {
	// Setup environment argument
	evtJSON, err := json.Marshal(&evt)
	if err != nil {
		return err
	}

	evtFlag := flag.Lookup("evt")
	if evtFlag == nil {
		flag.String("evt", string(evtJSON), "Lambda Environment Arguments")
	} else {
		flag.Set("evt", string(evtJSON))
	}

	// Setup context argument
	ctxObj, _ := lambdacontext.FromContext(*ctx)
	lambdaContext := map[string]interface{}{
		"logStreamName":   lambdacontext.LogStreamName,
		"logGroupName":    lambdacontext.LogGroupName,
		"functionName":    lambdacontext.FunctionName,
		"functionVersion": lambdacontext.FunctionVersion,
		"awsRequestId":    ctxObj.AwsRequestID,
		"memoryLimitInMB": lambdacontext.MemoryLimitInMB,
	}
	ctxJSON, err := json.Marshal(lambdaContext)
	if err != nil {
		return err
	}

	ctxFlag := flag.Lookup("ctx")
	if ctxFlag == nil {
		flag.String("ctx", string(ctxJSON), "Lambda Context Arguments")
	} else {
		flag.Set("ctx", string(ctxJSON))
	}

	return nil
}

func main() {
	lambda.Start(Handle)
}
